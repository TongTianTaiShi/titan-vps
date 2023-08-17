package client

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"github.com/LMF709268224/titan-vps/api"

	"github.com/filecoin-project/go-jsonrpc"

	"github.com/LMF709268224/titan-vps/lib/rpcenc"
)

// NewTransaction creates a new http jsonrpc client.
func NewTransaction(ctx context.Context, addr string, requestHeader http.Header) (api.Transaction, jsonrpc.ClientCloser, error) {
	pushURL, err := getPushURL(addr)
	if err != nil {
		return nil, nil, err
	}

	var res api.TransactionStruct
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "titan",
		api.GetInternalStructs(&res),
		requestHeader,
		rpcenc.ReaderParamEncoder(pushURL),
		jsonrpc.WithErrors(api.RPCErrors),
	)

	return &res, closer, err
}

func getPushURL(addr string) (string, error) {
	pushURL, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	switch pushURL.Scheme {
	case "ws":
		pushURL.Scheme = "http"
	case "wss":
		pushURL.Scheme = "https"
	}
	///rpc/v0 -> /rpc/streams/v0/push

	pushURL.Path = path.Join(pushURL.Path, "../streams/v0/push")
	return pushURL.String(), nil
}

// NewBasis creates a new http jsonrpc client for basis
func NewBasis(ctx context.Context, addr string, requestHeader http.Header, opts ...jsonrpc.Option) (api.Basis, jsonrpc.ClientCloser, error) {
	pushURL, err := getPushURL(addr)
	if err != nil {
		return nil, nil, err
	}

	var res api.BasisStruct
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "titan",
		api.GetInternalStructs(&res), requestHeader,
		append([]jsonrpc.Option{
			rpcenc.ReaderParamEncoder(pushURL),
			jsonrpc.WithErrors(api.RPCErrors),
		}, opts...)...)

	return &res, closer, err
}

// NewCommonRPCV0 creates a new http jsonrpc client.
func NewCommonRPCV0(ctx context.Context, addr string, requestHeader http.Header) (api.Common, jsonrpc.ClientCloser, error) {
	var res api.CommonStruct
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "titan",
		api.GetInternalStructs(&res), requestHeader)

	return &res, closer, err
}
