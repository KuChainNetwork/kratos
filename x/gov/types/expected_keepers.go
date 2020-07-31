package types

import (
	"time"

	accountExported "github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/KuChainNetwork/kuchain/x/gov/external"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ParamSubspace defines the expected Subspace interface for parameters (noalias)
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
}

// SupplyKeeper defines the expected supply keeper for module accounts (noalias)
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) external.SupplyModuleAccount

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, external.SupplyModuleAccount)

	ModuleCoinsToPower(ctx sdk.Context, recipientModule string, amt Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr AccountID, amt Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr AccountID, recipientModule string, amt Coins) error
	BurnCoins(ctx sdk.Context, name AccountID, amt Coins) error
}

// StakingKeeper expected staking keeper (Validator and Delegator sets) (noalias)
type StakingKeeper interface {
	// iterate through bonded validators by operator address, execute func for each validator
	IterateBondedValidatorsByPower(
		sdk.Context, func(index int64, validator external.StakingValidatorI) (stop bool),
	)

	TotalBondedTokens(sdk.Context) sdk.Int // total bonded tokens within the validator set
	IterateDelegations(
		ctx sdk.Context, delegator AccountID,
		fn func(index int64, delegation external.StakingDelegationI) (stop bool),
	)

	Validator(sdk.Context, AccountID) external.StakingValidatorI
	JailByAccount(ctx sdk.Context, account AccountID)
	UnjailByAccount(ctx sdk.Context, account AccountID)
	SlashByValidatorAccount(ctx sdk.Context, valAccount AccountID, infractionHeight int64, slashFactor sdk.Dec)
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	AccountAuther
	GetAccount(sdk.Context, AccountID) accountExported.Account
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	AssetTransfer
	GetAllBalances(ctx sdk.Context, addr AccountID) Coins
	GetCoinPowers(ctx sdk.Context, account AccountID) Coins
	SpendableCoins(ctx sdk.Context, account AccountID) Coins
}

type DistributionKeeper interface {
	CanDistribution(ctx sdk.Context) (bool, time.Time)

	SetStartNotDistributionTimePoint(ctx sdk.Context, t time.Time)
}
