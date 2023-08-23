package main

import (
	"os"
	"time"

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

var (
	log               = logging.Logger("main")
	HeartbeatInterval = 10 * time.Second
)

const (
	// FlagMallRepo Flag
	FlagMallRepo = "mall-repo"
)

var AdvanceBlockCmd *cli.Command

func main() {
	types.RunningNodeType = types.NodeMall

	liblog.SetupLogLevels()

	local := []*cli.Command{
		initCmd,
		runCmd,
	}

	local = append(local, lcli.MallCMDs...)

	if AdvanceBlockCmd != nil {
		local = append(local, AdvanceBlockCmd)
	}

	interactiveDef := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	app := &cli.App{
		Name:                 "mall",
		Usage:                "Titan Edge Cloud Computing mall Service",
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
				Name:    FlagMallRepo,
				EnvVars: []string{"TITAN_MALL_PATH"},
				Hidden:  true,
				Value:   "~/.vpsmall",
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
	app.Metadata["repoType"] = repo.Mall

	lcli.RunApp(app)
}

var initCmd = &cli.Command{
	Name:  "init",
	Usage: "Initialize a mall repo",
	Action: func(cctx *cli.Context) error {
		log.Info("Initializing mall service")
		repoPath := cctx.String(FlagMallRepo)
		r, err := repo.NewFS(repoPath)
		if err != nil {
			return err
		}

		ok, err := r.Exists()
		if err != nil {
			return err
		}
		if ok {
			return xerrors.Errorf("repo at '%s' is already initialized", cctx.String(FlagMallRepo))
		}

		if err := r.Init(repo.Mall); err != nil {
			return err
		}

		return nil
	},
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Start mall service",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "http-server-timeout",
			Value: "30s",
		},
	},

	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		log.Info("Starting mall node")

		repoPath := cctx.String(FlagMallRepo)
		r, err := repo.NewFS(repoPath)
		if err != nil {
			return err
		}

		ok, err := r.Exists()
		if err != nil {
			return err
		}
		if !ok {
			if err := r.Init(repo.Mall); err != nil {
				return err
			}
		}

		lr, err := r.Lock(repo.Mall)
		if err != nil {
			return err
		}

		cfg, err := lr.Config()
		if err != nil {
			return err
		}

		bCfg := cfg.(*config.MallCfg)

		err = lr.Close()
		if err != nil {
			return err
		}

		shutdownChan := make(chan struct{})

		var bAPI api.Mall
		stop, err := node.New(cctx.Context,
			node.Mall(&bAPI),
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

		// Instantiate the mall handler.
		h, err := node.MallHandler(bAPI, true, bCfg, serverOptions...)
		if err != nil {
			return xerrors.Errorf("failed to instantiate rpc handler: %s", err.Error())
		}

		// Serve the RPC.
		rpcStopper, err := node.ServeRPC(h, "mall", bCfg.API.ListenAddress)
		if err != nil {
			return xerrors.Errorf("failed to start json-rpc endpoint: %s", err.Error())
		}

		log.Info("titan mall listen with:", bCfg.API.ListenAddress)

		// Monitor for shutdown.
		finishCh := node.MonitorShutdown(shutdownChan,
			node.ShutdownHandler{Component: "rpc server", StopFunc: rpcStopper},
			node.ShutdownHandler{Component: "node", StopFunc: stop},
		)
		// gin router
		// srv, err := web.NewServer(bCfg)
		// if err != nil {
		// 	log.Fatalf("create api server: %v\n", err)
		// }
		// go srv.Run()
		<-finishCh // fires when shutdown is complete.
		return nil
	},
}
