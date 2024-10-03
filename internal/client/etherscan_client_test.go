package client

import (
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/winQe/uniswap-fee-tracker/internal/types"
	"golang.org/x/time/rate"
)

// mockServerConfig holds the configuration for the mock HTTP server.
type mockServerConfig struct {
	expectedParams map[string]string
	responseBody   string
}

// createMockServer initializes a mock HTTP server based on the provided configuration.
func createMockServer(config mockServerConfig) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// Verify expected query parameters
		for key, expectedValue := range config.expectedParams {
			actualValue := query.Get(key)
			assert.Equal(nil, expectedValue, actualValue, fmt.Sprintf("Parameter %s mismatch", key))
		}

		// Respond with the specified JSON body
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, config.responseBody)
	}))
}

// initializeEtherscanClient sets up the EtherscanClient with a mock HTTP client.
func initializeEtherscanClient(mockServer *httptest.Server, apiKey string, poolAddress string) *EtherscanClient {
	rateLimiter := rate.NewLimiter(5, 5)
	dailyLimiter := rate.NewLimiter(1.15, 1000)
	rateLimitedClient := NewRateLimitedClient(rateLimiter, dailyLimiter)

	// Override the httpClient to use the mock server's client
	rateLimitedClient.httpClient = mockServer.Client()

	// Initialize EtherscanClient with the mock server's URL and mock RateLimitedClient
	return &EtherscanClient{
		RateLimitedClient: rateLimitedClient,
		baseURL:           mockServer.URL,
		apiKey:            apiKey,
		poolAddress:       poolAddress,
	}
}

func TestGetTransactionReceipt(t *testing.T) {
	expectedParams := map[string]string{
		"module": "proxy",
		"action": "eth_getTransactionByHash",
		"txhash": "0x003c8127556d023655168023988401be7cc46570be7713d42e8a9558c2ab1ae6",
		"apikey": "test-api-key",
	}

	// JSON response from actual API call with emitted logs
	sampleJSON := `{
        "jsonrpc": "2.0",
        "id": 1,
        "result": {
            "blockHash": "0x21ab72deeb4bb490bb3a6dc8ef46892e146a0c61b691354f5fa16c9dbf90b85f",
            "blockNumber": "0x13e5af1",
            "contractAddress": null,
            "cumulativeGasUsed": "0x1d9bc",
            "effectiveGasPrice": "0x16b86486ae",
            "from": "0x3d9aae030b9661e3605b3acb5d0385ede221a0cc",
            "gasUsed": "0x1d9bc",
            "status": "0x1",
            "to": "0x68d3a973e7272eb388022a5c6518d9b2a2e66fbf",
            "transactionHash": "0x003c8127556d023655168023988401be7cc46570be7713d42e8a9558c2ab1ae6",
            "transactionIndex": "0x0",
            "type": "0x2"
        }
    }`

	config := mockServerConfig{
		expectedParams: expectedParams,
		responseBody:   sampleJSON,
	}
	mockServer := createMockServer(config)
	defer mockServer.Close()

	client := initializeEtherscanClient(mockServer, "test-api-key", "")

	receipt, err := client.GetTransactionReceipt("0x003c8127556d023655168023988401be7cc46570be7713d42e8a9558c2ab1ae6")
	assert.NoError(t, err, "Expected no error from GetTransactionReceipt")

	// Define the expected ReceiptData
	expectedBlockNumber := uint64(0x13e5af1) // 20863729
	expectedGasUsed := uint64(0x1d9bc)       // 121276
	expectedGasPriceWei := new(big.Int)
	expectedGasPriceWei.SetString("0x16b86486ae", 0) // 97582876334

	assert.Equal(t, expectedBlockNumber, receipt.BlockNumber, "BlockNumber does not match")
	assert.Equal(t, "0x003c8127556d023655168023988401be7cc46570be7713d42e8a9558c2ab1ae6", receipt.Hash, "Transaction hash does not match")
	assert.Equal(t, expectedGasUsed, receipt.GasUsed, "GasUsed does not match")
	assert.Equal(t, expectedGasPriceWei, receipt.GasPriceWei, "GasPriceWei does not match")
}

