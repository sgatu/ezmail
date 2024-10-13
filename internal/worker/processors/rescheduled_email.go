package processors

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/domain/services"
)

func InitRescheduledEmailProcessor() func(rCtx *RunningContext) ([]string, Processor) {
	return func(rCtx *RunningContext) ([]string, Processor) {
		return []string{events.EVENT_TYPE_RESCHEDULED_EMAIL}, &RescheduledEmailProcessor{
			eventBus:          rCtx.EventBus,
			emailStoreService: rCtx.EmailStoreService,
			emailer:           rCtx.EmailerService,
			schEvtRepo:        rCtx.ScheduledEventsRepo,
			rc:                rCtx.Rc,
		}
	}
}

type RescheduledEmailProcessor struct {
	eventBus          events.EventBus
	emailStoreService services.EmailStoreService
	emailer           services.Emailer
	rc                *RescheduleConfig
	schEvtRepo        events.ScheduledEventRepository
}

func (rep *RescheduledEmailProcessor) Process(ctx context.Context, evt events.Event) error {
	evtP, ok := evt.(*events.RescheduledEmailEvent)
	if !ok {
		slog.Warn(fmt.Sprintf("Invalid event received by RescheduledEmailProcessor. Type = %s", evt.GetType()))
		return nil
	}
	email, err := rep.emailStoreService.PrepareEmail(ctx, evtP.Id)
	if err != nil {
		return err
	}
	err = rep.emailer.SendEmail(ctx, email)
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending rescheduled email with id %d", email.Id))
		if rep.rc != nil && evtP.Tries <= rep.rc.Retries && rep.schEvtRepo != nil {
			nextRun := time.Now().Add(time.Duration(rep.rc.RetryTimeMs) * time.Millisecond)
			evtP.When = nextRun
			evtP.Tries += 1
			errReschedule := rep.schEvtRepo.Push(ctx, nextRun, evtP)
			if errReschedule != nil {
				err = fmt.Errorf("could not reschedule email due do: %w, original error %w", errReschedule, err)
			}
		}
		return err
	}
	err = rep.emailStoreService.MarkEmailAsSent(ctx, evtP.Id)
	if err != nil {
		slog.Warn(fmt.Sprintf("Could not mark email as sent. Id: %d", email.Id))
	}

	return nil
}
