package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MarlonHeiber/ratelimiter/config"
	"github.com/MarlonHeiber/ratelimiter/limiter"
	"github.com/MarlonHeiber/ratelimiter/middleware"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	strategy := limiter.NewRedisLimiter(config.RedisAddr, config.RedisPassword, config.RedisDB)
	rl := limiter.NewRateLimiter(strategy)

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	handler := middleware.RateLimitMiddleware(rl)(mux)
	log.Println("Server Running at localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
