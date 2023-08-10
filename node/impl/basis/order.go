package basis

import (
	"context"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/filecoinbridge"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

func (m *Basis) CreateOrder(ctx context.Context, req types.CreateOrderReq) (string, error) {
	id, err := m.SaveVpsInstance(&req.Vps)
	if err != nil {
		return "", err
	}

	info := &types.OrderRecord{
		VpsID: id,
		User:  req.User,
		Value: "10000000000",
	}

	err = m.OrderMgr.CreatedOrder(info)
	if err != nil {
		return "", err
	}

	return info.To, nil
}

func (m *Basis) PaymentCompleted(ctx context.Context, req types.PaymentCompletedReq) (string, error) {
	record, err := m.LoadOrderRecord(req.OrderID)
	if err != nil {
		return "", err
	}

	if record.State != types.WaitingPayment {
		return "", xerrors.Errorf("Invalid order status %d", record.State)
	}

	cfg, err := m.GetBasisConfigFunc()
	if err != nil {
		return "", err
	}

	tx := req.TransactionID
	log.Debugf("tx:%s \n", tx)
	var cid cid.Cid
	err = filecoinbridge.EthGetMessageCidByTransactionHash(&cid, tx, cfg.LotusHTTPSAddr)
	if err != nil {
		return "", err
	}

	log.Debugf("cid:%s \n", cid.String())

	var msg filecoinbridge.Message
	err = filecoinbridge.ChainGetMessage(&msg, cid, cfg.LotusHTTPSAddr)
	if err != nil {
		return "", err
	}

	var info filecoinbridge.Lookup
	err = filecoinbridge.StateSearchMsg(&info, cid, cfg.LotusHTTPSAddr)
	if err != nil {
		return "", err
	}

	log.Debugf("Height:%d,ExitCode:%d,GasUsed:%d \n", info.Height, info.Receipt.ExitCode, info.Receipt.GasUsed)

	if info.Receipt.ExitCode == 0 {
		m.Notify.Pub(&types.FvmTransferWatch{
			TxHash: tx,
			From:   msg.From.String(),
			To:     msg.To.String(),
			Value:  msg.Value.String(),
		}, types.EventFvmTransferWatch.String())
	}

	return "", nil
}

func (m *Basis) CancelOrder(ctx context.Context, orderID string) error {
	return m.OrderMgr.CancelOrder(orderID)
}
