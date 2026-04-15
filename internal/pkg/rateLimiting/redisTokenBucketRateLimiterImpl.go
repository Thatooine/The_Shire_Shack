package rateLimiting

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/rateLimiting"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/mennanov/limiters"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisRateLimiterImpl struct {
	redisClient *redis.Client
	locker      limiters.DistLocker
}

func NewRedisRateLimiterImpl(redisClient *redis.Client) *RedisRateLimiterImpl {
	pool := goredis.NewPool(redisClient)
	return &RedisRateLimiterImpl{
		redisClient: redisClient,
		locker:      limiters.NewLockRedis(pool, "rate_limiter_lock"),
	}
}

func (r *RedisRateLimiterImpl) TokenStateBackend(_ context.Context, key string, ttl time.Duration) (*rateLimiting.TokenStateBackendResponse, error) {

	stateBackend := limiters.NewTokenBucketRedis(r.redisClient, key, ttl, false)
	return &rateLimiting.TokenStateBackendResponse{
		StateBackend: stateBackend,
	}, nil
}

func (r *RedisRateLimiterImpl) TokenBucket(_ context.Context, request rateLimiting.TokenBucketRequest, stateBackend limiters.TokenBucketStateBackend) *rateLimiting.TokenBucketResponse {
	tokenBucket := limiters.NewTokenBucket(
		request.Capacity,
		request.RefillRate,
		r.locker,
		stateBackend,
		limiters.NewSystemClock(),
		nil,
	)

	return &rateLimiting.TokenBucketResponse{
		TokenBucket: tokenBucket,
	}
}

func (r *RedisRateLimiterImpl) Limit(ctx context.Context, tokenBucket *limiters.TokenBucket) (*rateLimiting.LimitResponse, error) {
	wait, err := tokenBucket.Limit(ctx)
	if err != nil {
		if errors.Is(err, limiters.ErrLimitExhausted) {
			log.Ctx(ctx).Warn().Msg("rate limit exhausted")
			return nil, limiters.ErrLimitExhausted
		}
		log.Ctx(ctx).Error().Err(err).Msg("rate limiting failed")
		return nil, fmt.Errorf("rate limiting failed: %w", err)
	}

	return &rateLimiting.LimitResponse{
		TimeToRetry: wait,
	}, nil
}
