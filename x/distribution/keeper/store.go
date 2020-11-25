package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/store"
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
)

// get the delegator withdraw address, defaulting to the delegator address
func (k Keeper) GetDelegatorWithdrawAddr(ctx sdk.Context, delAddr AccountID) AccountID {
	store := store.NewStore(ctx, k.storeKey)
	b := store.Get(types.GetDelegatorWithdrawAddrKey(delAddr))
	if b == nil {
		return delAddr
	}
	return chainType.NewAccountIDFromByte(b)
}

// set the delegator withdraw address
func (k Keeper) SetDelegatorWithdrawAddr(ctx sdk.Context, delAddr, withdrawAddr AccountID) {
	store := store.NewStore(ctx, k.storeKey)
	store.Set(types.GetDelegatorWithdrawAddrKey(delAddr), withdrawAddr.StoreKey())
}

// delete a delegator withdraw addr
func (k Keeper) DeleteDelegatorWithdrawAddr(ctx sdk.Context, delID, withdrawID AccountID) {
	store := store.NewStore(ctx, k.storeKey)
	store.Delete(types.GetDelegatorWithdrawAddrKey(delID))
}

// iterate over delegator withdraw addrs
func (k Keeper) IterateDelegatorWithdrawAddID(ctx sdk.Context, handler func(del AccountID, addr AccountID) (stop bool)) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DelegatorWithdrawAddrPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addID, _ := chainType.NewAccountIDFromStr(string(iter.Value()))
		del := types.GetDelegatorWithdrawInfoAddressUseAccountID(iter.Key())
		if handler(del, addID) {
			break
		}
	}
}

// get the global fee pool distribution info
func (k Keeper) GetFeePool(ctx sdk.Context) (feePool types.FeePool) {
	store := store.NewStore(ctx, k.storeKey)
	b := store.Get(types.FeePoolKey)
	if b == nil {
		panic("Stored fee pool should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &feePool)
	return
}

// set the global fee pool distribution info
func (k Keeper) SetFeePool(ctx sdk.Context, feePool types.FeePool) {
	store := store.NewStore(ctx, k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&feePool)
	store.Set(types.FeePoolKey, b)
}

// GetPreviousProposerConsAddr returns the proposer consensus address for the
// current block.
// by cancer
func (k Keeper) GetPreviousProposerConsAddr(ctx sdk.Context) sdk.ConsAddress {
	store := store.NewStore(ctx, k.storeKey)
	bz := store.Get(types.ProposerKey)
	if bz == nil {
		panic("previous proposer not set")
	}

	addrValue := gogotypes.BytesValue{}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &addrValue)

	return sdk.ConsAddress(addrValue.Value)
}

// set the proposer public key for this block
func (k Keeper) SetPreviousProposerConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := store.NewStore(ctx, k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&gogotypes.BytesValue{Value: consAddr.Bytes()})
	store.Set(types.ProposerKey, bz)
}

// get the starting info associated with a delegator
func (k Keeper) GetDelegatorStartingInfo(ctx sdk.Context, val AccountID, del AccountID) (period types.DelegatorStartingInfo) {
	store := store.NewStore(ctx, k.storeKey)
	b := store.Get(types.GetDelegatorStartingInfoKey(val, del))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &period)
	return
}

// set the starting info associated with a delegator
func (k Keeper) SetDelegatorStartingInfo(ctx sdk.Context, val AccountID, del AccountID, period types.DelegatorStartingInfo) {
	store := store.NewStore(ctx, k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&period)
	store.Set(types.GetDelegatorStartingInfoKey(val, del), b)
}

// check existence of the starting info associated with a delegator
func (k Keeper) HasDelegatorStartingInfo(ctx sdk.Context, val AccountID, del AccountID) bool {
	store := store.NewStore(ctx, k.storeKey)
	return store.Has(types.GetDelegatorStartingInfoKey(val, del))
}

// delete the starting info associated with a delegator
func (k Keeper) DeleteDelegatorStartingInfo(ctx sdk.Context, val AccountID, del AccountID) {
	store := store.NewStore(ctx, k.storeKey)
	store.Delete(types.GetDelegatorStartingInfoKey(val, del))
}

// iterate over delegator starting infos
func (k Keeper) IterateDelegatorStartingInfos(ctx sdk.Context, handler func(val AccountID, del AccountID, info types.DelegatorStartingInfo) (stop bool)) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DelegatorStartingInfoPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var info types.DelegatorStartingInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &info)
		val, del := types.GetDelegatorStartingInfoAddresses(iter.Key())
		if handler(val, del, info) {
			break
		}
	}
}

