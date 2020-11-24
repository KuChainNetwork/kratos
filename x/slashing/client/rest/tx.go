package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	ctx := txutil.NewKuCLICtx(cliCtx)
	r.HandleFunc(
		"/slashing/unjail",
		unjailRequestHandlerFn(ctx),
	).Methods("POST")
}

// Unjail TX body
type UnjailReq struct {
	BaseReq      chainTypes.BaseReq `json:"base_req" yaml:"base_req"`
	ValidatorAcc string             `json:"validator_acc" yaml:"validator_acc"`
}

// FIX HERE
func unjailRequestHandlerFn(cliCtx txutil.KuCLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UnjailReq
		if !chainTypes.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		valAccountID, err := chainTypes.NewAccountIDFromStr(req.ValidatorAcc)

		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		valAccAddress, err := txutil.QueryAccountAuth(cliCtx, valAccountID)
		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account %s auth error : %s", valAccountID, err.Error()))
			return
		}

		msg := types.NewKuMsgUnjail(valAccAddress, valAccountID)
		err = msg.ValidateBasic()
		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
