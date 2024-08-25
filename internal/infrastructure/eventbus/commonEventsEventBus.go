package eventbus

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sgatu/ezmail/internal/domain/models/events"
)

type CommonEventsEventBus struct {
	redisConnection *redis.Conn
}

func (ce *CommonEventsEventBus) Push(ctx context.Context, event events.Event, queue string) error {
	eventData, err := event.Serialize()
	if err != nil {
		return err
	}
	result := ce.redisConnection.XAdd(ctx, &redis.XAddArgs{
		Stream: queue,
		Values: [2]string{"payload", eventData},
	})
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
