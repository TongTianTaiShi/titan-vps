package main

import (
	"fmt"
	"os"

	"github.com/LMF709268224/titan-vps/node/orders"
	gen "github.com/whyrusleeping/cbor-gen"
)

func main() {
	err := gen.WriteMapEncodersToFile("../cbor_gen.go", "orders",
		orders.OrderInfo{},
		orders.GoodsInfo{},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
