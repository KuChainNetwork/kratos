package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	rest "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CreateDexReq struct {
	BaseReq     rest.BaseReq `yaml:"base_req" json:"base_req"`
	Creator     string       `yaml:"creator" json:"creator"`
	Stakings    string       `yaml:"stakings" json:"stakings"`
	Description string       `yaml:"description" json:"description"`
}

type DestroyDexReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Creator string       `json:"creator" yml:"creator"`
}

type CreateCurrencyReq struct {
	BaseReq       rest.BaseReq        `json:"base_req" yaml:"base_req"`
	Creator       string              `json:"creator" yaml:"creator"`
	Base          types.BaseCurrency  `json:"base" yaml:"base"`
	Quote         types.QuoteCurrency `json:"quote" yaml:"quote"`
	DomainAddress string              `json:"domain_address" yaml:"domain_address"`
}

type UpdateCurrencyReq struct {
	BaseReq rest.BaseReq        `json:"base_req" yaml:"base_req"`
	Creator string              `json:"creator" yaml:"creator"`
	Base    types.BaseCurrency  `json:"base" yaml:"base"`
	Quote   types.QuoteCurrency `json:"quote" yaml:"quote"`
}

type ShutdownCurrencyReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	Creator   string       `json:"creator" yaml:"creator"`
	BaseCode  string       `json:"base_code" yaml:"base_code"`
	QuoteCode string       `json:"quote_code" yaml:"quote_code"`
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/dex/create",
		createDexHandlerFn(cliCtx),
	).Methods(http.MethodPost)
	r.HandleFunc(
		"/dex/destroy",
		destroyDexHandlerFn(cliCtx),
	).Methods(http.MethodPost)
	r.HandleFunc(
		"/dex/currency/create",
		createCurrencyHandlerFn(cliCtx),
	).Methods(http.MethodPost)
	r.HandleFunc(
		"/dex/currency/update",
		updateCurrencyHandlerFn(cliCtx),
	).Methods(http.MethodPost)
	r.HandleFunc(
		"/dex/currency/shutdown",
		shutdownCurrencyHandlerFn(cliCtx),
	).Methods(http.MethodPost)
}

// createDexHandlerFn returns the create dex handler
func createDexHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusBadRequest
		var err error
		defer func() {
			if nil != err {
				rest.WriteErrorResponse(w, statusCode, err.Error())
			}
		}()
		var body []byte
		if body, err = ioutil.ReadAll(r.Body); nil != err {
			return
		}
		var req CreateDexReq
		if err = cliCtx.Codec.UnmarshalJSON(body, &req); nil != err {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		var creatorName chainTypes.Name
		if creatorName, err = chainTypes.NewName(req.Creator); nil != err {
			return
		}
		var addr chainTypes.AccAddress
		if addr, err = sdk.AccAddressFromBech32(req.BaseReq.From); nil != err {
			return
		}
		var stakings types.Coins
		if stakings, err = chainTypes.ParseCoins(req.Stakings); nil != err {
			return
		}
		var creatorAccountID chainTypes.AccountID
		if creatorAccountID, err = chainTypes.NewAccountIDFromStr(req.Creator); nil != err {
			return
		}
		if types.MaxDexDescriptorLen < len(req.Description) {
			err = types.ErrDexDescTooLong
			return
		}
		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(creatorAccountID)
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{
			types.NewMsgCreateDex(addr, creatorName, stakings, []byte(req.Description)),
		})
	}
}

// destroyDexHandlerFn returns the destroy dex handler
func destroyDexHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusBadRequest
		var err error
		defer func() {
			if nil != err {
				rest.WriteErrorResponse(w, statusCode, err.Error())
			}
		}()
		var body []byte
		if body, err = ioutil.ReadAll(r.Body); nil != err {
			return
		}
		var req DestroyDexReq
		if err = cliCtx.Codec.UnmarshalJSON(body, &req); nil != err {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		var name chainTypes.Name
		if name, err = chainTypes.NewName(req.Creator); nil != err {
			return
		}
		var addr chainTypes.AccAddress
		if addr, err = sdk.AccAddressFromBech32(req.BaseReq.From); nil != err {
			return
		}
		var creatorAccountID chainTypes.AccountID
		if creatorAccountID, err = chainTypes.NewAccountIDFromStr(req.Creator); nil != err {
			return
		}
		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(creatorAccountID)
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{
			types.NewMsgDestroyDex(addr, name),
		})
	}
}

