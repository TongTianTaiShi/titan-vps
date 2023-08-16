package exchange

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/LMF709268224/titan-vps/lib/trxbridge/hexutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smirkcat/hdwallet"
	"golang.org/x/crypto/sha3"
)

func TestWatch(t *testing.T) {
	client, err := getGrpcClient("47.252.19.181:50051")
	if err != nil {
		fmt.Println("getGrpcClient err:", err.Error())
		return
	}

	prikey := ""
	toAddr := "TNXS7Xybbq8ZKiueGWomNNoUWqGhHCT1qe"
	valueStr := "23456789000000"
	privateKey, err := hdwallet.GetPrivateKeyByHexString(prikey)
	if err != nil {
		fmt.Println("GetPrivateKeyByHexString err:", err.Error())
		return
	}

	toAddress := common.HexToAddress(toAddr)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress))

	amount := new(big.Int)
	amount.SetString(valueStr, 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	tHash, err := client.TransferContract(privateKey, "TLBaRhANQoJFTqre9Nf1mjuwNWjCJeYqUL", data, 0)
	if err != nil {
		fmt.Println("Transfer err:", err.Error())
		return
	}

	fmt.Println("Transfer hash:", tHash)
}

func TestCreateAddr(t *testing.T) {
	for i := 0; i < 10; i++ {

		privateKey, err := hdwallet.NewPrivateKey("")
		if err != nil {
			fmt.Println("NewPrivateKey err : ", err)
			return
		}

		address := hdwallet.PrikeyToAddressTron(privateKey)
		prikey := hdwallet.PrikeyToHexString(privateKey)

		fmt.Println("\"", address, "\", #", prikey)
	}
}
