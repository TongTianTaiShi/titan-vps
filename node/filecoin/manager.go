package filecoin

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/filecoinfvm"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/filecoin-project/pubsub"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/crypto/sha3"
	"golang.org/x/xerrors"
)

var log = logging.Logger("fvm")

const (
	contractorAddr = "0x08906C7e01Bfb25483D9d411f1Fc18Df6b70a2F8"
	wsAddr         = "wss://wss.calibration.node.glif.io/apigw/lotus/rpc/v0"
	httpsAddr      = "https://api.calibration.node.glif.io/rpc/v1"
	privateKeyStr  = "3c3633bfaa3f8cfc2df9169d763eda6a8fb06d632e553f969f9dd2edd64dd11b"

	payeeAddress = "0xeb549F0B9887F4150dbD3bD0A257d99d5E316dBA" // need []string
)

// Manager is the node manager responsible for managing the online nodes
type Manager struct {
	notify *pubsub.PubSub
}

// NewManager creates a new instance of the node manager
func NewManager(pb *pubsub.PubSub) *Manager {
	manager := &Manager{
		notify: pb,
	}

	go manager.watchTransfer()

	return manager
}

func (m *Manager) watchTransfer() error {
	client, err := ethclient.Dial(wsAddr)
	if err != nil {
		return xerrors.Errorf("Dial err:%s", err.Error())
	}

	tokenAddress := common.HexToAddress(contractorAddr)

	myAbi, err := filecoinfvm.NewAbi(tokenAddress, client)
	if err != nil {
		return xerrors.Errorf("NewAbi err:%s", err.Error())
	}

	sink := make(chan *filecoinfvm.AbiTransfer)

	sub, err := myAbi.WatchTransfer(nil, sink, nil, nil)
	if err != nil {
		return xerrors.Errorf("Transfer err:%s", err.Error())
	}

	log.Debugf("tx sent: %s \n", sub)

	for {
		select {
		case err := <-sub.Err():
			if err != nil {
				log.Debugln(time.Now().Format("2006-01-02 15:04:05"), " err:", err)
				sub, err = myAbi.WatchTransfer(nil, sink, nil, nil)
				if err != nil {
					return xerrors.Errorf("Transfer err:%s", err.Error())
				}
			}
		case tr := <-sink:
			log.Infof("from:%s,to:%s,value:%d, RawTxHash:%s,RawBlockNumber:%d, Removed:%v \n", tr.From.String(), tr.To.String(), tr.Value, tr.Raw.TxHash.String(), tr.Raw.BlockNumber, tr.Raw.Removed)
			if !tr.Raw.Removed {
				m.notify.Pub(tr, types.EventTransfer.String())
			}
		}
	}
}

func (m *Manager) Check(tx string) error {
	return chainGetMessage(tx)
}

func (m *Manager) Transfer(toAddr, valueStr string) (string, error) {
	client, err := ethclient.Dial(httpsAddr)
	if err != nil {
		return "", xerrors.Errorf("Dial err:%s", err.Error())
	}

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", xerrors.Errorf("HexToECDSA err:%s", err.Error())
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", xerrors.New("publicKey err:")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	toAddress := common.HexToAddress(toAddr)
	tokenAddress := common.HexToAddress(contractorAddr)
	transferFnSignature := []byte("transfer(address,uint256)")

	myAbi, err := filecoinfvm.NewAbi(tokenAddress, client)
	if err != nil {
		return "", xerrors.Errorf("NewAbi err:%s", err.Error())
	}

	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	log.Debugln(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	log.Debugln(hexutil.Encode(paddedAddress))

	amount := new(big.Int)
	amount.SetString(valueStr, 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	log.Debugln(hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", xerrors.Errorf("NetworkID err:%s", err.Error())
	}

	signer := etypes.LatestSignerForChainID(chainID)
	to := &bind.TransactOpts{
		Signer: func(address common.Address, transaction *etypes.Transaction) (*etypes.Transaction, error) {
			return etypes.SignTx(transaction, signer, privateKey)
		},
		From:    fromAddress,
		Context: context.Background(),
		// GasLimit: gasLimit,
	}

	signedTx, err := myAbi.Transfer(to, toAddress, amount)
	if err != nil {
		return "", xerrors.Errorf("Transfer err:%s", err.Error())
	}

	log.Infof("tx sent: %s \n", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}
