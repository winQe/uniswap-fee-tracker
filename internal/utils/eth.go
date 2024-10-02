package utils

import (
	"math/big"
	"regexp"
	"strings"
)

func SanitizeTransactionHash(hash string) string {
	// Remove any whitespace
	hash = strings.TrimSpace(hash)

	// Ensure the hash is in the correct format (0x followed by 64 hexadecimal characters)
	regex := regexp.MustCompile(`^0x[a-fA-F0-9]{64}$`)
	if !regex.MatchString(hash) {
		return ""
	}

	return hash
}

func ConvertToETH(gasPriceWei *big.Int) float64 {
	ethValue := new(big.Float).SetInt(gasPriceWei)
	ethValue.Mul(ethValue, big.NewFloat(1e-18))

	result, _ := ethValue.Float64()
	return result
}
