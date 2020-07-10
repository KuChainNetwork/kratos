package types

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
)

// query endpoints supported by the auth Querier
const (
	QueryCoin            = "coin"
	QueryCoins           = "coins"
	QueryCoinPower       = "coinpower"
	QueryCoinPowers      = "coinpowers"
	QueryCoinStat        = "coinstate"
	QueryCoinDescription = "coindesc"
	QueryCoinLocked      = "coinslocked"
)

// QueryCoinParams defines the params for querying coin.
type QueryCoinParams struct {
	AccountID types.AccountID
	Creator   types.Name
	Symbol    types.Name
}

// NewQueryCoinParams creates a new instance of QueryCoinParams.
func NewQueryCoinParams(accountID types.AccountID, creator, symbol types.Name) QueryCoinParams {
	return QueryCoinParams{
		AccountID: accountID,
		Creator:   creator,
		Symbol:    symbol,
	}
}

// QueryCoinsParams defines the params for querying coin.
type QueryCoinsParams struct {
	AccountID types.AccountID
}

// NewQueryCoinsParams creates a new instance of QueryCoinParams.
func NewQueryCoinsParams(accountID types.AccountID) QueryCoinsParams {
	return QueryCoinsParams{
		AccountID: accountID,
	}
}

// QueryCoinPowerParams defines the params for querying coin.
type QueryCoinPowerParams struct {
	AccountID types.AccountID
	Creator   types.Name
	Symbol    types.Name
}

// NewQueryCoinPowerParams creates a new instance of QueryCoinParams.
func NewQueryCoinPowerParams(accountID types.AccountID, creator, symbol types.Name) QueryCoinPowerParams {
	return QueryCoinPowerParams{
		AccountID: accountID,
		Creator:   creator,
		Symbol:    symbol,
	}
}

// QueryCoinPowersParams defines the params for querying coin.
type QueryCoinPowersParams struct {
	AccountID types.AccountID
}

// NewQueryCoinPowersParams creates a new instance of QueryCoinParams.
func NewQueryCoinPowersParams(accountID types.AccountID) QueryCoinPowersParams {
	return QueryCoinPowersParams{
		AccountID: accountID,
	}
}

// QueryCoinStatParams defines the params for querying coin stat.
type QueryCoinStatParams struct {
	Creator types.Name
	Symbol  types.Name
}

// NewQueryCoinStatParams creates a new instance of QueryCoinStatParams.
func NewQueryCoinStatParams(creator, symbol types.Name) QueryCoinStatParams {
	return QueryCoinStatParams{
		Creator: creator,
		Symbol:  symbol,
	}
}

// QueryCoinDescParams defines the params for querying coin desc.
type QueryCoinDescParams struct {
	Creator types.Name
	Symbol  types.Name
}

// NewQueryCoinDescParams creates a new instance of QueryCoinDescParams.
func NewQueryCoinDescParams(creator, symbol types.Name) QueryCoinDescParams {
	return QueryCoinDescParams{
		Creator: creator,
		Symbol:  symbol,
	}
}

// QueryLockedCoinsParams defines the params for querying coin.
type QueryLockedCoinsParams struct {
	AccountID types.AccountID
}

// NewQueryLockedCoinsParamscreates a new instance of QueryCoinParams.
func NewQueryLockedCoinsParams(accountID types.AccountID) QueryLockedCoinsParams {
	return QueryLockedCoinsParams{
		AccountID: accountID,
	}
}

type LockedCoins struct {
	Coins             types.Coins `json:"coins" yaml:"coins"`
	UnlockBlockHeight int64       `json:"unlock_block_height" yaml:"unlock_block_height"`
}

type QueryLockedCoinsResponse struct {
	LockedCoins types.Coins   `json:"coins"`
	Locks       []LockedCoins `json:"locks"`
}
