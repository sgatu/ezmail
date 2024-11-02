package test

type EventInvalid struct{}

func (ei *EventInvalid) GetType() string {
	return "invalid"
}

func (ei *EventInvalid) Serialize() (string, error) {
	return "", nil
}
