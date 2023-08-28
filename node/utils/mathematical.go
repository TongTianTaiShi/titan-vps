package utils

import (
	"fmt"
	"math/big"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"
)

// BigIntAdd add
func BigIntAdd(numstr string, value string) (string, error) {
	n, _ := new(big.Int).SetString(numstr, 10)
	m, _ := new(big.Int).SetString(value, 10)

	if n == nil {
		return "0", &api.ErrWeb{Code: terrors.EncodingError.Int(), Message: fmt.Sprintf("BigIntReduce err num is %s", numstr)}
	}

	if m == nil {
		return "0", &api.ErrWeb{Code: terrors.EncodingError.Int(), Message: fmt.Sprintf("BigIntReduce err num is %s", value)}
	}

	m.Add(n, m)
	return m.String(), nil
}

// BigIntReduce reduce
func BigIntReduce(numstr string, value string) (string, error) {
	n, _ := new(big.Int).SetString(numstr, 10)
	m, _ := new(big.Int).SetString(value, 10)

	if n == nil {
		return "0", &api.ErrWeb{Code: terrors.EncodingError.Int(), Message: fmt.Sprintf("BigIntReduce err num is %s", numstr)}
	}

	if m == nil {
		return "0", &api.ErrWeb{Code: terrors.EncodingError.Int(), Message: fmt.Sprintf("BigIntReduce err num is %s", value)}
	}

	if n.Cmp(m) < 0 {
		return "0", &api.ErrWeb{Code: terrors.InsufficientBalance.Int(), Message: terrors.InsufficientBalance.String()}
	}
	z := new(big.Int).Sub(n, m)

	return z.String(), nil
}
