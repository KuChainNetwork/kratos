package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	dexTypes "github.com/KuChainNetwork/kuchain/x/dex/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

func getCoinsFromKVStore(ctx sdk.Context, cdc *codec.Codec, storeKey types.StoreKey, key []byte) Coins {
	store := ctx.KVStore(storeKey)
	bz := store.Get(key)
	if bz == nil {
		return Coins{}
	}

	res := Coins{}
	if err := cdc.UnmarshalBinaryBare(bz, &res); err != nil {
		panic(errors.Wrap(err, "get coins from data unmarshal"))
	}

	return res
}

func setCoinsToKVStore(ctx sdk.Context, cdc *codec.Codec, storeKey types.StoreKey, key []byte, amt Coins) {
	store := ctx.KVStore(storeKey)
	if amt.IsZero() {
		if store.Has(key) {
			store.Delete(key)
		}
		return
	}

	bz, err := cdc.MarshalBinaryBare(amt)
	if err != nil {
		panic(errors.Wrap(err, "marshal dex error"))
	}

	store.Set(key, bz)
}

func (k DexKeeper) GetSigInSumForDex(ctx sdk.Context, dex AccountID) Coins {
	key := dexTypes.GenStoreKey(dexTypes.DexSigSumStoreKeyPrefix, dex.Bytes())
	return getCoinsFromKVStore(ctx, k.cdc, k.key, key)
}

func (k DexKeeper) GetSigInForDex(ctx sdk.Context, account, dex AccountID) Coins {
	key := dexTypes.GenStoreKey(dexTypes.DexSigInStoreKeyPrefix, dex.Bytes(), account.Bytes())
	return getCoinsFromKVStore(ctx, k.cdc, k.key, key)
}

func (k DexKeeper) updateSigInSumForDex(ctx sdk.Context, isSub bool, dex AccountID, amt Coins) error {
	key := dexTypes.GenStoreKey(dexTypes.DexSigSumStoreKeyPrefix, dex.Bytes())
	curr := getCoinsFromKVStore(ctx, k.cdc, k.key, key)
	newCoins := curr.Add(amt...)
	if isSub {
		n, isNegative := curr.SafeSub(amt)
		if isNegative {
			return dexTypes.ErrDexSigInChangeToNegative
		}

		newCoins = n
	}

	setCoinsToKVStore(ctx, k.cdc, k.key, key, newCoins)
	return nil
}

func (k DexKeeper) updateSigIn(ctx sdk.Context, isSub bool, id, dex AccountID, amt Coins) (Coins, error) {
	key := dexTypes.GenStoreKey(dexTypes.DexSigInStoreKeyPrefix, dex.Bytes(), id.Bytes())

	curr := getCoinsFromKVStore(ctx, k.cdc, k.key, key)
	newCoins := curr.Add(amt...)
	if isSub {
		n, isNegative := curr.SafeSub(amt)
		if isNegative {
			return Coins{}, dexTypes.ErrDexSigInChangeToNegative
		}

		newCoins = n
	}

	setCoinsToKVStore(ctx, k.cdc, k.key, key, newCoins)

	if err := k.updateSigInSumForDex(ctx, isSub, dex, amt); err != nil {
		return Coins{}, sdkerrors.Wrap(err, "updateSigInSumForDex error")
	}

	return newCoins, nil
}

func (k DexKeeper) SigIn(ctx sdk.Context, id, dex AccountID, amt Coins) error {
	if _, ok := k.getDex(ctx, dex.MustName()); !ok {
		return errors.Wrapf(dexTypes.ErrDexNotExists, "dex %s not exists to sigin", dex.String())
	}

	// check user balance
	if balance, err := k.assetKeeper.GetCoins(ctx, id); nil != err {
		return errors.Wrapf(err, "GetCoins error")
	} else if !balance.IsAllGTE(amt) {
		return errors.Wrapf(dexTypes.ErrDexSigInAmountNotEnough, "user sigIn amount not enough")
	}

	// update sigIn state
	curr, err := k.updateSigIn(ctx, false, id, dex, amt)
	if err != nil {
		return errors.Wrapf(err, "updateSigIn error")
	}

	// change asset state
	if err := k.assetKeeper.Approve(ctx, id, dex, curr, true); err != nil {
		return errors.Wrapf(err, "asset Approve error")
	}

	return nil
}

