package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockRateCache is a mock implementation of the RateCache interface.
type MockRateCache struct {
	mock.Mock
}

// GetRate mocks the GetRate method of RateCache.
func (m *MockRateCache) GetRate(timestamp time.Time) (float64, error) {
	args := m.Called(timestamp)
	return args.Get(0).(float64), args.Error(1)
}

// StoreRate mocks the StoreRate method of RateCache.
func (m *MockRateCache) StoreRate(timestamp time.Time, value float64) error {
	args := m.Called(timestamp, value)
	return args.Error(0)
}
