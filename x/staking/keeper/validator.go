package keeper

import (
	"bytes"
	"fmt"
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/KuChainNetwork/kuchain/chain/store"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Cache the amino decoding of validators, as it can be the case that repeated slashing calls
// cause many calls to GetValidator, which were shown to throttle the state machine in our
// simulation. Note this is quite biased though, as the simulator does more slashes than a
// live chain should, however we require the slashing to be fast as noone pays gas for it.
type cachedValidator struct {
	val       types.Validator
	marshaled string // marshaled amino bytes for the validator object (not operator address)
}

func newCachedValidator(val types.Validator, marshaled string) cachedValidator {
	return cachedValidator{
		val:       val,
		marshaled: marshaled,
	}
}

// get a single validator
func (k Keeper) GetValidator(ctx sdk.Context, acc types.AccountID) (validator types.Validator, found bool) {
	store := store.NewStore(ctx, k.storeKey)
	value := store.Get(types.GetValidatorKey(acc))
	if value == nil {
		return validator, false
	}

	// If these amino encoded bytes are in the cache, return the cached validator
	strValue := string(value)
	if val, ok := k.validatorCache[strValue]; ok {
		valToReturn := val.val
		// Doesn't mutate the cache's value
		valToReturn.OperatorAccount = acc
		return valToReturn, true
	}

	// amino bytes weren't found in cache, so amino unmarshal and add it to the cache
	validator = types.MustUnmarshalValidator(k.cdc, value)
	cachedVal := newCachedValidator(validator, strValue)
	k.validatorCache[strValue] = newCachedValidator(validator, strValue)
	k.validatorCacheList.PushBack(cachedVal)

	// if the cache is too big, pop off the last element from it
	if k.validatorCacheList.Len() > aminoCacheSize {
		valToRemove := k.validatorCacheList.Remove(k.validatorCacheList.Front()).(cachedValidator)
		delete(k.validatorCache, valToRemove.marshaled)
	}

	validator = types.MustUnmarshalValidator(k.cdc, value)
	return validator, true
}

func (k Keeper) mustGetValidator(ctx sdk.Context, acc types.AccountID) types.Validator {
	validator, found := k.GetValidator(ctx, acc)
	if !found {
		panic(fmt.Sprintf("validator record not found for address: %X\n", acc))
	}
	return validator
}

// get a single validator by consensus address
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAcc sdk.ConsAddress) (validator types.Validator, found bool) {
	store := store.NewStore(ctx, k.storeKey)
	opAccBytes := store.Get(types.GetValidatorByConsAddrKey(consAcc))
	if opAccBytes == nil {
		return validator, false
	}

	return k.GetValidator(ctx, types.NewAccountIDFromByte(opAccBytes))
}

func (k Keeper) mustGetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) types.Validator {
	validator, found := k.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		panic(fmt.Errorf("validator with consensus-Address %s not found", consAddr))
	}
	return validator
}

// set the main record holding validator details
func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) {
	store := store.NewStore(ctx, k.storeKey)
	bz := types.MustMarshalValidator(k.cdc, validator)
	store.Set(types.GetValidatorKey(validator.OperatorAccount), bz)
}

// validator index
func (k Keeper) SetValidatorByConsAddr(ctx sdk.Context, validator types.Validator) {
	store := store.NewStore(ctx, k.storeKey)
	store.Set(types.GetValidatorByConsAddrKey(validator.GetConsAccount()), validator.OperatorAccount.StoreKey())
}

// validator index
func (k Keeper) SetValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	// jailed validators are not kept in the power index
	if validator.Jailed {
		return
	}
	store := store.NewStore(ctx, k.storeKey)
	store.Set(types.GetValidatorsByPowerIndexKey(validator), validator.OperatorAccount.StoreKey())
}

// validator index
func (k Keeper) DeleteValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	store := store.NewStore(ctx, k.storeKey)
	store.Delete(types.GetValidatorsByPowerIndexKey(validator))
}

