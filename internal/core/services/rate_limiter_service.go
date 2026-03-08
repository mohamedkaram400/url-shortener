package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)


type RateLimiterService struct {
	Redis *redis.Client
}

func NewRateLimiterService(redis *redis.Client) *RateLimiterService {
	return &RateLimiterService{Redis: redis}
}

func (r *RateLimiterService) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	count, err := r.Redis.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		r.Redis.Expire(ctx, key, window)
	}

	if count > int64(limit) {
		return false, nil
	}

	return true, nil
}