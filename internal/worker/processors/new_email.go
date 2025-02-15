package processors

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/domain/services"
)

func InitNewEmailProcessor() func(rCtx *RunningContext) ([]string, Processor) {
	return func(rCtx *RunningContext) ([]string, Processor) {
		return []string{events.EVENT_TYPE_NEW_EMAIL}, &NewEmailProcessor{
			eventBus:          rCtx.EventBus,
			emailStoreService: rCtx.EmailStoreService,
			emailer:           rCtx.EmailerService,
			schEvtRepo:        rCtx.ScheduledEventsRepo,
			rc:                rCtx.ResC,
		}
	}
}

type NewEmailProcessor struct {
	eventBus          events.EventBus
	emailStoreService services.EmailStoreService
	emailer           services.Emailer
	rc                *RescheduleConfig
	schEvtRepo        events.ScheduledEventRepository
}

func (nep *NewEmailProcessor) Process(ctx context.Context, evt events.Event) error {
	evtP, ok := evt.(*events.NewEmailEvent)
	if !ok {
		slog.Warn(fmt.Sprintf("Invalid event received by NewEmailProcessor. Type = %s", evt.GetType()))
		return nil
	}
	email, err := nep.emailStoreService.PrepareEmail(ctx, evtP.Id)
	if err != nil {
		return err
	}
	err = nep.emailer.SendEmail(ctx, email)
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending email with id %d", email.Id))
		if nep.rc != nil && nep.schEvtRepo != nil {

			nextRun := time.Now().Add(time.Duration(nep.rc.RetryTimeMs) * time.Millisecond)
			rescheduledEvent := events.CreateRescheduledEmailEvent(email.Id, nextRun)
			errReschedule := nep.schEvtRepo.Push(ctx, rescheduledEvent.When, rescheduledEvent)
			if errReschedule != nil {
				err = fmt.Errorf("could not reschedule email due do: %w, original error %w", errReschedule, err)
			}
		}
		return err
	}
	err = nep.emailStoreService.MarkEmailAsSent(ctx, email.Id)
	if err != nil {
		slog.Warn(fmt.Sprintf("Could not mark email as sent. Id: %d, err %s", email.Id, err.Error()))
	}
	return nil
}
