package orders

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/aliyun"
	"github.com/LMF709268224/titan-vps/lib/filecoinbridge"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/LMF709268224/titan-vps/node/transaction"
	"github.com/filecoin-project/go-statemachine"
	"github.com/filecoin-project/pubsub"
	"github.com/ipfs/go-datastore"
	"golang.org/x/xerrors"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("orders")

const (
	checkOrderInterval = 10 * time.Second
	orderTimeoutTime   = 10 * time.Minute
)

// Manager manager order
type Manager struct {
	stateMachineWait   sync.WaitGroup
	orderStateMachines *statemachine.StateGroup
	*db.SQLDB

	notify *pubsub.PubSub

	ongoingOrders map[string]*types.OrderRecord
	orderLock     *sync.Mutex

	cfg  config.MallCfg
	tMgr *transaction.Manager
}

// NewManager returns a new manager instance
func NewManager(ds datastore.Batching, sdb *db.SQLDB, pb *pubsub.PubSub, getCfg dtypes.GetMallConfigFunc, fm *transaction.Manager) (*Manager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	m := &Manager{
		SQLDB:         sdb,
		notify:        pb,
		ongoingOrders: make(map[string]*types.OrderRecord),
		orderLock:     &sync.Mutex{},
		cfg:           cfg,
		tMgr:          fm,
	}

	// state machine initialization
	m.stateMachineWait.Add(1)
	m.orderStateMachines = statemachine.New(ds, m, OrderInfo{})

	return m, nil
}

// Start initializes and starts the order state machine and associated tickers
func (m *Manager) Start(ctx context.Context) {
	if err := m.initStateMachines(ctx); err != nil {
		log.Errorf("restartStateMachines err: %s", err.Error())
	}

	// go m.subscribeEvents()
	go m.checkOrdersTimeout()
}

func (m *Manager) checkOrdersTimeout() {
	ticker := time.NewTicker(checkOrderInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		for _, orderRecord := range m.ongoingOrders {
			orderID := orderRecord.OrderID
			addr := orderRecord.To

			info, err := m.LoadOrderRecord(orderID)
			if err != nil {
				log.Errorf("checkOrderTimeout LoadOrderRecord %s , %s err:%s", addr, orderID, err.Error())
				continue
			}

			log.Debugf("checkout %s , %s ", addr, orderID)

			if info.State.Int() != Done.Int() && info.CreatedTime.Add(orderTimeoutTime).Before(time.Now()) {

				height := m.getHeight()

				err = m.orderStateMachines.Send(OrderHash(orderID), OrderTimeOut{Height: height})
				if err != nil {
					log.Errorf("checkOrderTimeout Send %s , %s err:%s", addr, orderID, err.Error())
					continue
				}
			}
		}
	}
}

func (m *Manager) getOrderIDByToAddress(to string) (string, bool) {
	for _, orderRecord := range m.ongoingOrders {
		if orderRecord.To == to {
			return orderRecord.OrderID, true
		}
	}

	return "", false
}

func (m *Manager) subscribeEvents() {
	subTransfer := m.notify.Sub(types.EventFvmTransferWatch.String())
	defer m.notify.Unsub(subTransfer)

	for {
		select {
		case u := <-subTransfer:
			tr := u.(*types.FvmTransferWatch)

			if orderID, exist := m.getOrderIDByToAddress(tr.To); exist {
				err := m.orderStateMachines.Send(OrderHash(orderID), PaymentResult{
					&PaymentInfo{
						TxHash: tr.TxHash,
						From:   tr.From,
						To:     tr.To,
						Value:  tr.Value,
					},
				})
				if err != nil {
					log.Errorf("subscribeNodeEvents Send %s err:%s", orderID, err.Error())
					continue
				}
			}
		}
	}
}

// Terminate stops the order state machine
func (m *Manager) Terminate(ctx context.Context) error {
	return m.orderStateMachines.Stop(ctx)
}

// CancelOrder cancel vps order
func (m *Manager) CancelOrder(orderID string) error {
	height := m.getHeight()

	return m.orderStateMachines.Send(OrderHash(orderID), OrderCancel{Height: height})
}

// CreatedOrder create vps order
func (m *Manager) CreatedOrder(req *types.OrderRecord) error {
	m.stateMachineWait.Wait()
	req.CreatedHeight = m.getHeight()

	// create order task
	return m.orderStateMachines.Send(OrderHash(req.OrderID), CreateOrder{orderInfoFrom(req)})
}

func (m *Manager) addOrder(req *types.OrderRecord) error {
	m.orderLock.Lock()
	defer m.orderLock.Unlock()

	if _, exist := m.ongoingOrders[req.UserID]; exist {
		return xerrors.New("user have order")
	}

	m.ongoingOrders[req.UserID] = req

	return nil
}

