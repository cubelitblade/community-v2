package events

import "errors"

// ErrUnknownEvent is returned when an unrecognized routing key is received.
var ErrUnknownEvent = errors.New("unknown event type")
