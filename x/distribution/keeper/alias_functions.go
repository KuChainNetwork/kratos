package keeper

import (
	chainType "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/distribution/types"
	supplyexported "github.com/KuChain-io/kuchain/x/supply/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	AccountID = chainType.AccountID
)

// get outstanding rewards , old by cancer
func (k Keeper) GetValidatorOutstandingRewardsCoins(ctx sdk.Context, val AccountID) sdk.DecCoins {
	return k.GetValidatorOutstandingRewards(ctx, val).Rewards
}

// get the community coins
func (k Keeper) GetFeePoolCommunityCoins(ctx sdk.Context) sdk.DecCoins {
	return k.GetFeePool(ctx).CommunityPool
}

// GetDistributionAccount returns the distribution ModuleAccount
func (k Keeper) GetDistributionAccount(ctx sdk.Context) supplyexported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}
