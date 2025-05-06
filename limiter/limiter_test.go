package limiter_test

import (
	"testing"
	"time"

	"github.com/MarlonHeiber/ratelimiter/config"
	"github.com/MarlonHeiber/ratelimiter/limiter"
)

type mockStrategy struct {
	calls      map[string]int
	blocked    map[string]time.Time
	lastAccess map[string]time.Time
}

func newMockStrategy() *mockStrategy {
	return &mockStrategy{
		calls:      make(map[string]int),
		blocked:    make(map[string]time.Time),
		lastAccess: make(map[string]time.Time),
	}
}

func (m *mockStrategy) Allow(key string, limit int, blockSeconds int) (bool, int, error) {
	now := time.Now()

	if until, ok := m.blocked[key]; ok {
		if now.Before(until) {
			return false, m.calls[key], nil
		}
		delete(m.blocked, key)
		m.calls[key] = 0
	}

	m.calls[key]++
	m.lastAccess[key] = now

	if m.calls[key] > limit {
		m.blocked[key] = now.Add(time.Duration(blockSeconds) * time.Second)
		return false, m.calls[key], nil
	}

	return true, m.calls[key], nil
}

// TestLimiter_Allow tests the Allow method of the RateLimiter
func TestLimiter_ByIP(t *testing.T) {
	mock := newMockStrategy()
	rl := limiter.NewRateLimiter(mock)

	config.DefaultRateLimit = 3
	config.BlockDuration = 1

	ip := "127.0.0.1"

	for i := 0; i < 3; i++ {
		ok, _, _, _, _ := rl.Allow("", ip)
		if !ok {
			t.Fatalf("expected request %d to be allowed", i+1)
		}
	}

	ok, _, _, _, _ := rl.Allow("", ip)
	if ok {
		t.Fatal("expected request to be blocked after limit exceeded")
	}
}

// TestLimiter_ByToken tests the Allow method of the RateLimiter with a token
func TestLimiter_ByTokenOverridesIP(t *testing.T) {
	mock := newMockStrategy()
	rl := limiter.NewRateLimiter(mock)

	config.DefaultRateLimit = 2
	config.BlockDuration = 1
	config.TokenLimits = make(map[string]int)
	config.TokenLimits["Example123"] = 5

	ip := "127.0.0.1"
	token := "Example123"

	for i := 0; i < 5; i++ {
		ok, _, _, _, _ := rl.Allow(token, ip)
		if !ok {
			t.Fatalf("expected token request %d to be allowed", i+1)
		}
	}

	ok, _, _, _, _ := rl.Allow(token, ip)
	if ok {
		t.Fatal("expected token to be blocked after exceeding limit")
	}
}

// TestLimiter_AllowWithDifferentIPs tests the block and the ubblock after blockduration time.
func TestLimiter_UnlockAfterBlockDuration(t *testing.T) {
	mock := newMockStrategy()
	rl := limiter.NewRateLimiter(mock)

	config.DefaultRateLimit = 2
	config.BlockDuration = 1
	ip := "10.0.0.1"

	for i := 0; i < 3; i++ {
		rl.Allow("", ip)
	}

	time.Sleep(1100 * time.Millisecond)

	ok, _, _, _, _ := rl.Allow("", ip)
	if !ok {
		t.Fatal("expected to allow request after block duration")
	}
}
