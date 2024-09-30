package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestRateLimitedClientGet(t *testing.T) {
	// Mock server, just returns OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	limiter := rate.NewLimiter(10, 1) // 10 requests per second, burst of 1
	client := NewRateLimitedClient(limiter)

	start := time.Now()
	for i := 0; i < 3; i++ {
		resp, err := client.get(server.URL)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}
	}
	duration := time.Since(start)

	if duration < 200*time.Millisecond {
		t.Errorf("Requests were not rate-limited as expected")
	}
}
