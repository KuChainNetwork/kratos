package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

// RegisterTxRoutes registers all transaction routes on the provided router.
func RegisterTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", QueryTxsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", BroadcastTxRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/encode", EncodeTxRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/decode", DecodeTxRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/sign_msg/encode", EncodeMsgRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/fee", QueryTxFeeAndGasConsumed(cliCtx)).Methods("POST")
}
