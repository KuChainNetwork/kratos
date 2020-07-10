package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Unjail calls the staking Unjail function to unjail a validator if the
// jailed period has concluded
func (k Keeper) Unjail(ctx sdk.Context, valAccountID types.AccountID) error {
	validator := k.sk.Validator(ctx, valAccountID)
	if validator == nil {
		return types.ErrNoValidatorForAddress
	}

	// cannot be unjailed if not jailed
	if !validator.IsJailed() {
		return types.ErrValidatorNotJailed
	}

	consAddr := sdk.ConsAddress(validator.GetConsPubKey().Address())

	info, found := k.GetValidatorSigningInfo(ctx, consAddr)
	if !found {
		return types.ErrNoValidatorForAddress
	}

	// cannot be unjailed if tombstoned
	if info.Tombstoned {
		return types.ErrValidatorJailed
	}

	// cannot be unjailed until out of jail
	if ctx.BlockHeader().Time.Before(info.JailedUntil) {
		return types.ErrValidatorJailed
	}

	k.sk.Unjail(ctx, consAddr)
	return nil
}
