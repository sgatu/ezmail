package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Event interface {
	Serialize() (string, error)
	GetType() string
}

type EventBus interface {
	Push(ctx context.Context, e Event) error
}
type BusReader interface {
	Read(ctx context.Context, limit int32) ([]EventWrapper, error)
	Commit(ctx context.Context, commitInfo interface{}) error
}
type TypedEvent struct {
	Type string `json:"type"`
}
type EventWrapper struct {
	Event Event
	Id    string
}

const (
	EVENT_TYPE_NEW_EMAIL             = "new_email"
	EVENT_TYPE_RESCHEDULED_EMAIL     = "rescheduled_email"
	EVENT_TYPE_DOMAIN_REGISTER       = "domain_register"
	EVENT_TYPE_REFRESH_DOMAIN_STATUS = "refresh_domain_status"
)

type ScheduledEventRepository interface {
	Push(ctx context.Context, when time.Time, evt Event) error
	GetNext(ctx context.Context) (Event, error)
	RemoveNext(ctx context.Context) error
}

func RetrieveTypedEvent(eventData []byte) (Event, error) {
	var tEvent TypedEvent
	err := json.Unmarshal(eventData, &tEvent)
	if err != nil {
		return nil, err
	}
	var resultEvent Event
	switch tEvent.Type {
	case EVENT_TYPE_NEW_EMAIL:
		resultEvent = &NewEmailEvent{}
	case EVENT_TYPE_RESCHEDULED_EMAIL:
		resultEvent = &RescheduledEmailEvent{}
	case EVENT_TYPE_DOMAIN_REGISTER:
		resultEvent = &DomainRegisterEvent{}
	case EVENT_TYPE_REFRESH_DOMAIN_STATUS:
		resultEvent = &RefreshDomainEvent{}
	default:
		err = fmt.Errorf("could not parse event - unknown type")
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(eventData, resultEvent)
	if err != nil {
		return nil, err
	}
	return resultEvent, nil
}
