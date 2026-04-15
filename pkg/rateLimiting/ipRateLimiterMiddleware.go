package rateLimiting

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/mennanov/limiters"
	"github.com/rs/zerolog/log"
)

type ipRateLimiterMiddleware struct {
	next        http.Handler
	rateLimiter RedisTokenBucketRateLimiter
	refillRate  time.Duration
	capacity    int64
	ttl         time.Duration
}

// NewIpRateLimiterMiddleware returns a gorilla/mux-compatible middleware that
// rate limits requests by client IP address using a Redis-backed token bucket.
func NewIpRateLimiterMiddleware(
	rateLimiter RedisTokenBucketRateLimiter,
	capacity int64,
	refillRate time.Duration,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &ipRateLimiterMiddleware{
			next:        next,
			rateLimiter: rateLimiter,
			refillRate:  refillRate,
			capacity:    capacity,
			ttl:         time.Duration(int64(refillRate) * capacity),
		}
	}
}

func (m *ipRateLimiterMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	key := "ip_token_bucket:" + ip

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
			log.Ctx(ctx).Warn().Str("ip", ip).Msg("ip rate limit exceeded")
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
