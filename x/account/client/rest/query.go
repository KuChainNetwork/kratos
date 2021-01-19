package rest

import (
	"fmt"
	"net/http"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/client/utils"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/account/{name}",
		getAccountHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/account/auth/{auth}",
		getAuthHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/accounts/{auth}",
		getAccountsByAuthHandlerFn(cliCtx),
	).Methods("GET")
}

func getAccountHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		accGetter := types.NewAccountRetriever(cliCtx)

		key, err := chainTypes.NewAccountIDFromStr(name)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx, ok := utils.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		acc, height, err := accGetter.GetAccountWithHeight(key)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		utils.PostProcessResponse(w, cliCtx, acc)
	}
}

func getAuthHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		auth := vars["auth"]

		cliCtx, ok := utils.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		accGetter := types.NewAccountRetriever(cliCtx)

		key, err := chainTypes.AccAddressFromBech32(auth)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("new acc-address error %v", err.Error()))
			return
		}

		data, height, err := accGetter.GetAddAuthWithHeight(key)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		utils.PostProcessResponse(w, cliCtx, data)
	}
}

func getAccountsByAuthHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		auth := vars["auth"]

		cliCtx, ok := utils.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryAccountsByAuthParams(auth)
		bz, err := cliCtx.Codec().MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAccountsByAuth)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var result types.Accounts
		if err = cliCtx.Codec().UnmarshalJSON(res, &result); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		utils.PostProcessResponse(w, cliCtx, result)
	}
}
