package events

import (
	"encoding/json"
	"time"
)

func CreateNewEmailEvent(id int64) *NewEmailEvent {
	return &NewEmailEvent{
		Id:   id,
		When: time.Now().UTC(),
		Type: "new_email",
	}
}

type NewEmailEvent struct {
	When time.Time `json:"when"`
	Type string    `json:"type"`
	Id   int64     `json:"id,string"`
}

func (nee *NewEmailEvent) Serialize() (string, error) {
	result, err := json.Marshal(nee)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (nee *NewEmailEvent) GetType() string {
	return nee.Type
}
