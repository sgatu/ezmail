package thirdparty

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type BaseRedisClient interface {
	ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd
	XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd
	XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type RedisClient struct {
	*redis.Client
}
