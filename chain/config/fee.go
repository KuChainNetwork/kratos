package config

import (
	"github.com/KuChainNetwork/kuchain/chain/types/coin"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	feePriceMiniLimit coin.DecCoins
)

func SetFeePriceMiniLimit(f coin.DecCoins) {
	feePriceMiniLimit = f
}

// GetFeeMiniLimit get fee gas price mini limit
func GetFeePriceMiniLimit() sdk.DecCoins {
	return feePriceMiniLimit.ToSDK()
}
