package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	DefaultRateLimit int
	BlockDuration    int
	TokenLimits      map[string]int
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
)

func LoadConfig() error {
	_ = godotenv.Load()

	DefaultRateLimit, _ = strconv.Atoi(getEnv("RATE_LIMIT_DEFAULT", "5"))
	BlockDuration, _ = strconv.Atoi(getEnv("BLOCK_DURATION_SECONDS", "300"))

	RedisAddr = getEnv("REDIS_ADDR", "localhost:6379")
	RedisPassword = getEnv("REDIS_PASSWORD", "")
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	RedisDB = db

	TokenLimits = make(map[string]int)
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "RATE_LIMIT_TOKEN_") {
			parts := strings.SplitN(e, "=", 2)
			token := strings.TrimPrefix(parts[0], "RATE_LIMIT_TOKEN_")
			limit, _ := strconv.Atoi(parts[1])
			TokenLimits[token] = limit
		}
	}

	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
