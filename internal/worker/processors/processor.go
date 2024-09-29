package processors

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/models/events"
)

type Processor interface {
	Process(ctx context.Context, evt events.Event) error
}
