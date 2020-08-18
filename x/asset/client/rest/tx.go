package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	rest "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TransferReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	From    string       `json:"from" yaml:"from"`
	To      string       `json:"to" yaml:"to"`
	Amount  string       `json:"amount" yaml:"amount"`
}

type CreateReq struct {
	BaseReq       rest.BaseReq `json:"base_req" yaml:"base_req"`
	Creator       string       `json:"creator" yaml:"creator"`
	Symbol        string       `json:"symbol" yaml:"symbol"`
	MaxSupply     string       `json:"max_supply" yaml:"max_supply"`
	CanIssue      string       `json:"can_issue" yaml:"can_issue"`
	CanLock       string       `json:"can_lock" yaml:"can_lock"`
	CanBurn       string       `json:"can_burn" yaml:"can_burn"`
	IssueToHeight string       `json:"issue_to_height" yaml:"issue_to_height"`
	InitSupply    string       `json:"init_supply" yaml:"init_supply"`
	Desc          string       `json:"desc" yaml:"desc"`
}

type IssueReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Creator string       `json:"creator" yaml:"creator"`
	Symbol  string       `json:"symbol" yaml:"symbol"`
	Amount  string       `json:"amount" yaml:"amount"`
}

type BurnReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Account string       `json:"account" yaml:"account"`
	Amount  string       `json:"amount" yaml:"amount"`
}

type LockReq struct {
	BaseReq           rest.BaseReq `json:"base_req" yaml:"base_req"`
	Account           string       `json:"account" yaml:"account"`
	UnlockBlockHeight string       `json:"unlock_block_height" yaml:"unlock_block_height"`
	Amount            string       `json:"amount" yaml:"amount"`
}

type UnlockReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Account string       `json:"account" yaml:"account"`
	Amount  string       `json:"amount" yaml:"amount"`
}

type ExerciseReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Account string       `json:"account" yaml:"account"`
	Amount  string       `json:"amount" yaml:"amount"`
}

func TransferRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TransferReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		from, err := chainTypes.NewAccountIDFromStr(req.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(from)
		authAddress, err := txutil.QueryAccountAuth(ctx, from)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		to, err := chainTypes.NewAccountIDFromStr(req.To)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		amount, err := chainTypes.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgTransfer(authAddress, from, to, amount)
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{msg})
	}
}

func CreateRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		creator, err := chainTypes.NewName(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		creatorID := types.NewAccountIDFromName(creator)

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(creatorID)
		auth, err := txutil.QueryAccountAuth(ctx, creatorID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		symbol, err := chainTypes.NewName(req.Symbol)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		maxSupply, err := chainTypes.ParseCoin(req.MaxSupply)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		isCanIssue := req.CanIssue == "1"
		isCanLock := req.CanLock == "1"
		isCanBurn := req.CanBurn == "1"
		issueToHeight, err := strconv.ParseInt(req.IssueToHeight, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		initSupply, err := chainTypes.ParseCoin(req.InitSupply)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("init supply parse error, %s", err.Error()))
			return
		}

		if chainTypes.CoinDenom(creator, symbol) != maxSupply.Denom {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("coin denom should equal %s != %s",
				chainTypes.CoinDenom(creator, symbol), maxSupply.Denom))
			return
		}

		if maxSupply.Denom != initSupply.Denom {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("init coin denom should equal %s != %s",
				initSupply.Denom, maxSupply.Denom))
			return
		}

		if len(req.Desc) > types.CoinDescriptionLen {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("coin description should be less than %d", types.CoinDescriptionLen))
			return
		}

		msg := types.NewMsgCreate(auth, creator, symbol, maxSupply, isCanIssue, isCanLock, isCanBurn, issueToHeight, initSupply, []byte(req.Desc))
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{msg})
	}
}

func IssueRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IssueReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		creator, err := chainTypes.NewName(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		creatorID := types.NewAccountIDFromName(creator)

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(creatorID)
		auth, err := txutil.QueryAccountAuth(ctx, creatorID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		symbol, err := chainTypes.NewName(req.Symbol)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		amount, err := chainTypes.ParseCoin(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if chainTypes.CoinDenom(creator, symbol) != amount.Denom {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("coin denom should equal %s != %s",
				chainTypes.CoinDenom(creator, symbol), amount.Denom))
			return
		}

		msg := types.NewMsgIssue(auth, creator, symbol, amount)
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{msg})
	}
}

func BurnRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BurnReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		acc, err := chainTypes.NewAccountIDFromStr(req.Account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(acc)
		auth, err := txutil.QueryAccountAuth(ctx, acc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		amount, err := chainTypes.ParseCoin(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgBurn(auth, acc, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{msg})
	}
}

func LockRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LockReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		account, err := types.NewAccountIDFromStr(req.Account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("account parse error, %s", err.Error()))
			return
		}

		unlockBlockHeight, err := strconv.Atoi(req.UnlockBlockHeight)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("unlock block height parse error, %s", err.Error()))
			return
		}

		amount, err := chainTypes.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("amount parse error, %s", err.Error()))
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(account)
		auth, err := txutil.QueryAccountAuth(ctx, account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account auth error, %s", err.Error()))
			return
		}

		msg := types.NewMsgLockCoin(auth, account, amount, int64(unlockBlockHeight))
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{msg})
	}
}

func UnlockRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LockReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		account, err := types.NewAccountIDFromStr(req.Account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("account parse error, %s", err.Error()))
			return
		}

		amount, err := chainTypes.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("amount parse error, %s", err.Error()))
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(account)
		auth, err := txutil.QueryAccountAuth(ctx, account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account auth error, %s", err.Error()))
			return
		}

		msg := types.NewMsgUnlockCoin(auth, account, amount)
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{msg})
	}
}

func ExerciseRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ExerciseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		account, err := types.NewAccountIDFromStr(req.Account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("account parse error, %s", err.Error()))
			return
		}

		amount, err := chainTypes.ParseCoin(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("amount parse error, %s", err.Error()))
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(account)
		auth, err := txutil.QueryAccountAuth(ctx, account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account auth error, %s", err.Error()))
			return
		}

		msg := types.NewMsgExercise(auth, account, amount)
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{msg})
	}
}
