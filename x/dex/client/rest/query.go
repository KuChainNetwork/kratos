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
	r.HandleFunc("/dex/sigIn/{creator}/{account}",
		getSigInStatusHandlerFn(cliCtx),
	).Methods(http.MethodGet)
	r.HandleFunc(
		"/dex/symbol/{creator}/{baseCreator}/{baseCode}/{quoteCreator}/{quoteCode}",
		getSymbolHandlerFn(cliCtx),
	).Methods(http.MethodGet)
}

// getCreatorHandlerFn function returns the get dex REST handler.
func getCreatorHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["creator"]
		creator, err := chainTypes.NewName(name)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		getter := types.NewDexRetriever(cliCtx)
		dex, _, err := getter.GetDexWithHeight(creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, dex)
	}
}

// getSigInStatusHandlerFn
func getSigInStatusHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		creatorStr := vars["creator"]
		accountStr := vars["account"]
		creator, err := chainTypes.NewAccountIDFromStr(creatorStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		account, err := chainTypes.NewAccountIDFromStr(accountStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		getter := types.NewDexRetriever(cliCtx)
		coins, _, err := getter.GetSigInWithHeight(account, creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, coins)
	}
}

// getSymbolHandlerFn returns get symbaol REST handler
func getSymbolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["creator"]
		baseCreator := vars["baseCreator"]
		baseCode := vars["baseCode"]
		quoteCreator := vars["quoteCreator"]
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
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		getter := types.NewDexRetriever(cliCtx)
		var dex *types.Dex
		if dex, _, err = getter.GetDexWithHeight(creator); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		baseCreatorName, err := chainTypes.NewName(baseCreator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		baseCodeName, err := chainTypes.NewName(baseCode)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		baseCode = types.CoinDenom(baseCreatorName, baseCodeName)
		quoteCreatorName, err := chainTypes.NewName(quoteCreator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		quoteCodeName, err := chainTypes.NewName(quoteCode)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		quoteCode = types.CoinDenom(quoteCreatorName, quoteCodeName)
		currency, ok := dex.Symbol(baseCreator, baseCode, quoteCreator, quoteCode)
		if !ok {
			rest.WriteErrorResponse(w,
				http.StatusNotFound,
				fmt.Sprintf("%s/%s not exists", baseCode, quoteCode))
			return
		}
		rest.PostProcessResponse(w, cliCtx, currency)
	}
}
