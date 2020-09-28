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

	// update sigIn state
	curr, err := k.updateSigIn(ctx, false, id, dex, amt)
	if err != nil {
		return errors.Wrapf(err, "updateSigIn error")
	}

	// change asset state
	if err := k.assetKeeper.Approve(ctx, id, dex, curr, true); err != nil {
		return errors.Wrapf(err, "asset Approve error")
	}

	if err := k.assetKeeper.LockCoins(ctx, id, -1, amt); err != nil {
		return errors.Wrapf(err, "asset lock coins error")
	}

	return nil
}

func (k DexKeeper) SigOut(ctx sdk.Context, isTimeout bool, id, dex AccountID, amt Coins) error {
	if _, ok := k.getDex(ctx, dex.MustName()); !ok {
		return errors.Wrapf(dexTypes.ErrDexNotExists, "dex %s not exists to sigin", dex.String())
	}

	// TODO: imp isTimeout

	// update sigIn state
	curr, err := k.updateSigIn(ctx, true, id, dex, amt)
	if err != nil {
		return errors.Wrapf(err, "updateSigIn error")
	}

	// change asset state
	if err := k.assetKeeper.Approve(ctx, id, dex, curr, true); err != nil {
		return errors.Wrapf(err, "asset Approve error")
	}

	if err := k.assetKeeper.UnLockFreezedCoins(ctx, id, amt); err != nil {
		return errors.Wrapf(err, "asset lock coins error")
	}

	return nil
}
