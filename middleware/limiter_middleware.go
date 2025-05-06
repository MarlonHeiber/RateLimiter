package middleware

import (
	"log"
	"net"
	"net/http"

	"github.com/MarlonHeiber/ratelimiter/limiter"
)

func RateLimitMiddleware(rl *limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			token := r.Header.Get("API_KEY")

			allowed, key, limit, count, err := rl.Allow(token, ip)
			if err != nil {
				log.Printf("[RateLimiter] ERRO: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			log.Printf("[RateLimiter] IP=%s | Token=%s | Key=%s | Limit=%d | Count=%d | Allowed=%v",
				ip, token, key, limit, count, allowed)

			if !allowed {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
