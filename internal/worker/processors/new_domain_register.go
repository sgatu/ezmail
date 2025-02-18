package processors

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sgatu/ezmail/internal/domain/models/events"
)

func InitNewDomainRegisterProcessor() func(rCtx *RunningContext) ([]string, Processor) {
	return func(rCtx *RunningContext) ([]string, Processor) {
		return []string{events.EVENT_TYPE_DOMAIN_REGISTER}, &NewDomainRegisterProcessor{
			sch:  rCtx.ScheduledEventsRepo,
			refC: rCtx.RefC,
		}
	}
}

type NewDomainRegisterProcessor struct {
	sch  events.ScheduledEventRepository
	refC *RefreshConfig
}

func (ndrp *NewDomainRegisterProcessor) Process(ctx context.Context, evt events.Event) error {
	evtP, ok := evt.(*events.DomainRegisterEvent)
	if !ok {
		slog.Warn(fmt.Sprintf("Invalid event received by NewDomainRegisterProcessor. Type = %s", evt.GetType()))
		return nil
	}
	ndrp.sch.Push(
		ctx,
		time.Now().Add(time.Duration(ndrp.refC.RetryDelaySec)*time.Second).UTC(),
		events.NewRefreshDomainEvent(evtP.DomainId, ndrp.refC.MaxRetries, ndrp.refC.RetryDelaySec),
	)
	return nil
}
