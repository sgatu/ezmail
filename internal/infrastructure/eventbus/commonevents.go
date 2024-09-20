package eventbus

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/sgatu/ezmail/internal/domain/models/events"
)

type CommonEventsEventBus struct {
	redisConnection *redis.Client
}

func NewCommonEventsEventBus(redisConn *redis.Client) *CommonEventsEventBus {
	return &CommonEventsEventBus{
		redisConnection: redisConn,
	}
}

func (ce *CommonEventsEventBus) Push(ctx context.Context, event events.Event, queue string) error {
	eventData, err := event.Serialize()
	if err != nil {

		slog.Warn(fmt.Sprintf("Could not serialize event. Type = %s, Err = %s", event.GetType(), err))
		return err
	}
	result := ce.redisConnection.XAdd(ctx, &redis.XAddArgs{
		Stream: queue,
		Values: []string{"payload", eventData},
	})
	if result.Err() != nil {
		slog.Warn(fmt.Sprintf("Could not send event to queue. Type = %s, Err = %s, Queue = %s", event.GetType(), err, queue))
		return result.Err()
	}
	return nil
}

type CommonEventsBusReader struct {
	redisConnection *redis.Client
}
