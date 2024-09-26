package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sgatu/ezmail/internal/worker/processors"
)

type executor struct {
	runningCtx *processors.RunningContext
	processors map[string][]processors.Processor
	running    bool
}

func NewExecutor(
	runningCtx *processors.RunningContext,
	processorsList ...func(rCtx *processors.RunningContext) ([]string, processors.Processor),
) *executor {
	e := executor{
		runningCtx: runningCtx,
		processors: make(map[string][]processors.Processor),
		running:    true,
	}
	for _, pDef := range processorsList {
		types, proc := pDef(e.runningCtx)
		for _, t := range types {
			if val, ok := e.processors[t]; ok {
				e.processors[t] = append(val, proc)
			} else {
				e.processors[t] = append(make([]processors.Processor, 0), proc)
			}
		}
	}
	return &e
}

func (e *executor) Run() {
	e.running = true
	go func() {
		for e.running {
			time.Sleep(50 * time.Millisecond)
			ctx := context.Background()
			msgs, err := e.runningCtx.BusReader.Read(ctx, 1)
			if err != nil {
				slog.Warn(fmt.Errorf("could not retrieve messages from queue due to %w", err).Error())
				continue
			}
			if len(msgs) == 0 {
				continue
			}
			event := msgs[0]
			processors, ok := e.processors[event.GetType()]
			if !ok {
				continue
			}
			for _, proc := range processors {
				_ = proc.Process(ctx, event) // ignore error for now, processor should deal with it and at least log it
			}

		}
	}()
}

func (e *executor) Stop() {
	e.running = false
}