func (m *Manager) removeOrder(userID string) {
	m.orderLock.Lock()
	defer m.orderLock.Unlock()

	delete(m.ongoingOrders, userID)
}

func (m *Manager) createAliyunInstance(vpsInfo *types.CreateInstanceReq) (*types.CreateInstanceResponse, error) {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret

	priceUnit := vpsInfo.PeriodUnit
	period := vpsInfo.Period
	regionID := vpsInfo.RegionId
	if priceUnit == "Year" {
		priceUnit = "Month"
		period = period * 12
	}

	var securityGroupID string

	group, err := aliyun.DescribeSecurityGroups(regionID, k, s)
	if err == nil && len(group) > 0 {
		securityGroupID = group[0]
	} else {
		securityGroupID, err = aliyun.CreateSecurityGroup(regionID, k, s)
		if err != nil {
			log.Errorf("CreateSecurityGroup err: %s", err.Error())
			return nil, xerrors.New(*err.Data)
		}
	}
	log.Debugln("securityGroupID:", securityGroupID, " , DryRun:", vpsInfo.DryRun)
	result, err := aliyun.CreateInstance(k, s, vpsInfo, vpsInfo.DryRun)
	if err != nil {
		log.Errorf("CreateInstance err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}
	address, err := aliyun.AllocatePublicIPAddress(regionID, k, s, result.InstanceID)
	if err != nil {
		log.Errorf("AllocatePublicIpAddress err: %s", err.Error())
	} else {
		result.PublicIpAddress = address
	}

	err = aliyun.AuthorizeSecurityGroup(regionID, k, s, securityGroupID)
	if err != nil {
		log.Errorf("AuthorizeSecurityGroup err: %s", err.Error())
	}
	//randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	//keyPairName := "KeyPair" + fmt.Sprintf("%06d", randNew.Intn(1000000))
	//keyInfo, err := aliyun.CreateKeyPair(regionID, k, s, keyPairName)
	//if err != nil {
	//	log.Errorf("CreateKeyPair err: %s", err.Error())
	//} else {
	//	result.PrivateKey = keyInfo.PrivateKeyBody
	//}
	//var instanceIds []string
	//instanceIds = append(instanceIds, result.InstanceID)
	//_, err = aliyun.AttachKeyPair(regionID, k, s, keyPairName, instanceIds)
	//if err != nil {
	//	log.Errorf("AttachKeyPair err: %s", err.Error())
	//}
	go func() {
		time.Sleep(1 * time.Minute)

		err = aliyun.StartInstance(regionID, k, s, result.InstanceID)
		if err != nil {
			log.Infoln("StartInstance err:", err)
		}
	}()
	info := &types.MyInstance{
		OrderID:            vpsInfo.OrderID,
		UserID:             vpsInfo.UserID,
		InstanceId:         result.InstanceID,
		Price:              vpsInfo.TradePrice,
		InternetChargeType: vpsInfo.InternetChargeType,
		Location:           vpsInfo.RegionId,
	}
	saveErr := m.SaveMyInstancesInfo(info)
	if err != nil {
		log.Errorf("SaveMyInstancesInfo:%v", saveErr)
	}
	var instanceIds []string
	instanceIds = append(instanceIds, result.InstanceID)
	instanceInfo, err := aliyun.DescribeInstances(regionID, k, s, instanceIds)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
	}
	//todo
	fmt.Println(instanceInfo.Body.Instances.Instance[0])
	fmt.Println(instanceInfo.Body.Instances.Instance[0].PublicIpAddress.IpAddress[0])
	ip := instanceInfo.Body.Instances.Instance[0].PublicIpAddress.IpAddress[0]
	securityGroupId := instanceInfo.Body.Instances.Instance[0].SecurityGroupIds.SecurityGroupId[0]
	instanceDetailsInfo := &types.CreateInstanceReq{
		IpAddress:       *ip,
		InstanceId:      result.InstanceID,
		SecurityGroupId: *securityGroupId,
		OrderID:         vpsInfo.OrderID,
		UserID:          vpsInfo.UserID,
	}
	errU := m.UpdateVpsInstance(instanceDetailsInfo)
	if errU != nil {
		log.Errorf("UpdateVpsInstance:%v", errU)
	}
	return result, nil
}

func (m *Manager) getHeight() int64 {
	var msg filecoinbridge.TipSet
	err := filecoinbridge.ChainHead(&msg, m.cfg.LotusHTTPSAddr)
	if err != nil {
		log.Errorf("ChainHead err:%s", err.Error())
		return 0
	}

	return msg.Height
}
