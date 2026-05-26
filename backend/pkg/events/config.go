// Package events provides generic RabbitMQ publish and consume infrastructure.
package events

import "time"

// Config holds RabbitMQ connection and topology settings.
type Config struct {
	URL            string
	ExchangeName   string
	QueueName      string
	ConsumerTag    string
	Bindings       []string
	MessageTimeout time.Duration
	PrefetchCount  int
	DeadLetterExchangeName string
	DeadLetterQueueName    string
	MaxDeliveryCount       int
}
