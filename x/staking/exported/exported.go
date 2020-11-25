package exported

import (
	"time"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegationI delegation bond for a delegated proof of stake system
type DelegationI interface {
	// FIXME: delete Addr
	GetDelegatorAccountID() types.AccountID // delegator sdk.AccAddress for the bond
	GetValidatorAccountID() types.AccountID // validator operator address
	GetShares() sdk.Dec                     // amount of validator's shares held in this delegation
}

// ValidatorI expected validator functions
type ValidatorI interface {
	IsJailed() bool                                         // whether the validator is jailed
	GetMoniker() string                                     // moniker of the validator
	GetStatus() BondStatus                                  // status of the validator
	IsBonded() bool                                         // check if has a bonded status
	IsUnbonded() bool                                       // check if has status unbonded
	IsUnbonding() bool                                      // check if has status unbonding
	InvalidExRate() bool                                    // check if invalid ex rate
	GetOperator() sdk.ValAddress                            // operator address to receive/return validators coins
	GetOperatorAccountID() types.AccountID                  // operator account to receive/return validators coins
	GetConsPubKey() crypto.PubKey                           // validation consensus pubkey
	GetConsAddr() sdk.ConsAddress                           // validation consensus address
	GetTokens() sdk.Int                                     // validation tokens
	GetBondedTokens() sdk.Int                               // validator bonded tokens
	GetConsensusPower() int64                               // validation power in tendermint
	GetCommission() sdk.Dec                                 // validator commission rate
	GetCommissionMaxRate() sdk.Dec                          // validator commission max rate
	GetMinSelfDelegation() sdk.Int                          // validator minimum self delegation
	GetDelegatorShares() sdk.Dec                            // total outstanding delegator shares
	TokensFromShares(sdk.Dec) sdk.Dec                       // token worth of provided delegator shares
	TokensFromSharesTruncated(sdk.Dec) sdk.Dec              // token worth of provided delegator shares, truncated
	TokensFromSharesRoundUp(sdk.Dec) sdk.Dec                // token worth of provided delegator shares, rounded up
	SharesFromTokens(amt sdk.Int) (sdk.Dec, error)          // shares worth of delegator's bond
	SharesFromTokensTruncated(amt sdk.Int) (sdk.Dec, error) // truncated shares worth of delegator's bond
	CommissionValidateNewRate(sdk.Dec, time.Time) error     // validate new rate
}

// Event Hooks
// These can be utilized to communicate between a staking keeper and another
// keeper which must take particular actions when validators/delegators change
// state. The second keeper must implement this interface, which then the
// staking keeper can call.

// StakingHooks event hooks for staking validator object (noalias)
type StakingHooks interface {
	AfterValidatorCreated(ctx sdk.Context, valAddr types.AccountID)                           // Must be called when a validator is created
	BeforeValidatorModified(ctx sdk.Context, valAddr types.AccountID)                         // Must be called when a validator's state changes
	AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr types.AccountID) // Must be called when a validator is deleted

	AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr types.AccountID)         // Must be called when a validator is bonded
	AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr types.AccountID) // Must be called when a validator begins unbonding

	BeforeDelegationCreated(ctx sdk.Context, delAddr types.AccountID, valAddr types.AccountID)        // Must be called when a delegation is created
	BeforeDelegationSharesModified(ctx sdk.Context, delAddr types.AccountID, valAddr types.AccountID) // Must be called when a delegation's shares are modified
	BeforeDelegationRemoved(ctx sdk.Context, delAddr types.AccountID, valAddr types.AccountID)        // Must be called when a delegation is removed
	AfterDelegationModified(ctx sdk.Context, delAddr types.AccountID, valAddr types.AccountID)
	BeforeValidatorSlashed(ctx sdk.Context, valAddr types.AccountID, fraction sdk.Dec)
}
