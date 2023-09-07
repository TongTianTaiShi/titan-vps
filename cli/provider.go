package cli

import (
	"fmt"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/urfave/cli/v2"
)

var providerCmds = &cli.Command{
	Name: "provider",
	Subcommands: []*cli.Command{
		getEmailVerifyCmd,
		setInvitationCodeCmd,
		loginCmd,
	},
}

var setInvitationCodeCmd = &cli.Command{
	Name:  "set_invitation",
	Usage: "set invitation code",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "code",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		code := cctx.String("code")

		return api.SetInvitationCode(ctx, code)
	},
}

var getEmailVerifyCmd = &cli.Command{
	Name:  "get_email_verify",
	Usage: "get email verify",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "id",
		},
		&cli.IntFlag{
			Name: "type",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		id := cctx.String("id")
		loginType := cctx.Int("type")
		code, err := api.GetVerifyMessage(ctx, id, types.LoginType(loginType))
		if err != nil {
			fmt.Println(id, types.LoginType(loginType).String(), err)
			return err
		} else if loginType == 2 {
			fmt.Println(id, types.LoginType(loginType).String(), "Success")
			return nil
		} else {
			fmt.Println(id, types.LoginType(loginType).String(), code)
		}
		return nil
	},
}

var loginCmd = &cli.Command{
	Name:  "login",
	Usage: "login",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "email",
		},
		&cli.StringFlag{
			Name: "verify",
		},
		&cli.StringFlag{
			Name: "inv",
		},
		&cli.StringFlag{
			Name: "passwd",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := ReqContext(cctx)

		api, closer, err := GetMallAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		email := cctx.String("email")
		inv := cctx.String("inv")
		verify := cctx.String("verify")
		passwd := cctx.String("passwd")
		req := &types.AccountRequest{
			Type:           2,
			VerifyCode:     verify,
			Email:          email,
			InvitationCode: inv,
			Password:       passwd,
		}
		response, err := api.LoginAccount(ctx, req)
		if err != nil {
			return err
		} else {
			fmt.Println(response)
		}

		return nil
	},
}
