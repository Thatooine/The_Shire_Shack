package rateLimiting

import (
	"context"
	"time"

	"github.com/mennanov/limiters"
)

// TokenBucketRateLimiter defines the interface for rate limiting operations.
type TokenBucketRateLimiter interface {
	TokenBucket(ctx context.Context, request TokenBucketRequest, stateBackend limiters.TokenBucketStateBackend) *TokenBucketResponse
	Limit(ctx context.Context, tokenBucket *limiters.TokenBucket) (*LimitResponse, error)
}

type TokenBucketRequest struct {
	RefillRate time.Duration
	Capacity   int64
}

type TokenBucketResponse struct {
	TokenBucket *limiters.TokenBucket
}

type LimitResponse struct {
	TimeToRetry time.Duration
}
