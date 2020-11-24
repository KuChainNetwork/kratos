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
	balance, err := a.assetKeeper.GetCoins(ctx, id)
	if nil != err {
		return errors.Wrapf(err, "GetCoins error")
	}
	approve, err := a.assetKeeper.GetApproveCoins(ctx, id, dex)
	if err != nil {
		return errors.Wrapf(err, "GetApproveCoins error")
	}
	currAmt := amt
	if approve != nil {
		currAmt = approve.Amount.Add(amt...)
	}
	if !balance.IsAllGTE(currAmt) {
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

func (a DexKeeper) Deal(ctx sdk.Context, dex, from, to AccountID,
	amtFrom, amtTo, feeFrom, feeTo Coins) error {
	if _, ok := a.getDex(ctx, dex.MustName()); !ok {
		return errors.Wrapf(dexTypes.ErrDexNotExists, "dex %s not exists to sigin", dex.String())
	}

	// update sigIn state to FromAccount, should sub amtFrom(include fee), and add gotted(amtTo-toFee)
	approveAddForFrom := amtTo.Sub(feeTo)
	approveNowForFrom, err := a.assetKeeper.GetApproveCoins(ctx, from, dex)
	if err != nil {
		return errors.Wrapf(err, "get approve data error acc: %s, dex: %s", from, dex)
	}

	if approveNowForFrom != nil && !approveNowForFrom.Amount.IsZero() {
		approveAddForFrom = approveAddForFrom.Add(approveNowForFrom.Amount...)
	}

	if err := a.assetKeeper.Approve(ctx, from, dex, approveAddForFrom, true); err != nil {
		return errors.Wrapf(err, "asset Approve error")
	}

	if _, err := a.updateSigIn(ctx, true, from, dex, amtFrom); err != nil {
		return errors.Wrapf(err, "updateSigIn %s %s by %s error", dex, from, amtFrom)
	}

	if _, err := a.updateSigIn(ctx, false, from, dex, amtTo.Sub(feeTo)); err != nil {
		return errors.Wrapf(err, "updateSigIn %s %s by %s error", dex, from, amtFrom)
	}

	// update sigIn state to ToAccount, should sub amtTo(include fee), and add gotted(amtFrom-fromFee)
	approveAddForTo := amtFrom.Sub(feeFrom)
	approveNowForTo, err := a.assetKeeper.GetApproveCoins(ctx, to, dex)
	if err != nil {
		return errors.Wrapf(err, "get approve data error acc: %s, dex: %s", to, dex)
	}

	if approveNowForTo != nil && !approveNowForTo.Amount.IsZero() {
		approveAddForTo = approveAddForTo.Add(approveNowForTo.Amount...)
	}

	if err := a.assetKeeper.Approve(ctx, to, dex, approveAddForTo, true); err != nil {
		return errors.Wrapf(err, "asset Approve error")
	}

	if _, err := a.updateSigIn(ctx, true, to, dex, amtTo); err != nil {
		return errors.Wrapf(err, "updateSigIn %s %s by %s error", dex, to, amtTo)
	}

	if _, err := a.updateSigIn(ctx, false, to, dex, amtFrom.Sub(feeFrom)); err != nil {
		return errors.Wrapf(err, "updateSigIn %s %s by %s error", dex, from, amtFrom)
	}

	return nil
}