// get historical rewards for a particular period
func (k Keeper) GetValidatorHistoricalRewards(ctx sdk.Context, val AccountID, period uint64) (rewards types.ValidatorHistoricalRewards) {
	store := store.NewStore(ctx, k.storeKey)
	b := store.Get(types.GetValidatorHistoricalRewardsKey(val, period))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
	return
}

// set historical rewards for a particular period
func (k Keeper) SetValidatorHistoricalRewards(ctx sdk.Context, val AccountID, period uint64, rewards types.ValidatorHistoricalRewards) {
	store := store.NewStore(ctx, k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&rewards)
	store.Set(types.GetValidatorHistoricalRewardsKey(val, period), b)
}

// iterate over historical rewards
func (k Keeper) IterateValidatorHistoricalRewards(ctx sdk.Context,
	handler func(val AccountID, period uint64, rewards types.ValidatorHistoricalRewards) (stop bool)) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorHistoricalRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards types.ValidatorHistoricalRewards
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &rewards)
		addr, period := types.GetValidatorHistoricalRewardsAddressPeriod(iter.Key())
		if handler(addr, period, rewards) {
			break
		}
	}
}

// delete a historical reward
func (k Keeper) DeleteValidatorHistoricalReward(ctx sdk.Context, val AccountID, period uint64) {
	store := store.NewStore(ctx, k.storeKey)
	store.Delete(types.GetValidatorHistoricalRewardsKey(val, period))
}

// delete historical rewards for a validator
func (k Keeper) DeleteValidatorHistoricalRewards(ctx sdk.Context, val AccountID) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetValidatorHistoricalRewardsPrefix(val))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// delete all historical rewards
func (k Keeper) DeleteAllValidatorHistoricalRewards(ctx sdk.Context) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorHistoricalRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// historical reference count (used for testcases)
func (k Keeper) GetValidatorHistoricalReferenceCount(ctx sdk.Context) (count uint64) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorHistoricalRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards types.ValidatorHistoricalRewards
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &rewards)
		count += uint64(rewards.ReferenceCount)
	}
	return
}

// get current rewards for a validator
func (k Keeper) GetValidatorCurrentRewards(ctx sdk.Context, val AccountID) (rewards types.ValidatorCurrentRewards) {
	store := store.NewStore(ctx, k.storeKey)
	b := store.Get(types.GetValidatorCurrentRewardsKey(val))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
	return
}

// set current rewards for a validator
func (k Keeper) SetValidatorCurrentRewards(ctx sdk.Context, val AccountID, rewards types.ValidatorCurrentRewards) {
	store := store.NewStore(ctx, k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&rewards)
	store.Set(types.GetValidatorCurrentRewardsKey(val), b)
}

// delete current rewards for a validator
func (k Keeper) DeleteValidatorCurrentRewards(ctx sdk.Context, val AccountID) {
	store := store.NewStore(ctx, k.storeKey)
	store.Delete(types.GetValidatorCurrentRewardsKey(val))
}

// iterate over current rewards
func (k Keeper) IterateValidatorCurrentRewards(ctx sdk.Context, handler func(val AccountID, rewards types.ValidatorCurrentRewards) (stop bool)) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorCurrentRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards types.ValidatorCurrentRewards
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &rewards)
		addr := types.GetValidatorCurrentRewardsAddress(iter.Key())
		if handler(addr, rewards) {
			break
		}
	}
}

// get accumulated commission for a validator
func (k Keeper) GetValidatorAccumulatedCommission(ctx sdk.Context, val AccountID) (commission types.ValidatorAccumulatedCommission) {
	store := store.NewStore(ctx, k.storeKey)
	b := store.Get(types.GetValidatorAccumulatedCommissionKey(val))
	if b == nil {
		return types.ValidatorAccumulatedCommission{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &commission)
	return
}

// set accumulated commission for a validator
func (k Keeper) SetValidatorAccumulatedCommission(ctx sdk.Context, val AccountID, commission types.ValidatorAccumulatedCommission) {
	var bz []byte

	store := store.NewStore(ctx, k.storeKey)
	if commission.Commission.IsZero() {
		bz = k.cdc.MustMarshalBinaryLengthPrefixed(&types.ValidatorAccumulatedCommission{})
	} else {
		bz = k.cdc.MustMarshalBinaryLengthPrefixed(&commission)
	}

	store.Set(types.GetValidatorAccumulatedCommissionKey(val), bz)
}

// delete accumulated commission for a validator
func (k Keeper) DeleteValidatorAccumulatedCommission(ctx sdk.Context, val AccountID) {
	store := store.NewStore(ctx, k.storeKey)
	store.Delete(types.GetValidatorAccumulatedCommissionKey(val))
}

// iterate over accumulated commissions
func (k Keeper) IterateValidatorAccumulatedCommissions(ctx sdk.Context,
	handler func(val AccountID, commission types.ValidatorAccumulatedCommission) (stop bool)) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorAccumulatedCommissionPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var commission types.ValidatorAccumulatedCommission
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &commission)
		addr := types.GetValidatorAccumulatedCommissionAddress(iter.Key())
		if handler(addr, commission) {
			break
		}
	}
}

