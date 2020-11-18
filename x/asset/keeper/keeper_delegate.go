package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SpendableCoins return account spendable coins
func (a AssetKeeper) SpendableCoins(ctx sdk.Context, id types.AccountID) Coins {
	res, err := a.getCoins(ctx, id)
	if err != nil {
		return Coins{}
	}

	lockeds, err := a.getCoinsLocked(ctx, id)
	if err != nil || lockeds == nil {
		return res
	}

	spendable, isNegative := res.SafeSub(lockeds)
	if isNegative {
		return Coins{}
	}

	return spendable
}
