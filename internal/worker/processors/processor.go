package processors

import "github.com/sgatu/ezmail/internal/domain/models/events"

type Processor interface {
	Process(evt events.Event) error
}