// get validator outstanding rewards
func (k Keeper) GetValidatorOutstandingRewards(ctx sdk.Context, val AccountID) (rewards types.ValidatorOutstandingRewards) {
	store := store.NewStore(ctx, k.storeKey)
	bz := store.Get(types.GetValidatorOutstandingRewardsKey(val))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rewards)
	return
}

// set validator outstanding rewards
func (k Keeper) SetValidatorOutstandingRewards(ctx sdk.Context, val AccountID, rewards types.ValidatorOutstandingRewards) {
	store := store.NewStore(ctx, k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&rewards)
	store.Set(types.GetValidatorOutstandingRewardsKey(val), b)
}

// delete validator outstanding rewards
func (k Keeper) DeleteValidatorOutstandingRewards(ctx sdk.Context, val AccountID) {
	store := store.NewStore(ctx, k.storeKey)
	store.Delete(types.GetValidatorOutstandingRewardsKey(val))
}

// iterate validator outstanding rewards
func (k Keeper) IterateValidatorOutstandingRewards(ctx sdk.Context, handler func(val AccountID, rewards types.ValidatorOutstandingRewards) (stop bool)) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorOutstandingRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		rewards := types.ValidatorOutstandingRewards{}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &rewards)
		addr := types.GetValidatorOutstandingRewardsAddress(iter.Key())
		if handler(addr, rewards) {
			break
		}
	}
}

// get slash event for height
func (k Keeper) GetValidatorSlashEvent(ctx sdk.Context, val AccountID, height, period uint64) (event types.ValidatorSlashEvent, found bool) {
	store := store.NewStore(ctx, k.storeKey)
	b := store.Get(types.GetValidatorSlashEventKey(val, height, period))
	if b == nil {
		return types.ValidatorSlashEvent{}, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &event)
	return event, true
}

// set slash event for height
func (k Keeper) SetValidatorSlashEvent(ctx sdk.Context, val AccountID, height, period uint64, event types.ValidatorSlashEvent) {
	store := store.NewStore(ctx, k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&event)
	store.Set(types.GetValidatorSlashEventKey(val, height, period), b)
}

// iterate over slash events between heights, inclusive
func (k Keeper) IterateValidatorSlashEventsBetween(ctx sdk.Context, val AccountID, startingHeight uint64, endingHeight uint64,
	handler func(height uint64, event types.ValidatorSlashEvent) (stop bool)) {
	store := store.NewStore(ctx, k.storeKey)
	iter := store.Iterator(
		types.GetValidatorSlashEventKeyPrefix(val, startingHeight),
		types.GetValidatorSlashEventKeyPrefix(val, endingHeight+1),
	)
	defer iter.Close()
	ctx.Logger().Debug("IterateValidatorSlashEventsBetween", "iter.Valid()", iter.Valid(), "val", val)

	for ; iter.Valid(); iter.Next() {
		var event types.ValidatorSlashEvent
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &event)
		_, height := types.GetValidatorSlashEventAddressHeight(iter.Key())
		if handler(height, event) {
			break
		}
	}
}

// iterate over all slash events
func (k Keeper) IterateValidatorSlashEvents(ctx sdk.Context, handler func(val AccountID, height uint64, event types.ValidatorSlashEvent) (stop bool)) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorSlashEventPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var event types.ValidatorSlashEvent
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &event)
		val, height := types.GetValidatorSlashEventAddressHeight(iter.Key())
		if handler(val, height, event) {
			break
		}
	}
}

// delete slash events for a particular validator
func (k Keeper) DeleteValidatorSlashEvents(ctx sdk.Context, val AccountID) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetValidatorSlashEventPrefix(val))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// delete all slash events
func (k Keeper) DeleteAllValidatorSlashEvents(ctx sdk.Context) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorSlashEventPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
