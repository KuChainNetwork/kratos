package gov

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, bk types.BankKeeper, supplyKeeper types.SupplyKeeper, k Keeper, data GenesisState) {
	k.SetProposalID(ctx, data.StartingProposalID)
	k.SetDepositParams(ctx, data.DepositParams)
	k.SetVotingParams(ctx, data.VotingParams)
	k.SetTallyParams(ctx, data.TallyParams)

	// check if the deposits pool account exists
	moduleAcc := k.GetGovernanceAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	var totalDeposits types.Coins
	for _, deposit := range data.Deposits {
		k.SetDeposit(ctx, deposit)
		totalDeposits = totalDeposits.Add(deposit.Amount...)
	}

	for _, vote := range data.Votes {
		k.SetVote(ctx, vote)
	}

	for _, proposal := range data.Proposals {
		switch proposal.Status {
		case StatusDepositPeriod:
			k.InsertInactiveProposalQueue(ctx, proposal.ProposalID, proposal.DepositEndTime)
		case StatusVotingPeriod:
			k.InsertActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)
		default:
		}
		k.SetProposal(ctx, proposal)
	}

	if bk.GetCoinPowers(ctx, moduleAcc.GetID()).IsZero() && bk.GetAllBalances(ctx, moduleAcc.GetID()).IsZero() {
		supplyKeeper.SetModuleAccount(ctx, moduleAcc)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingProposalID, _ := k.GetProposalID(ctx)
	depositParams := k.GetDepositParams(ctx)
	votingParams := k.GetVotingParams(ctx)
	tallyParams := k.GetTallyParams(ctx)
	proposals := k.GetProposals(ctx)

	var proposalsDeposits Deposits
	var proposalsVotes Votes
	for _, proposal := range proposals {
		deposits := k.GetDeposits(ctx, proposal.ProposalID)
		proposalsDeposits = append(proposalsDeposits, deposits...)

		votes := k.GetVotes(ctx, proposal.ProposalID)
		proposalsVotes = append(proposalsVotes, votes...)
	}

	pvs := make([]types.PunishValidator, 0, 64)
	k.IterateAllPunishValidators(ctx, func(v types.PunishValidator) bool {
		pvs = append(pvs, v)
		return false
	})

	return GenesisState{
		StartingProposalID: startingProposalID,
		Deposits:           proposalsDeposits,
		Votes:              proposalsVotes,
		Proposals:          proposals,
		DepositParams:      depositParams,
		VotingParams:       votingParams,
		TallyParams:        tallyParams,
		PunishValidators:   pvs,
	}
}
