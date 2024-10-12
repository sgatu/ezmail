package redis

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sgatu/ezmail/internal/domain/models/events"
)

type RedisScheduledEventRepository struct {
	conn *redis.Client
}

func NewRedisScheduledEventRepository(client *redis.Client) *RedisScheduledEventRepository {
	return &RedisScheduledEventRepository{
		conn: client,
	}
}

const (
	redisKey = "q_scheduled_events"
)

func (repo *RedisScheduledEventRepository) Push(ctx context.Context, when time.Time, evt events.Event) error {
	evtData, err := evt.Serialize()
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Scheduling evt %s for %d", evt.GetType(), when.Unix()))
	result := repo.conn.ZAdd(ctx, redisKey, redis.Z{Score: float64(when.Unix()), Member: evtData})
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (repo *RedisScheduledEventRepository) GetNext(ctx context.Context) (events.Event, error) {
	nextOne, err := getNextOne(ctx, repo.conn)
	if err != nil {
		return nil, err
	}
	if nextOne == nil {
		return nil, nil
	}
	return events.RetrieveTypedEvent([]byte(*nextOne))
}

func (repo *RedisScheduledEventRepository) RemoveNext(ctx context.Context) error {
	nextOne, err := getNextOne(ctx, repo.conn)
	if err != nil {
		return err
	}
	if nextOne == nil {
		return nil
	}
	resultRem := repo.conn.ZRem(ctx, redisKey, *nextOne)
	return resultRem.Err()
}

func getNextOne(ctx context.Context, conn *redis.Client) (*string, error) {
	result := conn.ZRangeByScore(ctx, redisKey, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   strconv.FormatInt(time.Now().UTC().Unix(), 10),
		Count: 1,
	})
	slog.Info(fmt.Sprintf("Checking next evt at %s", strconv.FormatInt(time.Now().UTC().Unix(), 10)))
	if result.Err() != nil {
		return nil, result.Err()
	}
	values := result.Val()
	if len(values) == 0 {
		return nil, nil
	}
	return &values[0], nil
}
