package orders

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
	"time"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"
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

var USDRateInfo struct {
	USDRate float32
	ET      time.Time
}

// Manager manager order
type Manager struct {
	stateMachineWait   sync.WaitGroup
	orderStateMachines *statemachine.StateGroup
	*db.SQLDB

	notification *pubsub.PubSub

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
		notification:  pb,
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
	go m.CronGetInstanceDefaultInfo()
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
	subTransfer := m.notification.Sub(types.EventFvmTransferWatch.String())
	defer m.notification.Unsub(subTransfer)

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
	if exist, _ := m.orderStateMachines.Has(OrderHash(orderID)); !exist {
		return &api.ErrWeb{Code: terrors.NotFoundOrder.Int(), Message: terrors.NotFoundOrder.String()}
	}

	height := m.getHeight()

	err := m.orderStateMachines.Send(OrderHash(orderID), OrderCancel{Height: height})
	if err != nil {
		return &api.ErrWeb{Code: terrors.StateMachinesError.Int(), Message: err.Error()}
	}

	return nil
}

// PaymentCompleted cancel vps order
func (m *Manager) PaymentCompleted(orderID string) error {
	if exist, _ := m.orderStateMachines.Has(OrderHash(orderID)); !exist {
		return &api.ErrWeb{Code: terrors.NotFoundOrder.Int(), Message: terrors.NotFoundOrder.String()}
	}

	err := m.orderStateMachines.Send(OrderHash(orderID), PaymentResult{})
	if err != nil {
		return &api.ErrWeb{Code: terrors.StateMachinesError.Int(), Message: err.Error()}
	}

	return nil
}

// CreatedOrder create vps order
func (m *Manager) CreatedOrder(req *types.OrderRecord) error {
	m.stateMachineWait.Wait()
	req.CreatedHeight = m.getHeight()

	m.addOrder(req)

	// create order task
	err := m.orderStateMachines.Send(OrderHash(req.OrderID), CreateOrder{orderInfoFrom(req)})
	if err != nil {
		return &api.ErrWeb{Code: terrors.StateMachinesError.Int(), Message: err.Error()}
	}

	return nil
}

func (m *Manager) addOrder(req *types.OrderRecord) {
	m.orderLock.Lock()
	defer m.orderLock.Unlock()

	// if _, exist := m.ongoingOrders[req.OrderID]; exist {
	// 	return xerrors.New("user have order")
	// }

	m.ongoingOrders[req.OrderID] = req

	return
}

