package utils

import (
	"fmt"
	"math/big"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"
)

// AddBigInt adds two big integers represented as strings.
func AddBigInt(numstr string, value string) (string, error) {
	// Convert input strings to big.Int
	n, nOk := new(big.Int).SetString(numstr, 10)
	m, mOk := new(big.Int).SetString(value, 10)

	// Check for conversion errors
	if !nOk || n == nil {
		return "0", &api.ErrWeb{Code: terrors.EncodingError.Int(), Message: fmt.Sprintf("AddBigInt error: invalid num %s", numstr)}
	}

	if !mOk || m == nil {
		return "0", &api.ErrWeb{Code: terrors.EncodingError.Int(), Message: fmt.Sprintf("AddBigInt error: invalid value %s", value)}
	}

	// Perform addition
	m.Add(n, m)
	return m.String(), nil
}

// ReduceBigInt subtracts one big integer represented as a string from another.
func ReduceBigInt(numstr string, value string) (string, error) {
	// Convert input strings to big.Int
	n, nOk := new(big.Int).SetString(numstr, 10)
	m, mOk := new(big.Int).SetString(value, 10)

	// Check for conversion errors
	if !nOk || n == nil {
		return "0", &api.ErrWeb{Code: terrors.EncodingError.Int(), Message: fmt.Sprintf("ReduceBigInt error: invalid num %s", numstr)}
	}

	if !mOk || m == nil {
		return "0", &api.ErrWeb{Code: terrors.EncodingError.Int(), Message: fmt.Sprintf("ReduceBigInt error: invalid value %s", value)}
	}

	// Check if n is less than m
	if n.Cmp(m) < 0 {
		return "0", &api.ErrWeb{Code: terrors.InsufficientBalance.Int(), Message: terrors.InsufficientBalance.String()}
	}

	// Perform subtraction
	result := new(big.Int).Sub(n, m)
	return result.String(), nil
}
