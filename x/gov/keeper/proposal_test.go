package keeper_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

type validProposal struct{}

func (validProposal) GetTitle() string       { return "title" }
func (validProposal) GetDescription() string { return "description" }
func (validProposal) ProposalRoute() string  { return types.RouterKey }
func (validProposal) ProposalType() string   { return types.ProposalTypeText }
func (validProposal) String() string         { return "" }
func (validProposal) ValidateBasic() error   { return nil }

type invalidProposalTitle1 struct{ validProposal }

func (invalidProposalTitle1) GetTitle() string { return "" }

type invalidProposalTitle2 struct{ validProposal }

func (invalidProposalTitle2) GetTitle() string { return strings.Repeat("1234567890", 100) }

type invalidProposalDesc1 struct{ validProposal }

func (invalidProposalDesc1) GetDescription() string { return "" }

type invalidProposalDesc2 struct{ validProposal }

func (invalidProposalDesc2) GetDescription() string { return strings.Repeat("1234567890", 1000) }

type invalidProposalRoute struct{ validProposal }

func (invalidProposalRoute) ProposalRoute() string { return "nonexistingroute" }

type invalidProposalValidation struct{ validProposal }

func (invalidProposalValidation) ValidateBasic() error {
	return errors.New("invalid proposal")
}

func registerTestCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(validProposal{}, "test/validproposal", nil)
	cdc.RegisterConcrete(invalidProposalTitle1{}, "test/invalidproposalt1", nil)
	cdc.RegisterConcrete(invalidProposalTitle2{}, "test/invalidproposalt2", nil)
	cdc.RegisterConcrete(invalidProposalDesc1{}, "test/invalidproposald1", nil)
	cdc.RegisterConcrete(invalidProposalDesc2{}, "test/invalidproposald2", nil)
	cdc.RegisterConcrete(invalidProposalRoute{}, "test/invalidproposalr", nil)
	cdc.RegisterConcrete(invalidProposalValidation{}, "test/invalidproposalv", nil)
}

func TestProposalAll(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("TestGetSetProposal", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)
		proposalID := proposal.ProposalID
		keeper.SetProposal(ctx, proposal)

		gotProposal, ok := keeper.GetProposal(ctx, proposalID)
		require.True(t, ok)
		require.True(t, ProposalEqual(keeper, proposal, gotProposal))
	})
	Convey("TestGetSetProposal", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		tp := TestProposal
		proposal, err := keeper.SubmitProposal(ctx, tp)
		require.NoError(t, err)

		require.True(t, proposal.VotingStartTime.Equal(time.Time{}))

		keeper.ActivateVotingPeriod(ctx, proposal)

		require.True(t, proposal.VotingStartTime.Equal(ctx.BlockHeader().Time))

		proposal, ok := keeper.GetProposal(ctx, proposal.ProposalID)
		require.True(t, ok)

		activeIterator := keeper.ActiveProposalQueueIterator(ctx, proposal.VotingEndTime)
		require.True(t, activeIterator.Valid())

		proposalID := types.GetProposalIDFromBytes(activeIterator.Value())
		require.Equal(t, proposalID, proposal.ProposalID)
		activeIterator.Close()
	})
	Convey("TestGetSetProposal", t, func() {
		proposalID := uint64(1)
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.GovKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		status := []types.ProposalStatus{types.StatusDepositPeriod, types.StatusVotingPeriod}

		addr1 := Accdel[30]

		for _, s := range status {
			for i := 0; i < 50; i++ {
				p := types.NewProposal(TestProposal, proposalID, time.Now(), time.Now())
				p.Status = s

				if i%2 == 0 {
					d := types.NewDeposit(proposalID, addr1, nil)
					v := types.NewVote(proposalID, addr1, types.OptionYes)
					keeper.SetDeposit(ctx, d)
					keeper.SetVote(ctx, v)
				}

				keeper.SetProposal(ctx, p)
				proposalID++
			}
		}

		testCases := []struct {
			params             types.QueryProposalsParams
			expectedNumResults int
		}{
			{types.NewQueryProposalsParams(1, 50, types.StatusNil, chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID()), 50},
			{types.NewQueryProposalsParams(1, 50, types.StatusDepositPeriod, chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID()), 50},
			{types.NewQueryProposalsParams(1, 50, types.StatusVotingPeriod, chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID()), 50},
			{types.NewQueryProposalsParams(1, 25, types.StatusNil, chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID()), 25},
			{types.NewQueryProposalsParams(2, 25, types.StatusNil, chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID()), 25},
			{types.NewQueryProposalsParams(1, 50, types.StatusRejected, chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID()), 0},
			{types.NewQueryProposalsParams(1, 50, types.StatusNil, addr1, chainTypes.EmptyAccountID()), 50},
			{types.NewQueryProposalsParams(1, 50, types.StatusNil, chainTypes.EmptyAccountID(), addr1), 50},
			{types.NewQueryProposalsParams(1, 50, types.StatusNil, addr1, addr1), 50},
			{types.NewQueryProposalsParams(1, 50, types.StatusDepositPeriod, addr1, addr1), 25},
			{types.NewQueryProposalsParams(1, 50, types.StatusDepositPeriod, chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID()), 50},
			{types.NewQueryProposalsParams(1, 50, types.StatusVotingPeriod, chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID()), 50},
		}

		for _, tc := range testCases {
			proposals := keeper.GetProposalsFiltered(ctx, tc.params)
			require.Len(t, proposals, tc.expectedNumResults)

			for _, p := range proposals {
				if len(tc.params.ProposalStatus.String()) != 0 {
					require.Equal(t, tc.params.ProposalStatus, p.Status)
				}
			}
		}
	})
}
