package mock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type MockRedisClient struct {
	zAddReturn          *redis.IntCmd
	zRemReturn          *redis.IntCmd
	zRangeByScoreReturn *redis.StringSliceCmd
	xAddReturn          *redis.StringCmd
	xReadReturn         *redis.XStreamSliceCmd
	getReturn           *redis.StringCmd
	setReturn           *redis.StatusCmd

	zRangeLastRequest struct {
		ctx context.Context
		opt *redis.ZRangeBy
		key string
	}
	xAddLastRequest struct {
		ctx context.Context
		a   *redis.XAddArgs
	}
	zAddLastRequest struct {
		ctx     context.Context
		key     string
		members []redis.Z
	}
	zRemLastRequest struct {
		ctx     context.Context
		key     string
		members []interface{}
	}
	xReadLastRequest struct {
		ctx context.Context
		a   *redis.XReadArgs
	}
	getLastRequest struct {
		ctx context.Context
		key string
	}
	setLastRequest struct {
		value      interface{}
		ctx        context.Context
		key        string
		expiration time.Duration
	}

	ZAddCalls          int
	ZRemCalls          int
	ZRangeByScoreCalls int
	XAddCalls          int
	XReadCalls         int
	GetCalls           int
	SetCalls           int
}

func (mr *MockRedisClient) XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd {
	mr.XReadCalls++
	mr.xReadLastRequest = struct {
		ctx context.Context
		a   *redis.XReadArgs
	}{
		ctx: ctx,
		a:   a,
	}
	return mr.xReadReturn
}

func (mr *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	mr.GetCalls++
	mr.getLastRequest = struct {
		ctx context.Context
		key string
	}{
		ctx: ctx,
		key: key,
	}
	return mr.getReturn
}

func (mr *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	mr.SetCalls++
	mr.setLastRequest = struct {
		value      interface{}
		ctx        context.Context
		key        string
		expiration time.Duration
	}{
		ctx:        ctx,
		key:        key,
		value:      value,
		expiration: expiration,
	}
	return mr.setReturn
}

func (mr *MockRedisClient) XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd {
	mr.XAddCalls++
	mr.xAddLastRequest = struct {
		ctx context.Context
		a   *redis.XAddArgs
	}{
		ctx: ctx,
		a:   a,
	}
	return mr.xAddReturn
}

func (mr *MockRedisClient) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	mr.ZAddCalls++
	mr.zAddLastRequest = struct {
		ctx     context.Context
		key     string
		members []redis.Z
	}{
		ctx:     ctx,
		key:     key,
		members: members,
	}
	return mr.zAddReturn
}

func (mr *MockRedisClient) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	mr.ZRemCalls++
	mr.zRemLastRequest = struct {
		ctx     context.Context
		key     string
		members []interface{}
	}{
		ctx:     ctx,
		key:     key,
		members: members,
	}
	return mr.zRemReturn
}

func (mr *MockRedisClient) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	mr.ZRangeByScoreCalls++
	mr.zRangeLastRequest = struct {
		ctx context.Context
		opt *redis.ZRangeBy
		key string
	}{
		ctx: ctx,
		key: key,
		opt: opt,
	}
	return mr.zRangeByScoreReturn
}

func (mr *MockRedisClient) SetZAddResult(ret *redis.IntCmd) {
	mr.zAddReturn = ret
}

func (mr *MockRedisClient) SetZRemResult(ret *redis.IntCmd) {
	mr.zRemReturn = ret
}

func (mr *MockRedisClient) SetZRangeByScoreResult(ret *redis.StringSliceCmd) {
	mr.zRangeByScoreReturn = ret
}

func (mr *MockRedisClient) SetXAddResult(ret *redis.StringCmd) {
	mr.xAddReturn = ret
}
