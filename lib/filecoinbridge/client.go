package filecoinbridge

// "github.com/ethereum/go-ethereum/accounts/abi/bind"
// "github.com/ethereum/go-ethereum/common"
// etypes "github.com/ethereum/go-ethereum/core/types"
// "github.com/ethereum/go-ethereum/crypto"
// "github.com/ethereum/go-ethereum/ethclient"
// "golang.org/x/xerrors"

// GrpcClient GrpcClient
type GrpcClient struct {
	Url        string
	Contractor string
}

// NewGrpcClient NewGrpcClient
func NewGrpcClient(url, contractorAddr string) *GrpcClient {
	client := new(GrpcClient)
	client.Url = url
	client.Contractor = contractorAddr
	return client
}

// GetBalance get balance
// func (g *GrpcClient) GetBalance(addr string) (*big.Int, error) {
// 	client, err := ethclient.Dial(g.Url)
// 	if err != nil {
// 		return big.NewInt(0), xerrors.Errorf("Dial err:%s", err.Error())
// 	}

// 	tokenAddress := common.HexToAddress(g.Contractor)

// 	myAbi, err := NewFvm(tokenAddress, client)
// 	if err != nil {
// 		return big.NewInt(0), xerrors.Errorf("NewAbi err:%s", err.Error())
// 	}

// 	return myAbi.BalanceOf(nil, common.HexToAddress(addr))
// }

// func (g *GrpcClient) Mint(privateKeyStr, toAddr, valueStr string) (string, error) {
// 	client, err := ethclient.Dial(g.Url)
// 	if err != nil {
// 		return "", xerrors.Errorf("Dial err:%s", err.Error())
// 	}

// 	tokenAddress := common.HexToAddress(g.Contractor)

// 	myAbi, err := NewFvm(tokenAddress, client)
// 	if err != nil {
// 		return "", xerrors.Errorf("NewAbi err:%s", err.Error())
// 	}

// 	privateKey, err := crypto.HexToECDSA(privateKeyStr)
// 	if err != nil {
// 		return "", xerrors.Errorf("HexToECDSA err:%s", err.Error())
// 	}

// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		return "", xerrors.New("publicKey err:")
// 	}

// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
// 	amount := new(big.Int)
// 	amount.SetString(valueStr, 10)

// 	chainID, err := client.NetworkID(context.Background())
// 	if err != nil {
// 		return "", xerrors.Errorf("NetworkID err:%s", err.Error())
// 	}

// 	signer := etypes.LatestSignerForChainID(chainID)
// 	opt := &bind.TransactOpts{
// 		Signer: func(address common.Address, transaction *etypes.Transaction) (*etypes.Transaction, error) {
// 			return etypes.SignTx(transaction, signer, privateKey)
// 		},
// 		From:    fromAddress,
// 		Context: context.Background(),
// 		// GasLimit: gasLimit,
// 	}

// 	tr, err := myAbi.Mint(opt, common.HexToAddress(toAddr), amount)
// 	if err != nil {
// 		return "", xerrors.Errorf("Mint err:%s", err.Error())
// 	}

// 	return tr.Hash().Hex(), nil
// }

// Transfer transfer to
// func (g *GrpcClient) Transfer(privateKeyStr, toAddr, valueStr string) (string, error) {
// 	client, err := ethclient.Dial(g.Url)
// 	if err != nil {
// 		return "", xerrors.Errorf("Dial err:%s", err.Error())
// 	}

// 	privateKey, err := crypto.HexToECDSA(privateKeyStr)
// 	if err != nil {
// 		return "", xerrors.Errorf("HexToECDSA err:%s", err.Error())
// 	}

// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		return "", xerrors.New("publicKey err:")
// 	}

// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
// 	toAddress := common.HexToAddress(toAddr)
// 	tokenAddress := common.HexToAddress(g.Contractor)

// 	myAbi, err := NewFvm(tokenAddress, client)
// 	if err != nil {
// 		return "", xerrors.Errorf("NewAbi err:%s", err.Error())
// 	}

// 	amount := new(big.Int)
// 	amount.SetString(valueStr, 10)

// 	chainID, err := client.NetworkID(context.Background())
// 	if err != nil {
// 		return "", xerrors.Errorf("NetworkID err:%s", err.Error())
// 	}

// 	signer := etypes.LatestSignerForChainID(chainID)
// 	to := &bind.TransactOpts{
// 		Signer: func(address common.Address, transaction *etypes.Transaction) (*etypes.Transaction, error) {
// 			return etypes.SignTx(transaction, signer, privateKey)
// 		},
// 		From:    fromAddress,
// 		Context: context.Background(),
// 		// GasLimit: gasLimit,
// 	}

// 	signedTx, err := myAbi.Transfer(to, toAddress, amount)
// 	if err != nil {
// 		return "", xerrors.Errorf("Transfer err:%s", err.Error())
// 	}

// 	return signedTx.Hash().Hex(), nil
// }
