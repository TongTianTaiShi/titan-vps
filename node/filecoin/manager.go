package filecoin

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/filecoinfvm"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/filecoin-project/pubsub"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/crypto/sha3"
	"golang.org/x/xerrors"
)

var log = logging.Logger("filecoin")

// Manager is the node manager responsible for managing the online nodes
type Manager struct {
	notify *pubsub.PubSub

	cfg config.BasisCfg
}

// NewManager creates a new instance of the node manager
func NewManager(pb *pubsub.PubSub, getCfg dtypes.GetBasisConfigFunc) (*Manager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		notify: pb,
		cfg:    cfg,
	}

	go manager.watchTransfer()

	return manager, nil
}

func (m *Manager) watchTransfer() error {
	client, err := ethclient.Dial(m.cfg.LotusWsAddr)
	if err != nil {
		return xerrors.Errorf("Dial err:%s", err.Error())
	}

	tokenAddress := common.HexToAddress(m.cfg.ContractorAddr)

	myAbi, err := filecoinfvm.NewAbi(tokenAddress, client)
	if err != nil {
		return xerrors.Errorf("NewAbi err:%s", err.Error())
	}

	sink := make(chan *filecoinfvm.AbiTransfer)

	sub, err := myAbi.WatchTransfer(nil, sink, nil, nil)
	if err != nil {
		return xerrors.Errorf("Transfer err:%s", err.Error())
	}

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
				m.notify.Pub(&types.FvmTransfer{
					ID:    tr.Raw.TxHash.Hex(),
					From:  tr.From.Hex(),
					To:    tr.To.Hex(),
					Value: tr.Value.Int64(),
				}, types.EventTransfer.String())
			}
		}
	}
}

func (m *Manager) GetHeight() int64 {
	var msg tipSet
	err := filecoinfvm.ChainHead(&msg, m.cfg.LotusHTTPSAddr)
	if err != nil {
		log.Errorf("ChainHead err:%s", err.Error())
		return 0
	}

	return msg.Height
}

// GetBalance get balance
func (m *Manager) GetBalance(addr string) (*big.Int, error) {
	client, err := ethclient.Dial(m.cfg.LotusHTTPSAddr)
	if err != nil {
		return big.NewInt(0), xerrors.Errorf("Dial err:%s", err.Error())
	}

	tokenAddress := common.HexToAddress(m.cfg.ContractorAddr)

	myAbi, err := filecoinfvm.NewAbi(tokenAddress, client)
	if err != nil {
		return big.NewInt(0), xerrors.Errorf("NewAbi err:%s", err.Error())
	}

	return myAbi.BalanceOf(nil, common.HexToAddress(addr))
}

// CheckMessage check
func (m *Manager) CheckMessage(tx string) error {
	log.Debugf("tx:%s \n", tx)
	var cid cid.Cid
	err := filecoinfvm.EthGetMessageCidByTransactionHash(&cid, tx, m.cfg.LotusHTTPSAddr)
	if err != nil {
		return err
	}

	log.Debugf("cid:%s \n", cid.String())

	var msg message
	err = filecoinfvm.ChainGetMessage(&msg, cid, m.cfg.LotusHTTPSAddr)
	if err != nil {
		return err
	}

	var info lookup
	err = filecoinfvm.StateSearchMsg(&info, cid, m.cfg.LotusHTTPSAddr)
	if err != nil {
		return err
	}

	log.Debugf("Height:%d,ExitCode:%d,GasUsed:%d \n", info.Height, info.Receipt.ExitCode, info.Receipt.GasUsed)

	if info.Receipt.ExitCode == 0 {
		m.notify.Pub(&types.FvmTransfer{
			ID:    tx,
			From:  msg.From.String(),
			To:    msg.To.String(),
			Value: msg.Value.Int64(),
		}, types.EventTransfer.String())
	}

	return nil
}

func (m *Manager) Mint(toAddr string) (string, error) {
	client, err := ethclient.Dial(m.cfg.LotusHTTPSAddr)
	if err != nil {
		return "", xerrors.Errorf("Dial err:%s", err.Error())
	}

	tokenAddress := common.HexToAddress(m.cfg.ContractorAddr)

	myAbi, err := filecoinfvm.NewAbi(tokenAddress, client)
	if err != nil {
		return "", xerrors.Errorf("NewAbi err:%s", err.Error())
	}

	privateKey, err := crypto.HexToECDSA(m.cfg.PrivateKeyStr)
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
	transferFnSignature := []byte("transfer(address,uint256)")

	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	log.Debugln(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	log.Debugln(hexutil.Encode(paddedAddress))

	amount := big.NewInt(9000000000000000000)
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

	tr, err := myAbi.Mint(to, common.HexToAddress(toAddr), amount)
	if err != nil {
		return "", xerrors.Errorf("Mint err:%s", err.Error())
	}

	log.Infoln(tr)

	return tr.Hash().Hex(), nil
}

// Transfer transfer to
func (m *Manager) Transfer(toAddr, valueStr string) (string, error) {
	client, err := ethclient.Dial(m.cfg.LotusHTTPSAddr)
	if err != nil {
		return "", xerrors.Errorf("Dial err:%s", err.Error())
	}

	privateKey, err := crypto.HexToECDSA(m.cfg.PrivateKeyStr)
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
	tokenAddress := common.HexToAddress(m.cfg.ContractorAddr)
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
