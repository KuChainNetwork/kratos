package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/gov/external"
	"github.com/KuChainNetwork/kuchain/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: Break into several smaller functions for clarity

// Tally iterates over the votes and updates the tally of a proposal based on the voting power of the
// voters
func (keeper Keeper) Tally(ctx sdk.Context, proposal types.Proposal) (passes bool, burnDeposits bool, tallyResults types.TallyResult, punishBp []AccountID, punish bool, vetobp []AccountID) {
	results := make(map[types.VoteOption]sdk.Dec)
	results[types.OptionYes] = sdk.ZeroDec()
	results[types.OptionAbstain] = sdk.ZeroDec()
	results[types.OptionNo] = sdk.ZeroDec()
	results[types.OptionNoWithVeto] = sdk.ZeroDec()

	totalVotingPower := sdk.ZeroDec()
	currValidators := make(map[string]types.ValidatorGovInfo)
	// fetch all the bonded validators, insert them into currValidators
	keeper.sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator external.StakingValidatorI) (stop bool) {
		currValidators[validator.GetOperatorAccountID().String()] = types.NewValidatorGovInfo(
			validator.GetOperatorAccountID(),
			validator.GetBondedTokens(),
			validator.GetDelegatorShares(),
			sdk.ZeroDec(),
			types.OptionEmpty,
		)

		return false
	})

	keeper.IterateVotes(ctx, proposal.ProposalID, func(vote types.Vote) bool {
		//if validator, just record it in the map
		valAddrStr := vote.Voter.String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Option
			currValidators[valAddrStr] = val
		}

		keeper.deleteVote(ctx, vote.ProposalID, vote.Voter)
		return false
	})

	var punishValidators []AccountID
	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if val.Vote == types.OptionEmpty {
			punishValidators = append(punishValidators, val.Address)
			continue
		}

		if val.Vote == types.OptionNoWithVeto {
			vetobp = append(vetobp, val.Address)
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		fractionAfterDeductions := sharesAfterDeductions.Quo(val.DelegatorShares)
		votingPower := fractionAfterDeductions.MulInt(val.BondedTokens)

		results[val.Vote] = results[val.Vote].Add(votingPower)
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	tallyParams := keeper.GetTallyParams(ctx)
	tallyResults = types.NewTallyResultFromMap(results)

	// TODO: Upgrade the spec to cover all of these cases & remove pseudocode.
	// If there is no staked coins, the proposal fails
	if keeper.sk.TotalBondedTokens(ctx).IsZero() {
		return false, false, tallyResults, punishValidators, false, vetobp
	}

	// If there is not enough quorum of votes, the proposal fails
	percentVoting := totalVotingPower.Quo(keeper.sk.TotalBondedTokens(ctx).ToDec())
	if percentVoting.LT(tallyParams.Quorum) {
		return false, true, tallyResults, punishValidators, false, vetobp
	}

	// If no one votes (everyone abstains), proposal fails
	if totalVotingPower.Sub(results[types.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, false, tallyResults, punishValidators, false, vetobp
	}

	// If more than 1/3 of voters veto, proposal fails
	if results[types.OptionNoWithVeto].Quo(totalVotingPower).GT(tallyParams.Veto) {
		return false, true, tallyResults, punishValidators, true, vetobp
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if results[types.OptionYes].Quo(totalVotingPower.Sub(results[types.OptionAbstain])).GT(tallyParams.Threshold) {
		return true, false, tallyResults, punishValidators, false, vetobp
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails
	return false, false, tallyResults, punishValidators, false, vetobp
}

func (keeper Keeper) EmergencyPass(ctx sdk.Context, proposalID uint64) (passes bool, tallyResults types.TallyResult) {
	results := make(map[types.VoteOption]sdk.Dec)
	results[types.OptionYes] = sdk.ZeroDec()
	results[types.OptionAbstain] = sdk.ZeroDec()
	results[types.OptionNo] = sdk.ZeroDec()
	results[types.OptionNoWithVeto] = sdk.ZeroDec()

	currValidators := make(map[string]types.ValidatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators      TotalBondedTokens
	keeper.sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator external.StakingValidatorI) (stop bool) {
		currValidators[validator.GetOperatorAccountID().String()] = types.NewValidatorGovInfo(
			validator.GetOperatorAccountID(),
			validator.GetBondedTokens(),
			validator.GetDelegatorShares(),
			sdk.ZeroDec(),
			types.OptionEmpty,
		)

		return false
	})

	keeper.IterateVotes(ctx, proposalID, func(vote types.Vote) bool {
		//if validator, just record it in the map
		valAddrStr := vote.Voter.String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Option
			currValidators[valAddrStr] = val
		}
		return false
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if val.Vote == types.OptionEmpty {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		fractionAfterDeductions := sharesAfterDeductions.Quo(val.DelegatorShares)
		votingPower := fractionAfterDeductions.MulInt(val.BondedTokens)

		results[val.Vote] = results[val.Vote].Add(votingPower)
	}

	tallyParams := keeper.GetTallyParams(ctx)
	tallyResults = types.NewTallyResultFromMap(results)

	// TODO: Upgrade the spec to cover all of these cases & remove pseudocode.
	// If there is no staked coins, the proposal fails
	if keeper.sk.TotalBondedTokens(ctx).IsZero() {
		return false, tallyResults
	}

	// If more than 2/3 of the votes are in favor, proposal passes
	if results[types.OptionYes].Quo(keeper.sk.TotalBondedTokens(ctx).ToDec()).GT(tallyParams.Emergency) {
		return true, tallyResults
	}
	return false, tallyResults
}

//add jail information
func (keeper Keeper) Jail(ctx sdk.Context, validatorAccount AccountID, proposalID uint64) {
	punishValdator := types.NewPunishValidator(validatorAccount, ctx.BlockHeader().Height, ctx.BlockHeader().Time.Add(keeper.DowntimeJailDuration(ctx)), proposalID)
	keeper.SetPunishValidator(ctx, punishValdator)
	keeper.sk.JailByAccount(ctx, validatorAccount)
}

func (keeper Keeper) SlashValidator(ctx sdk.Context, validatorAccount AccountID) {
	keeper.sk.SlashByValidatorAccount(ctx, validatorAccount, ctx.BlockHeader().Height, keeper.GetSlashFraction(ctx))
}

func (keeper Keeper) Slash(ctx sdk.Context) {
	keeper.distrKeeper.SetStartNotDistributionTimePoint(ctx, ctx.BlockHeader().Time)
}
