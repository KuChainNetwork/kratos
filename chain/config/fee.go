package config

import "github.com/KuChainNetwork/kuchain/chain/types"

var (
	feePriceMiniLimit types.Coins
)

func SetFeePriceMiniLimit(f types.Coins) {
	feePriceMiniLimit = f
}

// GetFeeMiniLimit get fee gas price mini limit
func GetFeePriceMiniLimit() types.Coins {
	return feePriceMiniLimit
}
