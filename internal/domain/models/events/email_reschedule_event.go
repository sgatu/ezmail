package events

import (
	"encoding/json"
	"time"
)

func CreateRescheduledEmailEvent(id int64, when time.Time) *RescheduledEmailEvent {
	return &RescheduledEmailEvent{
		Id:    id,
		When:  when,
		Tries: 1,
		Type:  EVENT_TYPE_RESCHEDULED_EMAIL,
	}
}

type RescheduledEmailEvent struct {
	When  time.Time `json:"when"`
	Type  string    `json:"type"`
	Tries int8      `json:"tries"`
	Id    int64     `json:"id,string"`
}

func (ree *RescheduledEmailEvent) Serialize() (string, error) {
	result, err := json.Marshal(ree)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (ree *RescheduledEmailEvent) GetType() string {
	return ree.Type
}
