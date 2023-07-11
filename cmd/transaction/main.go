package main

import (
	"fmt"
	"os"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/build"
	lcli "github.com/LMF709268224/titan-vps/cli"
	cliutil "github.com/LMF709268224/titan-vps/cli/util"
	liblog "github.com/LMF709268224/titan-vps/lib/log"
	"github.com/LMF709268224/titan-vps/node"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/repo"
	"github.com/filecoin-project/go-jsonrpc"
	logging "github.com/ipfs/go-log/v2"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

var log = logging.Logger("main")

var AdvanceBlockCmd *cli.Command

const (
	// FlagTransactionRepo Flag
	FlagTransactionRepo = "transaction-repo"
)

func main() {
	types.RunningNodeType = types.NodeTransaction

	liblog.SetupLogLevels()

	local := []*cli.Command{
		initCmd,
		runCmd,
	}

	if AdvanceBlockCmd != nil {
		local = append(local, AdvanceBlockCmd)
	}

	interactiveDef := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	app := &cli.App{
		Name:                 "transaction",
		Usage:                "Titan Edge Cloud Computing Transaction Service",
		Version:              build.UserVersion(),
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				// examined in the Before above
				Name:        "color",
				Usage:       "use color in display output",
				DefaultText: "depends on output being a TTY",
			},
			&cli.StringFlag{
				Name:    FlagTransactionRepo,
				EnvVars: []string{"TITAN_TRANSACTION_PATH"},
				Hidden:  true,
				Value:   "~/.vpstransaction",
			},
			&cli.BoolFlag{
				Name:  "interactive",
				Usage: "setting to false will disable interactive functionality of commands",
				Value: interactiveDef,
			},
			&cli.BoolFlag{
				Name:  "force-send",
				Usage: "if true, will ignore pre-send checks",
			},
			cliutil.FlagVeryVerbose,
		},
		After: func(c *cli.Context) error {
			if r := recover(); r != nil {
				panic(r)
			}
			return nil
		},

		Commands: append(local, lcli.Commands...),
	}

	app.Setup()
	app.Metadata["repoType"] = repo.Transaction

	lcli.RunApp(app)
}

var initCmd = &cli.Command{
	Name:  "init",
	Usage: "Initialize a transaction repo",
	Action: func(cctx *cli.Context) error {
		log.Info("Initializing transaction service")
		repoPath := cctx.String(FlagTransactionRepo)
		r, err := repo.NewFS(repoPath)
		if err != nil {
			return err
		}

		ok, err := r.Exists()
		if err != nil {
			return err
		}
		if ok {
			return xerrors.Errorf("repo at '%s' is already initialized", cctx.String(FlagTransactionRepo))
		}

		if err := r.Init(repo.Transaction); err != nil {
			return err
		}

		return nil
	},
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Start transaction service",
	Flags: []cli.Flag{},

	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		log.Info("Starting transaction service")

		repoPath := cctx.String(FlagTransactionRepo)
		r, err := repo.NewFS(repoPath)
		if err != nil {
			return err
		}

		ok, err := r.Exists()
		if err != nil {
			return err
		}
		if !ok {
			if err := r.Init(repo.Transaction); err != nil {
				return err
			}
		}

		lr, err := r.Lock(repo.Transaction)
		if err != nil {
			return err
		}

		cfg, err := lr.Config()
		if err != nil {
			return err
		}

		tCfg := cfg.(*config.TransactionCfg)

		err = lr.Close()
		if err != nil {
			return err
		}

		shutdownChan := make(chan struct{})

		var tAPI api.Transaction
		stop, err := node.New(cctx.Context,
			node.Transaction(&tAPI),
			node.Base(),
			node.Repo(r),
		)
		if err != nil {
			return xerrors.Errorf("creating node: %w", err)
		}

		// Populate JSON-RPC options.
		serverOptions := []jsonrpc.ServerOption{jsonrpc.WithServerErrors(api.RPCErrors)}
		if maxRequestSize := cctx.Int("api-max-req-size"); maxRequestSize != 0 {
			serverOptions = append(serverOptions, jsonrpc.WithMaxRequestSize(int64(maxRequestSize)))
		}

		// Instantiate the transaction handler.
		h, err := node.TransactionHandler(tAPI, true, serverOptions...)
		if err != nil {
			return fmt.Errorf("failed to instantiate rpc handler: %s", err.Error())
		}

		// Serve the RPC.
		rpcStopper, err := node.ServeRPC(h, "transaction", tCfg.API.ListenAddress)
		if err != nil {
			return fmt.Errorf("failed to start json-rpc endpoint: %s", err.Error())
		}

		log.Info("transaction listen with:", tCfg.API.ListenAddress)

		// Monitor for shutdown.
		finishCh := node.MonitorShutdown(shutdownChan,
			node.ShutdownHandler{Component: "rpc server", StopFunc: rpcStopper},
			node.ShutdownHandler{Component: "node", StopFunc: stop},
		)
		<-finishCh // fires when shutdown is complete.
		return nil
	},
}
