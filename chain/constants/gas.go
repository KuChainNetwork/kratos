package constants

import "github.com/KuChain-io/kuchain/chain/types"

const (
	MinGasPriceString        = "0.0001kratos/kts"
	GasTxSizePrice    uint64 = 5
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
