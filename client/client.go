/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package client

import (
	"fmt"
	netRpc "net/rpc"
	"net/url"

	"github.com/Oneledger/protocol/rpc"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/node"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var (
	ErrEmptyQuery    = errors.New("empty query path")
	ErrEmptyResponse = errors.New("empty response")
)

type Context struct {
	rpcClient     rpcclient.Client
	broadcastMode string
	proof         bool
	oltClient     *netRpc.Client
}

/*
	Generators
*/
func NewLocalContext(node *node.Node) (cliCtx Context) {
	tmClient := rpcclient.NewLocal(node)

	cliCtx = Context{
		rpcClient:     tmClient,
		broadcastMode: BroadcastAsync,
		proof:         false,
	}
	return
}

func NewContext(rpcAddress, sdkAddress string) (cliCtx Context, err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Debug("Ignoring rpcClient Panic", "r", r)
			// debug.PrintStack()
		}
	}()

	// tm rpc Context
	var tmRPCClient = rpcclient.NewHTTP(rpcAddress, "/websocket")
	/*
		// check status of rpc; return client if everything fine
		_, err = tmRPCClient.Status()
		if err == nil {
			logger.Debug("rpcClient is running")

			cliCtx = Context{
				rpcClient:     tmRPCClient,
				broadcastMode: BroadcastCommit,
			}
			return
		}

	*/

	// try starting tmRPCClient client
	err = tmRPCClient.Start()
	if err != nil {
		err = fmt.Errorf("rpcClient is unavailable", "address", rpcAddress, err)
		return
	}

	u, err := url.Parse(sdkAddress)
	if err != nil {
		return
	}

	client, err := netRpc.DialHTTPPath("tcp", u.Host, rpc.Path)
	if err != nil {
		return
	}

	cliCtx = Context{
		rpcClient:     tmRPCClient,
		broadcastMode: BroadcastCommit,
		oltClient:     client,
	}

	return
}

/*
	Context Methods
*/

// Send a very specific query to return parse response.value
func (ctx Context) Query(serviceMethod string, args interface{}, response interface{}) error {

	err := ctx.oltClient.Call(serviceMethod, args, response)
	if err != nil {
		logger.Error("error while calling client server", "err", err)
	}

	return err
}

func (ctx Context) Tx(hash []byte, prove bool) (res *ctypes.ResultTx) {

	result, err := ctx.rpcClient.Tx(hash, prove)
	if err != nil {
		logger.Error("TxSearch Error", "err", err)
		return nil
	}

	logger.Debug("TxSearch", "hash", hash, "prove", prove, "result", result)
	return result
}

func (ctx Context) Block(height int64) (res *ctypes.ResultBlock) {

	h := blockHeightConvert(height)

	result, err := ctx.rpcClient.Block(h)
	if err != nil {
		logger.Error("error in getting a block at height", "height", height, "err", err)
		return nil
	}

	return result
}

func (ctx Context) Search(query string, prove bool, page, perPage int) (res *ctypes.ResultTxSearch) {

	result, err := ctx.rpcClient.TxSearch(query, prove, page, perPage)
	if err != nil {
		logger.Error("TxSearch Error", "err", err)
		return nil
	}

	logger.Debug("TxSearch", "query", query, "prove", prove, "result", result)

	return result
}

// abciQuery is a thin wrapper over tendermint rpc client's abciQuery
func (ctx Context) abciQuery(path string, packet []byte) (res *ctypes.ResultABCIQuery, err error) {

	if len(path) < 1 {
		return nil, ErrEmptyQuery
	}

	var response *ctypes.ResultABCIQuery
	response, err = ctx.rpcClient.ABCIQuery(path, packet)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, ErrEmptyResponse
	}

	return response, nil
}

func blockHeightConvert(height int64) (h *int64) {

	// Pass nil if given 0 to return the latest block
	if height != 0 {
		h = &height
	}
	return
}
