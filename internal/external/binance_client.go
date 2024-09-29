package external

import (
	"fmt"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

// BinanceClient is the client for interacting with the Binance Spot API.
type BinanceClient struct {
	*RateLimitedClient
	baseURL string
}

// NewBinanceClient initializes a new Binance Spot API client with rate limits.
// Binance has a limit of 6000 weights per minute, and each Kline API call costs 2 weights.
// This allows for approximately 3000 calls per minute or 50 calls per second.
//
// https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#limits
// https://developers.binance.com/docs/binance-spot-api-docs/web-socket-api#ip-limits
func NewBinanceClient() *BinanceClient {
	// 50 requests per second, with a burst of 30 requests.
	rateLimit := rate.NewLimiter(50, 30)

	return &BinanceClient{
		RateLimitedClient: NewRateLimitedClient(rateLimit),
		baseURL:           "https://api.binance.com/api/v3",
	}
}
