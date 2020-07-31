package keeper

import (
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	supplyexported "github.com/KuChainNetwork/kuchain/x/supply/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Coins    = chainTypes.Coins
	Coin     = chainTypes.Coin
	DecCoins = chainTypes.DecCoins
	DecCoin  = chainTypes.DecCoin
)

var (
	NewDec = chainTypes.NewDec
)

type (
	AccountID = chainTypes.AccountID
)

// get outstanding rewards , old by cancer
func (k Keeper) GetValidatorOutstandingRewardsCoins(ctx sdk.Context, val AccountID) DecCoins {
	return k.GetValidatorOutstandingRewards(ctx, val).Rewards
}

// get the community coins
func (k Keeper) GetFeePoolCommunityCoins(ctx sdk.Context) DecCoins {
	return k.GetFeePool(ctx).CommunityPool
}

// GetDistributionAccount returns the distribution ModuleAccount
func (k Keeper) GetDistributionAccount(ctx sdk.Context) supplyexported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}
