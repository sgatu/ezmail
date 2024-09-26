package redis

import "github.com/redis/go-redis/v9"

var NewClient = redis.NewClient

type Options = redis.Options
