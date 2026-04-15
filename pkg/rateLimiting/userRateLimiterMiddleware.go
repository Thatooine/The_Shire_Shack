package rateLimiting

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/authentication"
	"github.com/mennanov/limiters"
	"github.com/rs/zerolog/log"
)

type userRateLimiterMiddleware struct {
	next        http.Handler
	rateLimiter RedisTokenBucketRateLimiter
	refillRate  time.Duration
	capacity    int64
	ttl         time.Duration
}

// NewUserRateLimiterMiddleware returns a gorilla/mux-compatible middleware that
// rate limits authenticated requests by UserID using a Redis-backed token bucket.
func NewUserRateLimiterMiddleware(
	rateLimiter RedisTokenBucketRateLimiter,
	capacity int64,
	refillRate time.Duration,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &userRateLimiterMiddleware{
			next:        next,
			rateLimiter: rateLimiter,
			refillRate:  refillRate,
			capacity:    capacity,
			ttl:         time.Duration(int64(refillRate) * capacity),
		}
	}
}

func (m *userRateLimiterMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claim, ok := authentication.LoginClaimFromContext(ctx)
	if !ok {
		// No login claim — let the auth middleware handle rejection.
		m.next.ServeHTTP(w, r)
		return
	}

	key := "token_bucket:" + claim.UserID

	stateBackend, err := m.rateLimiter.TokenStateBackend(ctx, key, m.ttl)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create token state backend")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	tokenBucket := m.rateLimiter.TokenBucket(
		ctx,
		TokenBucketRequest{
			RefillRate: m.refillRate,
			Capacity:   m.capacity,
		}, stateBackend.StateBackend)

	if _, err := m.rateLimiter.Limit(ctx, tokenBucket.TokenBucket); err != nil {
		if errors.Is(err, limiters.ErrLimitExhausted) {
			log.Ctx(ctx).Warn().Str("userID", claim.UserID).Msg("rate limit exceeded")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "too many requests"})
			return
		}
		log.Ctx(ctx).Error().Err(err).Msg("rate limiting error")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	m.next.ServeHTTP(w, r)
}
