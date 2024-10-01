package client

import (
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestGetTransactionReceipt(t *testing.T) {
	// JSON response from actual API calls
	// https://api.etherscan.io/api?module=proxy&action=eth_getTransactionReceipt&txhash=0x003c8127556d023655168023988401be7cc46570be7713d42e8a9558c2ab1ae6&apikey=MRW44A62CVQTZU28FBCKT248ZKV5KHHC6Z
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

	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request parameters
		expectedModule := "proxy"
		expectedAction := "eth_getTransactionByHash"
		query := r.URL.Query()
		assert.Equal(t, expectedModule, query.Get("module"), "Module parameter mismatch")
		assert.Equal(t, expectedAction, query.Get("action"), "Action parameter mismatch")
		assert.Equal(t, "test-api-key", query.Get("apikey"), "API key mismatch")
		assert.Equal(t, "0x003c8127556d023655168023988401be7cc46570be7713d42e8a9558c2ab1ae6", query.Get("txhash"), "Transaction hash mismatch")

		// Return the sample JSON response
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, sampleJSON)
	}))
	defer mockServer.Close()

	rateLimiter := rate.NewLimiter(5, 5)
	dailyLimiter := rate.NewLimiter(1.15, 1000)
	rateLimitedClient := NewRateLimitedClient(rateLimiter, dailyLimiter)

	// Override the httpClient to use the mock server
	rateLimitedClient.httpClient = mockServer.Client()

	// Initialize EtherscanClient with the mock server's URL and mock RateLimitedClient
	client := &EtherscanClient{
		RateLimitedClient: rateLimitedClient,
		baseURL:           mockServer.URL,
		apiKey:            "test-api-key",
	}

	// Call the GetTransactionReceipt method with a sample transaction hash
	receipt, err := client.GetTransactionReceipt("0x003c8127556d023655168023988401be7cc46570be7713d42e8a9558c2ab1ae6")
	assert.NoError(t, err, "Expected no error from GetTransactionReceipt")

	// Define the expected ReceiptData
	expectedBlockNumber := uint64(0x13e5af1) // 20863729
	expectedGasUsed := uint64(0x1d9bc)       // 121276
	expectedGasPriceWei := new(big.Int)
	expectedGasPriceWei.SetString("0x16b86486ae", 0) // 0x16b86486ae = 97582876334

	// Assertions to verify the correctness of the parsed data
	assert.Equal(t, expectedBlockNumber, receipt.BlockNumber, "BlockNumber does not match")
	assert.Equal(t, "0x003c8127556d023655168023988401be7cc46570be7713d42e8a9558c2ab1ae6", receipt.Hash, "Transaction hash does not match")
	assert.Equal(t, expectedGasUsed, receipt.GasUsed, "GasUsed does not match")
	assert.Equal(t, expectedGasPriceWei, receipt.GasPriceWei, "GasPriceWei does not match")
}

