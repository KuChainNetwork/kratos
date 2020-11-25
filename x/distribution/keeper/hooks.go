package keeper

import (
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ types.StakingTypesStakingHooks = Hooks{}

// Create new distribution hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// initialize validator distribution record
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valID chainType.AccountID) {
	val := h.k.stakingKeeper.Validator(ctx, valID)
	h.k.initializeValidator(ctx, val)
}

// cleanup for after validator is removed
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valID chainType.AccountID) {
	// fetch outstanding
	outstanding := h.k.GetValidatorOutstandingRewardsCoins(ctx, valID)

	// force-withdraw commission
	commission := h.k.GetValidatorAccumulatedCommission(ctx, valID).Commission
	if !commission.IsZero() {
		// subtract from outstanding
		outstanding = outstanding.Sub(commission)

		// split into integral & remainder
		coins, remainder := commission.TruncateDecimal()

		// remainder to community pool
		feePool := h.k.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
		h.k.SetFeePool(ctx, feePool)

		// add to validator account
		if !coins.IsZero() {
			withdrawAddr := h.k.GetDelegatorWithdrawAddr(ctx, valID)
			err := h.k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins)
			if err != nil {
				panic(err)
			}
		}
	}

	// add outstanding to community pool
	feePool := h.k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(outstanding...)
	h.k.SetFeePool(ctx, feePool)

	// delete outstanding
	h.k.DeleteValidatorOutstandingRewards(ctx, valID)

	// remove commission record
	h.k.DeleteValidatorAccumulatedCommission(ctx, valID)

	// clear slashes
	h.k.DeleteValidatorSlashEvents(ctx, valID)

	// clear historical rewards
	h.k.DeleteValidatorHistoricalRewards(ctx, valID)

	// clear current rewards
	h.k.DeleteValidatorCurrentRewards(ctx, valID)
}

// increment period
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, _ chainType.AccountID, valID chainType.AccountID) {
	val := h.k.stakingKeeper.Validator(ctx, valID)
	h.k.IncrementValidatorPeriod(ctx, val)
}

// withdraw delegation rewards (which also increments period)
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr chainType.AccountID, valAId chainType.AccountID) {
	val := h.k.stakingKeeper.Validator(ctx, valAId)
	del := h.k.stakingKeeper.Delegation(ctx, delAddr, valAId)
	if _, err := h.k.withdrawDelegationRewards(ctx, val, del); err != nil {
		panic(err)
	}
}

// create new delegation period record
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delID chainType.AccountID, valID chainType.AccountID) {
	h.k.initializeDelegation(ctx, valID, delID)
}

// record the slash event
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valID chainType.AccountID, fraction sdk.Dec) {
	h.k.updateValidatorSlashFraction(ctx, valID, fraction)
}

// nolint - unused hooks
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ chainType.AccountID)                 {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ chainType.AccountID) {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ chainType.AccountID) {
}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ chainType.AccountID, _ chainType.AccountID) {}