// createCurrencyHandlerFn returns the create currency handler
func createCurrencyHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusBadRequest
		var err error
		defer func() {
			if nil != err {
				rest.WriteErrorResponse(w, statusCode, err.Error())
			}
		}()
		var body []byte
		if body, err = ioutil.ReadAll(r.Body); nil != err {
			return
		}
		var req CreateCurrencyReq
		if err = cliCtx.Codec.UnmarshalJSON(body, &req); nil != err {
			return
		}
		if 0 >= len(req.Base.Code) ||
			0 >= len(req.Base.Name) ||
			0 >= len(req.Base.FullName) ||
			0 >= len(req.Base.IconUrl) ||
			0 >= len(req.Base.TxUrl) ||
			0 >= len(req.Quote.Code) ||
			0 >= len(req.Quote.Name) ||
			0 >= len(req.Quote.FullName) ||
			0 >= len(req.Quote.IconUrl) ||
			0 >= len(req.Quote.TxUrl) ||
			0 >= len(req.DomainAddress) {
			err = errors.Errorf("incorrect request fields")
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		var name chainTypes.Name
		if name, err = chainTypes.NewName(req.Creator); nil != err {
			return
		}
		var addr chainTypes.AccAddress
		if addr, err = sdk.AccAddressFromBech32(req.BaseReq.From); nil != err {
			return
		}
		var creatorAccountID chainTypes.AccountID
		if creatorAccountID, err = chainTypes.NewAccountIDFromStr(req.Creator); nil != err {
			return
		}
		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(creatorAccountID)
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{
			types.NewMsgCreateCurrency(addr,
				name,
				&req.Base,
				&req.Quote,
				req.DomainAddress,
			),
		})
	}
}

// updateCurrencyHandlerFn returns the update currency handler
func updateCurrencyHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusBadRequest
		var err error
		defer func() {
			if nil != err {
				rest.WriteErrorResponse(w, statusCode, err.Error())
			}
		}()
		var body []byte
		if body, err = ioutil.ReadAll(r.Body); nil != err {
			return
		}
		var req UpdateCurrencyReq
		if err = cliCtx.Codec.UnmarshalJSON(body, &req); nil != err {
			return
		}
		if 0 >= len(req.Base.Code) &&
			0 >= len(req.Base.Name) &&
			0 >= len(req.Base.FullName) &&
			0 >= len(req.Base.IconUrl) &&
			0 >= len(req.Base.TxUrl) &&
			0 >= len(req.Quote.Code) &&
			0 >= len(req.Quote.Name) &&
			0 >= len(req.Quote.FullName) &&
			0 >= len(req.Quote.IconUrl) &&
			0 >= len(req.Quote.TxUrl) {
			err = errors.Errorf("incorrect request fields")
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		var name chainTypes.Name
		if name, err = chainTypes.NewName(req.Creator); nil != err {
			return
		}
		var addr chainTypes.AccAddress
		if addr, err = sdk.AccAddressFromBech32(req.BaseReq.From); nil != err {
			return
		}
		var creatorAccountID chainTypes.AccountID
		if creatorAccountID, err = chainTypes.NewAccountIDFromStr(req.Creator); nil != err {
			return
		}
		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(creatorAccountID)
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{
			types.NewMsgUpdateCurrency(addr,
				name,
				&req.Base,
				&req.Quote,
			),
		})
	}
}

// updateCurrencyHandlerFn returns the shutdown currency handler
func shutdownCurrencyHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusBadRequest
		var err error
		defer func() {
			if nil != err {
				rest.WriteErrorResponse(w, statusCode, err.Error())
			}
		}()
		var body []byte
		if body, err = ioutil.ReadAll(r.Body); nil != err {
			return
		}
		var req ShutdownCurrencyReq
		if err = cliCtx.Codec.UnmarshalJSON(body, &req); nil != err {
			return
		}
		if 0 >= len(req.BaseCode) ||
			0 >= len(req.QuoteCode) {
			err = errors.Errorf("incorrect request fields")
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		var name chainTypes.Name
		if name, err = chainTypes.NewName(req.Creator); nil != err {
			return
		}
		var addr chainTypes.AccAddress
		if addr, err = sdk.AccAddressFromBech32(req.BaseReq.From); nil != err {
			return
		}
		var creatorAccountID chainTypes.AccountID
		if creatorAccountID, err = chainTypes.NewAccountIDFromStr(req.Creator); nil != err {
			return
		}
		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(creatorAccountID)
		txutil.WriteGenerateStdTxResponse(w, ctx, req.BaseReq, []sdk.Msg{
			types.NewMsgShutdownCurrency(addr,
				name,
				req.BaseCode,
				req.QuoteCode,
			),
		})
	}
}
