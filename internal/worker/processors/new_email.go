package processors

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/domain/services"
)

type NewEmailProcessor struct {
	eventBus          events.EventBus
	emailStoreService services.EmailStoreService
	emailer           services.Emailer
	rescheduleFailed  bool
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
		if nep.rescheduleFailed {
			rescheduledEvent := events.CreateRescheduledEmailEvent(email.Id, time.Now().Add(5*time.Minute))
			rescheduledEvent.Id = 1
		}
		return err
	}
	return nil
}
