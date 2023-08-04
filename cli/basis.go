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
	describeRegionsCmd,
}

var orderCmds = &cli.Command{
	Name:  "order",
	Usage: "Manage order",
	Subcommands: []*cli.Command{
		createOrderCmd,
		cancelOrderCmd,
		paymentCompletedCmd,
	},
}

var userCmds = &cli.Command{
	Name:  "user",
	Usage: "Manage user",
	Subcommands: []*cli.Command{
		signCodeCmd,
		loginCmd,
		logoutCmd,
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

var createOrderCmd = &cli.Command{
	Name:  "create",
	Usage: "create order",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "user",
			Usage: "node type: edge 1, update 6",
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

		user := cctx.String("user")

		address, err := api.CreateOrder(ctx, types.CreateOrderReq{Vps: "123456", User: user})
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
			Name:  "hash",
			Usage: "node type: edge 1, update 6",
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

		hash := cctx.String("hash")

		return api.CancelOrder(ctx, hash)
	},
}

var paymentCompletedCmd = &cli.Command{
	Name:  "pc",
	Usage: "payment completed",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "hash",
			Usage: "node type: edge 1, update 6",
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

		hash := cctx.String("hash")

		str, err := api.PaymentCompleted(ctx, types.PaymentCompletedReq{OrderID: hash})
		if err != nil {
			return err
		}

		fmt.Println(str)
		return nil
	},
}

var signCodeCmd = &cli.Command{
	Name:  "sc",
	Usage: "user get code",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "",
			Usage: "",
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

		userId := cctx.String("user_id")

		str, err := api.SignCode(ctx, userId)

		if err != nil {
			return err
		}

		fmt.Println(str)
		return nil
	},
}

var loginCmd = &cli.Command{
	Name:  "login",
	Usage: "user login",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "",
			Usage: "",
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

		userId := cctx.String("user_id")
		publicKey := cctx.String("public_key")
		signature := cctx.String("signature")

		str, err := api.Login(ctx, &types.UserReq{
			UserId:    userId,
			PublicKey: publicKey,
			Signature: signature,
		})
		if err != nil {
			return err
		}

		fmt.Println(str)
		return nil
	},
}

var logoutCmd = &cli.Command{
	Name:  "logout",
	Usage: "user logout",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "",
			Usage: "",
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

		userId := cctx.String("user_id")
		token := cctx.String("token")

		err = api.Logout(ctx, &types.UserReq{
			UserId: userId,
			Token:  token,
		})
		if err != nil {
			return err
		}
		return nil
	},
}
