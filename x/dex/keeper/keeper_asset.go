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

func (a DexKeeper) GetSigInSumForDex(ctx sdk.Context, dex AccountID) Coins {
	key := dexTypes.GenStoreKey(dexTypes.DexSigSumStoreKeyPrefix, dex.Bytes())
	return getCoinsFromKVStore(ctx, a.cdc, a.key, key)
}

func (a DexKeeper) GetSigInForDex(ctx sdk.Context, account, dex AccountID) Coins {
	key := dexTypes.GenStoreKey(dexTypes.DexSigInStoreKeyPrefix, dex.Bytes(), account.Bytes())
	return getCoinsFromKVStore(ctx, a.cdc, a.key, key)
}

func (a DexKeeper) updateSigInSumForDex(ctx sdk.Context, isSub bool, dex AccountID, amt Coins) error {
	key := dexTypes.GenStoreKey(dexTypes.DexSigSumStoreKeyPrefix, dex.Bytes())
	curr := getCoinsFromKVStore(ctx, a.cdc, a.key, key)
	newCoins := curr.Add(amt...)
	if isSub {
		n, isNegative := curr.SafeSub(amt)
		if isNegative {
			return dexTypes.ErrDexSigInChangeToNegative
		}

		newCoins = n
	}

	setCoinsToKVStore(ctx, a.cdc, a.key, key, newCoins)
	return nil
}

func (a DexKeeper) updateSigIn(ctx sdk.Context, isSub bool, id, dex AccountID, amt Coins) (Coins, error) {
	key := dexTypes.GenStoreKey(dexTypes.DexSigInStoreKeyPrefix, dex.Bytes(), id.Bytes())

	curr := getCoinsFromKVStore(ctx, a.cdc, a.key, key)
	newCoins := curr.Add(amt...)
	if isSub {
		n, isNegative := curr.SafeSub(amt)
		if isNegative {
			return Coins{}, dexTypes.ErrDexSigInChangeToNegative
		}

		newCoins = n
	}

	setCoinsToKVStore(ctx, a.cdc, a.key, key, newCoins)

	if err := a.updateSigInSumForDex(ctx, isSub, dex, amt); err != nil {
		return Coins{}, sdkerrors.Wrap(err, "updateSigInSumForDex error")
	}

	return newCoins, nil
}

func (a DexKeeper) SigIn(ctx sdk.Context, id, dex AccountID, amt Coins) error {
	if _, ok := a.getDex(ctx, dex.MustName()); !ok {
		return errors.Wrapf(dexTypes.ErrDexNotExists, "dex %s not exists to sigin", dex.String())
	}

	// check user balance
	if balance, err := a.assetKeeper.GetCoins(ctx, id); nil != err {
		return errors.Wrapf(err, "GetCoins error")
	} else if !balance.IsAllGTE(amt) {
		return errors.Wrapf(dexTypes.ErrDexSigInAmountNotEnough, "user sigIn amount not enough")
	}

	// update sigIn state
	curr, err := a.updateSigIn(ctx, false, id, dex, amt)
	if err != nil {
		return errors.Wrapf(err, "updateSigIn error")
	}

	// change asset state
	if err := a.assetKeeper.Approve(ctx, id, dex, curr, true); err != nil {
		return errors.Wrapf(err, "asset Approve error")
	}

	return nil
}

func (a DexKeeper) GetSigOutReqHeight(ctx sdk.Context, id AccountID) (int64, bool) {
	key := dexTypes.DexSigOutReqStoreKey(id)

	store := ctx.KVStore(a.key)

	if !store.Has(key) {
		return 0, false
	}

	bz := store.Get(key)
	if bz == nil {
		return 0, false
	}

	var res int64
	if err := a.cdc.UnmarshalBinaryBare(bz, &res); err != nil {
		panic(errors.Wrap(err, "get height error from data unmarshal"))
	}

	return res, res > 0
}

func (a DexKeeper) setSigOutReqHeightToStore(ctx sdk.Context, id AccountID, height int64) {
	key := dexTypes.DexSigOutReqStoreKey(id)

	store := ctx.KVStore(a.key)

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

	bz, err := a.cdc.MarshalBinaryBare(height)
	if err != nil {
		panic(errors.Wrap(err, "marshal dex error"))
	}

	store.Set(key, bz)
}

