// nolint
package keeper

import (
	"time"

	"github.com/KuChainNetwork/kuchain/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// FIXME: use AccountID type

func (k Keeper) AfterValidatorBonded(ctx sdk.Context, acc sdk.ConsAddress, _ types.AccountID) {
	// Update the signing info start height or create a new signing info
	_, found := k.GetValidatorSigningInfo(ctx, acc)
	if !found {
		signingInfo := types.NewValidatorSigningInfo(
			acc,
			ctx.BlockHeight(),
			0,
			time.Unix(0, 0),
			false,
			0,
		)
		k.SetValidatorSigningInfo(ctx, acc, signingInfo)
	}
}

// When a validator is created, add the address-pubkey relation.
func (k Keeper) AfterValidatorCreated(ctx sdk.Context, valAddr types.AccountID) {
	validator := k.sk.Validator(ctx, valAddr)
	k.AddPubkey(ctx, validator.GetConsPubKey())
}

// When a validator is removed, delete the address-pubkey relation.
func (k Keeper) AfterValidatorRemoved(ctx sdk.Context, address sdk.ConsAddress) {
	k.deleteAddrPubkeyRelation(ctx, crypto.Address(address))
}

//_________________________________________________________________________________________

// Hooks wrapper struct for slashing keeper
type Hooks struct {
	k Keeper
}

var _ types.StakingHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr types.AccountID) {
	h.k.AfterValidatorBonded(ctx, consAddr, valAddr)
}

// Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, _ types.AccountID) {
	h.k.AfterValidatorRemoved(ctx, consAddr)
}

// Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr types.AccountID) {
	h.k.AfterValidatorCreated(ctx, valAddr)
}

// nolint - unused hooks
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ types.AccountID)   {}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ types.AccountID)                           {}
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ types.AccountID, _ types.AccountID)        {}
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ types.AccountID, _ types.AccountID) {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ types.AccountID, _ types.AccountID)        {}
func (h Hooks) AfterDelegationModified(_ sdk.Context, _ types.AccountID, _ types.AccountID)        {}
func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ types.AccountID, _ sdk.Dec)                 {}
