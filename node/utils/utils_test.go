package utils

import (
	"fmt"
	"testing"
)

func TestRate(t *testing.T) {
	r := GetUSDRate()
	fmt.Println("GetUSDRate :", r)
}

func TestBigIntReduce(t *testing.T) {
	s, e := BigIntReduce("123456", "548955.52")
	fmt.Println("BigIntReduce :", s)
	fmt.Println("BigIntReduce :", e)
}