func (k DexKeeper) GetSigOutReqHeight(ctx sdk.Context, id AccountID) (int64, bool) {
	key := dexTypes.DexSigOutReqStoreKey(id)

	store := ctx.KVStore(k.key)

	if !store.Has(key) {
		return 0, false
	}

	bz := store.Get(key)
	if bz == nil {
		return 0, false
	}

	var res int64
	if err := k.cdc.UnmarshalBinaryBare(bz, &res); err != nil {
		panic(errors.Wrap(err, "get height error from data unmarshal"))
	}

	return res, res > 0
}

func (k DexKeeper) setSigOutReqHeightToStore(ctx sdk.Context, id AccountID, height int64) {
	key := dexTypes.DexSigOutReqStoreKey(id)

	store := ctx.KVStore(k.key)

	// cleanup sigout req
	if height == 0 {
		if store.Has(key) {
			store.Delete(key)
		}
		return
	}

	if store.Has(key) {
		panic(errors.Errorf("cannot update req height"))
	}

	bz, err := k.cdc.MarshalBinaryBare(height)
	if err != nil {
		panic(errors.Wrap(err, "marshal dex error"))
	}

	store.Set(key, bz)
}

func (k DexKeeper) SigOut(ctx sdk.Context, isTimeout bool, id, dex AccountID, amt Coins) error {
	if _, ok := k.getDex(ctx, dex.MustName()); !ok {
		return errors.Wrapf(dexTypes.ErrDexNotExists, "dex %s not exists to sigin", dex.String())
	}

	if isTimeout {
		reqHeight, ok := k.GetSigOutReqHeight(ctx, id)

		// if not has req, set req height
		if !ok {
			k.setSigOutReqHeightToStore(ctx, id, ctx.BlockHeight())
			// just return
			return nil
		}

		// is has req, so check if current height >= (reqHeight + unlockHeight)
		// has req, but not ok, so error
		if ctx.BlockHeight() >= (reqHeight + SigOutByUserUnlockHeight) {
			// cleanup req
			k.setSigOutReqHeightToStore(ctx, id, 0)
		} else {
			return errors.Wrapf(dexTypes.ErrDexSigOutByUserNoUnlock,
				"sigOut by user need wait to %d but %d", reqHeight+SigOutByUserUnlockHeight, ctx.BlockHeight())
		}
	} else {
		// if has req, but sigout by dex
		if _, ok := k.GetSigOutReqHeight(ctx, id); ok {
			// cleanup req
			k.setSigOutReqHeightToStore(ctx, id, 0)
		}
	}

	// update sigIn state
	curr, err := k.updateSigIn(ctx, true, id, dex, amt)
	if err != nil {
		return errors.Wrapf(err, "updateSigIn error")
	}

	// change asset state
	if err := k.assetKeeper.Approve(ctx, id, dex, curr, true); err != nil {
		return errors.Wrapf(err, "asset Approve error")
	}

	return nil
}

func (k DexKeeper) Deal(ctx sdk.Context, dex, from, to AccountID, amtFrom, amtTo Coins) error {
	if _, ok := k.getDex(ctx, dex.MustName()); !ok {
		return errors.Wrapf(dexTypes.ErrDexNotExists, "dex %s not exists to sigin", dex.String())
	}

	// update sigIn state
	if _, err := k.updateSigIn(ctx, true, from, dex, amtFrom); err != nil {
		return errors.Wrapf(err, "updateSigIn %s %s by %s error", dex, from, amtFrom)
	}

	if _, err := k.updateSigIn(ctx, true, to, dex, amtTo); err != nil {
		return errors.Wrapf(err, "updateSigIn %s %s by %s error", dex, to, amtTo)
	}

	return nil
}
