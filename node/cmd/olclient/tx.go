package main

import (
	gcontext "context"
	"encoding/hex"
	"github.com/Oneledger/protocol/node/serialize"
	"os"

	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/sdk/pb"
	"github.com/spf13/cobra"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type txRequest struct {
	hash  string
	proof bool
}

var txArgs = &txRequest{}

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Look up a particular transaction by its hash",
	Run:   requestTx,
}

func init() {
	RootCmd.AddCommand(txCmd)
	txCmd.Flags().StringVar(&txArgs.hash, "hash", "", "Get a transaction by its hex-encoded hash")
	txCmd.Flags().BoolVar(&txArgs.proof, "proof", false, "Include proof in the transaction")
}

func requestTx(cmd *cobra.Command, args []string) {
	client := comm.NewSDKClient()
	ctx := gcontext.Background()
	hashBytes, err := hex.DecodeString(txArgs.hash)
	if err != nil {
		shared.Console.Error("Failed to decode hash", err)
		os.Exit(1)
	}
	req := &pb.TxRequest{
		Hash:  hashBytes,
		Proof: txArgs.proof,
	}

	reply, err := client.Tx(ctx, req)
	if err != nil {
		shared.Console.Error("Internal error", err)
		os.Exit(1)
	}

	tx := &ctypes.ResultTx{}
	err = serialize.JSONSzr.Deserialize(reply.Results, tx)
	if err != nil {
		shared.Console.Error("Deserialize failed", err)
		os.Exit(1)
	}

	shared.Console.Info(tx)

}
