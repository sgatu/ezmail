package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sgatu/ezmail/internal/domain/models/events"
)

type RedisScheduledEventRepository struct {
	conn *redis.Client
}

const (
	redisKey = "q_scheduled_events"
)

func (repo *RedisScheduledEventRepository) Push(ctx context.Context, when time.Time, evt events.Event) error {
	evtData, err := evt.Serialize()
	if err != nil {
		return err
	}
	result := repo.conn.ZAdd(ctx, redisKey, redis.Z{Score: float64(when.Unix()), Member: evtData})
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (repo *RedisScheduledEventRepository) GetNext(ctx context.Context) (events.Event, error) {
	result := repo.conn.ZRangeByScore(ctx, redisKey, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   strconv.FormatInt(time.Now().Unix(), 10),
		Count: 1,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	values := result.Val()
	if len(values) == 0 {
		return nil, nil
	}
	return events.RetrieveTypedEvent([]byte(values[0]))
}

func (repo *RedisScheduledEventRepository) RemoveNext(ctx context.Context) error {
	result := repo.conn.ZRangeByScore(ctx, redisKey, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   strconv.FormatInt(time.Now().Unix(), 10),
		Count: 1,
	})
	if result.Err() != nil {
		return result.Err()
	}
	values := result.Val()
	if len(values) == 0 {
		return nil
	}
	resultRem := repo.conn.ZRem(ctx, redisKey, values[0])
	return resultRem.Err()
}
