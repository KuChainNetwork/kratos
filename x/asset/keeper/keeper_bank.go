package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BlacklistedAddr checks if a given address is blacklisted (i.e restricted from
// receiving funds)
func (a AssetKeeper) BlacklistedAddr(addr sdk.AccAddress) bool {
	// TODO: node black list
	return false
}

// GetAllBalances get all coins for a account
func (a AssetKeeper) GetAllBalances(ctx sdk.Context, ID types.AccountID) Coins {
	res, _ := a.GetCoins(ctx, ID)
	return res
}

// GetBalance get coins balance for account
func (a AssetKeeper) GetBalance(ctx sdk.Context, ID types.AccountID, denom string) Coin {
	creator, symbol, err := types.CoinAccountsFromDenom(denom)
	if err != nil {
		panic(err)
	}

	res, _ := a.GetCoin(ctx, ID, creator, symbol)
	return res
}
