// noalias
// DONTCOVER
package types

import (
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/KuChainNetwork/kuchain/x/slashing/external"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper expected account keeper
type AccountKeeper interface {
	chainTypes.AccountAuther
	GetAccount(sdk.Context, chainTypes.AccountID) exported.Account
	IterateAccounts(ctx sdk.Context, process func(exported.Account) (stop bool))
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	chainTypes.AssetTransfer

	SpendableCoins(ctx sdk.Context, addr chainTypes.AccountID) chainTypes.Coins
}

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	WithKeyTable(table external.ParamsKeyTable) external.ParamsSubspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps external.ParamSet)
	SetParamSet(ctx sdk.Context, ps external.ParamSet)
}

// StakingKeeper expected staking keeper
type StakingKeeper interface {
	// iterate through validators by operator address, execute func for each validator
	IterateValidators(sdk.Context,
		func(index int64, validator external.StakingValidatorl) (stop bool))

	Validator(sdk.Context, AccountID) external.StakingValidatorl                 // get a particular validator by operator address
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) external.StakingValidatorl // get a particular validator by consensus address

	// Slash slash the validator and delegators of the validator, specifying offense height, offense power, and slash fraction
	Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec)
	Jail(sdk.Context, sdk.ConsAddress)   // jail a validator
	Unjail(sdk.Context, sdk.ConsAddress) // unjail a validator

	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, AccountID, AccountID) external.StakingDelegatel

	// MaxValidators returns the maximum amount of bonded validators
	MaxValidators(sdk.Context) uint32
	GetAllValidatorInterfaces(ctx sdk.Context) []external.StakingValidatorl
}

// StakingHooks event hooks for staking validator object (noalias)
type StakingHooks interface {
	AfterValidatorCreated(ctx sdk.Context, valAddr AccountID)                           // Must be called when a validator is created
	AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr AccountID) // Must be called when a validator is deleted
	AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr AccountID)  // Must be called when a validator is bonded
}
