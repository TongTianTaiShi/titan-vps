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
func BigIntReduce(numstr string, value string) (string, bool) {
	n, _ := new(big.Int).SetString(numstr, 10)
	m, _ := new(big.Int).SetString(value, 10)

	if n.Cmp(m) < 0 {
		return "0", false
	}
	z := new(big.Int).Sub(n, m)

	return z.String(), true
}
