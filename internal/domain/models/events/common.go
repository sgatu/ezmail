package events

import "context"

type Event interface {
	Serialize() (string, error)
	GetType() string
}

type EventBus interface {
	Push(ctx context.Context, e Event, queue string) error
}
type BusReader interface {
	Read(ctx context.Context, limit int32) error
	Commit(ctx context.Context, commitInfo interface{}) error
}
