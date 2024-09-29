package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/winQe/uniswap-fee-tracker/internal/client"
	"github.com/winQe/uniswap-fee-tracker/internal/mocks"
)

func TestPriceManager_GetETHUSDT(t *testing.T) {
	timestamp := time.Now()

	t.Run("cache hits", func(t *testing.T) {
		mockCache := new(mocks.MockRateCache)
		mockClient := new(mocks.MockPriceClient)

		// Simulate cache hit, external API shouldn't be called
		mockCache.On("GetRate", "eth-usdt").Return(4800.75, nil)
		mockClient.AssertNotCalled(t, "GetETHUSDT")

		priceManager := NewPriceManager(mockCache, mockClient)
		price, err := priceManager.GetETHUSDT(timestamp)

		assert.NoError(t, err)
		assert.Equal(t, 4800.75, price)

		mockCache.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("cache miss, valid external API response", func(t *testing.T) {
		mockCache := new(mocks.MockRateCache)
		mockClient := new(mocks.MockPriceClient)

		// Simulate cache miss, valid external API response, and storing in cache
		mockCache.On("GetRate", "eth-usdt").Return(0.0, errors.New("cache miss"))
		mockClient.On("GetETHUSDT", timestamp).Return(&client.KlineData{ClosePrice: 1850.00}, nil)
		mockCache.On("StoreRate", "eth-usdt", 1850.00).Return(nil)

		priceManager := NewPriceManager(mockCache, mockClient)
		price, err := priceManager.GetETHUSDT(timestamp)

		assert.NoError(t, err)
		assert.Equal(t, 1850.00, price)

		mockCache.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("cache miss, external API error", func(t *testing.T) {
		mockCache := new(mocks.MockRateCache)
		mockClient := new(mocks.MockPriceClient)

		// Simulate cache miss, external API API failure
		mockCache.On("GetRate", "eth-usdt").Return(0.0, errors.New("cache miss"))
		mockClient.On("GetETHUSDT", timestamp).Return((*client.KlineData)(nil), errors.New("external API error"))

		priceManager := NewPriceManager(mockCache, mockClient)
		price, err := priceManager.GetETHUSDT(timestamp)

		assert.Error(t, err)
		assert.Equal(t, 0.0, price)
		assert.Contains(t, err.Error(), "external API error")

		mockCache.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("cache miss, store error", func(t *testing.T) {
		mockCache := new(mocks.MockRateCache)
		mockClient := new(mocks.MockPriceClient)

		// Simulate cache miss, valid external API response, but cache store fails
		mockCache.On("GetRate", "eth-usdt").Return(0.0, errors.New("cache miss"))
		mockClient.On("GetETHUSDT", timestamp).Return(&client.KlineData{ClosePrice: 1850.00}, nil)
		mockCache.On("StoreRate", "eth-usdt", 1850.00).Return(errors.New("could not store in cache"))

		priceManager := NewPriceManager(mockCache, mockClient)
		price, err := priceManager.GetETHUSDT(timestamp)

		assert.NoError(t, err)
		assert.Equal(t, 1850.00, price)

		mockCache.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})
}
