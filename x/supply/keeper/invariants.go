package keeper

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/x/supply/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterInvariants register all supply invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "total-supply", TotalSupply(k))
}

// AllInvariants runs all invariants of the supply module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return TotalSupply(k)(ctx)
	}
}

// TotalSupply checks that the total supply reflects all the coins held in accounts
func TotalSupply(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var expectedTotal types.Coins
		supply := k.GetSupply(ctx)

		k.bk.IterateAllCoins(ctx, func(_ types.AccountID, balance types.Coins) bool {
			expectedTotal = expectedTotal.Add(balance...)
			return false
		})

		k.bk.IterateAllCoinPowers(ctx, func(_ types.AccountID, balance types.Coins) bool {
			expectedTotal = expectedTotal.Add(balance...)
			return false
		})

		broken := !expectedTotal.IsEqual(supply.GetTotal())

		return sdk.FormatInvariant(types.ModuleName, "total supply",
			fmt.Sprintf(
				"\tsum of accounts coins: %v\n"+
					"\tsupply.Total:          %v\n",
				expectedTotal, supply.GetTotal())), broken
	}
}
