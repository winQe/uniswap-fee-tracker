package external

import (
	"fmt"
	"net/url"

	"golang.org/x/time/rate"
)

// Client for Etherscan API
type EtherscanClient struct {
	*RateLimitedClient
	baseUrl string
	apiKey  string
}

// NewEtherscanClient initializes Etherscan with Free Plan API Limits
func NewEtherscanClient(apiKey string) *EtherscanClient {
	// 5 API calls per second
	secondLimiter := rate.NewLimiter(5, 5) // 5 requests per second, burst of 5

	// 100,000 API calls per day = ~1.15 requests per second
	dailyLimiter := rate.NewLimiter(1.15, 1000) // 1.15 requests per second, burst of 1,000

	return &EtherscanClient{
		RateLimitedClient: NewRateLimitedClient(secondLimiter, dailyLimiter),
		baseUrl:           "https://api.etherscan.io/api",
		apiKey:            apiKey,
	}
}

func (e *EtherscanClient) GetTransactionByHash(hash string) error {
	params := url.Values{}

	// https://docs.etherscan.io/api-endpoints/geth-parity-proxy#eth_gettransactionbyhash
	params.Add("module", "proxy")
	params.Add("action", "eth_getTransactionByHash")
	params.Add("txhash", hash)
	params.Add("apikey", e.apiKey)

	txUrl := fmt.Sprintf("%s?%s", e.baseUrl, params.Encode())
	_, err := e.Get(txUrl)
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	// TODO: Parse the response and return the gas fee
	return nil
}
