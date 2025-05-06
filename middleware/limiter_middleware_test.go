package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MarlonHeiber/ratelimiter/limiter"
	"github.com/MarlonHeiber/ratelimiter/middleware"
)

type alwaysAllowStrategy struct{}

func (s *alwaysAllowStrategy) Allow(key string, limit int, blockSeconds int) (bool, int, error) {
	return true, 1, nil
}

type alwaysBlockStrategy struct{}

func (s *alwaysBlockStrategy) Allow(key string, limit int, blockSeconds int) (bool, int, error) {
	return false, 100, nil
}

// TestMiddleware_AllowRequest tests the RateLimitMiddleware with a strategy that always allows requests
func TestMiddleware_AllowRequest(t *testing.T) {
	rl := limiter.NewRateLimiter(&alwaysAllowStrategy{})

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	w := httptest.NewRecorder()

	handler := middleware.RateLimitMiddleware(rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

// TestMiddleware_BlockRequest tests the RateLimitMiddleware with a strategy that always blocks requests
// and returns a 429 status code
func TestMiddleware_BlockRequest(t *testing.T) {
	rl := limiter.NewRateLimiter(&alwaysBlockStrategy{})

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	w := httptest.NewRecorder()

	handler := middleware.RateLimitMiddleware(rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, got %d", w.Code)
	}

	expectedMsg := "you have reached the maximum number of requests or actions allowed within a certain time frame\n"
	if w.Body.String() != expectedMsg {
		t.Fatalf("unexpected error message: %s", w.Body.String())
	}
}
