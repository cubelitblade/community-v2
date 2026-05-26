package events

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	dialTimeout        = 10 * time.Second
	backoffMaxInterval = 30 * time.Second
	defaultTimeout     = 30 * time.Second
)

// Consumer consumes domain events from RabbitMQ and dispatches them to a MessageHandler.
type Consumer struct {
	cfg     *Config
	handler MessageHandler
	logger  *slog.Logger

	cancel  context.CancelFunc
	mu      sync.Mutex
	running bool
	wg      sync.WaitGroup
}

// closeConn closes an AMQP connection, ignoring any error.
func closeConn(conn *amqp.Connection, logger *slog.Logger) {
	if conn == nil {
		return
	}

	if err := conn.Close(); err != nil {
		logger.Warn("Failed to close connection", "error", err)
	}
}

// closeChannel closes an AMQP channel, ignoring any error.
func closeChannel(channel *amqp.Channel, logger *slog.Logger) {
	if channel == nil {
		return
	}

	if err := channel.Close(); err != nil {
		logger.Warn("Failed to close channel", "error", err)
	}
}

// NewConsumer creates a RabbitMQ Consumer.
func NewConsumer(cfg *Config, handler MessageHandler, logger *slog.Logger) *Consumer {
	if cfg == nil {
		panic("nil config")
	}

	if handler == nil {
		panic("nil handler")
	}

	if logger == nil {
		panic("nil logger")
	}

	return &Consumer{
		cfg:     cfg,
		handler: handler,
		logger:  logger,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return fmt.Errorf("consumer already running")
	}
	c.running = true
	c.mu.Unlock()

	ctx, c.cancel = context.WithCancel(ctx)
	c.wg.Add(1)
	go c.run(ctx)
	return nil
}

func (c *Consumer) Stop() {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		c.logger.Warn("Consumer not running, stop is ignored")
		return
	}
	c.running = false
	c.mu.Unlock()

	if c.cancel != nil {
		c.cancel()
	}

	c.wg.Wait()

	c.logger.Info("Consumer stopped")
}

func (c *Consumer) run(ctx context.Context) {
	defer c.wg.Done()

	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 1 * time.Second
	b.MaxInterval = backoffMaxInterval
	b.MaxElapsedTime = 5 * time.Minute

	for {
		err := c.connectAndConsume(ctx)
		if err != nil {
			c.logger.Error("consumer loop exited", "error", err)
		}

		if ctx.Err() != nil {
			c.logger.Info("consumer loop exiting due to cancellation")
			return
		}

		b.Reset()
		for {
			wait := b.NextBackOff()
			if wait == backoff.Stop {
				c.logger.Error("reconnect retries exhausted")
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(wait):
			}

			err := c.connectAndConsume(ctx)
			if err == nil {
				c.logger.Info("reconnected successfully")
				break
			}
			c.logger.Error("reconnect failed, will retry", "error", err)

			if ctx.Err() != nil {
				return
			}
		}
	}
}

func (c *Consumer) connectAndConsume(ctx context.Context) error {
	conn, err := amqp.DialConfig(c.cfg.URL, amqp.Config{
		Dial: amqp.DefaultDial(dialTimeout),
	})
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}
	defer closeConn(conn, c.logger)

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("conn.Channel: %w", err)
	}
	defer closeChannel(channel, c.logger)

	if err := c.declareTopology(channel); err != nil {
		return err
	}

	deliveries, err := channel.Consume(
		c.cfg.QueueName,
		c.cfg.ConsumerTag,
		false, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("start consuming: %w", err)
	}

	c.logger.Info("connected and consuming",
		"queue", c.cfg.QueueName,
	)

	c.consumeLoop(ctx, channel, deliveries)
	return nil
}

func (c *Consumer) consumeLoop(ctx context.Context, channel *amqp.Channel, deliveries <-chan amqp.Delivery) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("consumer panicked", "recover", r)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("context canceled")
			return
		case d, ok := <-deliveries:
			if !ok {
				c.logger.Info("deliveries channel closed")
				return
			}
			c.handleDelivery(ctx, channel, d)
		}
	}
}

