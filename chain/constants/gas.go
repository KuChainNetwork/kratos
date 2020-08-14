package constants

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
)

var (
	MinGasPriceString           = "0.01" + DefaultBondDenom
	GasTxSizePrice       uint64 = 5
	EstimatedGasConsumed uint64 = 24000
)

var (
	MinGasPrice types.DecCoins
)

func init() {
	if minGasPrice, err := types.ParseDecCoins(MinGasPriceString); err != nil {
		panic(err)
	} else {
		MinGasPrice = minGasPrice
	}
}
