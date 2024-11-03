package eventbus

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/thirdparty"
)

type RedisEventBus struct {
	redisConnection thirdparty.BaseRedisClient
	eventsTopic     string
	maxLen          int64
}

func NewRedisEventBus(redisConn thirdparty.BaseRedisClient, maxLen int64, eventsTopic string) *RedisEventBus {
	return &RedisEventBus{
		redisConnection: redisConn,
		maxLen:          maxLen,
		eventsTopic:     eventsTopic,
	}
}

func (ce *RedisEventBus) Push(ctx context.Context, event events.Event) error {
	eventData, err := event.Serialize()
	if err != nil {

		slog.Warn(fmt.Sprintf("Could not serialize event. Type = %s, Err = %s", event.GetType(), err))
		return err
	}
	result := ce.redisConnection.XAdd(ctx, &redis.XAddArgs{
		Stream: ce.eventsTopic,
		Values: []string{"payload", eventData},
		MaxLen: ce.maxLen,
	})
	if result.Err() != nil {
		slog.Warn(fmt.Sprintf("Could not send event to queue. Type = %s, Err = %s, Queue = %s", event.GetType(), err, ce.eventsTopic))
		return result.Err()
	}
	return nil
}
