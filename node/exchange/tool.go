package exchange

import (
	"github.com/LMF709268224/titan-vps/lib/filecoinbridge"
	"github.com/LMF709268224/titan-vps/lib/trxbridge"
)

// NewGrpcClient creates and starts a gRPC client for a given address.
func getGrpcClient(addr string) (*trxbridge.GrpcClient, error) {
	node := trxbridge.NewGrpcClient(addr)
	err := node.Start()
	if err != nil {
		return nil, err
	}

	return node, nil
}

// GetTronHeight retrieves the current block height of a Tron node at the specified address.
func getTronHeight(addr string) int64 {
	client, err := getGrpcClient(addr)
	if err != nil {
		log.Errorln("getGrpcClient err :", err.Error())
		return 0
	}

	block, err := client.GetNowBlock()
	if err != nil {
		log.Errorln("GetNowBlock err :", err.Error())
		return 0
	}

	return block.GetBlockHeader().RawData.Number
}

// GetFilecoinHeight retrieves the current block height of a Filecoin node at the specified address.
func getFilecoinHeight(addr string) int64 {
	var msg filecoinbridge.TipSet
	err := filecoinbridge.ChainHead(&msg, addr)
	if err != nil {
		log.Errorf("ChainHead err:%s", err.Error())
		return 0
	}

	return msg.Height
}
