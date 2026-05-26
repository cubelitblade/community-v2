package events

import "context"

// MessageHandler processes a single event delivery.
type MessageHandler interface {
	Handle(ctx context.Context, routingKey string, body []byte) error
}