// validator index
func (k Keeper) HasValidatorByPowerIndex(ctx sdk.Context, power []byte) bool {
	store := store.NewStore(ctx, k.storeKey)
	return store.Has(power)
}

// validator index
func (k Keeper) SetNewValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	store := store.NewStore(ctx, k.storeKey)
	store.Set(types.GetValidatorsByPowerIndexKey(validator), validator.OperatorAccount.StoreKey())
}

// Update the tokens of an existing validator, update the validators power index key
func (k Keeper) AddValidatorTokensAndShares(ctx sdk.Context, validator types.Validator,
	tokensToAdd sdk.Int) (valOut types.Validator, addedShares sdk.Dec) {
	k.DeleteValidatorByPowerIndex(ctx, validator)
	validator, addedShares = validator.AddTokensFromDel(tokensToAdd)
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	return validator, addedShares
}

// Update the tokens of an existing validator, update the validators power index key
func (k Keeper) RemoveValidatorTokensAndShares(ctx sdk.Context, validator types.Validator,
	sharesToRemove sdk.Dec) (valOut types.Validator, removedTokens sdk.Int) {
	k.DeleteValidatorByPowerIndex(ctx, validator)
	validator, removedTokens = validator.RemoveDelShares(sharesToRemove)
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	return validator, removedTokens
}

// Update the tokens of an existing validator, update the validators power index key
func (k Keeper) RemoveValidatorTokens(ctx sdk.Context,
	validator types.Validator, tokensToRemove sdk.Int) types.Validator {
	k.DeleteValidatorByPowerIndex(ctx, validator)
	validator = validator.RemoveTokens(tokensToRemove)
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	return validator
}

// UpdateValidatorCommission attempts to update a validator's commission rate.
// An error is returned if the new commission rate is invalid.
func (k Keeper) UpdateValidatorCommission(ctx sdk.Context,
	validator types.Validator, newRate sdk.Dec) (types.Commission, error) {
	commission := validator.Commission
	blockTime := ctx.BlockHeader().Time

	if err := commission.ValidateNewRate(newRate, blockTime); err != nil {
		return commission, err
	}

	commission.Rate = newRate
	commission.UpdateTime = blockTime

	return commission, nil
}

// remove the validator record and associated indexes
// except for the bonded validator index which is only handled in ApplyAndReturnTendermintUpdates
func (k Keeper) RemoveValidator(ctx sdk.Context, address types.AccountID) {
	// first retrieve the old validator record
	validator, found := k.GetValidator(ctx, address)
	if !found {
		return
	}

	if !validator.IsUnbonded() {
		panic("cannot call RemoveValidator on bonded or unbonding validators")
	}
	if validator.Tokens.IsPositive() {
		panic("attempting to remove a validator which still contains tokens")
	}
	if validator.Tokens.IsPositive() {
		panic("validator being removed should never have positive tokens")
	}

	valConsAddr := validator.GetConsAccount()

	// delete the old validator record
	store := store.NewStore(ctx, k.storeKey)
	store.Delete(types.GetValidatorKey(address))
	store.Delete(types.GetValidatorByConsAddrKey(valConsAddr))
	store.Delete(types.GetValidatorsByPowerIndexKey(validator))

	// call hooks
	k.AfterValidatorRemoved(ctx, validator.GetConsAccount(), validator.OperatorAccount)
}

// get groups of validators

// get the set of all validators with no limits, used during genesis dump
func (k Keeper) GetAllValidators(ctx sdk.Context) (validators []types.Validator) {
	store := store.NewStore(ctx, k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}

	return validators
}

// return a given amount of all the validators
func (k Keeper) GetValidators(ctx sdk.Context, maxRetrieve uint32) (validators []types.Validator) {
	store := store.NewStore(ctx, k.storeKey)
	validators = make([]types.Validator, maxRetrieve)

	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators[i] = validator
		i++
	}

	return validators[:i] // trim if the array length < maxRetrieve
}

