package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type LockedCoins = types.LockedCoins

type accountLockedCoins struct {
	ID      types.AccountID `json:"id" yaml:"id"`
	Lockeds []LockedCoins   `json:"lockeds" yaml:"lockeds"`
}

func (a AssetKeeper) getCoinsLockedStat(ctx sdk.Context, id types.AccountID) (accountLockedCoins, error) {
	store := ctx.KVStore(a.key)
	bz := store.Get(types.CoinLockedStatStoreKey(id))
	res := accountLockedCoins{
		ID: id,
	}

	if bz == nil {
		return res, nil
	}

	if err := a.cdc.UnmarshalBinaryBare(bz, &res); err != nil {
		return res, sdkerrors.Wrap(err, "get coins locked state unmarshal")
	}

	return res, nil
}

func (a AssetKeeper) setCoinsLockedStat(ctx sdk.Context, id types.AccountID, stat accountLockedCoins) error {
	store := ctx.KVStore(a.key)
	bz, err := a.cdc.MarshalBinaryBare(stat)
	if err != nil {
		return sdkerrors.Wrap(err, "set coins locked state marshal error")
	}
	store.Set(types.CoinLockedStatStoreKey(id), bz)
	return nil
}

func (a AssetKeeper) setCoinsLocked(ctx sdk.Context, account types.AccountID, coin types.Coins) error {
	store := ctx.KVStore(a.key)
	bz, err := a.cdc.MarshalBinaryBare(coin)
	if err != nil {
		return sdkerrors.Wrap(err, "set coins locked marshal error")
	}

	key := types.CoinLockedStoreKey(account)
	if bz == nil {
		ctx.Logger().Debug("set coins", "account", account, "coin", coin)
		if store.Has(key) {
			store.Delete(key)
		}
		return nil
	}

	store.Set(types.CoinLockedStoreKey(account), bz)
	return nil
}

func (a AssetKeeper) getCoinsLocked(ctx sdk.Context, account types.AccountID) (types.Coins, error) {
	store := ctx.KVStore(a.key)
	bz := store.Get(types.CoinLockedStoreKey(account))
	if bz == nil {
		return types.Coins{}, nil
	}

	var coins types.Coins

	if err := a.cdc.UnmarshalBinaryBare(bz, &coins); err != nil {
		return types.Coins{}, sdkerrors.Wrap(err, "get coins locked unmarshal")
	}

	return coins, nil
}

// LockCoins lock coins for a account
func (a AssetKeeper) LockCoins(ctx sdk.Context, account types.AccountID, unlockBlockHeight int64, coins types.Coins) error {
	if coins.IsZero() {
		return nil
	}

	currentCoins, err := a.getCoins(ctx, account)
	if err != nil {
		return sdkerrors.Wrap(err, "LockCoins: get coins in lock coins")
	}

	coinLocked, err := a.getCoinsLocked(ctx, account)
	if err != nil {
		return sdkerrors.Wrap(err, "LockCoins: get coins locked")
	}

	coinsLockedAll := coinLocked.Add(coins...)

	if !currentCoins.IsAllGTE(coinsLockedAll) {
		return types.ErrAssetLockCoinsNoEnough
	}

	if unlockBlockHeight <= ctx.BlockHeight() {
		return types.ErrAssetLockUnlockBlockHeightErr
	}

	stat, err := a.getCoinsLockedStat(ctx, account)
	if err != nil {
		return sdkerrors.Wrap(err, "LockCoins: get coins locked stat")
	}

	stat.Lockeds = append(stat.Lockeds, LockedCoins{
		Coins:             coins,
		UnlockBlockHeight: unlockBlockHeight,
	})

	if err := a.setCoinsLockedStat(ctx, account, stat); err != nil {
		return sdkerrors.Wrap(err, "LockCoins: set coins locked stat")
	}

	if err := a.setCoinsLocked(ctx, account, coinsLockedAll); err != nil {
		return sdkerrors.Wrap(err, "LockCoins: add coins locked")
	}

	return nil
}

func (a AssetKeeper) UnLockCoins(ctx sdk.Context, account types.AccountID, coins types.Coins) error {
	coinLocked, err := a.getCoinsLocked(ctx, account)
	if err != nil {
		return sdkerrors.Wrap(err, "UnlockCoins: get coins locked")
	}

	stat, err := a.getCoinsLockedStat(ctx, account)
	if err != nil {
		return sdkerrors.Wrap(err, "UnlockCoins: get coins locked stat")
	}

	height := ctx.BlockHeight()
	coinsCanUnLocked := types.Coins{}

	newStat := accountLockedCoins{
		ID:      account,
		Lockeds: make([]LockedCoins, 0, len(stat.Lockeds)),
	}

	for _, l := range stat.Lockeds {
		if l.UnlockBlockHeight <= height {
			coinsCanUnLocked = coinsCanUnLocked.Add(l.Coins...)
		} else {
			newStat.Lockeds = append(newStat.Lockeds, l)
		}
	}

	if !coinsCanUnLocked.IsEqual(coins) {
		return sdkerrors.Wrapf(types.ErrAssetUnLockCoins, "unlock should be %s", coinsCanUnLocked.String())
	}

	newCoinsLocked, isNegative := coinLocked.SafeSub(coinsCanUnLocked)
	if isNegative {
		return sdkerrors.Wrapf(types.ErrAssetUnLockCoins, "unlock sum be %s >= %s",
			coinsCanUnLocked.String(), coinLocked.String())
	}

	err = a.setCoinsLocked(ctx, account, newCoinsLocked)
	if err != nil {
		return sdkerrors.Wrap(err, "UnlockedCoins")
	}

	err = a.setCoinsLockedStat(ctx, account, newStat)
	return sdkerrors.Wrap(err, "UnlockedCoins")

}

// GetLockCoins get locked data
func (a AssetKeeper) GetLockCoins(ctx sdk.Context, account types.AccountID) (types.Coins, []LockedCoins, error) {
	lockedStat, err := a.getCoinsLockedStat(ctx, account)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "get lock stat")
	}

	all, err := a.getCoinsLocked(ctx, account)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "get lock coins")
	}

	return all, lockedStat.Lockeds, nil
}

// CheckIsCanUseCoins check if the account can use this coins
func (a AssetKeeper) CheckIsCanUseCoins(ctx sdk.Context, account types.AccountID, coins types.Coins) error {
	currentCoins, err := a.getCoins(ctx, account)
	if err != nil {
		return sdkerrors.Wrap(err, "CheckIsCanUseCoins: get coins in lock coins")
	}

	return a.checkIsCanUseCoins(ctx, account, coins, currentCoins)
}

func (a AssetKeeper) checkIsCanUseCoins(ctx sdk.Context, account types.AccountID, coins, currentCoins types.Coins) error {
	coinLocked, err := a.getCoinsLocked(ctx, account)
	if err != nil {
		return sdkerrors.Wrap(err, "CheckIsCanUseCoins: get coins locked")
	}

	if currentCoins.IsAllGTE(coinLocked.Add(coins...)) {
		return nil
	} else {
		return types.ErrAssetCoinsLocked
	}
}
