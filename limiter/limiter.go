package limiter

import "github.com/MarlonHeiber/ratelimiter/config"

type Strategy interface {
	Allow(key string, limit int, blockSeconds int) (bool, int, error)
}

type RateLimiter struct {
	strategy Strategy
}

func NewRateLimiter(s Strategy) *RateLimiter {
	return &RateLimiter{strategy: s}
}

func (r *RateLimiter) Allow(token, ip string) (bool, string, int, int, error) {
	key := ip
	limit := config.DefaultRateLimit

	if token != "" {
		if val, ok := config.TokenLimits[token]; ok {
			key = "token:" + token
			limit = val
		}
	}

	allowed, count, err := r.strategy.Allow(key, limit, config.BlockDuration)
	return allowed, key, limit, count, err
}
