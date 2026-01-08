package middlewares

import (
	"net/http"
	"time"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/core/cache/redis"
	"github.com/go-chi/httprate"
	httprateredis "github.com/go-chi/httprate-redis"
)

type RateLimiterFactory struct {
	redis *redis.Client
}

func NewRateLimiterFactory(redis *redis.Client) *RateLimiterFactory {
	return &RateLimiterFactory{
		redis: redis,
	}
}

func (r *RateLimiterFactory) NewRateLimiter(reqLimit int, perDuration time.Duration) func(http.Handler) http.Handler {
	return httprate.Limit(
		reqLimit,
		perDuration,
		httprate.WithKeyByIP(),
		httprateredis.WithRedisLimitCounter(&httprateredis.Config{
			Client: r.redis.PureClient(),
		}),
		httprate.WithLimitHandler(rateLimiterHandler),
	)
}

func rateLimiterHandler(w http.ResponseWriter, r *http.Request) {
	errors.HandleError(r.Context(), w, domainErr.NewTooManyRequestsError("too many requests, try later", nil))
}
