package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/gov/types"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestVotes(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("TestVotes", t, func() {
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

		var invalidOption types.VoteOption = 0x10

		require.Error(t, keeper.AddVote(ctx, proposalID, TestAddrs[0], types.OptionYes), "proposal not on voting period")
		require.Error(t, keeper.AddVote(ctx, 10, TestAddrs[0], types.OptionYes), "invalid proposal ID")

		proposal.Status = types.StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		require.Error(t, keeper.AddVote(ctx, proposalID, TestAddrs[0], invalidOption), "invalid option")

		// Test first vote
		require.NoError(t, keeper.AddVote(ctx, proposalID, TestAddrs[0], types.OptionAbstain))
		vote, found := keeper.GetVote(ctx, proposalID, TestAddrs[0])
		require.True(t, found)
		require.Equal(t, TestAddrs[0], vote.Voter)
		require.Equal(t, proposalID, vote.ProposalID)
		require.Equal(t, types.OptionAbstain, vote.Option)

		// Test change of vote
		require.NoError(t, keeper.AddVote(ctx, proposalID, TestAddrs[0], types.OptionYes))
		vote, found = keeper.GetVote(ctx, proposalID, TestAddrs[0])
		require.True(t, found)
		require.Equal(t, TestAddrs[0], vote.Voter)
		require.Equal(t, proposalID, vote.ProposalID)
		require.Equal(t, types.OptionYes, vote.Option)

		// Test second vote
		require.NoError(t, keeper.AddVote(ctx, proposalID, TestAddrs[1], types.OptionNoWithVeto))
		vote, found = keeper.GetVote(ctx, proposalID, TestAddrs[1])
		require.True(t, found)
		require.Equal(t, TestAddrs[1], vote.Voter)
		require.Equal(t, proposalID, vote.ProposalID)
		require.Equal(t, types.OptionNoWithVeto, vote.Option)

		// Test vote iterator
		// NOTE order of deposits is determined by the addresses
		votes := keeper.GetAllVotes(ctx)
		require.Len(t, votes, 2)
		require.Equal(t, votes, keeper.GetVotes(ctx, proposalID))
		require.Equal(t, TestAddrs[0], votes[1].Voter)
		require.Equal(t, proposalID, votes[1].ProposalID)
		require.Equal(t, types.OptionYes, votes[1].Option)
		require.Equal(t, TestAddrs[1], votes[0].Voter)
		require.Equal(t, proposalID, votes[0].ProposalID)
		require.Equal(t, types.OptionNoWithVeto, votes[0].Option)
	})
}
