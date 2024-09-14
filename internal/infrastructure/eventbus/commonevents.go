package eventbus

import (
	"context"
	"fmt"

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
		fmt.Println("err1")

		fmt.Println(err)
		return err
	}
	result := ce.redisConnection.XAdd(ctx, &redis.XAddArgs{
		Stream: queue,
		Values: []string{"payload", eventData},
	})
	if result.Err() != nil {
		fmt.Println(result.Err())
		return result.Err()
	}
	return nil
}
