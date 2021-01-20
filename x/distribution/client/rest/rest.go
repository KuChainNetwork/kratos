package rest

import (
	"net/http"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	rest "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	govRest "github.com/KuChainNetwork/kuchain/x/gov/client/rest"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
)

// RegisterRoutes register distribution REST routes.
func RegisterRoutes(cliCtx client.Context, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r, queryRoute)
}

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func ProposalRESTHandler(cliCtx client.Context) govRest.ProposalRESTHandler {
	return govRest.ProposalRESTHandler{
		SubRoute: "community_pool_spend",
		Handler:  postProposalHandlerFn(cliCtx),
	}
}

func postProposalHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CommunityPoolSpendProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec(), &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewCommunityPoolSpendProposal(req.Title, req.Description, req.Recipient, req.Amount)
		msg := types.GovTypesNewKuMsgSubmitProposal(req.ProposerAccAddress, content, req.Deposit, req.Proposer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
