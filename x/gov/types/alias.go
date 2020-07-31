package types

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
)

type (
	AccountID     = types.AccountID
	Coin          = types.Coin
	Coins         = types.Coins
	Name          = types.Name
	AccountAuther = types.AccountAuther
	AssetTransfer = types.AssetTransfer
	KuMsg         = types.KuMsg
)

var (
	MustName             = types.MustName
	NewCoin              = types.NewCoin
	NewCoins             = types.NewCoins
	NewAccountIDFromByte = types.NewAccountIDFromByte
)
