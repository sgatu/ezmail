package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/sgatu/ezmail/internal/domain/models/events"
)

type Scheduler struct {
	schedRepo events.ScheduledEventRepository
	evBus     events.EventBus
	wg        *sync.WaitGroup
	running   bool
}

func NewScheduler(repo events.ScheduledEventRepository, evBus events.EventBus, wg *sync.WaitGroup) *Scheduler {
	return &Scheduler{
		schedRepo: repo,
		running:   false,
		evBus:     evBus,
		wg:        wg,
	}
}

func (sched *Scheduler) Run() {
	sched.running = true
	sched.wg.Add(1)
	slog.Info("Starting scheduler")
	go func() {
		var cancel context.CancelFunc
		var ctx context.Context
		for sched.running {
			if cancel != nil {
				cancel()
			}
			time.Sleep(5000 * time.Millisecond)
			ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Millisecond)
			next, err := sched.schedRepo.GetNext(ctx)
			if err != nil {
				if !errors.Is(err, os.ErrDeadlineExceeded) {
					slog.Warn(fmt.Errorf("could not get next scheduled event due to %w", err).Error())
				}
				continue
			}
			if next == nil {
				continue
			}
			slog.Info(fmt.Sprintf("Found scheduled event %s, sending to queue", next.GetType()))
			sched.evBus.Push(ctx, next)
			sched.schedRepo.RemoveNext(ctx)
		}
		if cancel != nil {
			cancel()
		}
		sched.wg.Done()
	}()
}

func (sched *Scheduler) Stop() {
	slog.Info("Shutting down scheduler")
	sched.running = false
}
