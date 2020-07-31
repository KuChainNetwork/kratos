package rest

import (
	"fmt"
	"net/http"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	rest "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/params/external"
	"github.com/KuChainNetwork/kuchain/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PostProposalParamsReq struct {
	BaseReq        rest.BaseReq           `json:"base_req" yaml:"base_req"`
	Title          string                 `json:"title" yaml:"title"`                     // Title of the proposal
	Description    string                 `json:"description" yaml:"description"`         // Description of the proposal
	InitialDeposit string                 `json:"initial_deposit" yaml:"initial_deposit"` // Coins to add to the proposal's deposit
	ProposerAcc    string                 `json:"proposer_acc" yaml:"proposer_acc"`       // account of the proposer
	ParamChanges   []proposal.ParamChange `json:"param_changes" yaml:"param_changes"`
}

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the param
// change REST handler with a given sub-route.
func ProposalRESTHandler(cliCtx context.CLIContext) external.GovProposalRESTHandler {
	ctx := txutil.NewKuCLICtx(cliCtx)
	return external.GovProposalRESTHandler{
		SubRoute: "param_change",
		Handler:  postProposalParamsHandlerFn(ctx),
	}
}

func postProposalParamsHandlerFn(cliCtx txutil.KuCLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostProposalParamsReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		ProposalAccount, err := chainTypes.NewAccountIDFromStr(req.ProposerAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("proposer account id error, %v", err))
			return
		}

		content := proposal.NewParameterChangeProposal(req.Title, req.Description, req.ParamChanges)
		deposit, err := chainTypes.ParseCoins(req.InitialDeposit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// Get proposal address
		authAccAddress, err := txutil.QueryAccountAuth(cliCtx, ProposalAccount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account %s auth error, %v", ProposalAccount, err))
			return
		}
		msg := external.GovNewMsgSubmitProposal(authAccAddress, content, deposit, ProposalAccount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx = cliCtx.WithFromAccount(ProposalAccount)

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
