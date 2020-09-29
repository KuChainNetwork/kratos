package rest

import (
	"io/ioutil"
	"net/http"
	"time"

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

type CreateSymbolReq struct {
	BaseReq       rest.BaseReq        `json:"base_req" yaml:"base_req"`
	Creator       string              `json:"creator" yaml:"creator"`
	Base          types.BaseCurrency  `json:"base" yaml:"base"`
	Quote         types.QuoteCurrency `json:"quote" yaml:"quote"`
	DomainAddress string              `json:"domain_address" yaml:"domain_address"`
}

type UpdateSymbolReq struct {
	BaseReq rest.BaseReq        `json:"base_req" yaml:"base_req"`
	Creator string              `json:"creator" yaml:"creator"`
	Base    types.BaseCurrency  `json:"base" yaml:"base"`
	Quote   types.QuoteCurrency `json:"quote" yaml:"quote"`
}

type ShutdownSymbolReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	Creator   string       `json:"creator" yaml:"creator"`
	BaseCode  string       `json:"base_code" yaml:"base_code"`
	QuoteCode string       `json:"quote_code" yaml:"quote_code"`
}

type UpdateDexReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Creator     string       `json:"creator" yaml:"creator"`
	Description string       `json:"description" yaml:"description"`
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
		"/dex/symbol/create",
		createSymbolHandlerFn(cliCtx),
	).Methods(http.MethodPost)
	r.HandleFunc(
		"/dex/symbol/update",
		updateSymbolHandlerFn(cliCtx),
	).Methods(http.MethodPost)
	r.HandleFunc("/dex/symbol/pause",
		pauseSymbolHandlerFn(cliCtx),
	).Methods(http.MethodPost)
	r.HandleFunc("/dex/symbol/restore",
		restoreSymbolHandlerFn(cliCtx),
	).Methods(http.MethodPost)
	r.HandleFunc(
		"/dex/symbol/shutdown",
		shutdownSymbolHandlerFn(cliCtx),
	).Methods(http.MethodPost)
	r.HandleFunc(
		"/dex/update",
		updateDexHandlerFn(cliCtx),
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

// createSymbolHandlerFn returns the create currency handler
func createSymbolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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
		var req CreateSymbolReq
		if err = cliCtx.Codec.UnmarshalJSON(body, &req); nil != err {
			return
		}
		if !req.Base.Validate() || !req.Quote.Validate() || 0 >= len(req.DomainAddress) {
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
			types.NewMsgCreateSymbol(addr,
				name,
				&req.Base,
				&req.Quote,
				req.DomainAddress,
				time.Time{}, // use server time
			),
		})
	}
}

// updateSymbolHandlerFn returns the update currency handler
func updateSymbolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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
		var req UpdateSymbolReq
		if err = cliCtx.Codec.UnmarshalJSON(body, &req); nil != err {
			return
		}
		if 0 >= len(req.Base.Code) || 0 >= len(req.Quote.Code) ||
			(req.Base.Empty(false) && req.Quote.Empty(false)) {
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
			types.NewMsgUpdateSymbol(addr,
				name,
				&req.Base,
				&req.Quote,
			),
		})
	}
}

// pauseSymbolHandlerFn returns the pause currency handler
func pauseSymbolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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
		var req ShutdownSymbolReq
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
			types.NewMsgPauseSymbol(addr,
				name,
				req.BaseCode,
				req.QuoteCode,
			),
		})
	}
}

// restoreSymbolHandlerFn returns the restore currency handler
func restoreSymbolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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
		var req ShutdownSymbolReq
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
			types.NewMsgRestoreSymbol(addr,
				name,
				req.BaseCode,
				req.QuoteCode,
			),
		})
	}
}

// shutdownSymbolHandlerFn returns the shutdown currency handler
func shutdownSymbolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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
		var req ShutdownSymbolReq
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
			types.NewMsgShutdownSymbol(addr,
				name,
				req.BaseCode,
				req.QuoteCode,
			),
		})
	}
}

// updateDexHandlerFn returns the shutdown currency handler
func updateDexHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if nil != err {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			}
		}()
		var body []byte
		if body, err = ioutil.ReadAll(r.Body); nil != err {
			return
		}

		var req UpdateDexReq
		if err = cliCtx.Codec.UnmarshalJSON(body, &req); nil != err {
			return
		}
		if types.MaxDexDescriptorLen < len(req.Description) {
			err = types.ErrDexDescTooLong

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

			types.NewMsgUpdateDexDescription(addr, name, []byte(req.Description)),
		})
	}
}
