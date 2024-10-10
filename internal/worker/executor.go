package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/sgatu/ezmail/internal/worker/processors"
)

type executor struct {
	wg         *sync.WaitGroup
	runningCtx *processors.RunningContext
	processors map[string][]processors.Processor
	running    bool
}

func NewExecutor(
	runningCtx *processors.RunningContext,
	wg *sync.WaitGroup,
	processorsList ...func(rCtx *processors.RunningContext) ([]string, processors.Processor),
) *executor {
	e := executor{
		runningCtx: runningCtx,
		processors: make(map[string][]processors.Processor),
		running:    true,
		wg:         wg,
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
	e.wg.Add(1)
	go func() {
		anyProcessed := true
		var cancel context.CancelFunc
		var ctx context.Context
		for e.running {
			if cancel != nil {
				cancel()
			}
			if anyProcessed {
				time.Sleep(50 * time.Millisecond)
			} else {
				slog.Debug("Looking for messages, none in last loop, 5s sleep")
				time.Sleep(5000 * time.Millisecond)
			}
			anyProcessed = false
			ctx, cancel = context.WithTimeout(context.Background(), 1500*time.Millisecond)
			msgs, err := e.runningCtx.BusReader.Read(ctx, 1)
			if err != nil {
				slog.Warn(fmt.Errorf("could not retrieve messages from queue due to %w", err).Error())
				continue
			}
			if len(msgs) == 0 {
				continue
			}
			eventWrapper := msgs[0]
			processors, ok := e.processors[eventWrapper.Event.GetType()]
			if !ok {
				continue
			}
			anyProcessed = true
			err = nil
			var lastId *string = nil
			for _, proc := range processors {
				err = proc.Process(ctx, eventWrapper.Event)
				if err != nil {
					break
				}
				lastId = &eventWrapper.Id
			}
			if lastId != nil {
				e.runningCtx.BusReader.Commit(ctx, &lastId)
			} else {
				anyProcessed = false
			}
		}
		if cancel != nil {
			cancel()
		}
		e.wg.Done()
	}()
}

func (e *executor) Stop() {
	slog.Info("Shutting down executor...")
	e.running = false
}
