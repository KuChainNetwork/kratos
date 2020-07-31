package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/KuChainNetwork/kuchain/x/supply/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetBondedPool returns the bonded tokens pool's module account
func (k Keeper) GetBondedPool(ctx sdk.Context) (bondedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, types.BondedPoolName)
}

// GetNotBondedPool returns the not bonded tokens pool's module account
func (k Keeper) GetNotBondedPool(ctx sdk.Context) (notBondedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, types.NotBondedPoolName)
}

// bondedTokensToNotBonded transfers coins from the bonded to the not bonded pool within staking
func (k Keeper) bondedTokensToNotBonded(ctx sdk.Context, tokens sdk.Int) {
	coins := types.NewCoins(types.NewCoin(k.BondDenom(ctx), tokens))
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.BondedPoolName, types.NotBondedPoolName, coins)
	if err != nil {
		panic(err)
	}
}

// notBondedTokensToBonded transfers coins from the not bonded to the bonded pool within staking
func (k Keeper) notBondedTokensToBonded(ctx sdk.Context, tokens sdk.Int) {
	coins := types.NewCoins(types.NewCoin(k.BondDenom(ctx), tokens))
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.NotBondedPoolName, types.BondedPoolName, coins)
	if err != nil {
		panic(err)
	}
}

// burnBondedTokens removes coins from the bonded pool module account
func (k Keeper) burnBondedTokens(ctx sdk.Context, amt sdk.Int) error {
	if !amt.IsPositive() {
		// skip as no coins need to be burned
		return nil
	}
	coins := types.NewCoins(types.NewCoin(k.BondDenom(ctx), amt))
	return k.supplyKeeper.BurnCoins(ctx, types.BondedPoolAccountID, coins)
}

// burnNotBondedTokens removes coins from the not bonded pool module account
func (k Keeper) burnNotBondedTokens(ctx sdk.Context, amt sdk.Int) error {
	if !amt.IsPositive() {
		// skip as no coins need to be burned
		return nil
	}
	coins := types.NewCoins(types.NewCoin(k.BondDenom(ctx), amt))
	return k.supplyKeeper.BurnCoins(ctx, types.NotBondedPoolAccountID, coins)
}

// TotalBondedTokens total staking tokens supply which is bonded
func (k Keeper) TotalBondedTokens(ctx sdk.Context) sdk.Int {
	bondedPool := k.GetBondedPool(ctx)
	return k.bankKeeper.GetCoinPowerByDenomd(ctx, bondedPool.GetID(), k.BondDenom(ctx)).Amount
}

// TotalBondedTokens total staking tokens supply which is bonded
func (k Keeper) TotalNotBondedTokens(ctx sdk.Context) sdk.Int {
	notBondedPool := k.GetNotBondedPool(ctx)
	return k.bankKeeper.GetCoinPowerByDenomd(ctx, notBondedPool.GetID(), k.BondDenom(ctx)).Amount
}

// StakingTokenSupply staking tokens from the total supply
func (k Keeper) StakingTokenSupply(ctx sdk.Context) sdk.Int {
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(k.BondDenom(ctx))
}

// BondedRatio the fraction of the staking tokens which are currently bonded
func (k Keeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	stakeSupply := k.StakingTokenSupply(ctx)
	if stakeSupply.IsPositive() {
		return k.TotalBondedTokens(ctx).ToDec().QuoInt(stakeSupply)
	}

	return sdk.ZeroDec()
}
