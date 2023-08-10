package exchange

import (
	"github.com/LMF709268224/titan-vps/lib/filecoinbridge"
	"github.com/LMF709268224/titan-vps/lib/trxbridge"
)

// GetGrpcClient
func getGrpcClient(addr string) (*trxbridge.GrpcClient, error) {
	node := trxbridge.NewGrpcClient(addr)
	err := node.Start()
	if err != nil {
		return nil, err
	}

	return node, nil
}

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

func getFilecoinHeight(addr string) int64 {
	var msg filecoinbridge.TipSet
	err := filecoinbridge.ChainHead(&msg,addr)
	if err != nil {
		log.Errorf("ChainHead err:%s", err.Error())
		return 0
	}

	return msg.Height
}