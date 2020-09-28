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

func (k DexKeeper) updateSigInSumForDex(ctx sdk.Context, dex AccountID, amt Coins) error {
	key := dexTypes.GenStoreKey(dexTypes.DexSigSumStoreKeyPrefix, dex.Bytes())
	curr := getCoinsFromKVStore(ctx, k.cdc, k.key, key)
	newCoins := curr.Add(amt...)
	if newCoins.IsAnyNegative() {
		return dexTypes.ErrDexSigInChangeToNegative
	}
	setCoinsToKVStore(ctx, k.cdc, k.key, key, newCoins)
	return nil
}

func (k DexKeeper) updateSigIn(ctx sdk.Context, id, dex AccountID, amt Coins) error {
	key := dexTypes.GenStoreKey(dexTypes.DexSigInStoreKeyPrefix, dex.Bytes(), id.Bytes())

	curr := getCoinsFromKVStore(ctx, k.cdc, k.key, key)
	newCoins := curr.Add(amt...)
	if newCoins.IsAnyNegative() {
		return dexTypes.ErrDexSigInChangeToNegative
	}
	setCoinsToKVStore(ctx, k.cdc, k.key, key, newCoins)

	return sdkerrors.Wrap(k.updateSigInSumForDex(ctx, dex, amt), "updateSigInSumForDex error")
}

func (k DexKeeper) SigIn(ctx sdk.Context, id, dex AccountID, amt Coins) error {
	if _, ok := k.getDex(ctx, dex.MustName()); !ok {
		return errors.Wrapf(dexTypes.ErrDexNotExists, "dex %s not exists to sigin", dex.String())
	}

	// update sigIn state
	if err := k.updateSigIn(ctx, id, dex, amt); err != nil {
		return errors.Wrapf(err, "updateSigIn error")
	}

	// change asset state
	if err := k.assetKeeper.Approve(ctx, id, dex, amt, true); err != nil {
		return errors.Wrapf(err, "asset Approve error")
	}

	if err := k.assetKeeper.LockCoins(ctx, id, -1, amt); err != nil {
		return errors.Wrapf(err, "asset lock coins error")
	}

	return nil
}
