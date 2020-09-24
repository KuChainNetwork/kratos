package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

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

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/dex/create",
		createDexHandlerFn(cliCtx),
	).Methods(http.MethodPost)
	r.HandleFunc(
		"/dex/destroy",
		destroyDexHandlerFn(cliCtx),
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
