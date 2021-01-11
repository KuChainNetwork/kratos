package keeper

import (
	"fmt"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// register all distribution invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "nonnegative-outstanding",
		NonNegativeOutstandingInvariant(k))
	ir.RegisterRoute(types.ModuleName, "can-withdraw",
		CanWithdrawInvariant(k))
	ir.RegisterRoute(types.ModuleName, "reference-count",
		ReferenceCountInvariant(k))
	ir.RegisterRoute(types.ModuleName, "module-account",
		ModuleAccountInvariant(k))
}

// AllInvariants runs all invariants of the distribution module
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := CanWithdrawInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		res, stop = NonNegativeOutstandingInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		res, stop = ReferenceCountInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		return ModuleAccountInvariant(k)(ctx)
	}
}

// NonNegativeOutstandingInvariant checks that outstanding unWithdrawn fees are never negative
func NonNegativeOutstandingInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int
		var outstanding chainTypes.DecCoins

		k.IterateValidatorOutstandingRewards(ctx, func(addr AccountID, rewards types.ValidatorOutstandingRewards) (stop bool) {
			outstanding = rewards.Rewards
			if outstanding.IsAnyNegative() {
				count++
				msg += fmt.Sprintf("\t%v has negative outstanding coins: %v\n", addr, outstanding)
			}
			return false
		})
		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "nonnegative outstanding",
			fmt.Sprintf("found %d validators with negative outstanding rewards\n%s", count, msg)), broken
	}
}

// CanWithdrawInvariant checks that current rewards can be completely withdrawn
func CanWithdrawInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// cache, we don't want to write changes
		ctx, _ = ctx.CacheContext()

		var remaining chainTypes.DecCoins

		valDelegations := make(map[string][]types.AccountID)
		for _, del := range k.stakingKeeper.GetAllSDKDelegations(ctx) {
			val := del.GetValidator().String()
			valDelegations[val] = append(valDelegations[val], del.GetDelegator())
		}

		// iterate over all validators
		k.stakingKeeper.IterateValidators(ctx, func(_ int64, val types.StakingExportedValidatorI) (stop bool) {
			valID := val.GetOperator()
			_, _ = k.WithdrawValidatorCommission(ctx, valID)

			delegations, ok := valDelegations[val.GetOperator().String()]
			if ok {
				for _, del := range delegations {
					if _, err := k.WithdrawDelegationRewards(ctx, del, valID); err != nil {
						panic(err)
					}
				}
			}

			remaining = k.GetValidatorOutstandingRewardsCoins(ctx, val.GetOperator())
			if len(remaining) > 0 && remaining[0].Amount.IsNegative() {
				return true
			}

			return false
		})

		broken := len(remaining) > 0 && remaining[0].Amount.IsNegative()
		return sdk.FormatInvariant(types.ModuleName, "can withdraw",
			fmt.Sprintf("remaining coins: %v\n", remaining)), broken
	}
}

// ReferenceCountInvariant checks that the number of historical rewards records is correct
func ReferenceCountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		valCount := uint64(0)
		k.stakingKeeper.IterateValidators(ctx, func(_ int64, val types.StakingExportedValidatorI) (stop bool) {
			valCount++
			return false
		})
		dels := k.stakingKeeper.GetAllSDKDelegations(ctx)
		slashCount := uint64(0)
		k.IterateValidatorSlashEvents(ctx,
			func(_ AccountID, _ uint64, _ types.ValidatorSlashEvent) (stop bool) {
				slashCount++
				return false
			})

		// one record per validator (last tracked period), one record per
		// delegation (previous period), one record per slash (previous period)
		expected := valCount + uint64(len(dels)) + slashCount
		count := k.GetValidatorHistoricalReferenceCount(ctx)
		broken := count != expected

		return sdk.FormatInvariant(types.ModuleName, "reference count",
			fmt.Sprintf("expected historical reference count: %d = %v validators + %v delegations + %v slashes\n"+
				"total validator historical reference count: %d\n",
				expected, valCount, len(dels), slashCount, count)), broken
	}
}

// ModuleAccountInvariant checks that the coins held by the distr ModuleAccount
// is consistent with the sum of validator outstanding rewards
func ModuleAccountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var expectedCoins chainTypes.DecCoins
		k.IterateValidatorOutstandingRewards(ctx, func(_ AccountID, rewards types.ValidatorOutstandingRewards) (stop bool) {
			expectedCoins = expectedCoins.Add(rewards.Rewards...)
			return false
		})

		communityPool := k.GetFeePoolCommunityCoins(ctx)
		expectedInt, _ := expectedCoins.Add(communityPool...).TruncateDecimal()

		macc := k.GetDistributionAccount(ctx)
		balances := k.BankKeeper.GetCoinPowers(ctx, macc.GetID())

		broken := !balances.IsEqual(expectedInt)
		return sdk.FormatInvariant(
			types.ModuleName, "ModuleAccount coins",
			fmt.Sprintf("\texpected ModuleAccount coins:     %s\n"+
				"\tdistribution ModuleAccount coins: %s\n",
				expectedInt, balances,
			),
		), broken
	}
}
