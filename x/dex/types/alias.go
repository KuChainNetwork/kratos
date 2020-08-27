package types

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
)

type (
	AccountID = types.AccountID
	Name      = types.Name
	Coin      = types.Coin
	Coins     = types.Coins
	DecCoins  = types.DecCoins
)

var (
	CoinDenom             = types.CoinDenom
	CoinAccountsFromDenom = types.CoinAccountsFromDenom
	NewCoin               = types.NewCoin
	NewName               = types.NewName
	MustName              = types.MustName
	NewAccountIDFromName  = types.NewAccountIDFromName
	NewAccountIDFromStr   = types.NewAccountIDFromStr
	NewInt                = types.NewInt
)
