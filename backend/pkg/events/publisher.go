package events

import (
	"context"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher sends domain events to a RabbitMQ exchange.
type Publisher struct {
	cfg     *Config
	logger  *slog.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewPublisher creates a RabbitMQ Publisher, connects, and declares the exchange.
func NewPublisher(cfg *Config, logger *slog.Logger) (*Publisher, error) {
	if cfg == nil {
		panic("nil config")
	}

	if logger == nil {
		panic("nil logger")
	}

	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		closeConn(conn, logger)

		return nil, fmt.Errorf("open channel: %w", err)
	}

	if err := channel.ExchangeDeclare(cfg.ExchangeName, "topic", true, false, false, false, nil); err != nil {
		closeChannel(channel, logger)
		closeConn(conn, logger)

		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	return &Publisher{
		cfg:     cfg,
		logger:  logger,
		conn:    conn,
		channel: channel,
	}, nil
}

// Publish sends an event with the given routing key and body to the exchange.
func (p *Publisher) Publish(ctx context.Context, routingKey string, body []byte) error {
	err := p.channel.PublishWithContext(ctx, p.cfg.ExchangeName, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		return fmt.Errorf("publish %s: %w", routingKey, err)
	}

	return nil
}

// Close tears down the AMQP channel and connection.
func (p *Publisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return fmt.Errorf("close channel: %w", err)
		}
	}

	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			return fmt.Errorf("close connection: %w", err)
		}
	}

	return nil
}
