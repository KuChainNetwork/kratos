package rest

import (
	"fmt"
	"github.com/KuChain-io/kuchain/chain/client/txutil"
	chaintype "github.com/KuChain-io/kuchain/chain/types"
	rest "github.com/KuChain-io/kuchain/chain/types"
	govutils "github.com/KuChain-io/kuchain/x/gov/client/utils"
	"github.com/KuChain-io/kuchain/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"net/http"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, phs []ProposalRESTHandler) {
	propSubRtr := r.PathPrefix("/gov/proposals").Subrouter()
	for _, ph := range phs {
		propSubRtr.HandleFunc(fmt.Sprintf("/%s", ph.SubRoute), ph.Handler).Methods("POST")
	}

	kuCliCtx := txutil.NewKuCLICtx(cliCtx)

	r.HandleFunc("/gov/proposals", postProposalHandlerFn(kuCliCtx)).Methods("POST")
	r.HandleFunc("/gov/deposits", depositHandlerFn(kuCliCtx)).Methods("POST")
	r.HandleFunc("/gov/votes", voteHandlerFn(kuCliCtx)).Methods("POST")
}

func postProposalHandlerFn(cliCtx txutil.KuCLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		deposit, err := sdk.ParseCoins(req.InitialDeposit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposerAccount, err := chaintype.NewAccountIDFromStr(req.ProposerAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("proposer account id error, %v", err))
			return
		}

		content := types.ContentFromProposalType(req.Title, req.Description, types.ProposalTypeText)

		proposalAccAddress, err := txutil.QueryAccountAuth(cliCtx, proposerAccount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account %s auth error, %v", proposerAccount, err))
			return
		}

		msg := types.NewKuMsgSubmitProposal(proposalAccAddress, content, deposit, proposerAccount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func depositHandlerFn(cliCtx txutil.KuCLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DepositReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		if len(req.ProposalId) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "proposalId required but not specified")
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, req.ProposalId)
		if !ok {
			return
		}

		amount, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		depositor, err := chaintype.NewAccountIDFromStr(req.Depositor)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("depositor account id error, %v", err))
			return
		}

		// Get depositor address
		depositorAccAddress, err := txutil.QueryAccountAuth(cliCtx, depositor)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account %s auth error, %v", depositor, err))
			return
		}

		msg := types.NewKuMsgDeposit(depositorAccAddress, depositor, proposalID, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func voteHandlerFn(cliCtx txutil.KuCLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req VoteReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if len(req.ProposalId) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "proposalId required but not specified")
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, req.ProposalId)
		if !ok {
			return
		}

		voteOption, err := types.VoteOptionFromString(govutils.NormalizeVoteOption(req.Option))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		VoterAccount, err := chaintype.NewAccountIDFromStr(req.Voter)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("depositor account id error, %v", err))
			return
		}

		voterAccAddress, err := txutil.QueryAccountAuth(cliCtx, VoterAccount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account %s auth error, %v", VoterAccount, err))
			return
		}

		msg := types.NewKuMsgVote(voterAccAddress, VoterAccount, proposalID, voteOption)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
