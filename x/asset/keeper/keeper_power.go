package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IssueCoinPower add coin power to account, if a account has coin power, it can gen coins by power
func (a AssetKeeper) IssueCoinPower(ctx sdk.Context, id types.AccountID, amt Coins) (Coins, error) {
	for _, c := range amt {
		if err := a.issueCoinStat(ctx, c); err != nil {
			return Coins{}, sdkerrors.Wrapf(err, "issue %s state error", c)
		}
	}
	return a.addCoinPower(ctx, id, amt)
}

func (a AssetKeeper) addCoinPower(ctx sdk.Context, id types.AccountID, amt Coins) (Coins, error) {
	coins, err := a.getCoinsPower(ctx, id)
	if err != nil {
		return Coins{}, sdkerrors.Wrapf(err, "get %s coins powers error", id)
	}

	if amt.IsZero() {
		return coins, nil
	}

	added := coins.Add(amt...)

	if err := a.setCoinsPower(ctx, id, added); err != nil {
		return Coins{}, sdkerrors.Wrapf(err, "set coins powers in add %s error", id)
	}

	return added, nil
}

// BurnCoinPower sub coin power from account, if has no power return error
func (a AssetKeeper) BurnCoinPower(ctx sdk.Context, id types.AccountID, amt Coins) (Coins, error) {
	for _, c := range amt {
		if err := a.burnCoinStat(ctx, c); err != nil {
			return Coins{}, sdkerrors.Wrapf(err, "burn %s state error", c)
		}
	}
	return a.subCoinPower(ctx, id, amt)
}

func (a AssetKeeper) subCoinPower(ctx sdk.Context, id types.AccountID, amt Coins) (Coins, error) {
	coins, err := a.getCoinsPower(ctx, id)
	if err != nil {
		return Coins{}, sdkerrors.Wrapf(err, "get %s coins powers error", id)
	}

	if amt.IsZero() {
		return coins, nil
	}

	if amt.IsAnyNegative() {
		return Coins{}, sdkerrors.Wrapf(types.ErrAssetCoinNoEnough, "amt should not be negative")
	}

	ctx.Logger().Debug("sub coin power", "c", coins, "amt", amt)

	subed, hasNeg := coins.SafeSub(amt)
	if hasNeg {
		return Coins{}, sdkerrors.Wrapf(types.ErrAssetCoinNoEnough, "sub coin power error no enough")
	}

	if subed == nil {
		subed = NewCoins()
	}

	if err := a.setCoinsPower(ctx, id, subed); err != nil {
		return Coins{}, sdkerrors.Wrapf(err, "set coins powers in add %s error", id)
	}

	return subed, nil
}

// SendCoinPower sub coin power from account, if has no power return error
func (a AssetKeeper) SendCoinPower(ctx sdk.Context, from, to types.AccountID, amt Coins) error {
	// just return if from == to
	if from.Eq(to) {
		return nil
	}

	if _, err := a.subCoinPower(ctx, from, amt); err != nil {
		return sdkerrors.Wrapf(err, "get %s coins powers in send error", from)
	}

	if _, err := a.addCoinPower(ctx, to, amt); err != nil {
		return sdkerrors.Wrapf(err, "set coins power in add %s error", to)
	}

	if err := a.ak.EnsureAccount(ctx, to); err != nil {
		return sdkerrors.Wrapf(err, "ensure account %s error", to)
	}

	return nil
}

// CoinsToPower accounts coins to coin power, so that it can be send to module account
func (a AssetKeeper) CoinsToPower(ctx sdk.Context, from, to types.AccountID, amt Coins) error {
	if amt.IsZero() {
		return nil
	}

	// sub coin from account, this will not change the coin supply stat
	coins, err := a.getCoins(ctx, from)
	if err != nil {
		return sdkerrors.Wrap(err, "CoinsToPower: get from coin error")
	}

	subed, hasNeg := coins.SafeSub(amt)
	if hasNeg {
		return sdkerrors.Wrap(types.ErrAssetCoinNoEnough, "CoinsToPower: sub coins")
	}

	if err := a.checkIsCanUseCoins(ctx, from, amt, coins, false); err != nil {
		return sdkerrors.Wrap(err, "coinsToPower")
	}

	if err := a.setCoins(ctx, from, subed); err != nil {
		return sdkerrors.Wrap(err, "CoinsToPower: set coins")
	}

	if _, err := a.addCoinPower(ctx, to, amt); err != nil {
		return sdkerrors.Wrapf(err, "CoinsToPower: set coins power in add %s error", to)
	}

	return nil
}

// ExerciseCoinPower exercise coin power to get coins to account
func (a AssetKeeper) ExerciseCoinPower(ctx sdk.Context, id types.AccountID, amt types.Coin) error {
	ctx.Logger().Debug("exercise coin power", "id", id, "amount", amt)

	if amt.IsZero() {
		return nil
	}

	if _, err := a.subCoinPower(ctx, id, NewCoins(amt)); err != nil {
		return sdkerrors.Wrapf(err, "get %s coins powers in exercise error", id)
	}

	coins, err := a.getCoins(ctx, id)
	if err != nil {
		return sdkerrors.Wrapf(err, "get coins error")
	}

	return sdkerrors.Wrapf(a.setCoins(ctx, id, coins.Add(amt)), "set coins in exercise error")
}
