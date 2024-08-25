package events

import "encoding/json"

type UserCreatedEvent struct {
	EventType string `json:"event_type"`
	When      int64  `json:"when"`
	UserId    int64  `json:"user_id"`
}

func (uce *UserCreatedEvent) Serialize() (string, error) {
	data, err := json.Marshal(uce)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
