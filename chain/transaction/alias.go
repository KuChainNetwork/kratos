package transaction

import "github.com/KuChainNetwork/kuchain/chain/types"

type (
	StdTx        = types.StdTx
	StdSignMsg   = types.StdSignMsg
	StdSignature = types.StdSignature
)

var (
	NewStdTx  = types.NewStdTx
	NewStdFee = types.NewStdFee
)

type (
	Coins    = types.Coins
	Coin     = types.Coin
	DecCoins = types.DecCoins
	DecCoin  = types.DecCoin
)

var (
	NewDec  = types.NewDec
	NewCoin = types.NewCoin
)
