package transaction

import (
	"context"
	"crypto/ecdsa"

	"github.com/LMF709268224/titan-vps/lib/filecoinbridge"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/xerrors"
)

// func (m *Manager) watchFvmTransactions() error {
// 	client, err := ethclient.Dial(m.cfg.LotusWsAddr)
// 	if err != nil {
// 		return xerrors.Errorf("Dial err:%s", err.Error())
// 	}

// 	tokenAddress := common.HexToAddress(m.cfg.TitanContractorAddr)

// 	fAbi, err := filecoinbridge.NewFvm(tokenAddress, client)
// 	if err != nil {
// 		return xerrors.Errorf("NewAbi err:%s", err.Error())
// 	}

// 	sink := make(chan *filecoinbridge.AbiTransfer)

// 	sub, err := fAbi.WatchTransfer(nil, sink, nil, nil)
// 	if err != nil {
// 		return xerrors.Errorf("WatchTransfer err:%s", err.Error())
// 	}

// 	for {
// 		select {
// 		case err := <-sub.Err():
// 			if err != nil {
// 				log.Debugln(time.Now().Format("2006-01-02 15:04:05"), " err:", err)
// 				sub, err = fAbi.WatchTransfer(nil, sink, nil, nil)
// 				if err != nil {
// 					return xerrors.Errorf("WatchTransfer err:%s", err.Error())
// 				}
// 			}
// 		case tr := <-sink:
// 			log.Debugf("from:%s,to:%s,value:%d, RawTxHash:%s,RawBlockNumber:%d, Removed:%v \n", tr.From.String(), tr.To.Hex(), tr.Value, tr.Raw.TxHash.String(), tr.Raw.BlockNumber, tr.Raw.Removed)
// 			if !tr.Raw.Removed {
// 				m.notification.Pub(&types.FvmTransferWatch{
// 					TxHash: tr.Raw.TxHash.Hex(),
// 					From:   tr.From.Hex(),
// 					To:     tr.To.Hex(),
// 					Value:  tr.Value.String(),
// 				}, types.EventFvmTransferWatch.String())
// 			}
// 		}
// 	}
// }

func (m *Manager) SendMsg(info filecoinbridge.IpcOrderInfo) error {
	client, err := ethclient.Dial(m.cfg.LotusWsAddr)
	if err != nil {
		return xerrors.Errorf("Dial err:%s", err.Error())
	}

	tokenAddress := common.HexToAddress(m.cfg.TitanContractorAddr)

	fAbi, err := filecoinbridge.NewFvm(tokenAddress, client)
	if err != nil {
		return xerrors.Errorf("NewAbi err:%s", err.Error())
	}

	privateKey, err := crypto.HexToECDSA(m.cfg.PrivateKeyStr)
	if err != nil {
		return xerrors.Errorf("HexToECDSA err:%s", err.Error())
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return xerrors.New("publicKey err:")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return xerrors.Errorf("NetworkID err:%s", err.Error())
	}

	signer := etypes.LatestSignerForChainID(chainID)
	opt := &bind.TransactOpts{
		Signer: func(address common.Address, transaction *etypes.Transaction) (*etypes.Transaction, error) {
			return etypes.SignTx(transaction, signer, privateKey)
		},
		From:    fromAddress,
		Context: context.Background(),
		// GasLimit: gasLimit,
	}

	tr, err := fAbi.SetOrderInfo(opt, info)
	if err != nil {
		return xerrors.Errorf("SetOrderInfo err:%s", err.Error())
	}

	log.Infof("SetOrderInfo tr:%v", tr)

	return nil
}