func (a DexKeeper) SigOut(ctx sdk.Context, isTimeout bool, id, dex AccountID, amt Coins) error {
	if _, ok := a.getDex(ctx, dex.MustName()); !ok {
		return errors.Wrapf(dexTypes.ErrDexNotExists, "dex %s not exists to sigin", dex.String())
	}

	if isTimeout {
		reqHeight, ok := a.GetSigOutReqHeight(ctx, id)

		// if not has req, set req height
		if !ok {
			a.setSigOutReqHeightToStore(ctx, id, ctx.BlockHeight())
			// just return
			return nil
		}

		// is has req, so check if current height >= (reqHeight + unlockHeight)
		// has req, but not ok, so error
		if ctx.BlockHeight() >= (reqHeight + SigOutByUserUnlockHeight) {
			// cleanup req
			a.setSigOutReqHeightToStore(ctx, id, 0)
		} else {
			return errors.Wrapf(dexTypes.ErrDexSigOutByUserNoUnlock,
				"sigOut by user need wait to %d but %d", reqHeight+SigOutByUserUnlockHeight, ctx.BlockHeight())
		}
	} else if _, ok := a.GetSigOutReqHeight(ctx, id); ok {
		// cleanup req
		a.setSigOutReqHeightToStore(ctx, id, 0)
	}

	// update sigIn state
	curr, err := a.updateSigIn(ctx, true, id, dex, amt)
	if err != nil {
		return errors.Wrapf(err, "updateSigIn error")
	}

	// change asset state
	if err := a.assetKeeper.Approve(ctx, id, dex, curr, true); err != nil {
		return errors.Wrapf(err, "asset Approve error")
	}

	return nil
}

func (a DexKeeper) Deal(ctx sdk.Context, msgData dexTypes.MsgDexDealData) error {
	dex := msgData.Dex
	if _, ok := a.getDex(ctx, dex.MustName()); !ok {
		return errors.Wrapf(dexTypes.ErrDexNotExists, "dex %s not exists to sigin", dex.String())
	}

	// check user balance
	from := msgData.TransferData.From
	amountFrom := msgData.TransferData.FromAsset.Add(msgData.TransferData.FromFee...)
	fromKey := dexTypes.GenStoreKey(dexTypes.DexSigInStoreKeyPrefix, dex.Bytes(), from.Bytes())
	currDexFrom := getCoinsFromKVStore(ctx, a.cdc, a.key, fromKey)
	_, isNegative := currDexFrom.SafeSub(amountFrom)
	if isNegative {
		return errors.Wrapf(dexTypes.ErrDexDealAmountNotEnough,
			"dex deal user amount not enough")
	}

	to := msgData.TransferData.To
	amountTo := msgData.TransferData.ToAsset
	toKey := dexTypes.GenStoreKey(dexTypes.DexSigInStoreKeyPrefix, dex.Bytes(), to.Bytes())
	currDexTo := getCoinsFromKVStore(ctx, a.cdc, a.key, toKey)
	_, isNegative = currDexTo.SafeSub(amountTo)
	if isNegative {
		return errors.Wrapf(dexTypes.ErrDexDealAmountNotEnough,
			"dex deal user amount not enough")
	}

	// update sigIn state
	_, err := a.updateSigIn(ctx, true, from, dex, amountFrom)
	if err != nil {
		return errors.Wrapf(err, "updateSigIn %s %s by %s error", dex, from, amountFrom)
	}

	_, err = a.updateSigIn(ctx, true, to, dex, amountTo)
	if err != nil {
		return errors.Wrapf(err, "updateSigIn %s %s by %s error", dex, to, amountTo)
	}

	currFrom, err := a.updateSigIn(ctx, false, from, dex, amountTo)
	if err != nil {
		return errors.Wrapf(err, "updateSigIn %s %s by %s error", dex, from, amountTo)
	}

	assetFrom, isNegative := msgData.TransferData.FromAsset.SafeSub(msgData.TransferData.ToFee)
	if isNegative {
		return sdkerrors.Wrap(dexTypes.ErrDexDealAmountNotEnough, "dex deal user coins error")
	}
	currTo, err := a.updateSigIn(ctx, false, to, dex, assetFrom)
	if err != nil {
		return errors.Wrapf(err, "updateSigIn %s %s by %s error", dex, from, assetFrom)
	}

	// change asset state
	if err := a.assetKeeper.Approve(ctx, from, dex, currFrom, true); err != nil {
		return errors.Wrapf(err, "asset Approve error")
	}
	if err := a.assetKeeper.Approve(ctx, to, dex, currTo, true); err != nil {
		return errors.Wrapf(err, "asset Approve error")
	}

	// transfer from->dex->to
	err = a.assetKeeper.TransferDetail(ctx, from, dex, amountFrom, true)
	if err != nil {
		return errors.Wrapf(err, "deal transfer error")
	}
	err = a.assetKeeper.TransferDetail(ctx, dex, to, assetFrom, false)
	if err != nil {
		return errors.Wrapf(err, "deal transfer error")
	}

	// transfer to->dex->from
	err = a.assetKeeper.TransferDetail(ctx, to, dex, amountTo, true)
	if err != nil {
		return errors.Wrapf(err, "deal transfer error")
	}
	err = a.assetKeeper.TransferDetail(ctx, dex, from, amountTo, false)
	if err != nil {
		return errors.Wrapf(err, "deal transfer error")
	}

	return nil
}
