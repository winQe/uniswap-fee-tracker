package domain

import "time"

// PriceManagerInterface defines interface for price manager
type PriceManagerInterface interface {
	GetETHUSDT(timestamp time.Time) (float64, error)
}
