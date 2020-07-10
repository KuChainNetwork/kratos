package types

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	NewAccount(ctx sdk.Context, acc Account) Account

	GetAccount(sdk.Context, types.AccountID) Account // can return nil.
	SetAccount(ctx sdk.Context, acc exported.Account)
	IterateAccounts(ctx sdk.Context, cb func(account Account) (stop bool))
}

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	SendCoinPower(ctx sdk.Context, from, to types.AccountID, amt sdk.Coins) error
	IssueCoinPower(ctx sdk.Context, id types.AccountID, amt sdk.Coins) (sdk.Coins, error)
	BurnCoinPower(ctx sdk.Context, id types.AccountID, amt sdk.Coins) (sdk.Coins, error)
	CoinsToPower(ctx sdk.Context, from, to types.AccountID, amt sdk.Coins) error

	GetCoinsTotalSupply(ctx sdk.Context) types.Coins
	GetCoinTotalSupply(ctx sdk.Context, creator, symbol types.Name) types.Coin

	IterateAllCoins(ctx sdk.Context, cb func(address types.AccountID, balance sdk.Coins) (stop bool))
	IterateAllCoinPowers(ctx sdk.Context, cb func(address types.AccountID, balance sdk.Coins) (stop bool))
}
