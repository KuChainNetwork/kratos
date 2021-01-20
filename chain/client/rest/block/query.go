package block

import (
	"net/http"
	"strconv"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/client/utils"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func QueryDecodeBlockRequestHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				"couldn't parse block height. Assumed format is '/block/{height}/decode'.")
			return
		}

		chainHeight, err := rpc.GetChainHeight(cliCtx.Ctx())
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "failed to parse chain height")
			return
		}

		if height > chainHeight {
			rest.WriteErrorResponse(w, http.StatusNotFound, "requested block height is bigger then the chain length")
			return
		}

		output, err := getBlock(cliCtx, &height)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponseBare(w, cliCtx, output)
	}
}

func LatestDecodeBlockRequestHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := getBlock(cliCtx, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponseBare(w, cliCtx, output)
	}
}
