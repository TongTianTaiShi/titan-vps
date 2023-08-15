package utils

import "math/big"

// BigIntAdd add
func BigIntAdd(numstr string, value string) string {
	n, _ := new(big.Int).SetString(numstr, 10)
	m, _ := new(big.Int).SetString(value, 10)
	m.Add(n, m)
	return m.String()
}

// BigIntReduce reduce
func BigIntReduce(numstr string, value string) string {
	n, _ := new(big.Int).SetString(numstr, 10)
	m, _ := new(big.Int).SetString(value, 10)
	z := new(big.Int).Sub(n, m)
	return z.String()
}
