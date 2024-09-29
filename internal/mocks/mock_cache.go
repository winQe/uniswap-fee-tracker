package mocks

import "github.com/stretchr/testify/mock"

// MockRateCache is a mock implementation of the RateCache interface.
type MockRateCache struct {
	mock.Mock
}

// GetRate mocks the GetRate method of RateCache.
func (m *MockRateCache) GetRate(key string) (float64, error) {
	args := m.Called(key)
	return args.Get(0).(float64), args.Error(1)
}

// StoreRate mocks the StoreRate method of RateCache.
func (m *MockRateCache) StoreRate(key string, value float64) error {
	args := m.Called(key, value)
	return args.Error(0)
}
