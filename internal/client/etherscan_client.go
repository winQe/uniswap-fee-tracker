package client

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/time/rate"
)

// EtherscanClient is the client for interacting with the Etherscan API.
type EtherscanClient struct {
	*RateLimitedClient
	baseURL     string
	apiKey      string
	poolAddress string
}

// TransactionData represents the simplified transaction result from the API calls
type TransactionData struct {
	BlockNumber uint64
	Hash        string
	GasUsed     uint64
	GasPriceWei *big.Int
}

// receiptResponse represents API response for the transaction receipt.
type receiptResponse struct {
	ID     int            `json:"id"`
	Result receiptDetails `json:"result"`
}

// receiptDetails holds the core details of a transaction.
type receiptDetails struct {
	BlockNumber       string `json:"blockNumber"`
	Hash              string `json:"transactionHash"`
	GasUsed           string `json:"gasUsed"`
	EffectiveGasPrice string `json:"effectiveGasPrice"`
}

// tokenTxResponse represents the API response of tokenTx API call
type tokenTxResponse struct {
	Status  string           `json:"status"` // OK = 1
	Message string           `json:"message"`
	Result  []tokenTxDetails `json:"result"`
}

// tokenTxDetails holds the core detail of the transaction from tokenTx
type tokenTxDetails struct {
	BlockNumber  string `json:"blockNumber"`
	TimeStamp    string `json:"timeStamp"`
	Hash         string `json:"hash"`
	TokenName    string `json:"tokenName"`
	TokenSymbol  string `json:"tokenSymbol"`
	TokenDecimal string `json:"tokenDecimal"`
	GasPrice     string `json:"gasPrice"`
	GasUsed      string `json:"gasUsed"`
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

func (e *EtherscanClient) GetTransactionReceipt(hash string) (*TransactionData, error) {
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
	txData := &TransactionData{
		Hash:        txHash,
		BlockNumber: blockNumber,
		GasUsed:     gasUsed,
		GasPriceWei: gasPriceWei,
	}

	return txData, nil
}

func (e *EtherscanClient) listTransactions(offset *int, startBlock *uint64, endBlock *uint64) ([]TransactionData, error) {
	params := url.Values{}

	// Required parameters
	params.Add("module", "account")
	params.Add("action", "tokentx")
	params.Add("address", e.poolAddress)
	params.Add("apikey", e.apiKey)
	params.Add("sort", "desc")

	// Optional parameters
	if offset != nil {
		params.Add("offset", strconv.Itoa(*offset))
	}

	if startBlock != nil {
		params.Add("startblock", strconv.FormatUint(*startBlock, 10))
	}

	if endBlock != nil {
		params.Add("endblock", strconv.FormatUint(*endBlock, 10))
	}

	txURL := fmt.Sprintf("%s?%s", e.baseURL, params.Encode())
	resp, err := e.get(txURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var result tokenTxResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	// Check for success in the API response (status == 1)
	if result.Status != "1" {
		return nil, fmt.Errorf("Etherscan server error: %s", result.Message)
	}

	transactions, err := convertResponseToTransactionData(result)
	if err != nil {
		return nil, fmt.Errorf("error converting response to transactionData: %v", err)
	}

	return transactions, nil
}

// Convert a tokenTxResponse to a slice of TransactionData
func convertResponseToTransactionData(response tokenTxResponse) ([]TransactionData, error) {
	var transactions []TransactionData

	for _, detail := range response.Result {
		txData, err := convertToTransactionData(detail)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *txData)
	}

	return transactions, nil
}

// Convert tokenTxDetails to TransactionData
func convertToTransactionData(details tokenTxDetails) (*TransactionData, error) {
	blockNumber, err := strconv.ParseUint(details.BlockNumber, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting BlockNumber: %v", err)
	}

	gasUsed, err := strconv.ParseUint(details.GasUsed, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting GasUsed: %v", err)
	}

	// Convert GasPrice to big.Int
	gasPriceWei := new(big.Int)
	_, ok := gasPriceWei.SetString(details.GasPrice, 10)
	if !ok {
		return nil, fmt.Errorf("error converting GasPrice to big.Int")
	}

	txData := &TransactionData{
		BlockNumber: blockNumber,
		Hash:        details.Hash,
		GasUsed:     gasUsed,
		GasPriceWei: gasPriceWei,
	}

	return txData, nil
}
