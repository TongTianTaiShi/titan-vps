package filecoinbridge

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs/go-cid"
)

type (
	// lotus struct
	TipSet struct {
		Height int64
	}

	Message struct {
		Version uint64

		To   address.Address
		From address.Address

		Nonce uint64

		Value big.Int

		GasLimit   int64
		GasFeeCap  big.Int
		GasPremium big.Int

		Method uint64
		Params []byte
	}

	MessageReceipt struct {
		ExitCode int64
		GasUsed  int64
	}

	Lookup struct {
		Message   cid.Cid // Can be different than requested, in case it was replaced, but only gas values changed
		Receipt   MessageReceipt
		ReturnDec interface{}
		Height    int64
	}
)
