package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	blockRest "github.com/KuChainNetwork/kuchain/chain/client/rest/block"
)

func RegisterBlockRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/blocks/latest/decode", blockRest.LatestDecodeBlockRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/blocks/{height}/decode", blockRest.QueryDecodeBlockRequestHandlerFn(cliCtx)).Methods("GET")
}
