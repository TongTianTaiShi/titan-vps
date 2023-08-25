package cli

import (
	"fmt"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/urfave/cli/v2"
)

// MallCMDs Mall cmd
var MallCMDs = []*cli.Command{
	WithCategory("order", orderCmds),
	WithCategory("user", userCmds),
	WithCategory("vps", vpsCmds),
	WithCategory("admin", adminCmds),
}

var adminCmds = &cli.Command{
	Name:  "admin",
	Usage: "Manage admin",
	Subcommands: []*cli.Command{
		createAdminCmd,
		approveWithdrawalCmd,
		rejectWithdrawalCmd,
		getWithdrawalCmd,
		getAddressesCmd,
		supplementRechargeCmd,
	},
}

var vpsCmds = &cli.Command{
	Name:  "vps",
	Usage: "Manage vps",
	Subcommands: []*cli.Command{
		describeRegionsCmd,
		describeInstanceTypeCmd,
		describeImageCmd,
		describePriceCmd,
		createKeyPairCmd,
		getDeskCmd,
		UpdateDefaultInfoCmd,
		GetInstanceDefaultCmd,
		GetInstanceCpuCmd,
		GetInstanceRenewStatusCmd,
		GetInstanceMemoryCmd,
		describeInstancesCmd,
	},
}

var orderCmds = &cli.Command{
	Name:  "order",
	Usage: "Manage order",
	Subcommands: []*cli.Command{
		createOrderCmd,
		cancelOrderCmd,
		paymentCompletedCmd,
		listCmd,
	},
}

var userCmds = &cli.Command{
	Name:  "user",
	Usage: "Manage user",
	Subcommands: []*cli.Command{
		getBalanceCmd,
		getRechargeAddrCmd,
		withdrawCmd,
	},
}

var describeRegionsCmd = &cli.Command{
	Name:  "dr",
	Usage: "describe regions",

	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()

		list, err := api.DescribeRegions(ctx)
		if err != nil {
			return err
		}

		fmt.Println(list)
		return nil
	},
}

var describeImageCmd = &cli.Command{
	Name:  "dim",
	Usage: "describe regions",

	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()

		list, err := api.DescribeImages(ctx, "cn-hangzhou", "")
		if err != nil {
			return err
		}

		fmt.Println(list)
		return nil
	},
}

var createKeyPairCmd = &cli.Command{
	Name:  "ckp",
	Usage: "describe regions",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "instanceID",
			Usage: "region id",
			Value: "",
		},
	},
	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()
		oid := cctx.String("instanceID")
		list, err := api.CreateKeyPair(ctx, "cn-hangzhou", oid)
		if err != nil {
			return err
		}

		fmt.Println(list)
		return nil
	},
}

var getDeskCmd = &cli.Command{
	Name:  "gdc",
	Usage: "get  desk indo",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "it",
			Usage: "region id",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "regionID",
			Usage: "region id",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "sys",
			Usage: "region id",
			Value: "",
		},
	},
	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()
		instanceId := cctx.String("it")
		regionID := cctx.String("regionID")
		sys := cctx.String("sys")
		desk := &types.AvailableResourceReq{
			InstanceType:        instanceId,
			RegionId:            regionID,
			DestinationResource: sys,
		}

		list, err := api.DescribeAvailableResourceForDesk(ctx, desk)
		if err != nil {
			return err
		}
		for _, v := range list {
			fmt.Println(v.Value)
		}
		return nil
	},
}

var UpdateDefaultInfoCmd = &cli.Command{
	Name:  "tdc",
	Usage: "UpdateDefaultInfo",
	Flags: []cli.Flag{},
	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()
		err = api.UpdateInstanceDefaultInfo(ctx)
		if err != nil {
			return err
		}
		return nil
	},
}

var GetInstanceDefaultCmd = &cli.Command{
	Name:  "gidc",
	Usage: "Get InstanceDefault indo",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "c",
			Usage: "region id",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "m",
			Usage: "region id",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "p",
			Usage: "region id",
			Value: "",
		},
	},
	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()
		c := int32(cctx.Int64("c"))
		p := cctx.Int64("p")
		m := float32(cctx.Float64("m"))
		req := &types.InstanceTypeFromBaseReq{
			RegionId:         "cn-hangzhou",
			MemorySize:       m,
			CpuCoreCount:     c,
			CpuArchitecture:  "",
			InstanceCategory: "",
			Page:             p,
			Limit:            100,
		}
		list, err := api.GetInstanceDefaultInfo(ctx, req)
		if err != nil {
			return err
		}
		for _, data := range list.List {
			fmt.Println(data.CpuCoreCount)
			fmt.Println(data.MemorySize)
			fmt.Println(data.InstanceTypeId)
			fmt.Println(data.Price)
		}

		return nil
	},
}

