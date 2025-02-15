package events

import (
	"encoding/json"
	"time"
)

type RefreshDomainEvent struct {
	DomainId           int64
	Type               string
	MaxRetries         int64
	CurrentRetries     int64
	LastRetry          time.Time
	TimeBetweenRetries int64
}

func NewRefreshDomainEvent(domaintId int64, maxRetries int64, timeBetweenRetries int64) *RefreshDomainEvent {
	return &RefreshDomainEvent{
		DomainId:           domaintId,
		MaxRetries:         maxRetries,
		TimeBetweenRetries: timeBetweenRetries,
		LastRetry:          time.Now().UTC(),
		CurrentRetries:     0,
		Type:               EVENT_TYPE_REFRESH_DOMAIN_STATUS,
	}
}

func (rde *RefreshDomainEvent) PrepareNext() bool {
	if rde.CurrentRetries+1 >= rde.MaxRetries {
		return false
	}
	rde.CurrentRetries++
	rde.LastRetry = time.Now().UTC()
	return true
}

func (rde *RefreshDomainEvent) GetType() string {
	return rde.Type
}

func (rde *RefreshDomainEvent) Serialize() (string, error) {
	result, err := json.Marshal(rde)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