func TestListTransactions(t *testing.T) {
	// Response of actual api call
	// https://api.etherscan.io/api%20%20%20?module=account&action=tokentx&page=1&offset=100&startblock=20871328&endblock=20871331&sort=desc&address=0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640

	sampleListTransactionsJSON := `{
        "status": "1",
        "message": "OK",
        "result": [
            {
                "blockNumber": "20871331",
                "timeStamp": "1727793983",
                "hash": "0xf508343089e789298f09e941e7c76bc500809e3f203b17d4d5769e263fa4d3f1",
                "nonce": "32072",
                "blockHash": "0x45a90ab0d934a425dd3e222da2772b59f826d7fedb14b012b37ac49122e18f70",
                "from": "0x68d3a973e7272eb388022a5c6518d9b2a2e66fbf",
                "contractAddress": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                "to": "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
                "value": "28379090239168221039",
                "tokenName": "Wrapped Ether",
                "tokenSymbol": "WETH",
                "tokenDecimal": "18",
                "transactionIndex": "3",
                "gas": "240958",
                "gasPrice": "34768791303",
                "gasUsed": "121242",
                "cumulativeGasUsed": "564598",
                "input": "deprecated",
                "confirmations": "300"
            },
            {
                "blockNumber": "20871331",
                "timeStamp": "1727793983",
                "hash": "0xf508343089e789298f09e941e7c76bc500809e3f203b17d4d5769e263fa4d3f1",
                "nonce": "32072",
                "blockHash": "0x45a90ab0d934a425dd3e222da2772b59f826d7fedb14b012b37ac49122e18f70",
                "from": "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
                "contractAddress": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
                "to": "0x68d3a973e7272eb388022a5c6518d9b2a2e66fbf",
                "value": "72819772278",
                "tokenName": "USDC",
                "tokenSymbol": "USDC",
                "tokenDecimal": "6",
                "transactionIndex": "3",
                "gas": "240958",
                "gasPrice": "34768791303",
                "gasUsed": "121242",
                "cumulativeGasUsed": "564598",
                "input": "deprecated",
                "confirmations": "300"
            },
            {
                "blockNumber": "20871328",
                "timeStamp": "1727793947",
                "hash": "0x8a4ed869c6b0ba8ed9543ec13f634a8105523eed2848a699c0b2150ae694bfc8",
                "nonce": "65",
                "blockHash": "0x49adaad0e17eabe786a7044a7138dc036367815b4f7126602400883ef591b060",
                "from": "0x8449e4198a021e8a2a5537c0508430b8febf8efc",
                "contractAddress": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
                "to": "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
                "value": "1200000000",
                "tokenName": "USDC",
                "tokenSymbol": "USDC",
                "tokenDecimal": "6",
                "transactionIndex": "27",
                "gas": "467430",
                "gasPrice": "24087796268",
                "gasUsed": "341965",
                "cumulativeGasUsed": "2854901",
                "input": "deprecated",
                "confirmations": "303"
            },
            {
                "blockNumber": "20871328",
                "timeStamp": "1727793947",
                "hash": "0x8a4ed869c6b0ba8ed9543ec13f634a8105523eed2848a699c0b2150ae694bfc8",
                "nonce": "65",
                "blockHash": "0x49adaad0e17eabe786a7044a7138dc036367815b4f7126602400883ef591b060",
                "from": "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
                "contractAddress": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                "to": "0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad",
                "value": "467118523758354530",
                "tokenName": "Wrapped Ether",
                "tokenSymbol": "WETH",
                "tokenDecimal": "18",
                "transactionIndex": "27",
                "gas": "467430",
                "gasPrice": "24087796268",
                "gasUsed": "341965",
                "cumulativeGasUsed": "2854901",
                "input": "deprecated",
                "confirmations": "303"
            }
        ]
    }`

	config := mockServerConfig{
		expectedParams: map[string]string{
			"module":     "account",
			"action":     "tokentx",
			"address":    "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
			"apikey":     "test-api-key",
			"sort":       "desc",
			"offset":     "10",
			"startblock": "1000000",
			"endblock":   "2000000",
		},
		responseBody: sampleListTransactionsJSON,
	}
	mockServer := createMockServer(config)
	defer mockServer.Close()

	client := initializeEtherscanClient(mockServer, "test-api-key", "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640")

	// Define sample input parameters
	offset := 10
	startBlock := uint64(1000000)
	endBlock := uint64(2000000)
	transactions, err := client.ListTransactions(&offset, &startBlock, &endBlock, nil)
	assert.NoError(t, err, "Expected no error from listTransactions")

	// Parse the Unix timestamps from the JSON to time.Time
	timestamp1 := time.Unix(1727793983, 0).Local()
	timestamp2 := time.Unix(1727793983, 0).Local()
	timestamp3 := time.Unix(1727793947, 0).Local()
	timestamp4 := time.Unix(1727793947, 0).Local()

	expectedGasPrice1, _ := new(big.Int).SetString("34768791303", 10)
	expectedGasPrice2, _ := new(big.Int).SetString("34768791303", 10)
	expectedGasPrice3, _ := new(big.Int).SetString("24087796268", 10)
	expectedGasPrice4, _ := new(big.Int).SetString("24087796268", 10)

	expectedTransactions := []types.TransactionData{
		{
			BlockNumber: 20871331,
			Hash:        "0xf508343089e789298f09e941e7c76bc500809e3f203b17d4d5769e263fa4d3f1",
			GasUsed:     121242,
			GasPriceWei: expectedGasPrice1,
			Timestamp:   timestamp1,
		},
		{
			BlockNumber: 20871331,
			Hash:        "0xf508343089e789298f09e941e7c76bc500809e3f203b17d4d5769e263fa4d3f1",
			GasUsed:     121242,
			GasPriceWei: expectedGasPrice2,
			Timestamp:   timestamp2,
		},
		{
			BlockNumber: 20871328,
			Hash:        "0x8a4ed869c6b0ba8ed9543ec13f634a8105523eed2848a699c0b2150ae694bfc8",
			GasUsed:     341965,
			GasPriceWei: expectedGasPrice3,
			Timestamp:   timestamp3,
		},
		{
			BlockNumber: 20871328,
			Hash:        "0x8a4ed869c6b0ba8ed9543ec13f634a8105523eed2848a699c0b2150ae694bfc8",
			GasUsed:     341965,
			GasPriceWei: expectedGasPrice4,
			Timestamp:   timestamp4,
		},
	}

	// Assertions to verify the correctness of the parsed data
	assert.Equal(t, expectedTransactions, transactions, "Transaction data does not match expected values")
}
