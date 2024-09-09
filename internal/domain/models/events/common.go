package events

import "context"

type Event interface {
	Serialize() (string, error)
}
type TypedEvent struct {
	Type string `json:"type"`
}

func (te *TypedEvent) GetType() string {
	return te.Type
}

type EventBus interface {
	Push(ctx context.Context, e Event, queue string) error
}
