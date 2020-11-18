package keeper

import (
	"errors"

	"github.com/KuChainNetwork/kuchain/chain/store"
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetCoins get coins by account id
func (a AssetKeeper) GetCoins(ctx sdk.Context, account types.AccountID) (types.Coins, error) {
	return a.getCoins(ctx, account)
}

// GetCoin get coin by account id and coin demon
func (a AssetKeeper) GetCoin(ctx sdk.Context, account types.AccountID, creator, symbol types.Name) (types.Coin, error) {
	coins, err := a.getCoins(ctx, account)
	if err != nil {
		return types.Coin{}, sdkerrors.Wrapf(err, "get coin %s %s-%s", account, creator, symbol)
	}
	denomd := types.CoinDenom(creator, symbol)

	return types.NewCoin(denomd, coins.AmountOf(denomd)), nil
}

// GetCoinPowers get coin powers by account id
func (a AssetKeeper) GetCoinPowers(ctx sdk.Context, account types.AccountID) types.Coins {
	res, err := a.getCoinsPower(ctx, account)
	if err != nil {
		panic(err)
	}

	return res
}

// GetCoinPower get coin power by account id and coin demon
func (a AssetKeeper) GetCoinPower(ctx sdk.Context, account types.AccountID, creator, symbol types.Name) (types.Coin, error) {
	coins, err := a.getCoinsPower(ctx, account)
	if err != nil {
		return types.Coin{}, sdkerrors.Wrapf(err, "get coin power %s %s-%s", account, creator, symbol)
	}
	denomd := types.CoinDenom(creator, symbol)

	return types.NewCoin(denomd, coins.AmountOf(denomd)), nil
}

// GetCoinPower get coin power by account id and coin demon
func (a AssetKeeper) GetCoinPowerByDenomd(ctx sdk.Context, account types.AccountID, denomd string) types.Coin {
	coins, err := a.getCoinsPower(ctx, account)
	if err != nil {
		panic(err)
	}

	return types.NewCoin(denomd, coins.AmountOf(denomd))
}

// GetCoinDesc get coin description by coin demon
func (a AssetKeeper) GetCoinDesc(ctx sdk.Context, creator, symbol types.Name) (*types.CoinDescription, error) {
	return a.getDescription(ctx, creator, symbol)
}

// GetCoinStat get coin stat data by coin demon
func (a AssetKeeper) GetCoinStat(ctx sdk.Context, creator, symbol types.Name) (*types.CoinStat, error) {
	return a.getStat(ctx, creator, symbol)
}

// GetCoinsTotalSupply get all coin stat data
func (a AssetKeeper) GetCoinsTotalSupply(ctx sdk.Context) types.Coins {
	store := store.NewStore(ctx, a.key)
	iterator := sdk.KVStorePrefixIterator(store, types.GetKeyPrefix(types.CoinStatStoreKeyPrefix))

	res := types.Coins{}

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var stat types.CoinStat

		if err := a.cdc.UnmarshalBinaryBare(iterator.Value(), &stat); err != nil {
			panic(errors.New("unmarshal coins state in store error"))
		}
		// FIXME: No iterator all coins
		res = res.Add(stat.Supply)
	}

	return res
}

// GetCoinTotalSupply get coin stat data by coin demon
func (a AssetKeeper) GetCoinTotalSupply(ctx sdk.Context, creator, symbol types.Name) types.Coin {
	stat, err := a.getStat(ctx, creator, symbol)
	if err != nil {
		panic(sdkerrors.Wrapf(err, "get coin total supply %s/%s", creator, symbol))
	}
	return stat.Supply
}

// IterateAllCoins iterate all account 's coins
func (a AssetKeeper) IterateAllCoins(ctx sdk.Context, cb func(address types.AccountID, balance Coins) (stop bool)) {
	store := store.NewStore(ctx, a.key)
	iterator := sdk.KVStorePrefixIterator(store, types.GetKeyPrefix(types.CoinStoreKeyPrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var (
			coins types.Coins
		)

		id := types.AccountIDFromCoinStoreKey(iterator.Key())

		if err := a.cdc.UnmarshalBinaryBare(iterator.Value(), &coins); err != nil {
			panic(errors.New("unmarshal coins in store error"))
		}

		if cb(id, coins) {
			break
		}
	}
}

// GetCoin get coin by account id and coin demon
func (a AssetKeeper) GetApproveCoins(ctx sdk.Context, account, spender types.AccountID) (*ApproveData, error) {
	return a.getApprove(ctx, account, spender)
}

// IterateAllCoins iterate all account 's coins
func (a AssetKeeper) IterateAllCoinPowers(
	ctx sdk.Context,
	cb func(address types.AccountID, balance Coins) (stop bool)) {
	store := store.NewStore(ctx, a.key)
	iterator := sdk.KVStorePrefixIterator(store, types.GetKeyPrefix(types.CoinPowerStoreKeyPrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var (
			coins types.Coins
		)

		id := types.AccountIDFromCoinStoreKey(iterator.Key())

		if err := a.cdc.UnmarshalBinaryBare(iterator.Value(), &coins); err != nil {
			panic(errors.New("unmarshal coins in store error"))
		}

		if cb(id, coins) {
			break
		}
	}
}
