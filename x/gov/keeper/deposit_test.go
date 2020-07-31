package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/staking/exported"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestDeposits(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("TestDeposits", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		stakingKeeper := app.StakeKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID

		fourStake := chainTypes.NewCoins(chainTypes.NewCoin(stakingKeeper.BondDenom(ctx), exported.TokensFromConsensusPower(220)))
		fiveStake := chainTypes.NewCoins(chainTypes.NewCoin(stakingKeeper.BondDenom(ctx), exported.TokensFromConsensusPower(250)))

		require.True(t, proposal.TotalDeposit.IsEqual(chainTypes.NewCoins()))

		// Check no deposits at beginning
		deposit, found := keeper.GetDeposit(ctx, proposalID, TestAddrs[1])
		require.False(t, found)
		proposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		require.True(t, proposal.VotingStartTime.Equal(time.Time{}))

		// Check first deposit
		votingStarted, err := keeper.AddDeposit(ctx, proposalID, TestAddrs[0], fourStake)
		require.NoError(t, err)
		require.False(t, votingStarted)
		deposit, found = keeper.GetDeposit(ctx, proposalID, TestAddrs[0])
		require.True(t, found)
		require.Equal(t, fourStake, deposit.Amount)
		require.Equal(t, TestAddrs[0], deposit.Depositor)
		proposal, ok = keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		require.Equal(t, fourStake, proposal.TotalDeposit)

		// Check a second deposit from same address
		votingStarted, err = keeper.AddDeposit(ctx, proposalID, TestAddrs[0], fiveStake)
		require.NoError(t, err)
		require.False(t, votingStarted)
		deposit, found = keeper.GetDeposit(ctx, proposalID, TestAddrs[0])
		require.True(t, found)
		require.Equal(t, fourStake.Add(fiveStake...), deposit.Amount)
		require.Equal(t, TestAddrs[0], deposit.Depositor)
		proposal, ok = keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		require.Equal(t, fourStake.Add(fiveStake...), proposal.TotalDeposit)

		// Check third deposit from a new address
		votingStarted, err = keeper.AddDeposit(ctx, proposalID, TestAddrs[1], fourStake)
		require.NoError(t, err)
		require.True(t, votingStarted)
		deposit, found = keeper.GetDeposit(ctx, proposalID, TestAddrs[1])
		require.True(t, found)
		require.Equal(t, TestAddrs[1], deposit.Depositor)
		require.Equal(t, fourStake, deposit.Amount)
		proposal, ok = keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		require.Equal(t, fourStake.Add(fiveStake...).Add(fourStake...), proposal.TotalDeposit)

		// Check that proposal moved to voting period
		proposal, ok = keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		require.True(t, proposal.VotingStartTime.Equal(ctx.BlockHeader().Time))

		// Test deposit iterator
		// NOTE order of deposits is determined by the addresses
		deposits := keeper.GetAllDeposits(ctx)
		require.Len(t, deposits, 2)
		require.Equal(t, deposits, keeper.GetDeposits(ctx, proposalID))
		require.Equal(t, TestAddrs[1], deposits[0].Depositor)
		require.Equal(t, fourStake, deposits[0].Amount)
		require.Equal(t, TestAddrs[0], deposits[1].Depositor)
		require.Equal(t, fourStake.Add(fiveStake...), deposits[1].Amount)

		// Test Refund Deposits
		deposit, found = keeper.GetDeposit(ctx, proposalID, TestAddrs[1])
		require.True(t, found)
		require.Equal(t, fourStake, deposit.Amount)
		keeper.RefundDeposits(ctx, proposalID)
		deposit, found = keeper.GetDeposit(ctx, proposalID, TestAddrs[1])
		require.False(t, found)
	})
}