func (c *Consumer) handleDelivery(ctx context.Context, channel *amqp.Channel, delivery amqp.Delivery) {
	c.logger.Debug("received event", "routing_key", delivery.RoutingKey)

	timeout := c.cfg.MessageTimeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	msgCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := c.handler.Handle(msgCtx, delivery.RoutingKey, delivery.Body); err != nil {
		c.logger.Error("handle event failed", "routing_key", delivery.RoutingKey, "error", err)

		c.handleProcessingFailure(channel, delivery)
		return
	}
	if err := delivery.Ack(false); err != nil {
		c.logger.Warn("failed to ack", "routing_key", delivery.RoutingKey, "error", err)
	}
}

func (c *Consumer) publishForRetry(ch *amqp.Channel, delivery amqp.Delivery, retryCount int) error {
	headers := delivery.Headers
	if headers == nil {
		headers = amqp.Table{}
	}
	headers["x-retry-count"] = retryCount + 1

	return ch.Publish(
		c.cfg.ExchangeName,
		delivery.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  delivery.ContentType,
			Body:         delivery.Body,
			DeliveryMode: amqp.Persistent,
			Headers:      headers,
		},
	)
}

func (c *Consumer) declareTopology(channel *amqp.Channel) error {
	if err := channel.ExchangeDeclare(c.cfg.DeadLetterExchangeName, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare dead letter exchange: %w", err)
	}

	dlq, err := channel.QueueDeclare(c.cfg.DeadLetterQueueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare dead letter queue: %w", err)
	}

	if err := channel.QueueBind(dlq.Name, "#", c.cfg.DeadLetterExchangeName, false, nil); err != nil {
		return fmt.Errorf("bind dead letter queue: %w", err)
	}

	if err := channel.ExchangeDeclare(c.cfg.ExchangeName, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	queue, err := channel.QueueDeclare(c.cfg.QueueName, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange": c.cfg.DeadLetterExchangeName,
	})
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	for _, key := range c.cfg.Bindings {
		if err := channel.QueueBind(queue.Name, key, c.cfg.ExchangeName, false, nil); err != nil {
			return fmt.Errorf("bind %s: %w", key, err)
		}
	}

	if err := channel.Qos(c.cfg.PrefetchCount, 0, false); err != nil {
		return fmt.Errorf("set QoS: %w", err)
	}

	return nil
}

func (c *Consumer) handleProcessingFailure(channel *amqp.Channel, delivery amqp.Delivery) {
	retryCount := getRetryCount(delivery.Headers)

	maxDeliveryCount := c.cfg.MaxDeliveryCount
	if maxDeliveryCount <= 0 {
		maxDeliveryCount = 3
	}

	if retryCount < maxDeliveryCount {
		c.retryMessage(channel, delivery, retryCount)
		return
	}

	c.logger.Warn("message exceeded max retries, sending to DLQ", "routing_key", delivery.RoutingKey)
	if err := delivery.Nack(false, false); err != nil {
		c.logger.Warn("failed to nack", "routing_key", delivery.RoutingKey, "error", err)
	}
}

func (c *Consumer) retryMessage(channel *amqp.Channel, delivery amqp.Delivery, retryCount int) {
	if err := c.publishForRetry(channel, delivery, retryCount); err != nil {
		c.logger.Error("failed to republish for retry", "error", err)
		if nackErr := delivery.Nack(false, true); nackErr != nil {
			c.logger.Warn("failed to nack", "routing_key", delivery.RoutingKey, "error", nackErr)
		}
		return
	}
	if err := delivery.Ack(false); err != nil {
		c.logger.Warn("failed to ack after republish", "routing_key", delivery.RoutingKey, "error", err)
	}
}

func getRetryCount(headers amqp.Table) int {
	val, ok := headers["x-retry-count"]
	if !ok {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}
