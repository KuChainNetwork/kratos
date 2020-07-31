package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/gov/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKeeper(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("TestIncrementProposalNumber", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		tp := TestProposal
		keeper.SubmitProposal(ctx, tp)
		keeper.SubmitProposal(ctx, tp)
		keeper.SubmitProposal(ctx, tp)
		keeper.SubmitProposal(ctx, tp)
		keeper.SubmitProposal(ctx, tp)
		proposal6, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)

		require.Equal(t, uint64(6), proposal6.ProposalID)
	})
	Convey("TestProposalQueues", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		// create test proposals
		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)

		inactiveIterator := keeper.InactiveProposalQueueIterator(ctx, proposal.DepositEndTime)
		require.True(t, inactiveIterator.Valid())

		proposalID := types.GetProposalIDFromBytes(inactiveIterator.Value())
		require.Equal(t, proposalID, proposal.ProposalID)
		inactiveIterator.Close()

		keeper.ActivateVotingPeriod(ctx, proposal)

		proposal, ok := keeper.GetProposal(ctx, proposal.ProposalID)
		require.True(t, ok)

		activeIterator := keeper.ActiveProposalQueueIterator(ctx, proposal.VotingEndTime)
		require.True(t, activeIterator.Valid())
		proposalID, _ = types.SplitActiveProposalQueueKey(activeIterator.Key())
		require.Equal(t, proposalID, proposal.ProposalID)
		activeIterator.Close()
	})
}
