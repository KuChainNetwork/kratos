package rest

import (
	"net/http"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/client/utils"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func getCoinsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		account := vars["account"]

		key, err := chainTypes.NewAccountIDFromStr(account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := utils.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		accGetter := types.NewAssetRetriever(cliCtx)

		coin, height, err := accGetter.GetCoins(key)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		utils.PostProcessResponse(w, cliCtx, coin)
	}
}

func getCoinPowersHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		account := vars["account"]

		key, err := chainTypes.NewAccountIDFromStr(account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := utils.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		accGetter := types.NewAssetRetriever(cliCtx)

		coin, height, err := accGetter.GetCoinPowers(key)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		utils.PostProcessResponse(w, cliCtx, coin)
	}
}

func getCoinsLockedHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		account := vars["account"]

		cliCtx, ok := utils.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		accGetter := types.NewAssetRetriever(cliCtx)

		key, err := chainTypes.NewAccountIDFromStr(account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := accGetter.GetLockedCoins(key)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		utils.PostProcessResponse(w, cliCtx, res)
	}
}

func getCoinStatHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		cliCtx, ok := utils.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		accGetter := types.NewAssetRetriever(cliCtx)

		creator, err := chainTypes.NewName(vars["creator"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		symbol, err := chainTypes.NewName(vars["symbol"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := accGetter.GetCoinStat(creator, symbol)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		utils.PostProcessResponse(w, cliCtx, res)
	}
}
