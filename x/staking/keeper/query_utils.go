package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/store"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Return all validators that a delegator is bonded to. If maxRetrieve is supplied, the respective amount will be returned.
func (k Keeper) GetDelegatorValidators(
	ctx sdk.Context, delegatorAddr AccountID, maxRetrieve uint32,
) []types.Validator {

	validators := make([]types.Validator, maxRetrieve)

	store := store.NewStore(ctx, k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delegatorAddr)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())

		validator, found := k.GetValidator(ctx, delegation.ValidatorAccount)
		if !found {
			panic(types.ErrNoValidatorFound)
		}

		validators[i] = validator
		i++
	}

	return validators[:i] // trim
}

// return a validator that a delegator is bonded to
func (k Keeper) GetDelegatorValidator(
	ctx sdk.Context, delegatorAddr AccountID, validatorAddr AccountID,
) (validator types.Validator, err error) {

	delegation, found := k.GetDelegation(ctx, delegatorAddr, validatorAddr)
	if !found {
		return validator, types.ErrNoDelegation
	}

	validator, found = k.GetValidator(ctx, delegation.ValidatorAccount)
	if !found {
		panic(types.ErrNoValidatorFound)
	}

	return validator, nil
}

//_____________________________________________________________________________________

// return all delegations for a delegator
func (k Keeper) GetAllDelegatorDelegations(ctx sdk.Context, delegator AccountID) []types.Delegation {
	delegations := make([]types.Delegation, 0)

	store := store.NewStore(ctx, k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) //smallest to largest
	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
		i++
	}

	return delegations
}

// return all unbonding-delegations for a delegator
func (k Keeper) GetAllUnbondingDelegations(ctx sdk.Context, delegator AccountID) []types.UnbondingDelegation {
	unbondingDelegations := make([]types.UnbondingDelegation, 0)

	store := store.NewStore(ctx, k.storeKey)
	delegatorPrefixKey := types.GetUBDsKey(delegator.StoreKey())
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	for i := 0; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
		i++
	}

	return unbondingDelegations
}

// return all redelegations for a delegator
func (k Keeper) GetAllRedelegations(
	ctx sdk.Context, delegator AccountID, srcValAddress, dstValAddress AccountID,
) []types.Redelegation {

	store := store.NewStore(ctx, k.storeKey)
	delegatorPrefixKey := types.GetREDsKey(delegator.StoreKey())
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	srcValFilter := !(srcValAddress.Empty())
	dstValFilter := !(dstValAddress.Empty())

	redelegations := []types.Redelegation{}

	for ; iterator.Valid(); iterator.Next() {
		redelegation := types.MustUnmarshalRED(k.cdc, iterator.Value())
		if srcValFilter && !(srcValAddress.Eq(redelegation.ValidatorSrcAccount)) {
			continue
		}
		if dstValFilter && !(dstValAddress.Eq(redelegation.ValidatorDstAccount)) {
			continue
		}

		redelegations = append(redelegations, redelegation)
	}

	return redelegations
}
