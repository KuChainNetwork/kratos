package rest

import (
	"net/http"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	rest "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/client/common"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	// Withdraw all delegator rewards
	r.HandleFunc(
		"/distribution/delegators/rewards",
		withdrawDelegatorRewardsHandlerFn(cliCtx, queryRoute),
	).Methods("POST")

	// Withdraw delegation rewards
	r.HandleFunc(
		"/distribution/delegators_validator/rewards",
		withdrawDelegationRewardsHandlerFn(cliCtx),
	).Methods("POST")

	// Replace the rewards withdrawal address
	r.HandleFunc(
		"/distribution/delegators/withdraw_account",
		setDelegatorWithdrawalAddrHandlerFn(cliCtx),
	).Methods("POST")

	// Withdraw validator rewards and commission
	r.HandleFunc(
		"/distribution/validators/rewards",
		withdrawValidatorRewardsHandlerFn(cliCtx),
	).Methods("POST")

}

type (
	withdrawRewardsReq struct {
		BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
		DelegatorAcc string       `json:"delegator_acc" yaml:"delegator_acc"`
		ValidatorAcc string       `json:"validator_acc" yaml:"validator_acc"`
	}

	setWithdrawalAddrReq struct {
		BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
		DelegatorAcc string       `json:"delegator_acc" yaml:"delegator_acc"`
		WithdrawAcc  string       `json:"withdraw_acc" yaml:"withdraw_acc"`
	}

	fundCommunityPoolReq struct {
		BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
		Amount       string       `json:"amount" yaml:"amount"`
		DepositorAcc string       `json:"depositor_acc" yaml:"depositor_acc"`
	}
)

// Withdraw delegator rewards
func withdrawDelegatorRewardsHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawRewardsReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		delegatorAcc, err := chainTypes.NewAccountIDFromStr(req.DelegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(delegatorAcc)
		auth, err := txutil.QueryAccountAuth(ctx, delegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg, err := common.WithdrawAllDelegatorRewards(cliCtx, auth, queryRoute, delegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, txutil.NewKuCLICtx(cliCtx), req.BaseReq, msg)
	}
}

// Withdraw delegation rewards
func withdrawDelegationRewardsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawRewardsReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		delegatorAcc, err := chainTypes.NewAccountIDFromStr(req.DelegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		validatorAcc, err := chainTypes.NewAccountIDFromStr(req.ValidatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(delegatorAcc)
		auth, err := txutil.QueryAccountAuth(ctx, delegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgWithdrawDelegatorReward(auth, delegatorAcc, validatorAcc)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, txutil.NewKuCLICtx(cliCtx), req.BaseReq, []sdk.Msg{msg})
	}
}

// Replace the rewards withdrawal address
func setDelegatorWithdrawalAddrHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setWithdrawalAddrReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) { //bugs, x
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		delegatorAcc, err := chainTypes.NewAccountIDFromStr(req.DelegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		withdrawAcc, err := chainTypes.NewAccountIDFromStr(req.WithdrawAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(delegatorAcc)
		auth, err := txutil.QueryAccountAuth(ctx, delegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSetWithdrawAccountId(auth, delegatorAcc, withdrawAcc)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, txutil.NewKuCLICtx(cliCtx), req.BaseReq, []sdk.Msg{msg})
	}
}

// Withdraw validator rewards and commission
func withdrawValidatorRewardsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawRewardsReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		validatorAcc, err := chainTypes.NewAccountIDFromStr(req.ValidatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(validatorAcc)
		auth, err := txutil.QueryAccountAuth(ctx, validatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgWithdrawDelegatorReward(auth, validatorAcc, validatorAcc)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, txutil.NewKuCLICtx(cliCtx), req.BaseReq, []sdk.Msg{msg})
	}
}

func checkDelegatorAddressVar(w http.ResponseWriter, r *http.Request) (chainTypes.AccountID, bool) {
	accID, err := chainTypes.NewAccountIDFromStr(mux.Vars(r)["delegatorAddr"])
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return chainTypes.EmptyAccountID(), false
	}

	return accID, true
}

func checkValidatorAddressVar(w http.ResponseWriter, r *http.Request) (chainTypes.AccountID, bool) {
	addr, err := chainTypes.NewAccountIDFromStr(mux.Vars(r)["validatorAddr"])
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return chainTypes.EmptyAccountID(), false
	}

	return addr, true
}
