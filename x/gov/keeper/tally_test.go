package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KuChainNetwork/kuchain/x/staking"
	"github.com/KuChainNetwork/kuchain/x/staking/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/gov/types"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestTally(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("TestTallyNoOneVotes", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, powers)

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.False(t, passes)
		require.True(t, burnDeposits)
		require.True(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyNoOneVotes", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{2, 5, 0})

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		err = keeper.AddVote(ctx, proposalID, TestAddrs[0], types.OptionYes)
		require.Nil(t, err)

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, _, _, _, _ := keeper.Tally(ctx, proposal)
		require.False(t, passes)
		require.True(t, burnDeposits)
	})
	Convey("TestTallyOnlyValidatorsAllYes", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{5, 5, 5})

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr3, types.OptionYes))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.True(t, passes)
		require.False(t, burnDeposits)
		require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyOnlyValidators51No", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{5, 6, 0})

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionNo))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, _, _, _, _ := keeper.Tally(ctx, proposal)

		require.False(t, passes)
		require.False(t, burnDeposits)
	})
	Convey("TestTallyOnlyValidators51No", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{5, 6, 0})

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionNo))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionYes))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.True(t, passes)
		require.False(t, burnDeposits)
		require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyOnlyValidatorsVetoed", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{6, 6, 7})

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr3, types.OptionNoWithVeto))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.False(t, passes)
		require.True(t, burnDeposits)
		require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyOnlyValidatorsAbstainPasses", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{6, 6, 7})

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionAbstain))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionNo))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr3, types.OptionYes))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.True(t, passes)
		require.False(t, burnDeposits)
		require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyOnlyValidatorsAbstainFails", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{6, 6, 7})

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionAbstain))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr3, types.OptionNo))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.False(t, passes)
		require.False(t, burnDeposits)
		require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyOnlyValidatorsNonVoter", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{5, 6, 7})

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionNo))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.False(t, passes)
		require.False(t, burnDeposits)
		require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyDelgatorOverride", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{5, 6, 7})

		delTokens := exported.TokensFromConsensusPower(30)
		val1, found := stakingKeeper.GetValidator(ctx, valOpAddr1)
		require.True(t, found)

		_, err := stakingKeeper.Delegate(ctx, TestAddrs[0], delTokens, exported.Unbonded, val1, true)
		require.NoError(t, err)

		_ = staking.EndBlocker(ctx, *stakingKeeper)

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr3, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, TestAddrs[0], types.OptionNo))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.False(t, passes)
		require.False(t, burnDeposits)
		require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyDelgatorInherit", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{5, 6, 7})

		delTokens := exported.TokensFromConsensusPower(30)
		val3, found := stakingKeeper.GetValidator(ctx, valOpAddr3)
		require.True(t, found)

		_, err := stakingKeeper.Delegate(ctx, TestAddrs[0], delTokens, exported.Unbonded, val3, true)
		require.NoError(t, err)

		_ = staking.EndBlocker(ctx, *stakingKeeper)

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionNo))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionNo))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr3, types.OptionYes))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.True(t, passes)
		require.False(t, burnDeposits)
		require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyDelgatorMultipleOverride", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{5, 6, 7})

		delTokens := exported.TokensFromConsensusPower(10)
		val1, found := stakingKeeper.GetValidator(ctx, valOpAddr1)
		require.True(t, found)
		val2, found := stakingKeeper.GetValidator(ctx, valOpAddr2)
		require.True(t, found)

		_, err := stakingKeeper.Delegate(ctx, TestAddrs[0], delTokens, exported.Unbonded, val1, true)
		require.NoError(t, err)
		_, err = stakingKeeper.Delegate(ctx, TestAddrs[0], delTokens, exported.Unbonded, val2, true)
		require.NoError(t, err)

		_ = staking.EndBlocker(ctx, *stakingKeeper)

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr3, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, TestAddrs[0], types.OptionNo))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.True(t, passes)
		require.False(t, burnDeposits)
		require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyDelgatorMultipleInherit", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{25, 6, 7})

		delTokens := exported.TokensFromConsensusPower(10)
		val2, found := stakingKeeper.GetValidator(ctx, valOpAddr2)
		require.True(t, found)
		val3, found := stakingKeeper.GetValidator(ctx, valOpAddr3)
		require.True(t, found)

		_, err := stakingKeeper.Delegate(ctx, TestAddrs[0], delTokens, exported.Unbonded, val2, true)
		require.NoError(t, err)
		_, err = stakingKeeper.Delegate(ctx, TestAddrs[0], delTokens, exported.Unbonded, val3, true)
		require.NoError(t, err)

		_ = staking.EndBlocker(ctx, *stakingKeeper)

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionNo))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr3, types.OptionNo))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.False(t, passes)
		require.False(t, burnDeposits)
		require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyJailedValidator", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{25, 6, 7})

		delTokens := exported.TokensFromConsensusPower(10)
		val2, found := stakingKeeper.GetValidator(ctx, valOpAddr2)
		require.True(t, found)
		val3, found := stakingKeeper.GetValidator(ctx, valOpAddr3)
		require.True(t, found)

		_, err := stakingKeeper.Delegate(ctx, TestAddrs[0], delTokens, exported.Unbonded, val2, true)
		require.NoError(t, err)
		_, err = stakingKeeper.Delegate(ctx, TestAddrs[0], delTokens, exported.Unbonded, val3, true)
		require.NoError(t, err)

		_ = staking.EndBlocker(ctx, *stakingKeeper)

		stakingKeeper.Jail(ctx, sdk.ConsAddress(val2.GetConsPubKey().Address()))

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionNo))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr3, types.OptionNo))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.True(t, passes)
		require.False(t, burnDeposits)
		require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
	})
	Convey("TestTallyValidatorMultipleDelegations", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		stakingKeeper = stakingKeeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		createValidators(app, ctx, stakingKeeper, []int64{10, 10, 10})

		delTokens := exported.TokensFromConsensusPower(10)
		val2, found := stakingKeeper.GetValidator(ctx, valOpAddr2)
		require.True(t, found)

		_, err := stakingKeeper.Delegate(ctx, valAccAddr1, delTokens, exported.Unbonded, val2, true)
		require.NoError(t, err)

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionYes))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionNo))
		require.NoError(t, keeper.AddVote(ctx, proposalID, valAccAddr3, types.OptionYes))

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		passes, burnDeposits, tallyResults, _, _, _ := keeper.Tally(ctx, proposal)

		require.False(t, passes)
		require.False(t, burnDeposits)

		expectedYes := exported.TokensFromConsensusPower(20)
		expectedAbstain := exported.TokensFromConsensusPower(0)
		expectedNo := exported.TokensFromConsensusPower(20)
		expectedNoWithVeto := exported.TokensFromConsensusPower(0)
		expectedTallyResult := types.NewTallyResult(expectedYes, expectedAbstain, expectedNo, expectedNoWithVeto)

		require.True(t, tallyResults.Equals(expectedTallyResult))
	})
}
