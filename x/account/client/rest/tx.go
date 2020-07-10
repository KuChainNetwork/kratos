package rest

import (
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	rest "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type CreateAccountReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Creator     string       `json:"creator" yaml:"creator"`
	Account     string       `json:"account" yaml:"account"`
	AccountAuth string       `json:"account_auth" yaml:"account_auth"`
}

type UpdateAuthReq struct {
	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	Account        string       `json:"account" yaml:"account"`
	NewAccountAuth string       `json:"new_account_auth" yaml:"new_account_auth"`
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/account/create",
		createAccountHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/account/update_auth",
		updateAuthHandlerFn(cliCtx),
	).Methods("POST")
}

func createAccountHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateAccountReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		creator, err := chainTypes.NewAccountIDFromStr(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(creator)
		auth, err := txutil.QueryAccountAuth(ctx, creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		account, err := chainTypes.NewName(req.Account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		accountAuth, err := sdk.AccAddressFromBech32(req.AccountAuth)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreateAccount(auth, creator, account, accountAuth)
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{msg})
	}
}

func updateAuthHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateAuthReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		accountName, err := chainTypes.NewName(req.Account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		account := chainTypes.NewAccountIDFromName(accountName)

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(account)
		auth, err := txutil.QueryAccountAuth(ctx, account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		newAccountAuth, err := sdk.AccAddressFromBech32(req.NewAccountAuth)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgUpdateAccountAuth(auth, accountName, newAccountAuth)
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{msg})
	}
}