var GetInstanceCpuCmd = &cli.Command{
	Name:  "gicc",
	Usage: "get  desk indo",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "m",
			Usage: "region id",
			Value: "",
		},
	},
	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()
		m := float32(cctx.Float64("m"))
		req := &types.InstanceTypeFromBaseReq{
			RegionId:         "cn-hangzhou",
			MemorySize:       m,
			CpuArchitecture:  "",
			InstanceCategory: "",
		}
		list, err := api.GetInstanceCpuInfo(ctx, req)
		if err != nil {
			return err
		}
		for _, data := range list {
			fmt.Println(*data)
		}

		return nil
	},
}

var GetInstanceRenewStatusCmd = &cli.Command{
	Name:  "girs",
	Usage: "GetInstanceRenewStatus",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "i",
			Usage: "instance id",
			Value: "",
		},
	},
	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()
		m := cctx.String("i")
		req := types.SetRenewOrderReq{
			RegionID:   "cn-hangzhou",
			InstanceId: m,
		}
		list, err := api.GetRenewInstance(ctx, req)
		if err != nil {
			return err
		}
		fmt.Println(list)

		return nil
	},
}

var GetInstanceMemoryCmd = &cli.Command{
	Name:  "gimc",
	Usage: "get  desk indo",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "c",
			Usage: "region id",
			Value: "",
		},
	},
	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()
		c := int32(cctx.Int64("c"))
		req := &types.InstanceTypeFromBaseReq{
			RegionId:         "cn-hangzhou",
			CpuCoreCount:     c,
			CpuArchitecture:  "",
			InstanceCategory: "",
		}
		list, err := api.GetInstanceMemoryInfo(ctx, req)
		if err != nil {
			return err
		}
		for _, data := range list {
			fmt.Println(*data)
		}

		return nil
	},
}

var describeInstancesCmd = &cli.Command{
	Name:  "dicc",
	Usage: "describe regions",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "instanceID",
			Usage: "region id",
			Value: "",
		},
	},
	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()
		oid := cctx.String("instanceID")
		err = api.DescribeInstances(ctx, "cn-hangzhou", oid)
		if err != nil {
			return err
		}
		return nil
	},
}

var describePriceCmd = &cli.Command{
	Name:  "dpc",
	Usage: "describe price",
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()

		list, err := api.DescribePrice(ctx,
			&types.DescribePriceReq{
				RegionId:           "cn-hangzhou",
				InstanceType:       "ecs.s2.xlarge",
				PriceUnit:          "Week",
				Period:             1,
				Amount:             1,
				InternetChargeType: "PayByTraffic",
			})
		if err != nil {
			return err
		}

		fmt.Println(list)
		return nil
	},
}

var describeInstanceTypeCmd = &cli.Command{
	Name:  "dit",
	Usage: "describe instance type",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "rid",
			Usage: "region id",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "core",
			Usage: "core size",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "memory",
			Usage: "memory size",
			Value: "",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()
		regionID := cctx.String("rid")
		core := int32(cctx.Int("core"))
		memory := float32(cctx.Float64("memory"))
		list, err := api.DescribeInstanceType(ctx, &types.DescribeInstanceTypeReq{RegionId: regionID, CpuArchitecture: "X86", CpuCoreCount: core, MemorySize: memory, InstanceCategory: "General-purpose"})
		if err != nil {
			return err
		}
		fmt.Println(list)
		return nil
	},
}

var createOrderCmd = &cli.Command{
	Name:  "create",
	Usage: "create order",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()

		address, err := api.CreateOrder(ctx, types.CreateOrderReq{
			CreateInstanceReq: types.CreateInstanceReq{
				RegionId:                "cn-qingdao",
				ImageId:                 "aliyun_2_1903_x64_20G_alibase_20230731.vhd",
				PeriodUnit:              "Week",
				Period:                  1,
				InstanceType:            "ecs.n4.small",
				InternetChargeType:      "PayByTraffic",
				DryRun:                  true,
				InternetMaxBandwidthOut: 1,
				SystemDiskCategory:      "cloud_efficiency",
				SystemDiskSize:          40,
			},
			Amount: 1,
		})
		if err != nil {
			return err
		}

		fmt.Println(address)
		return nil
	},
}

