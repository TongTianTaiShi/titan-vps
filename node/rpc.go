package node

import (
	"context"
	"html/template"
	"net"
	"net/http"
	"time"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/metrics"
	"github.com/LMF709268224/titan-vps/metrics/proxy"
	"github.com/LMF709268224/titan-vps/node/handler"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"
	"go.opencensus.io/tag"
	"golang.org/x/xerrors"
)

var rpclog = logging.Logger("rpc")

// ServeRPC serves an HTTP handler over the supplied listen multiaddr.
//
// This function spawns a goroutine to run the server, and returns immediately.
// It returns the stop function to be called to terminate the endpoint.
//
// The supplied ID is used in tracing, by inserting a tag in the context.
func ServeRPC(h http.Handler, id string, addr string) (StopFunc, error) {
	// Start listening to the addr; if invalid or occupied, we will fail early.
	lst, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, xerrors.Errorf("could not listen: %w", err)
	}

	// Instantiate the server and start listening.
	srv := &http.Server{
		Handler:           h,
		ReadHeaderTimeout: 30 * time.Second,
		BaseContext: func(listener net.Listener) context.Context {
			ctx, _ := tag.New(context.Background(), tag.Upsert(metrics.APIInterface, id))
			return ctx
		},
	}

	go func() {
		err = srv.Serve(lst)
		if err != http.ErrServerClosed {
			rpclog.Warnf("rpc server failed: %s", err)
		}
	}()

	return srv.Shutdown, err
}

// TransactionHandler returns a transaction handler, to be mounted as-is on the server.
func TransactionHandler(a api.Transaction, permissioned bool, opts ...jsonrpc.ServerOption) (http.Handler, error) {
	m := mux.NewRouter()

	serveRpc := func(path string, hnd interface{}) {
		rpcServer := jsonrpc.NewServer(append(opts, jsonrpc.WithServerErrors(api.RPCErrors))...)
		rpcServer.Register("titan", hnd)

		var hand http.Handler = rpcServer
		if permissioned {
			hand = handler.New(a.AuthVerify, m.ServeHTTP)
		}

		m.Handle(path, hand)
	}

	fnapi := proxy.MetricedTransactionAPI(a)
	if permissioned {
		fnapi = api.PermissionedTransactionAPI(fnapi)
	}

	serveRpc("/rpc/v0", fnapi)
	m.PathPrefix("/").Handler(http.DefaultServeMux) // pprof

	return m, nil
}

// MallHandler returns handler, to be mounted as-is on the server.
func MallHandler(a api.Mall, permissioned bool, opts ...jsonrpc.ServerOption) (http.Handler, error) {
	m := mux.NewRouter()

	serveRPC := func(path string, hnd interface{}) {
		rpcServer := jsonrpc.NewServer(append(opts, jsonrpc.WithServerErrors(api.RPCErrors))...)
		rpcServer.Register("titan", hnd)

		var hand http.Handler = rpcServer
		if permissioned {
			hand = handler.New(a.AuthVerify, rpcServer.ServeHTTP)
		}

		m.Handle(path, hand)
	}

	wapi := proxy.MetricedMallAPI(a)
	if permissioned {
		wapi = api.PermissionedMallAPI(wapi)
	}

	serveRPC("/rpc/v0", wapi)
	m.HandleFunc("/rpc/index", homePage)
	m.PathPrefix("/").Handler(http.DefaultServeMux) // pprof

	return m, nil
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../../homepage.html"))
	tmpl.Execute(w, nil)
}
