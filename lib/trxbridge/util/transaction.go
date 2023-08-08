package util

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"time"

	"github.com/LMF709268224/titan-vps/lib/trxbridge/core"

	"github.com/ethereum/go-ethereum/crypto"

	"google.golang.org/protobuf/proto"
)

// SignTransaction
func SignTransaction(transaction *core.Transaction, key *ecdsa.PrivateKey) ([]byte, error) {
	transaction.GetRawData().Timestamp = time.Now().UnixNano() / 1000000
	rawData, err := proto.Marshal(transaction.GetRawData())
	if err != nil {
		return nil, err
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	contractList := transaction.GetRawData().GetContract()
	for range contractList {
		signature, err := crypto.Sign(hash, key)
		if err != nil {
			return nil, err
		}
		transaction.Signature = append(transaction.Signature, signature)
	}
	return hash, nil
}
