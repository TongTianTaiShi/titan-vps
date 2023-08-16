package cli

import (
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/urfave/cli/v2"
)

// BasisCMDs Basis cmd
var BasisCMDs = []*cli.Command{
	WithCategory("order", orderCmds),
	WithCategory("user", userCmds),
	WithCategory("vps", vpsCmds),
}

var vpsCmds = &cli.Command{
	Name:  "vps",
	Usage: "Manage vps",
	Subcommands: []*cli.Command{
		describeRegionsCmd,
		describeInstanceTypeCmd,
		describeImageCmd,
		describePriceCmd,
	},
}

var orderCmds = &cli.Command{
	Name:  "order",
	Usage: "Manage order",
	Subcommands: []*cli.Command{
		createOrderCmd,
		cancelOrderCmd,
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

		api, closer, err := GetBasisAPI(cctx)
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

		api, closer, err := GetBasisAPI(cctx)
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

var describePriceCmd = &cli.Command{
	Name:  "dpc",
	Usage: "describe price",
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetBasisAPI(cctx)
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

		api, closer, err := GetBasisAPI(cctx)
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

		api, closer, err := GetBasisAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()

		address, err := api.CreateOrder(ctx, types.CreateInstanceReq{
			RegionId:     "cn-qingdao",
			ImageId:      "aliyun_2_1903_x64_20G_alibase_20230704.vhd",
			PeriodUnit:   "week",
			Period:       1,
			InstanceType: "ecs.t5-lc1m1.small",
			DryRun:       true,
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

		api, closer, err := GetBasisAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		orderID := cctx.String("oid")

		return api.CancelOrder(ctx, orderID)
	},
}

var getBalanceCmd = &cli.Command{
	Name:  "balance",
	Usage: "get balance",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetBasisAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()

		str, err := api.GetBalance(ctx)
		if err != nil {
			return err
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

		api, closer, err := GetBasisAPI(cctx)
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

		api, closer, err := GetBasisAPI(cctx)
		if err != nil {
			return err
		}

		defer closer()

		withdrawAddr := cctx.String("addr")
		value := cctx.String("value")

		return api.Withdraw(ctx, withdrawAddr, value)
	},
}
