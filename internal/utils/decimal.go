package utils

import (
	"math/big"
	
	"github.com/shopspring/decimal"
)

// DecimalFromBigInt converts a big.Int with given exponent to decimal
func DecimalFromBigInt(value *big.Int, exp int32) decimal.Decimal {
	if value == nil {
		return decimal.Zero
	}
	return decimal.NewFromBigInt(value, exp)
}

// BigIntFromDecimal converts a decimal to big.Int with given exponent
func BigIntFromDecimal(value decimal.Decimal, exp int32) *big.Int {
	// Multiply by 10^(-exp) to get the integer representation
	multiplier := decimal.New(1, -exp)
	result := value.Mul(multiplier)
	return result.BigInt()
}

// DecimalFromString safely converts a string to decimal
func DecimalFromString(value string) decimal.Decimal {
	result, err := decimal.NewFromString(value)
	if err != nil {
		return decimal.Zero
	}
	return result
}