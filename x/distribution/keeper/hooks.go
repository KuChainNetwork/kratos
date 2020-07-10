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

var _ types.StakingTypesStakingHooks = Hooks{} //bugs , stacking interface

// Create new distribution hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// initialize validator distribution record
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valId chainType.AccountID) {
	val := h.k.stakingKeeper.Validator(ctx, valId)
	h.k.initializeValidator(ctx, val)
}

// cleanup for after validator is removed
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valId chainType.AccountID) {
	// fetch outstanding
	outstanding := h.k.GetValidatorOutstandingRewardsCoins(ctx, valId)

	// force-withdraw commission
	commission := h.k.GetValidatorAccumulatedCommission(ctx, valId).Commission
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

			//accAddr := sdk.AccAddress(valAddr)
			withdrawAddr := h.k.GetDelegatorWithdrawAddr(ctx, valId)
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
	h.k.DeleteValidatorOutstandingRewards(ctx, valId)

	// remove commission record
	h.k.DeleteValidatorAccumulatedCommission(ctx, valId)

	// clear slashes
	h.k.DeleteValidatorSlashEvents(ctx, valId)

	// clear historical rewards
	h.k.DeleteValidatorHistoricalRewards(ctx, valId)

	// clear current rewards
	h.k.DeleteValidatorCurrentRewards(ctx, valId)
}

// increment period
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, _ chainType.AccountID, valId chainType.AccountID) {
	val := h.k.stakingKeeper.Validator(ctx, valId)
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
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delId chainType.AccountID, valId chainType.AccountID) {
	h.k.initializeDelegation(ctx, valId, delId)
}

// record the slash event
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valId chainType.AccountID, fraction sdk.Dec) {
	h.k.updateValidatorSlashFraction(ctx, valId, fraction)
}

// nolint - unused hooks
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ chainType.AccountID)                 {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ chainType.AccountID) {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ chainType.AccountID) {
}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ chainType.AccountID, _ chainType.AccountID) {}