// get the current group of bonded validators sorted by power-rank
func (k Keeper) GetBondedValidatorsByPower(ctx sdk.Context) []types.Validator {
	maxValidators := k.MaxValidators(ctx)
	validators := make([]types.Validator, maxValidators)

	iterator := k.ValidatorsPowerStoreIterator(ctx)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxValidators); iterator.Next() {
		address := NewAccountIDFromByte(iterator.Value())
		validator := k.mustGetValidator(ctx, address)

		if validator.IsBonded() {
			validators[i] = validator
			i++
		}
	}
	return validators[:i] // trim
}

// returns an iterator for the current validator power store
func (k Keeper) ValidatorsPowerStoreIterator(ctx sdk.Context) sdk.Iterator {
	store := store.NewStore(ctx, k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.ValidatorsByPowerIndexKey)
}

// Last Validator Index

// Load the last validator power.
// Returns zero if the operator was not a validator last block.
func (k Keeper) GetLastValidatorPower(ctx sdk.Context, operator types.AccountID) (power int64) {
	store := store.NewStore(ctx, k.storeKey)
	bz := store.Get(types.GetLastValidatorPowerKey(operator))
	if bz == nil {
		return 0
	}

	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshalBinaryBare(bz, &intV)
	return intV.GetValue()
}

// Set the last validator power.
func (k Keeper) SetLastValidatorPower(ctx sdk.Context, operator types.AccountID, power int64) {
	store := store.NewStore(ctx, k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.Int64Value{Value: power})
	store.Set(types.GetLastValidatorPowerKey(operator), bz)
}

// Delete the last validator power.
func (k Keeper) DeleteLastValidatorPower(ctx sdk.Context, operator types.AccountID) {
	store := store.NewStore(ctx, k.storeKey)
	store.Delete(types.GetLastValidatorPowerKey(operator))
}

// returns an iterator for the consensus validators in the last block
func (k Keeper) LastValidatorsIterator(ctx sdk.Context) (iterator sdk.Iterator) {
	store := store.NewStore(ctx, k.storeKey)
	iterator = sdk.KVStorePrefixIterator(store, types.LastValidatorPowerKey)
	return iterator
}

// Iterate over last validator powers.
func (k Keeper) IterateLastValidatorPowers(ctx sdk.Context, handler func(operator types.AccountID, power int64) (stop bool)) {
	store := store.NewStore(ctx, k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.LastValidatorPowerKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		//	addr := sdk.ValAddress(iter.Key()[len(types.LastValidatorPowerKey):])
		addr := NewAccountIDFromByte(iter.Key()[len(types.LastValidatorPowerKey):])
		intV := &gogotypes.Int64Value{}

		k.cdc.MustUnmarshalBinaryBare(iter.Value(), intV)
		if handler(addr, intV.GetValue()) {
			break
		}
	}
}

// get the group of the bonded validators
func (k Keeper) GetLastValidators(ctx sdk.Context) (validators []types.Validator) {
	store := store.NewStore(ctx, k.storeKey)

	// add the actual validator power sorted store
	maxValidators := k.MaxValidators(ctx)
	validators = make([]types.Validator, maxValidators)

	iterator := sdk.KVStorePrefixIterator(store, types.LastValidatorPowerKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {
		// sanity check
		if i >= int(maxValidators) {
			panic("more validators than maxValidators found")
		}
		address := NewAccountIDFromByte(types.AddressFromLastValidatorPowerKey(iterator.Key()))
		validator := k.mustGetValidator(ctx, address)

		validators[i] = validator
		i++
	}
	return validators[:i] // trim
}

// Validator Queue

// gets a specific validator queue timeslice. A timeslice is a slice of ValAddresses corresponding to unbonding validators
// that expire at a certain time.
func (k Keeper) GetValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []types.AccountID {
	store := store.NewStore(ctx, k.storeKey)
	bz := store.Get(types.GetValidatorQueueTimeKey(timestamp))
	if bz == nil {
		return []types.AccountID{}
	}

	va := types.AccountIDs{}
	k.cdc.MustUnmarshalBinaryBare(bz, &va)
	return va.Addresses
}

// Sets a specific validator queue timeslice.
func (k Keeper) SetValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []types.AccountID) {
	store := store.NewStore(ctx, k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&types.AccountIDs{Addresses: keys})
	store.Set(types.GetValidatorQueueTimeKey(timestamp), bz)
}

