package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
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

// GetETHUSDT fetches latest ETH - USDT conversion rate for the given timestamp (1 min interval)
// It queries the Binance Kline API
func (b *BinanceClient) GetETHUSDT(timestamp time.Time) (*KlineData, error) {
	params := url.Values{}
	params.Add("symbol", "ETHUSDT")
	params.Add("interval", "1m")
	params.Add("startTime", fmt.Sprintf("%d", timestamp.UnixMilli()))

	klinesURL := fmt.Sprintf("%s/klines?%s", b.baseURL, params.Encode())

	res, err := b.httpClient.Get(klinesURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Format can be found here
	// https://binance-docs.github.io/apidocs/spot/en/#kline-candlestick-data
	var klines [][]interface{}
	if err := json.Unmarshal(body, &klines); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	// Check if klines is empty
	if len(klines) == 0 {
		return nil, fmt.Errorf("no kline data returned")
	}

	closePriceStr, ok := klines[0][4].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected format for close price")
	}
	closePrice, err := strconv.ParseFloat(closePriceStr, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting close price to float64: %v", err)
	}

	// Return struct instead of a single value for future-proofing
	return &KlineData{
		ClosePrice: closePrice,
	}, nil
}
