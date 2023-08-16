package transaction

import (
	"math/big"
	"strconv"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/trxbridge"
	"github.com/LMF709268224/titan-vps/lib/trxbridge/api"
	"github.com/LMF709268224/titan-vps/lib/trxbridge/core"
	"github.com/LMF709268224/titan-vps/lib/trxbridge/hexutil"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smirkcat/hdwallet"
	"golang.org/x/xerrors"
	"google.golang.org/protobuf/proto"
)

const checkBlockInterval = 3 * time.Second

// GetGrpcClient
func (m *Manager) getGrpcClient() (*trxbridge.GrpcClient, error) {
	node := trxbridge.NewGrpcClient(m.cfg.TrxHTTPSAddr)
	err := node.Start()
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (m *Manager) watchTronTransactions() {
	ticker := time.NewTicker(checkBlockInterval)
	defer ticker.Stop()

	client, err := m.getGrpcClient()
	if err != nil {
		log.Errorln("getGrpcClient err :", err.Error())
		return
	}

	startHeight := int64(39297600)
	limit := int64(50)
	heightStr := ""

	err = m.LoadConfigValue(db.ConfigTronHeight, &heightStr)
	if err == nil {
		i, err := strconv.ParseInt(heightStr, 10, 64)
		if err == nil {
			startHeight = i
		}
	}

	for {
		<-ticker.C

		block, err := client.GetNowBlock()
		if err != nil {
			log.Errorf("GetNowBlock err:%s", err.Error())
			continue
		}

		nowHeight := block.BlockHeader.RawData.Number
		endHeight := startHeight + limit
		if endHeight >= nowHeight {
			endHeight = nowHeight
		}
		if startHeight >= endHeight {
			continue
		}

		log.Debugf(" handleBlock height :%d, endHeight:%d \n", startHeight, endHeight)
		blockInfo, err := client.GetBlockByLimitNext(startHeight, endHeight)
		if err != nil {
			log.Errorf("GetBlockByLimitNext err:%s \n", err.Error())
			continue
		}
		m.handleBlocks(blockInfo)

		startHeight = endHeight
		str := strconv.FormatInt(startHeight, 10)
		err = m.SaveConfigValue(db.ConfigTronHeight, str)
		if err != nil {
			log.Errorf("SaveConfigValue err:%s \n", err.Error())
		}
	}
}

func (m *Manager) handleBlocks(blockInfo *api.BlockListExtention) {
	for _, v := range blockInfo.Block {
		err := m.handleBlock(v)
		if err != nil {
			log.Errorln(" handleBlock err :", err.Error())
		}
	}
}

func (m *Manager) handleBlock(blockExtention *api.BlockExtention) error {
	if blockExtention == nil || blockExtention.BlockHeader == nil {
		return xerrors.New("block is nil")
	}

	height := blockExtention.BlockHeader.RawData.Number

	for _, te := range blockExtention.Transactions {
		if len(te.Transaction.GetRet()) == 0 {
			continue
		}

		state := te.Transaction.GetRet()[0].ContractRet
		txID := hexutil.Encode(te.Txid)

		// userAddr := string(te.Transaction.RawData.Data)

		for _, contract := range te.Transaction.RawData.Contract {
			m.filterTransaction(contract, txID, height, state)
		}
	}

	return nil
}

func (m *Manager) filterTransaction(contract *core.Transaction_Contract, txID string, height int64, state core.Transaction_ResultContractResult) {
	if contract.Type == core.Transaction_Contract_TriggerSmartContract {
		// trc20
		unObj := &core.TriggerSmartContract{}
		err := proto.Unmarshal(contract.Parameter.GetValue(), unObj)
		if err != nil {
			// log.Errorf("parse trc20 err: %s", err.Error())
			return
		}

		contractAddress := hdwallet.EncodeCheck(unObj.GetContractAddress())

		if contractAddress != m.cfg.TrxContractorAddr {
			return
		}

		from := hdwallet.EncodeCheck(unObj.GetOwnerAddress())
		data := unObj.GetData()

		to, amount, isOk := m.decodeData(data)
		if !isOk {
			// log.Errorf("decodeData err: %s", txid)
			return
		}

		m.handleTransfer(txID, from, to, height, amount, state)
	}
}

func (m *Manager) decodeData(trc20 []byte) (to string, amount string, flag bool) {
	if len(trc20) >= 68 {
		if hexutil.Encode(trc20[:4]) != "a9059cbb" {
			return
		}
		trc20[15] = 65 // 0x41

		bb := common.TrimLeftZeroes(trc20[36:68])
		bu := new(big.Int).SetBytes(bb)
		amount = bu.String()

		to = hdwallet.EncodeCheck(trc20[15:36])
		flag = true
	}
	return
}

func (m *Manager) handleTransfer(txID, from, to string, height int64, amount string, state core.Transaction_ResultContractResult) {
	log.Infof("Transfer :%s,%s,%s,%s,%s,%s", txID, to, from, amount, state)

	if userID, ok := m.tronAddrs[to]; ok {
		m.notify.Pub(&types.TronTransferWatch{
			TxHash: txID,
			From:   from,
			To:     to,
			Value:  amount,
			State:  state,
			Height: height,
			UserID: userID,
		}, types.EventTronTransferWatch.String())
	}
}
