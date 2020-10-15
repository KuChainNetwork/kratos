package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterBlockRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/blocks/latest/decode", LatestDecodeBlockRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/blocks/{height}/decode", QueryDecodeBlockRequestHandlerFn(cliCtx)).Methods("GET")
}