var cancelOrderCmd = &cli.Command{
	Name:  "cancel",
	Usage: "cancel order",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "oid",
			Usage: "order id",
			Value: "",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		orderID := cctx.String("oid")

		return api.CancelUserOrder(ctx, orderID)
	},
}

var paymentCompletedCmd = &cli.Command{
	Name:  "payment",
	Usage: "payment order",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "oid",
			Usage: "order id",
			Value: "",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		orderID := cctx.String("oid")

		return api.PaymentUserOrder(ctx, orderID)
	},
}

var listCmd = &cli.Command{
	Name:  "list",
	Usage: "list order",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		infos, err := api.GetUserOrderRecords(ctx, 10, 0)
		if err != nil {
			return err
		}

		for _, info := range infos.List {
			fmt.Printf("%s Create:%v Expiration:%v \n", info.OrderID, info.CreatedTime, info.Expiration)
		}

		return nil
	},
}

var getBalanceCmd = &cli.Command{
	Name:  "balance",
	Usage: "get balance",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		bApi, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()

		str, err := bApi.GetBalance(ctx)
		if err != nil {
			if webErr, ok := err.(*api.ErrWeb); ok {
				fmt.Printf("web error code %d,message:%s \n", webErr.Code, webErr.Message)
			} else {
				fmt.Printf("web error message:%v \n", err)
			}
		}

		fmt.Println(str)
		return nil
	},
}

var getRechargeAddrCmd = &cli.Command{
	Name:  "gra",
	Usage: "get recharge address",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()

		str, err := api.GetRechargeAddress(ctx)
		if err != nil {
			return err
		}

		fmt.Println(str)
		return nil
	},
}

var withdrawCmd = &cli.Command{
	Name:  "withdraw",
	Usage: "withdraw",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "user withdraw address",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "value",
			Usage: "withdraw value",
			Value: "",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		bApi, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()

		withdrawAddr := cctx.String("addr")
		value := cctx.String("value")

		err = bApi.Withdraw(ctx, withdrawAddr, value)
		if webErr, ok := err.(*api.ErrWeb); ok {
			fmt.Printf("web error code %d,message:%s \n", webErr.Code, webErr.Message)
		} else {
			fmt.Printf("web error message:%v \n", err)
		}

		return err
	},
}

var approveWithdrawalCmd = &cli.Command{
	Name:  "aw",
	Usage: "approve withdrawal",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "oid",
			Usage: "user id",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "hash",
			Usage: "txHash",
			Value: "",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		oid := cctx.String("oid")
		hash := cctx.String("hash")

		return api.ApproveUserWithdrawal(ctx, oid, hash)
	},
}

var rejectWithdrawalCmd = &cli.Command{
	Name:  "rw",
	Usage: "reject withdrawal",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "oid",
			Usage: "user id",
			Value: "",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		oid := cctx.String("oid")

		return api.RejectUserWithdrawal(ctx, oid)
	},
}

var createAdminCmd = &cli.Command{
	Name:  "create",
	Usage: "create admin",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "uid",
			Usage: "user id",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "nick name",
			Value: "",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		uid := cctx.String("uid")
		name := cctx.String("name")

		return api.AddAdminUser(ctx, uid, name)
	},
}

var getWithdrawalCmd = &cli.Command{
	Name:  "list-withdrawal",
	Usage: "list withdrawals",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "uid",
			Usage: "user id",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "state",
			Usage: "order state",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "start",
			Usage: "start data",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "end",
			Usage: "end data",
			Value: "",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		uid := cctx.String("uid")
		state := cctx.String("state")
		start := cctx.String("start")
		end := cctx.String("end")
		info, err := api.GetWithdrawalRecords(ctx, &types.GetWithdrawRequest{
			Limit:     10,
			Offset:    0,
			UserID:    uid,
			State:     state,
			StartDate: start,
			EndDate:   end,
		})
		if err != nil {
			return err
		}

		for _, r := range info.List {
			fmt.Println(r)
		}

		return nil
	},
}

var getAddressesCmd = &cli.Command{
	Name:  "list-address",
	Usage: "list address",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		info, err := api.GetRechargeAddresses(ctx, 100, 0)
		if err != nil {
			return err
		}

		for _, r := range info.List {
			fmt.Println(r)
		}

		return nil
	},
}

var supplementRechargeCmd = &cli.Command{
	Name:  "supplement",
	Usage: "supplement recharge order",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "hash",
			Usage: "tx hash",
			Value: "",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		hash := cctx.String("hash")
		return api.SupplementRechargeOrder(ctx, hash)
	},
}
