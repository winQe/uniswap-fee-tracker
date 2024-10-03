package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/winQe/uniswap-fee-tracker/internal/domain"
)

// MockTransactionManager is a mock implementation of the TransactionManagerInterface
type MockTransactionManager struct {
	mock.Mock
}

// GetLatestBlockNumber mocks the GetLatestBlockNumber method
func (m *MockTransactionManager) GetLatestBlockNumber() (uint64, error) {
	args := m.Called()
	return args.Get(0).(uint64), args.Error(1)
}

// GetTransaction mocks the GetTransaction method
func (m *MockTransactionManager) GetTransaction(hash string) (*domain.TxWithPrice, error) {
	args := m.Called(hash)
	return args.Get(0).(*domain.TxWithPrice), args.Error(1)
}

// BatchProcessTransactions mocks the BatchProcessTransactions method
func (m *MockTransactionManager) BatchProcessTransactions(startBlock uint64, endBlock uint64, ctx context.Context) ([]domain.TxWithPrice, error) {
	args := m.Called(startBlock, endBlock, ctx)
	return args.Get(0).([]domain.TxWithPrice), args.Error(1)
}

// BatchProcessTransactionsByTimestamp mocks the BatchProcessTransactionsByTimestamp method
func (m *MockTransactionManager) BatchProcessTransactionsByTimestamp(startTime time.Time, endTime time.Time, ctx context.Context) ([]domain.TxWithPrice, error) {
	args := m.Called(startTime, endTime, ctx)
	return args.Get(0).([]domain.TxWithPrice), args.Error(1)
}
