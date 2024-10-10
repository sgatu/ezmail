package worker

import (
	"context"
	"fmt"
	"log/slog"
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

func NewScheduler(repo events.ScheduledEventRepository, wg *sync.WaitGroup) *Scheduler {
	return &Scheduler{
		schedRepo: repo,
		running:   false,
		wg:        wg,
	}
}

func (sched *Scheduler) Run() {
	sched.running = true
	sched.wg.Add(1)
	go func() {
		var cancel context.CancelFunc
		var ctx context.Context
		for sched.running {
			if cancel != nil {
				cancel()
			}
			time.Sleep(20 * time.Millisecond)
			ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Millisecond)
			next, err := sched.schedRepo.GetNext(ctx)
			if err != nil {
				slog.Warn(fmt.Errorf("could not get next scheduled event due to %w", err).Error())
				continue
			}
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
	sched.running = false
}
