package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(addr, password string, db int) *RedisLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisLimiter{client: rdb}
}

func (r *RedisLimiter) Allow(key string, limit int, blockSeconds int) (bool, int, error) {
	ctx := context.Background()
	countKey := fmt.Sprintf("rl:%s", key)
	blockKey := fmt.Sprintf("block:%s", key)

	blocked, err := r.client.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, 0, err
	}
	if blocked == 1 {
		return false, 0, nil
	}

	count, err := r.client.Incr(ctx, countKey).Result()
	if err != nil {
		return false, 0, err
	}
	if count == 1 {
		r.client.Expire(ctx, countKey, time.Second)
	}

	if int(count) > limit {
		r.client.Set(ctx, blockKey, "1", time.Duration(blockSeconds)*time.Second)
		return false, int(count), nil
	}

	return true, int(count), nil
}
