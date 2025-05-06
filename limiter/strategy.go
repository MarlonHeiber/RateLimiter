package limiter

type LimiterStrategy interface {
	Allow(key string, limit int, blockSeconds int) (bool, error)
}
