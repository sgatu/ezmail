package events

type Event interface {
	Serialize() (string, error)
}
type EventBus interface {
	Push(e Event, queue string) error
}
