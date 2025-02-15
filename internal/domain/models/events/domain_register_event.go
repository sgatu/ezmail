package events

import (
	"encoding/json"
	"time"
)

type DomainRegisterEvent struct {
	DomainId int64
	Type     string
	When     time.Time
}

func NewDomainRegisterEvent(domainId int64) *DomainRegisterEvent {
	return &DomainRegisterEvent{
		DomainId: domainId,
		Type:     EVENT_TYPE_DOMAIN_REGISTER,
		When:     time.Now().UTC(),
	}
}

func (dre *DomainRegisterEvent) Serialize() (string, error) {
	result, err := json.Marshal(dre)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (dre *DomainRegisterEvent) GetType() string {
	return dre.Type
}
