package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	blockRest "github.com/KuChainNetwork/kuchain/chain/client/rest/block"
	txRest "github.com/KuChainNetwork/kuchain/chain/client/rest/tx"
)

func RegisterBlockRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/blocks/latest/decode", blockRest.LatestDecodeBlockRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/blocks/{height}/decode", blockRest.QueryDecodeBlockRequestHandlerFn(cliCtx)).Methods("GET")
}

// RegisterTxRoutes registers all transaction routes on the provided router.
func RegisterTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/txs/{hash}", txRest.QueryTxRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", txRest.QueryTxsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", txRest.BroadcastTxRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/encode", txRest.EncodeTxRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/decode", txRest.DecodeTxRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/sign_msg/encode", txRest.EncodeMsgRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/fee", txRest.QueryTxFeeAndGasConsumed(cliCtx)).Methods("POST")
}
