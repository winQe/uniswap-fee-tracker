package client

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/time/rate"
)

// EtherscanClient is the client for interacting with the Etherscan API.
type EtherscanClient struct {
	*RateLimitedClient
	baseURL string
	apiKey  string
}

// receiptResponse represents the full response for the transaction receipt.
type receiptResponse struct {
	ID     int       `json:"id"`
	Result txDetails `json:"result"`
}

// txDetails holds the core details of a transaction.
type txDetails struct {
	BlockNumber       string `json:"blockNumber"`
	Hash              string `json:"transactionHash"`
	GasUsed           string `json:"gasUsed"`
	EffectiveGasPrice string `json:"effectiveGasPrice"`
}

// ReceiptData represents the simplified data for the transaction receipt, which is the return of this function
type ReceiptData struct {
	BlockNumber uint64
	Hash        string
	GasUsed     uint64
	GasPriceWei *big.Int
}

// NewEtherscanClient initializes Etherscan with Free Plan API Limits
func NewEtherscanClient(apiKey string) *EtherscanClient {
	// 5 API calls per second
	secondLimiter := rate.NewLimiter(5, 5) // 5 requests per second, burst of 5

	// 100,000 API calls per day = ~1.15 requests per second
	dailyLimiter := rate.NewLimiter(1.15, 1000) // 1.15 requests per second, burst of 1,000

	return &EtherscanClient{
		RateLimitedClient: NewRateLimitedClient(secondLimiter, dailyLimiter),
		baseURL:           "https://api.etherscan.io/api",
		apiKey:            apiKey,
	}
}

func (e *EtherscanClient) GetTransactionReceipt(hash string) (*ReceiptData, error) {
	params := url.Values{}

	// https://docs.etherscan.io/api-endpoints/geth-parity-proxy#eth_gettransactionbyhash
	params.Add("module", "proxy")
	params.Add("action", "eth_getTransactionByHash")
	params.Add("txhash", hash)
	params.Add("apikey", e.apiKey)

	txURL := fmt.Sprintf("%s?%s", e.baseURL, params.Encode())
	_, err := e.get(txURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	// Make the GET request
	resp, err := e.get(txURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var receipt receiptResponse
	err = json.NewDecoder(resp.Body).Decode(&receipt)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	// Extract values
	txHash := receipt.Result.Hash
	blockNumber, err := hexutil.DecodeUint64(receipt.Result.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("error converting block number: %v", err)
	}

	gasUsed, err := hexutil.DecodeUint64(receipt.Result.GasUsed)
	if err != nil {
		return nil, fmt.Errorf("error converting gas used: %v", err)
	}

	gasPriceWei, err := hexutil.DecodeBig(receipt.Result.EffectiveGasPrice)
	if err != nil {
		return nil, fmt.Errorf("error converting gas price %v", err)
	}

	// Construct the TransactionData and return
	txData := &ReceiptData{
		Hash:        txHash,
		BlockNumber: blockNumber,
		GasUsed:     gasUsed,
		GasPriceWei: gasPriceWei,
	}

	return txData, nil
}
