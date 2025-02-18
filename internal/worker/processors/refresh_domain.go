package processors

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/domain/services"
)

func InitRefreshDomainProcessor() func(rCtx *RunningContext) ([]string, Processor) {
	return func(rCtx *RunningContext) ([]string, Processor) {
		return []string{events.EVENT_TYPE_REFRESH_DOMAIN_STATUS}, &RefreshDomainProcessor{
			sch:          rCtx.ScheduledEventsRepo,
			diRepository: rCtx.DomainInfoRepository,
			identityMgr:  rCtx.IdentityManager,
		}
	}
}

type RefreshDomainProcessor struct {
	sch          events.ScheduledEventRepository
	diRepository domain.DomainInfoRepository
	identityMgr  services.IdentityManager
}

func (ndrp *RefreshDomainProcessor) Process(ctx context.Context, evt events.Event) error {
	evtP, ok := evt.(*events.RefreshDomainEvent)
	if !ok {
		slog.Warn(fmt.Sprintf("Invalid event received by RefreshDomainProcessor. Type = %s", evt.GetType()))
		return nil
	}
	di, err := ndrp.diRepository.GetDomainInfoById(ctx, evtP.DomainId)
	if err != nil {
		slog.Warn(fmt.Sprintf("Could not retrieve domain info by id %d, RefreshDomainProcessor", evtP.DomainId))
		return err
	}
	err = ndrp.identityMgr.RefreshIdentity(ctx, di)
	if err != nil {
		slog.Warn(fmt.Sprintf("Could not refresh domain status due to %s", err))
	}
	if err == nil {
		err = ndrp.diRepository.Save(ctx, di)
	}
	if !di.Validated {
		ok := evtP.PrepareNext()
		if ok {
			ndrp.sch.Push(
				ctx,
				time.Now().Add(time.Duration(evtP.TimeBetweenRetries)*time.Second).UTC(),
				evtP,
			)
		}
	}
	return err
}
