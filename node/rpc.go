package node

import (
	"context"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/metrics"
	"github.com/LMF709268224/titan-vps/metrics/proxy"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/handler"
	"github.com/xuri/excelize/v2"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"
	"go.opencensus.io/tag"
	"golang.org/x/xerrors"
)

var (
	rpclog  = logging.Logger("rpc")
	mallCfg *config.MallCfg
)

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
func MallHandler(a api.Mall, permissioned bool, cfg *config.MallCfg, opts ...jsonrpc.ServerOption) (http.Handler, error) {
	m := mux.NewRouter()
	mallCfg = cfg

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
	m.HandleFunc("/rpc/download/withdraw", downloadWithdrawFile)
	m.PathPrefix("/").Handler(http.DefaultServeMux) // pprof

	return m, nil
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../../homepage.html"))
	tmpl.Execute(w, nil)
}

func downloadWithdrawFile(w http.ResponseWriter, r *http.Request) {
	// token := r.URL.Query().Get("token")

	state := r.URL.Query().Get("state")
	userID := r.URL.Query().Get("user-id")
	startDate := r.URL.Query().Get("start-date")
	endDate := r.URL.Query().Get("end-date")

	fmt.Println("state:", state)
	fmt.Println("userID:", userID)
	fmt.Println("startDate:", startDate)
	fmt.Println("endDate:", endDate)

	statuses := make([]types.WithdrawState, 0)
	if state == "" {
		statuses = []types.WithdrawState{types.WithdrawCreate, types.WithdrawDone, types.WithdrawRefund}
	} else {
		s2, err := strconv.Atoi(state)
		if err != nil {
			http.Error(w, "state atoi", http.StatusInternalServerError)
			return
		}

		statuses = []types.WithdrawState{types.WithdrawState(s2)}
	}

	client, err := db.NewSQLDB(mallCfg.DatabaseAddress)
	if err != nil {
		http.Error(w, "NewSQLDB", http.StatusInternalServerError)
		return
	}

	rows, err := client.LoadWithdrawRecordRows(statuses, userID, startDate, endDate)
	if err != nil {
		http.Error(w, "LoadWithdrawRecordRows", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	file := excelize.NewFile()
	columns := []string{"OrderID", "UserID", "Value", "WithdrawAddr", "WithdrawHash", "CreatedTime", "State"}
	for i, colName := range columns {
		file.SetCellValue("Sheet1", string(rune('A'+i))+"1", colName)
	}

	rowIdx := 2
	for rows.Next() {
		info := &types.WithdrawRecord{}
		err = rows.StructScan(info)
		if err != nil {
			log.Errorf("asset StructScan err: %s", err.Error())
			continue
		}

		values := []string{info.OrderID, info.UserID, info.Value, info.WithdrawAddr, info.WithdrawHash, info.CreatedTime.String(), string(info.State)}
		for i, value := range values {
			file.SetCellValue("Sheet1", string(rune('A'+i))+strconv.Itoa(rowIdx), value)
		}
		rowIdx++
	}

	if rows.Err() != nil {
		http.Error(w, "Error fetching rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=Withdraw.xlsx")
	err = file.Write(w)
	if err != nil {
		http.Error(w, "Error writing excel file", http.StatusInternalServerError)
		return
	}
}
