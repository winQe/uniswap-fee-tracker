package client

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimitedClient defines a base structure for a rate-limited API client
type RateLimitedClient struct {
	httpClient   *http.Client
	rateLimiters []*rate.Limiter // Allows multiple rate limiter (some API has both limits per day and per second)
}

// NewRateLimitedClient initializes the base client with multiple rate limiters
func NewRateLimitedClient(rateLimits ...*rate.Limiter) *RateLimitedClient {
	return &RateLimitedClient{
		httpClient:   &http.Client{},
		rateLimiters: rateLimits,
	}
}

// Get sends a GET request with rate limits applied
func (c *RateLimitedClient) Get(url string) (*http.Response, error) {
	ctx := context.Background()

	// Apply the rate limits
	for _, limiter := range c.rateLimiters {
		if err := limiter.Wait(ctx); err != nil {
			return nil, err
		}
	}

	return c.httpClient.Get(url)
}
