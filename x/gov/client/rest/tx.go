package rest

import (
	"fmt"
	"net/http"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	govutils "github.com/KuChainNetwork/kuchain/x/gov/client/utils"
	"github.com/KuChainNetwork/kuchain/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
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
		if !chainTypes.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		deposit, err := chainTypes.ParseCoins(req.InitialDeposit)
		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposerAccount, err := chainTypes.NewAccountIDFromStr(req.ProposerAcc)
		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("proposer account id error, %v", err))
			return
		}

		content := types.ContentFromProposalType(req.Title, req.Description, types.ProposalTypeText)

		proposalAccAddress, err := txutil.QueryAccountAuth(cliCtx, proposerAccount)
		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account %s auth error, %v", proposerAccount, err))
			return
		}

		msg := types.NewKuMsgSubmitProposal(proposalAccAddress, content, deposit, proposerAccount)
		if err := msg.ValidateBasic(); err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func depositHandlerFn(cliCtx txutil.KuCLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DepositReq
		if !chainTypes.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		if len(req.ProposalID) == 0 {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, "ProposalID required but not specified")
			return
		}

		proposalID, ok := chainTypes.ParseUint64OrReturnBadRequest(w, req.ProposalID)
		if !ok {
			return
		}

		amount, err := chainTypes.ParseCoins(req.Amount)
		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		depositor, err := chainTypes.NewAccountIDFromStr(req.Depositor)
		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("depositor account id error, %v", err))
			return
		}

		// Get depositor address
		depositorAccAddress, err := txutil.QueryAccountAuth(cliCtx, depositor)
		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account %s auth error, %v", depositor, err))
			return
		}

		msg := types.NewKuMsgDeposit(depositorAccAddress, depositor, proposalID, amount)
		if err := msg.ValidateBasic(); err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func voteHandlerFn(cliCtx txutil.KuCLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req VoteReq
		if !chainTypes.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if len(req.ProposalID) == 0 {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, "ProposalID required but not specified")
			return
		}

		proposalID, ok := chainTypes.ParseUint64OrReturnBadRequest(w, req.ProposalID)
		if !ok {
			return
		}

		voteOption, err := types.VoteOptionFromString(govutils.NormalizeVoteOption(req.Option))
		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		VoterAccount, err := chainTypes.NewAccountIDFromStr(req.Voter)
		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("depositor account id error, %v", err))
			return
		}

		voterAccAddress, err := txutil.QueryAccountAuth(cliCtx, VoterAccount)
		if err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account %s auth error, %v", VoterAccount, err))
			return
		}

		msg := types.NewKuMsgVote(voterAccAddress, VoterAccount, proposalID, voteOption)
		if err := msg.ValidateBasic(); err != nil {
			chainTypes.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
