package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/dex/{creator}",
		getCreatorHandlerFn(cliCtx),
	).Methods("GET")
}

// getCreatorHandlerFn function returns the get dex REST handler.
func getCreatorHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["creator"]
		creator, err := chainTypes.NewName(name)
		if nil != err {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		getter := types.NewDexRetriever(cliCtx)
		dex, _, err := getter.GetDexWithHeight(creator)
		if nil != err {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, dex)
	}
}
