package types

import (
	"time"

	"github.com/KuChainNetwork/kuchain/chain/types"
	accExported "github.com/KuChainNetwork/kuchain/x/account/exported"
	supplyexported "github.com/KuChainNetwork/kuchain/x/supply/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias) by cancer
type AccountKeeperAccountID interface {
	types.AccountAuther
	GetAccount(ctx sdk.Context, Id AccountID) accExported.Account
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeperAccountID interface {
	types.AssetTransfer
	GetAllBalances(ctx sdk.Context, Id AccountID) Coins
	GetCoinPowers(ctx sdk.Context, Id AccountID) Coins
	SpendableCoins(ctx sdk.Context, Id AccountID) Coins
	CoinsToPower(ctx sdk.Context, from, to AccountID, amt Coins) error
}

// StakingKeeper expected staking keeper (noalias) by cancer
type StakingKeeperAccountID interface {
	// iterate through validators by operator address, execute func for each validator
	IterateValidators(sdk.Context,
		func(index int64, validator StakingExportedValidatorI) (stop bool))

	// iterate through bonded validators by operator address, execute func for each validator
	IterateBondedValidatorsByPower(sdk.Context,
		func(index int64, validator StakingExportedValidatorI) (stop bool))

	// iterate through the consensus validator set of the last block by operator address, execute func for each validator
	IterateLastValidators(sdk.Context,
		func(index int64, validator StakingExportedValidatorI) (stop bool))

	Validator(sdk.Context, AccountID) StakingExportedValidatorI                 // get a particular validator by operator address
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) StakingExportedValidatorI // get a particular validator by consensus address

	// slash the validator and delegators of the validator, specifying offence height, offence power, and slash fraction
	Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec)
	Jail(sdk.Context, sdk.ConsAddress)   // jail a validator
	Unjail(sdk.Context, sdk.ConsAddress) // unjail a validator

	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, AccountID, AccountID) StakingExportedDelegationI

	// MaxValidators returns the maximum amount of bonded validators
	MaxValidators(sdk.Context) uint32

	IterateDelegations(ctx sdk.Context, delegatorId AccountID,
		fn func(index int64, delegation StakingExportedDelegationI) (stop bool))

	GetLastTotalPower(ctx sdk.Context) sdk.Int
	GetLastValidatorPower(ctx sdk.Context, valId AccountID) int64

	GetAllSDKDelegations(ctx sdk.Context) []StakingDelegation
}

// StakingHooks event hooks for staking validator object (noalias) by cancer
type StakingHooksAccountID interface {
	AfterValidatorCreated(ctx sdk.Context, valId AccountID)                         // Must be called when a validator is created
	AfterValidatorRemoved(ctx sdk.Context, consId sdk.ConsAddress, valId AccountID) // Must be called when a validator is deleted

	BeforeDelegationCreated(ctx sdk.Context, delId AccountID, valId AccountID)        // Must be called when a delegation is created
	BeforeDelegationSharesModified(ctx sdk.Context, delId AccountID, valId AccountID) // Must be called when a delegation's shares are modified
	AfterDelegationModified(ctx sdk.Context, delId AccountID, valId AccountID)
	BeforeValidatorSlashed(ctx sdk.Context, valId AccountID, fraction sdk.Dec)
}

// SupplyKeeper defines the expected supply Keeper (noalias) by cancer
type SupplyKeeperAccountID interface {
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)

	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientId AccountID, amt Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderId AccountID, recipientModule string, amt Coins) error
}

type DistributionKeeper interface {
	CanDistribution(ctx sdk.Context) (bool, time.Time)

	SetStartNotDistributionTimePoint(ctx sdk.Context, t time.Time)
}
