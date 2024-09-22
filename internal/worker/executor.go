package worker

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/domain/services"
	"github.com/sgatu/ezmail/internal/worker/processors"
)

type executor struct {
	emailStoreService   services.EmailStoreService
	emailerService      services.Emailer
	scheduledEventsRepo events.ScheduledEventRepository
	busReader           events.BusReader
	processors          []processors.Processor
	running             bool
}

func NewExecutor(
	emailStoreService services.EmailStoreService,
	emailerService services.Emailer,
	schedulerEventsRepo events.ScheduledEventRepository,
	busReader events.BusReader,
	processorsList ...func() processors.Processor,
) *executor {
	e := executor{
		emailStoreService:   emailStoreService,
		emailerService:      emailerService,
		scheduledEventsRepo: schedulerEventsRepo,
		processors:          make([]processors.Processor, len(processorsList)),
	}
	for _, pDef := range processorsList {
		proc := pDef()
		e.processors = append(e.processors, proc)
	}
	return &e
}

func (e *executor) Run() {
	e.running = true
	go func() {
		for e.running {
			e.busReader.Read(context.Background(), 1)
		}
	}()
}

func (e *executor) Stop() {
	e.running = false
}
