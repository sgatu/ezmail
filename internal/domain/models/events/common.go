package events

import "context"

type Event interface {
	Serialize() (string, error)
}
type EventBus interface {
	Push(ctx context.Context, e Event, queue string) error
}
