package types

import (
	chaintype "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/gov/external"
	sdk "github.com/cosmos/cosmos-sdk/types"

	accountExported "github.com/KuChainNetwork/kuchain/x/account/exported"
	"time"
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

	ModuleCoinsToPower(ctx sdk.Context, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr chaintype.AccountID, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr chaintype.AccountID, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name chaintype.AccountID, amt sdk.Coins) error
}

// StakingKeeper expected staking keeper (Validator and Delegator sets) (noalias)
type StakingKeeper interface {
	// iterate through bonded validators by operator address, execute func for each validator
	IterateBondedValidatorsByPower(
		sdk.Context, func(index int64, validator external.StakingValidatorI) (stop bool),
	)

	TotalBondedTokens(sdk.Context) sdk.Int // total bonded tokens within the validator set
	IterateDelegations(
		ctx sdk.Context, delegator chaintype.AccountID,
		fn func(index int64, delegation external.StakingDelegationI) (stop bool),
	)

	Validator(sdk.Context, chaintype.AccountID) external.StakingValidatorI
	JailByAccount(ctx sdk.Context, account chaintype.AccountID)
	UnjailByAccount(ctx sdk.Context, account chaintype.AccountID)
	SlashByValidatorAccount(ctx sdk.Context, valAccount chaintype.AccountID, infractionHeight int64, slashFactor sdk.Dec)
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	chaintype.AccountAuther
	GetAccount(sdk.Context, chaintype.AccountID) accountExported.Account
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	chaintype.AssetTransfer
	GetAllBalances(ctx sdk.Context, addr chaintype.AccountID) sdk.Coins
	GetCoinPowers(ctx sdk.Context, account chaintype.AccountID) sdk.Coins
	SpendableCoins(ctx sdk.Context, account chaintype.AccountID) sdk.Coins
}

type DistributionKeeper interface {
	CanDistribution(ctx sdk.Context) (bool, time.Time)

	SetStartNotDistributionTimePoint(ctx sdk.Context, t time.Time)
}