// Deletes a specific validator queue timeslice.
func (k Keeper) DeleteValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time) {
	store := store.NewStore(ctx, k.storeKey)
	store.Delete(types.GetValidatorQueueTimeKey(timestamp))
}

// Insert an validator address to the appropriate timeslice in the validator queue
func (k Keeper) InsertValidatorQueue(ctx sdk.Context, val types.Validator) {
	timeSlice := k.GetValidatorQueueTimeSlice(ctx, val.UnbondingTime)
	timeSlice = append(timeSlice, val.OperatorAccount)
	k.SetValidatorQueueTimeSlice(ctx, val.UnbondingTime, timeSlice)
}

// Delete a validator address from the validator queue
func (k Keeper) DeleteValidatorQueue(ctx sdk.Context, val types.Validator) {
	timeSlice := k.GetValidatorQueueTimeSlice(ctx, val.UnbondingTime)
	newTimeSlice := []types.AccountID{}
	for _, addr := range timeSlice {
		if !bytes.Equal(addr.StoreKey(), val.OperatorAccount.StoreKey()) {
			newTimeSlice = append(newTimeSlice, addr)
		}
	}

	if len(newTimeSlice) == 0 {
		k.DeleteValidatorQueueTimeSlice(ctx, val.UnbondingTime)
	} else {
		k.SetValidatorQueueTimeSlice(ctx, val.UnbondingTime, newTimeSlice)
	}
}

// Returns all the validator queue timeslices from time 0 until endTime
func (k Keeper) ValidatorQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := store.NewStore(ctx, k.storeKey)
	return store.Iterator(types.ValidatorQueueKey, sdk.InclusiveEndBytes(types.GetValidatorQueueTimeKey(endTime)))
}

// Returns a concatenated list of all the timeslices before currTime, and deletes the timeslices from the queue
func (k Keeper) GetAllMatureValidatorQueue(ctx sdk.Context, currTime time.Time) (matureValsAddrs []types.AccountID) {
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	validatorTimesliceIterator := k.ValidatorQueueIterator(ctx, ctx.BlockHeader().Time)
	defer validatorTimesliceIterator.Close()

	for ; validatorTimesliceIterator.Valid(); validatorTimesliceIterator.Next() {
		timeslice := types.AccountIDs{}
		k.cdc.MustUnmarshalBinaryBare(validatorTimesliceIterator.Value(), &timeslice)

		matureValsAddrs = append(matureValsAddrs, timeslice.Addresses...)
	}

	return matureValsAddrs
}

// Unbonds all the unbonding validators that have finished their unbonding period
func (k Keeper) UnbondAllMatureValidatorQueue(ctx sdk.Context) {
	store := store.NewStore(ctx, k.storeKey)
	validatorTimesliceIterator := k.ValidatorQueueIterator(ctx, ctx.BlockHeader().Time)
	defer validatorTimesliceIterator.Close()

	for ; validatorTimesliceIterator.Valid(); validatorTimesliceIterator.Next() {
		timeslice := types.AccountIDs{}
		k.cdc.MustUnmarshalBinaryBare(validatorTimesliceIterator.Value(), &timeslice)

		for _, valAddr := range timeslice.Addresses {
			val, found := k.GetValidator(ctx, valAddr)
			if !found {
				panic("validator in the unbonding queue was not found")
			}

			if !val.IsUnbonding() {
				panic("unexpected validator in unbonding queue; status was not unbonding")
			}

			val = k.UnbondingToUnbonded(ctx, val)
			if val.GetDelegatorShares().IsZero() {
				k.RemoveValidator(ctx, val.OperatorAccount)
			}
		}

		store.Delete(validatorTimesliceIterator.Key())
	}
}
