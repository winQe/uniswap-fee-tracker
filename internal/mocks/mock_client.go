package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/winQe/uniswap-fee-tracker/internal/client"
	"github.com/winQe/uniswap-fee-tracker/internal/types"
)

// MockPriceClient is a mock implementation of the PriceClient interface.
type MockPriceClient struct {
	mock.Mock
}

// GetETHUSDT mocks the GetETHUSDT method of PriceClient.
func (m *MockPriceClient) GetETHUSDT(timestamp time.Time) (*client.KlineData, error) {
	args := m.Called(timestamp)
	return args.Get(0).(*client.KlineData), args.Error(1)
}

// MockTransactionClient is a mock implementation of the TransactionClient interface.
type MockTransactionClient struct {
	mock.Mock
}

// ListTransactions mocks the ListTransactions method.
func (m *MockTransactionClient) ListTransactions(batchSize *uint64, startBlock *uint64, endBlock *uint64, page *int) ([]types.TransactionData, error) {
	args := m.Called(batchSize, startBlock, endBlock, page)
	return args.Get(0).([]types.TransactionData), args.Error(1)
}
