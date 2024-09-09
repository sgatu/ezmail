package events

import (
	"encoding/json"
	"time"
)

type NewEmailEvent struct {
	When time.Time `json:"when"`
	TypedEvent
	Id int64 `json:"id"`
}

func (nee *NewEmailEvent) Serialize() (string, error) {
	result, err := json.Marshal(nee)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
