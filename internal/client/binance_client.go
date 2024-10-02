package client

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"golang.org/x/time/rate"
)

// KlineClient is the client for interacting with Binance Kline API using go-binance.
type KlineClient struct {
	binanceClient *binance.Client
	rateLimiter   *rate.Limiter
}

// NewKlineClient initializes a new KlineClient with rate limits.
// Binance has a limit of 6000 weights per minute, and each Kline API call costs 2 weights.
// This allows for approximately 3000 calls per minute or 50 calls per second.
//
// https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#limits
// https://developers.binance.com/docs/binance-spot-api-docs/web-socket-api#ip-limits
func NewKlineClient() *KlineClient {
	// Initialize the Binance client. No API key is required for public endpoints.
	binance.UseTestnet = true
	binanceClient := binance.NewClient("", "")

	// Set up a rate limiter: 50 requests per second with a burst of 30.
	rateLimiter := rate.NewLimiter(50, 30)

	return &KlineClient{
		binanceClient: binanceClient,
		rateLimiter:   rateLimiter,
	}
}

// GetETHUSDT fetches the latest ETH/USDT conversion rate for the given timestamp (1 min interval).
// It queries the Binance Kline API.
func (k *KlineClient) GetETHUSDT(timestamp time.Time) (*KlineData, error) {
	if timestamp.IsZero() {
		return nil, fmt.Errorf("timestamp is invalid")
	}
	// Respect the rate limit
	if err := k.rateLimiter.Wait(context.Background()); err != nil {
		return nil, fmt.Errorf("rate limiter error: %v", err)
	}

	// Prepare the Kline request
	klinesService := k.binanceClient.NewKlinesService()
	klinesService.Symbol("ETHUSDT").
		Interval("1m").
		StartTime(timestamp.UnixMilli()).
		Limit(1) // Fetch only the latest kline

	klines, err := klinesService.Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error fetching klines: %v", err)
	}

	// Check if klines are returned
	if len(klines) == 0 {
		return nil, fmt.Errorf("no kline data returned")
	}

	// Extract the close price from the first kline
	closePriceStr := klines[0].Close
	closePrice, err := strconv.ParseFloat(closePriceStr, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting close price to float64: %v", err)
	}

	// Return the structured KlineData
	return &KlineData{
		ClosePrice: closePrice,
	}, nil
}
