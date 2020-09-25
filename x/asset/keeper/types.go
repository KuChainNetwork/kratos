package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
)

type (
	AssetTransfer = types.AssetTransfer
	Context       = types.Context
)

var (
	CoinAccountsFromDenom = types.CoinAccountsFromDenom
	CoinDenom             = types.CoinDenom
)

type (
	AccountID = types.AccountID
	Coins     = types.Coins
	Coin      = types.Coin
	DecCoins  = types.DecCoins
	DecCoin   = types.DecCoin
)

var (
	NewDec        = types.NewDec
	NewCoins      = types.NewCoins
	NewInt64Coin  = types.NewInt64Coin
	NewInt64Coins = types.NewInt64Coins
)

type ApproveData struct {
	IsLock             bool      `json:"lock" yaml:"lock"`
	TimeOutBlockHeight int64     `json:"out_height" yaml:"out_height"`
	Controller         AccountID `json:"c" yaml:"c"`
	Amount             Coins     `json:"amount" yaml:"amount"`
}

func NewApproveData(amt Coins) *ApproveData {
	return &ApproveData{
		Amount: amt,
	}
}
