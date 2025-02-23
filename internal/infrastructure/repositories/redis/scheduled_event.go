package redis

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/thirdparty"
)

type RedisScheduledEventRepository struct {
	conn   thirdparty.BaseRedisClient
	schKey string
}

func NewRedisScheduledEventRepository(client thirdparty.BaseRedisClient, schedulingKey string) *RedisScheduledEventRepository {
	return &RedisScheduledEventRepository{
		conn:   client,
		schKey: schedulingKey,
	}
}

func (repo *RedisScheduledEventRepository) Push(ctx context.Context, when time.Time, evt events.Event) error {
	evtData, err := evt.Serialize()
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Scheduling evt %s for %d (%s)", evt.GetType(), when.Unix(), when), "Source", "RedisScheduledEventRepository")
	result := repo.conn.ZAdd(ctx, repo.schKey, redis.Z{Score: float64(when.Unix()), Member: evtData})
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (repo *RedisScheduledEventRepository) GetNext(ctx context.Context) (events.Event, error) {
	nextOne, err := repo.getNextOne(ctx)
	if err != nil {
		return nil, err
	}
	if nextOne == nil {
		return nil, nil
	}
	evt, err := events.RetrieveTypedEvent([]byte(*nextOne))
	if err == nil {
		slog.Info(fmt.Sprintf("Found scheduled event %s, at %s", evt.GetType(), time.Now().UTC()), "Source", "RedisScheduledEventRepository")
	}
	return evt, err
}

func (repo *RedisScheduledEventRepository) RemoveNext(ctx context.Context) error {
	nextOne, err := repo.getNextOne(ctx)
	if err != nil {
		return err
	}
	if nextOne == nil {
		return nil
	}
	resultRem := repo.conn.ZRem(ctx, repo.schKey, *nextOne)
	return resultRem.Err()
}

func (repo *RedisScheduledEventRepository) getNextOne(ctx context.Context) (*string, error) {
	result := repo.conn.ZRangeByScore(ctx, repo.schKey, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   strconv.FormatInt(time.Now().UTC().Unix(), 10),
		Count: 1,
	})
	slog.Debug(fmt.Sprintf("Checking next evt at %s", strconv.FormatInt(time.Now().UTC().Unix(), 10)), "Source", "RedisScheduledEventRepository")
	if result.Err() != nil {
		return nil, result.Err()
	}
	values := result.Val()
	if len(values) == 0 {
		return nil, nil
	}
	return &values[0], nil
}
