package client

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/winQe/uniswap-fee-tracker/internal/utils"
	"golang.org/x/time/rate"
)

// EtherscanClient is the client for interacting with the Etherscan API.
type EtherscanClient struct {
	*RateLimitedClient
	baseURL     string
	apiKey      string
	poolAddress string
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
	Status  string          `json:"status"` // OK = 1
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
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

// blockNumberResponse represents the API response for getting block number by timestamp.
type blockNumberResponse struct {
	Status  string `json:"status"`  // "1" indicates success
	Message string `json:"message"` // e.g., "OK"
	Result  string `json:"result"`  // Block number as a string
}

// NewEtherscanClient initializes Etherscan with Free Plan API Limits
func NewEtherscanClient(apiKey string, poolAddress string) *EtherscanClient {
	// 5 API calls per second
	secondLimiter := rate.NewLimiter(5, 5) // 5 requests per second, burst of 5

	// 100,000 API calls per day = ~1.15 requests per second
	dailyLimiter := rate.NewLimiter(1.15, 1000) // 1.15 requests per second, burst of 1,000

	return &EtherscanClient{
		RateLimitedClient: NewRateLimitedClient(secondLimiter, dailyLimiter),
		baseURL:           "https://api.etherscan.io/api",
		apiKey:            apiKey,
		poolAddress:       poolAddress,
	}
}

// GetTransactionReceipt fetches the transaction receipt based on the txHash
// TODO: Needs to verify whether transaction actually belongs in the WETH-USDT Pool
// TODO: Timestamp from eth_getBlockByNumber
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

// GetLatestTransaction fetches the latest transaction from the Uniswap V3 ETH-USDC pool.
func (e *EtherscanClient) GetLatestTransaction() (*TransactionData, error) {
	// Only the latest transaction
	offset := 1
	page := 1
	transactions, err := e.ListTransactions(&offset, nil, nil, &page)
	if err != nil || len(transactions) == 0 {
		return nil, fmt.Errorf("error fetching the latest transaction: %v", err)
	}
	return &transactions[0], nil
}

// listTransactions is a helper function to query transactions from the Uniswap V3 WETH-USDC pool based on optional parameters.
func (e *EtherscanClient) ListTransactions(offset *int, startBlock *uint64, endBlock *uint64, page *int) ([]TransactionData, error) {
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

	if page != nil {
		params.Add("page", strconv.Itoa(*page))
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

	// Check for success in the API response (status == "1")
	if result.Status != "1" {
		return nil, fmt.Errorf("Etherscan server error: %s", result.Message)
	}

	// Unmarshal the Result into []tokenTxDetails
	var transactionsDetails []tokenTxDetails
	if err := json.Unmarshal(result.Result, &transactionsDetails); err != nil {
		return nil, fmt.Errorf("error parsing transactions: %v", err)
	}

	// Convert to []TransactionData
	transactions, err := convertResponseToTransactionData(transactionsDetails)
	if err != nil {
		return nil, fmt.Errorf("error converting response to TransactionData: %v", err)
	}

	return transactions, nil
}

// Adjust convertResponseToTransactionData to accept tokenTxDetails as a parameter
func convertResponseToTransactionData(details []tokenTxDetails) ([]TransactionData, error) {
	var transactions []TransactionData

	for _, detail := range details {
		txData, err := convertToTransactionData(detail)
		if err != nil {
			// Log the error and skip the transaction
			fmt.Printf("Error converting transaction data: %v\n", err)
			continue
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

	unixTime, err := utils.ParseUnixTime(details.TimeStamp)
	if err != nil {
		return nil, fmt.Errorf("error converting timestamp: %v", err)
	}
	txTime := time.Unix(unixTime, 0)

	txData := &TransactionData{
		BlockNumber: blockNumber,
		Hash:        details.Hash,
		GasUsed:     gasUsed,
		GasPriceWei: gasPriceWei,
		Timestamp:   txTime,
	}

	return txData, nil
}

// GetBlockNumberByTimestamp fetches the block number closest(can be before of after) to the given timestamp.
func (e *EtherscanClient) GetBlockNumberByTimestamp(timestamp time.Time, before bool) (uint64, error) {
	closest := "after"
	if before {
		closest = "before"
	}

	// Prepare query parameters
	params := url.Values{}
	params.Add("module", "block")
	params.Add("action", "getblocknobytime")
	params.Add("timestamp", strconv.FormatInt(timestamp.Unix(), 10)) // base 10
	params.Add("closest", closest)
	params.Add("apikey", e.apiKey)

	blockURL := fmt.Sprintf("%s?%s", e.baseURL, params.Encode())

	resp, err := e.get(blockURL)
	if err != nil {
		return 0, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var blockResp blockNumberResponse
	if err := json.NewDecoder(resp.Body).Decode(&blockResp); err != nil {
		return 0, fmt.Errorf("error parsing JSON response: %v", err)
	}

	// Check if the API returned a successful status
	if blockResp.Status != "1" {
		return 0, fmt.Errorf("Etherscan API error: %s - %s", blockResp.Status, blockResp.Message)
	}

	blockNumber, err := strconv.ParseUint(blockResp.Result, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting block number: %v", err)
	}

	return blockNumber, nil
}
