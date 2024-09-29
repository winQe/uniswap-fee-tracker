package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/winQe/uniswap-fee-tracker/internal/client"
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
