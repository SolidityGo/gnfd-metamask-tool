package client

import (
	"github.com/cometbft/cometbft/rpc/client"
	chttp "github.com/cometbft/cometbft/rpc/client/http"
	libclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
)

type RPCClient struct {
	RpcClient client.Client
}

func httpClient(addr string) *chttp.HTTP {
	httpCli, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		panic(err)
	}
	cli, err := chttp.NewWithClient(addr, "/websocket", httpCli)
	if err != nil {
		panic(err)
	}
	return cli
}

func NewRPCClient(addr string) RPCClient {
	return RPCClient{
		RpcClient: httpClient(addr),
	}
}
