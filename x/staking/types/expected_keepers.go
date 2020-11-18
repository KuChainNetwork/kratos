package types

import (
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	accountExported "github.com/KuChainNetwork/kuchain/x/account/exported"
	stakingexported "github.com/KuChainNetwork/kuchain/x/staking/exported"
	"github.com/KuChainNetwork/kuchain/x/staking/external"
	supplyexported "github.com/KuChainNetwork/kuchain/x/supply/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DistributionKeeper expected distribution keeper (noalias)
type DistributionKeeper interface {
	GetFeePoolCommunityCoins(ctx sdk.Context) chainTypes.DecCoins
	GetValidatorOutstandingRewardsCoins(ctx sdk.Context, val chainTypes.AccountID) chainTypes.DecCoins
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	chainTypes.AccountAuther
	IterateAccounts(ctx sdk.Context, process func(accountExported.Account) (stop bool))
	GetAccount(sdk.Context, chainTypes.AccountID) accountExported.Account // only used for simulation
}

// AccountStatKeeper is interface for other modules to get account state.
type AccountStatKeeper interface {
	GetAccount(sdk.Context, chainTypes.AccountID) external.Account // can return nil.
	IterateAccounts(ctx sdk.Context, cb func(account external.Account) (stop bool))

	GetNextAccountNumber(ctx sdk.Context) uint64
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	chainTypes.AssetTransfer
	GetBalance(ctx sdk.Context, addr chainTypes.AccountID, denom string) Coin
	GetCoinPowers(ctx sdk.Context, account chainTypes.AccountID) Coins
	GetCoinPowerByDenomd(ctx sdk.Context, account chainTypes.AccountID, denomd string) Coin
	SpendableCoins(ctx sdk.Context, addr chainTypes.AccountID) Coins
}

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	GetSupply(ctx sdk.Context) supplyexported.SupplyI

	InitModuleAccount(ctx sdk.Context, moduleName string) error
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) supplyexported.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)

	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt Coins) error
	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr AccountID, amt Coins) error
	DelegateCoinsFromAccountToModule(ctx sdk.Context, recipientModule string, amt Coins) error

	BurnCoins(ctx sdk.Context, name chainTypes.AccountID, amt Coins) error
}

// ValidatorSet expected properties for the set of all validators (noalias)
type ValidatorSet interface {
	// iterate through validators by operator address, execute func for each validator
	IterateValidators(sdk.Context,
		func(index int64, validator stakingexported.ValidatorI) (stop bool))

	// iterate through bonded validators by operator address, execute func for each validator
	IterateBondedValidatorsByPower(sdk.Context,
		func(index int64, validator stakingexported.ValidatorI) (stop bool))

	// iterate through the consensus validator set of the last block by operator address, execute func for each validator
	IterateLastValidators(sdk.Context,
		func(index int64, validator stakingexported.ValidatorI) (stop bool))

	Validator(sdk.Context, AccountID) stakingexported.ValidatorI                 // get a particular validator by operator address
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingexported.ValidatorI // get a particular validator by consensus address
	TotalBondedTokens(sdk.Context) sdk.Int                                       // total bonded tokens within the validator set
	StakingTokenSupply(sdk.Context) sdk.Int                                      // total staking token supply

	// slash the validator and delegators of the validator, specifying offence height, offence power, and slash fraction
	Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec)
	Jail(sdk.Context, sdk.ConsAddress)   // jail a validator
	Unjail(sdk.Context, sdk.ConsAddress) // unjail a validator

	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, AccountID, AccountID) stakingexported.DelegationI

	// MaxValidators returns the maximum amount of bonded validators
	MaxValidators(sdk.Context) uint32
}

// DelegationSet expected properties for the set of all delegations for a particular (noalias)
type DelegationSet interface {
	GetValidatorSet() ValidatorSet // validator set for which delegation set is based upon

	// iterate through all delegations from one delegator by validator-AccAddress,
	//   execute func for each validator
	IterateDelegations(ctx sdk.Context, delegator AccountID,
		fn func(index int64, delegation stakingexported.DelegationI) (stop bool))
}

// Event Hooks
// These can be utilized to communicate between a staking keeper and another
// keeper which must take particular actions when validators/delegators change
// state. The second keeper must implement this interface, which then the
// staking keeper can call.

// StakingHooks event hooks for staking validator object (noalias)
type StakingHooks interface {
	AfterValidatorCreated(ctx sdk.Context, valAddr AccountID)                           // Must be called when a validator is created
	BeforeValidatorModified(ctx sdk.Context, valAddr AccountID)                         // Must be called when a validator's state changes
	AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr AccountID) // Must be called when a validator is deleted

	AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr AccountID)         // Must be called when a validator is bonded
	AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr AccountID) // Must be called when a validator begins unbonding

	BeforeDelegationCreated(ctx sdk.Context, delAddr AccountID, valAddr AccountID)        // Must be called when a delegation is created
	BeforeDelegationSharesModified(ctx sdk.Context, delAddr AccountID, valAddr AccountID) // Must be called when a delegation's shares are modified
	BeforeDelegationRemoved(ctx sdk.Context, delAddr AccountID, valAddr AccountID)        // Must be called when a delegation is removed
	AfterDelegationModified(ctx sdk.Context, delAddr AccountID, valAddr AccountID)
	BeforeValidatorSlashed(ctx sdk.Context, valAddr AccountID, fraction sdk.Dec)
}
