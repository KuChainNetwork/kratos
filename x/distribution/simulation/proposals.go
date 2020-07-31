package simulation

import (
	"math/rand"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/keeper"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
)

// OpWeightSubmitCommunitySpendProposal app params key for community spend proposal
const OpWeightSubmitCommunitySpendProposal = "op_weight_submit_community_spend_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper) []sim.WeightedProposalContent {
	return []sim.WeightedProposalContent{
		{
			AppParamsKey:       OpWeightSubmitCommunitySpendProposal,
			DefaultWeight:      simappparams.DefaultWeightCommunitySpendProposal,
			ContentSimulatorFn: SimulateCommunityPoolSpendProposalContent(k),
		},
	}
}

// SimulateCommunityPoolSpendProposalContent generates random community-pool-spend proposal content
func SimulateCommunityPoolSpendProposalContent(k keeper.Keeper) types.SimulationContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []types.SimulationAccount) govtypes.Content {
		simAccount, _ := types.SimulationRandomAcc(r, accs)

		balance := k.GetFeePool(ctx).CommunityPool
		if balance.Empty() {
			return nil
		}

		denomIndex := r.Intn(len(balance))
		amount, err := types.SimulationRandPositiveInt(r, balance[denomIndex].Amount.TruncateInt())
		if err != nil {
			return nil
		}

		aid, _ := chainTypes.NewAccountIDFromStr(string(simAccount.Address))

		return types.NewCommunityPoolSpendProposal(
			types.SimulationRandStringOfLength(r, 10),
			types.SimulationRandStringOfLength(r, 100),
			aid,
			chainTypes.NewCoins(chainTypes.NewCoin(balance[denomIndex].Denom, amount)),
		)
	}
}
