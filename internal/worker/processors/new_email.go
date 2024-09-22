package processors

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/domain/services"
)

type RescheduleConfig struct {
	Retries          int8
	RetryTimeSeconds int32
}
type NewEmailProcessor struct {
	eventBus            events.EventBus
	emailStoreService   services.EmailStoreService
	emailer             services.Emailer
	rescheduleConfig    *RescheduleConfig
	scheduledEventsRepo events.ScheduledEventRepository
}

func (nep *NewEmailProcessor) Process(evt events.Event) error {
	ctx := context.Background()
	evtP, ok := evt.(*events.NewEmailEvent)
	if !ok {
		slog.Warn(fmt.Sprintf("Invalid event received by NewEmailProcessor. Type = %s", evt.GetType()))
		return nil
	}
	email, err := nep.emailStoreService.PrepareEmail(ctx, evtP.Id)
	if err != nil {
		return err
	}
	err = nep.emailer.SendEmail(context.Background(), email)
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending email with id %d", email.Id))
		if nep.rescheduleConfig != nil {
			rescheduledEvent := events.CreateRescheduledEmailEvent(email.Id, time.Now().Add(time.Duration(nep.rescheduleConfig.RetryTimeSeconds)*time.Second))
			errReschedule := nep.scheduledEventsRepo.Push(ctx, rescheduledEvent.When, rescheduledEvent)
			if errReschedule != nil {
				err = fmt.Errorf("could not reschedule email due do: %w, original error %w", errReschedule, err)
			}
		}
		return err
	}
	return nil
}