func (m *Manager) removeOrder(orderID string) {
	m.orderLock.Lock()
	defer m.orderLock.Unlock()

	delete(m.ongoingOrders, orderID)
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
			return nil, xerrors.New(err.Error())
		}
	}
	log.Debugln("securityGroupID:", securityGroupID, " , DryRun:", vpsInfo.DryRun)
	result, err := aliyun.CreateInstance(k, s, vpsInfo, vpsInfo.DryRun)
	if err != nil {
		log.Errorf("CreateInstance err: %s", err.Error())
		return nil, xerrors.New(err.Error())
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
	var instanceIds []string
	instanceIds = append(instanceIds, result.InstanceID)
	instanceInfo, err := aliyun.DescribeInstances(regionID, k, s, instanceIds)
	if err != nil {
		log.Errorf("DescribeInstances err: %s", err.Error())
	}
	if len(instanceInfo.Body.Instances.Instance) > 0 {
		ip := instanceInfo.Body.Instances.Instance[0].PublicIpAddress.IpAddress[0]
		securityGroupId := ""
		if len(instanceInfo.Body.Instances.Instance) > 0 {
			securityGroupId = *instanceInfo.Body.Instances.Instance[0].SecurityGroupIds.SecurityGroupId[0]
		}
		OSType := instanceInfo.Body.Instances.Instance[0].OSType
		InstanceName := instanceInfo.Body.Instances.Instance[0].InstanceName
		BandwidthOut := instanceInfo.Body.Instances.Instance[0].InternetMaxBandwidthOut
		Cores := instanceInfo.Body.Instances.Instance[0].Cpu
		Memory := instanceInfo.Body.Instances.Instance[0].Memory
		instanceDetailsInfo := &types.CreateInstanceReq{
			IpAddress:       *ip,
			InstanceId:      result.InstanceID,
			SecurityGroupId: securityGroupId,
			OrderID:         vpsInfo.OrderID,
			UserID:          vpsInfo.UserID,
			OSType:          *OSType,
			Cores:           *Cores,
			Memory:          float32(*Memory),
		}
		errU := m.UpdateVpsInstance(instanceDetailsInfo)
		if errU != nil {
			log.Errorf("UpdateVpsInstance:%v", errU)
		}
		info := &types.MyInstance{
			OrderID:            vpsInfo.OrderID,
			UserID:             vpsInfo.UserID,
			InstanceId:         result.InstanceID,
			Price:              vpsInfo.TradePrice,
			InternetChargeType: vpsInfo.InternetChargeType,
			Location:           vpsInfo.RegionId,
			InstanceSystem:     *OSType,
			InstanceName:       *InstanceName,
			BandwidthOut:       *BandwidthOut,
		}
		saveErr := m.SaveMyInstancesInfo(info)
		if err != nil {
			log.Errorf("SaveMyInstancesInfo:%v", saveErr)
		}
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
func (m *Manager) CronGetInstanceDefaultInfo() {
	crontab := cron.New(cron.WithSeconds())
	var ctx context.Context
	task := func() {
		m.UpdateInstanceDefaultInfo(ctx)
	}
	//spec := "0 0 1,13 * * ?"
	spec := "*/60 * * * * ?"
	crontab.AddFunc(spec, task)
	crontab.Start()
	fmt.Println("start")
}
func (m *Manager) UpdateInstanceDefaultInfo(ctx context.Context) {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret
	regions, err := aliyun.DescribeRegions(k, s)
	if err != nil {
		log.Errorf("DescribePrice err:%v", err.Error())
		return
	}
	for _, region := range regions.Body.Regions.Region {
		instanceType := &types.DescribeInstanceTypeReq{
			RegionId:     *region.RegionId,
			CpuCoreCount: 0,
			MemorySize:   0,
		}
		instances, err := m.DescribeInstanceType(ctx, instanceType)
		if err != nil {
			log.Errorf("DescribeInstanceType err:%v", err.Error())
			continue
		}
		fmt.Println("start--------")
		for _, instance := range instances.InstanceTypes {
			time.Sleep(500 * time.Millisecond)
			images, err := m.DescribeImages(ctx, *region.RegionId, instance.InstanceTypeId)
			if err != nil {
				log.Errorf("DescribeImages err:%v", err.Error())
				continue
			}
			var disk = &types.AvailableResourceReq{
				InstanceType:        instance.InstanceTypeId,
				RegionId:            *region.RegionId,
				DestinationResource: "SystemDisk",
			}

			disks, err := m.DescribeAvailableResourceForDesk(ctx, disk)
			if err != nil {
				log.Errorf("DescribeAvailableResourceForDesk err:%v", err.Error())
				continue
			}
			for _, disk := range disks {
				priceReq := &types.DescribePriceReq{
					RegionId:                *region.RegionId,
					InstanceType:            instance.InstanceTypeId,
					PriceUnit:               "Week",
					ImageID:                 images[0].ImageId,
					InternetChargeType:      "PayByTraffic",
					SystemDiskCategory:      disk.Value,
					SystemDiskSize:          40,
					Period:                  1,
					Amount:                  1,
					InternetMaxBandwidthOut: 10,
				}
				price, err := aliyun.DescribePrice(k, s, priceReq)
				if err != nil {
					fmt.Println("get price fail")
					log.Errorf("DescribePrice err:%v", err.Error())
					continue
				}
				if USDRateInfo.USDRate == 0 || time.Now().After(USDRateInfo.ET) {
					UsdRate := aliyun.GetExchangeRate()
					USDRateInfo.USDRate = UsdRate
					USDRateInfo.ET = time.Now().Add(time.Hour)
				}
				if USDRateInfo.USDRate == 0 {
					USDRateInfo.USDRate = 7.2673
				}
				UsdRate := USDRateInfo.USDRate
				price.USDPrice = price.USDPrice / UsdRate
				fmt.Println(UsdRate)
				var info = &types.DescribeInstanceTypeFromBase{
					RegionId:               *region.RegionId,
					InstanceTypeId:         instance.InstanceTypeId,
					MemorySize:             instance.MemorySize,
					CpuArchitecture:        instance.CpuArchitecture,
					InstanceCategory:       instance.InstanceCategory,
					CpuCoreCount:           instance.CpuCoreCount,
					AvailableZone:          instance.AvailableZone,
					InstanceTypeFamily:     instance.InstanceTypeFamily,
					PhysicalProcessorModel: instance.PhysicalProcessorModel,
					Price:                  price.USDPrice,
				}
				saveErr := m.SaveInstancesInfo(info)
				if err != nil {
					log.Errorf("SaveMyInstancesInfo:%v", saveErr)
				}
			}

		}

	}
	return
}
func (m *Manager) DescribeInstanceType(ctx context.Context, instanceType *types.DescribeInstanceTypeReq) (*types.DescribeInstanceTypeResponse, error) {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret
	rsp, err := aliyun.DescribeInstanceTypes(k, s, instanceType)
	if err != nil {
		log.Errorf("DescribeInstanceTypes err: %s", err.Error())
		return nil, xerrors.New(err.Error())
	}
	AvailableResource, err := aliyun.DescribeAvailableResource(k, s, instanceType)
	if err != nil {
		log.Errorf("DescribeAvailableResource err: %s", err.Error())
		return nil, xerrors.New(err.Error())
	}
	rspDataList := &types.DescribeInstanceTypeResponse{
		NextToken: *rsp.Body.NextToken,
	}
	instanceTypes := make(map[string]int)
	AvailableZone := len(AvailableResource.Body.AvailableZones.AvailableZone)
	if AvailableZone < 0 {
		return rspDataList, nil
	}
	for _, data := range AvailableResource.Body.AvailableZones.AvailableZone {
		availableTypes := data.AvailableResources.AvailableResource
		if len(availableTypes) > 0 {
			for _, instanceTypeResource := range availableTypes {
				Resources := instanceTypeResource.SupportedResources.SupportedResource
				if len(Resources) > 0 {
					for _, Resource := range Resources {
						if *Resource.Status == "Available" {
							instanceTypes[*Resource.Value] = 1
						}
					}
				}
			}
		}
	}
	for _, data := range rsp.Body.InstanceTypes.InstanceType {
		if _, ok := instanceTypes[*data.InstanceTypeId]; !ok {
			continue
		}
		rspData := &types.DescribeInstanceType{
			InstanceCategory:       *data.InstanceCategory,
			InstanceTypeId:         *data.InstanceTypeId,
			MemorySize:             *data.MemorySize,
			CpuArchitecture:        *data.CpuArchitecture,
			InstanceTypeFamily:     *data.InstanceTypeFamily,
			CpuCoreCount:           *data.CpuCoreCount,
			AvailableZone:          AvailableZone,
			PhysicalProcessorModel: *data.PhysicalProcessorModel,
		}
		rspDataList.InstanceTypes = append(rspDataList.InstanceTypes, rspData)
	}
	return rspDataList, nil
}

func (m *Manager) DescribeImages(ctx context.Context, regionID, instanceType string) ([]*types.DescribeImageResponse, error) {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret

	rsp, err := aliyun.DescribeImages(regionID, k, s, instanceType)
	if err != nil {
		log.Errorf("DescribeImages err: %s", err.Error())
		return nil, xerrors.New(err.Error())
	}
	var rspDataList []*types.DescribeImageResponse
	for _, data := range rsp.Body.Images.Image {
		rspData := &types.DescribeImageResponse{
			ImageId:      *data.ImageId,
			ImageFamily:  *data.ImageFamily,
			ImageName:    *data.ImageName,
			Architecture: *data.Architecture,
			OSName:       *data.OSName,
			OSType:       *data.OSType,
			Platform:     *data.Platform,
		}
		rspDataList = append(rspDataList, rspData)
	}
	return rspDataList, nil
}

func (m *Manager) DescribeAvailableResourceForDesk(ctx context.Context, desk *types.AvailableResourceReq) ([]*types.AvailableResourceResponse, error) {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret
	rsp, err := aliyun.DescribeAvailableResourceForDesk(k, s, desk)
	if err != nil {
		log.Errorf("DescribeImages err: %s", err.Error())
		return nil, xerrors.New(err.Error())
	}
	var Category = map[string]int{
		"cloud":            1,
		"cloud_essd":       1,
		"cloud_ssd":        1,
		"cloud_efficiency": 1,
		"ephemeral_ssd":    1,
	}
	var rspDataList []*types.AvailableResourceResponse
	if len(rsp.Body.AvailableZones.AvailableZone) > 0 {
		AvailableResources := rsp.Body.AvailableZones.AvailableZone[0].AvailableResources.AvailableResource
		if len(AvailableResources) > 0 {
			systemDesk := AvailableResources[0].SupportedResources.SupportedResource
			if len(systemDesk) > 0 {
				for _, data := range systemDesk {
					if *data.Status == "Available" {
						if _, ok := Category[*data.Value]; ok {
							desk := &types.AvailableResourceResponse{
								Min:   *data.Min,
								Max:   *data.Max,
								Value: *data.Value,
								Unit:  *data.Unit,
							}
							rspDataList = append(rspDataList, desk)
						}
					}
				}
			}
		}
	}
	return rspDataList, nil
}
