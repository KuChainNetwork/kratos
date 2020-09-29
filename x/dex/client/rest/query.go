package rest

import (
	"fmt"
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
	).Methods(http.MethodGet)
	r.HandleFunc(
		"/dex/symbol/{creator}/{baseCode}/{quoteCode}",
		getSymbolHandlerFn(cliCtx),
	).Methods(http.MethodGet)
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

// getSymbolHandlerFn returns get symbaol REST handler
func getSymbolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["creator"]
		baseCode := vars["baseCode"]
		quoteCode := vars["quoteCode"]
		if 0 >= len(name) ||
			0 >= len(baseCode) ||
			0 >= len(quoteCode) {
			rest.WriteErrorResponse(w,
				http.StatusBadRequest,
				"creator, base code or quote code is empty")
			return
		}
		creator, err := chainTypes.NewName(name)
		if nil != err {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		getter := types.NewDexRetriever(cliCtx)
		var dex *types.Dex
		if dex, _, err = getter.GetDexWithHeight(creator); nil != err {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		currency, ok := dex.Symbol(baseCode, quoteCode)
		if !ok {
			rest.WriteErrorResponse(w,
				http.StatusNotFound,
				fmt.Sprintf("%s/%s not exists", baseCode, quoteCode))
			return
		}
		rest.PostProcessResponse(w, cliCtx, currency)
	}
}
