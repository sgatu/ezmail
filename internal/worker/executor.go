package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/sgatu/ezmail/internal/worker/processors"
)

type Executor struct {
	wg         *sync.WaitGroup
	runningCtx *processors.RunningContext
	processors map[string][]processors.Processor
	running    bool
}

func NewExecutor(
	runningCtx *processors.RunningContext,
	wg *sync.WaitGroup,
	processorsList ...func(rCtx *processors.RunningContext) ([]string, processors.Processor),
) *Executor {
	e := Executor{
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

func (e *Executor) Run() {
	e.running = true
	e.wg.Add(1)
	slog.Info("Starting executor")
	go func() {
		anyProcessed := true
		var cancel context.CancelFunc
		var ctx context.Context
		for e.running {
			if cancel != nil {
				cancel()
			}
			if anyProcessed {
				time.Sleep(1000 * time.Millisecond)
			} else {
				slog.Debug("Looking for messages, none in last loop, 5s sleep")
				time.Sleep(5000 * time.Millisecond)
			}
			anyProcessed = false
			ctx, cancel = context.WithTimeout(context.Background(), 1500*time.Millisecond)
			msgs, err := e.runningCtx.BusReader.Read(ctx, 1)
			if err != nil {
				if !errors.Is(err, os.ErrDeadlineExceeded) {
					slog.Warn(fmt.Errorf("could not retrieve messages from queue due to %w", err).Error())
				}
				continue
			}
			if len(msgs) == 0 {
				continue
			}
			eventWrapper := msgs[0]
			processors, ok := e.processors[eventWrapper.Event.GetType()]
			if !ok {
				slog.Warn(fmt.Sprintf("No processor found for %s\n", eventWrapper.Event.GetType()))
				// commit and ignore, no point on retrying the same event
				e.runningCtx.BusReader.Commit(ctx, eventWrapper.Id)
				continue
			}
			anyProcessed = true
			err = nil
			lastId := ""
			for _, proc := range processors {
				err = proc.Process(ctx, eventWrapper.Event)
				if err != nil {
					slog.Warn(fmt.Sprintf("Failed to process event %s/%s due to %s", eventWrapper.Id, eventWrapper.Event.GetType(), err.Error()))
					break
				}
				lastId = eventWrapper.Id
			}
			if len(lastId) != 0 {
				err = e.runningCtx.BusReader.Commit(ctx, lastId)
				if err != nil {
					slog.Error("Could not commit, closing executor")
					e.running = false
				}
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

func (e *Executor) Stop() {
	slog.Info("Shutting down executor")
	e.running = false
}
