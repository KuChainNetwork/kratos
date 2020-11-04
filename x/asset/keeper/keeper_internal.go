package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/store"
	"github.com/KuChainNetwork/kuchain/chain/types/coin"
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Some internal implements for asset keeper

func (a AssetKeeper) issueCoinStat(ctx sdk.Context, amount types.Coin) error {
	creator, symbol, err := CoinAccountsFromDenom(amount.Denom)
	if err != nil {
		return sdkerrors.Wrap(err, "issue coin error")
	}

	stat, err := a.getStat(ctx, creator, symbol)
	if stat == nil {
		return types.ErrAssetCoinNoExit
	}

	if err != nil {
		return sdkerrors.Wrap(err, "issue get stat error")
	}

	if !stat.Creator.Eq(creator) {
		return types.ErrAssetNoCreator
	}

	// check denom
	denom := types.CoinDenom(creator, symbol)
	if err := coin.ValidateDenom(denom); err != nil {
		return sdkerrors.Wrapf(types.ErrAssetDenom, "denom %s", denom)
	}

	if denom != amount.Denom {
		return sdkerrors.Wrap(types.ErrAssetDenom, "amount denom error")
	}

	newSupply := stat.Supply.Add(amount)
	maxSupply := stat.GetCurrentMaxSupplyLimit(ctx.BlockHeight())
	ctx.Logger().Debug("update state", "newSupply", newSupply, "maxSupply", stat.MaxSupply, "limit", maxSupply)

	// Core token no limit
	if (amount.Denom != constants.DefaultBondDenom) && (!maxSupply.IsGTE(newSupply)) {
		return types.ErrAssetIssueGTMaxSupply
	}

	stat.Supply = newSupply
	if err := a.setStat(ctx, stat); err != nil {
		return sdkerrors.Wrap(err, "set stat")
	}

	return nil
}

func (a AssetKeeper) burnCoinStat(ctx sdk.Context, amount types.Coin) error {
	creator, symbol, err := CoinAccountsFromDenom(amount.Denom)
	if err != nil {
		return sdkerrors.Wrap(err, "issue coin error")
	}

	stat, err := a.getStat(ctx, creator, symbol)
	if stat == nil {
		return types.ErrAssetCoinNoExit
	}

	if err != nil {
		return sdkerrors.Wrap(err, "issue get stat error")
	}

	if !stat.Creator.Eq(creator) {
		return types.ErrAssetNoCreator
	}

	// check denom
	denom := types.CoinDenom(creator, symbol)
	if err := coin.ValidateDenom(denom); err != nil {
		return sdkerrors.Wrapf(types.ErrAssetDenom, "denom %s", denom)
	}

	if denom != amount.Denom {
		return sdkerrors.Wrap(types.ErrAssetDenom, "amount denom error")
	}

	stat.Supply = stat.Supply.Sub(amount)
	if err := a.setStat(ctx, stat); err != nil {
		return sdkerrors.Wrap(err, "set stat")
	}

	return nil
}

func (a AssetKeeper) setCoins(ctx sdk.Context, account types.AccountID, coin types.Coins) error {
	store := store.NewStore(ctx, a.key)
	bz, err := a.cdc.MarshalBinaryBare(coin)
	if err != nil {
		return sdkerrors.Wrap(err, "set coins marshal error")
	}

	key := types.CoinStoreKey(account)

	if bz == nil {
		ctx.Logger().Debug("set coins", "account", account, "coin", coin)
		if store.Has(key) {
			store.Delete(key)
		}
		return nil
	}
	store.Set(key, bz)
	return nil
}

func (a AssetKeeper) getCoins(ctx sdk.Context, account types.AccountID) (types.Coins, error) {
	store := store.NewStore(ctx, a.key)
	bz := store.Get(types.CoinStoreKey(account))
	if bz == nil {
		return types.Coins{}, nil
	}

	var coins types.Coins

	if err := a.cdc.UnmarshalBinaryBare(bz, &coins); err != nil {
		return types.Coins{}, sdkerrors.Wrap(err, "get coins unmarshal")
	}

	return coins, nil
}

func (a AssetKeeper) setStat(ctx sdk.Context, stat *types.CoinStat) error {
	store := store.NewStore(ctx, a.key)
	bz, err := a.cdc.MarshalBinaryBare(*stat)
	if err != nil {
		return sdkerrors.Wrap(err, "set stat marshal error")
	}
	store.Set(types.CoinStatStoreKey(stat.Creator, stat.Symbol), bz)
	return nil
}

func (a AssetKeeper) getStat(ctx sdk.Context, creator, symbol types.Name) (*types.CoinStat, error) {
	store := store.NewStore(ctx, a.key)
	bz := store.Get(types.CoinStatStoreKey(creator, symbol))
	if bz == nil {
		return nil, types.ErrAssetCoinNoExit
	}

	var stat types.CoinStat

	if err := a.cdc.UnmarshalBinaryBare(bz, &stat); err != nil {
		return nil, sdkerrors.Wrap(err, "get stat unmarshal")
	}

	return &stat, nil
}

func (a AssetKeeper) setDescription(ctx sdk.Context, desc *types.CoinDescription) error {
	store := store.NewStore(ctx, a.key)
	bz, err := a.cdc.MarshalBinaryBare(*desc)
	if err != nil {
		return sdkerrors.Wrap(err, "set desc marshal error")
	}
	store.Set(types.CoinDescStoreKey(desc.Creator, desc.Symbol), bz)
	return nil
}

func (a AssetKeeper) getDescription(ctx sdk.Context, creator, symbol types.Name) (*types.CoinDescription, error) {
	store := store.NewStore(ctx, a.key)
	bz := store.Get(types.CoinDescStoreKey(creator, symbol))
	if bz == nil {
		return nil, types.ErrAssetCoinNoExit
	}

	var res types.CoinDescription

	if err := a.cdc.UnmarshalBinaryBare(bz, &res); err != nil {
		return nil, sdkerrors.Wrap(err, "get desc unmarshal")
	}

	return &res, nil
}

func (a AssetKeeper) setCoinsPower(ctx sdk.Context, account types.AccountID, coin types.Coins) error {
	store := store.NewStore(ctx, a.key)
	bz, err := a.cdc.MarshalBinaryBare(coin)
	if err != nil {
		return sdkerrors.Wrap(err, "set coins marshal error")
	}

	key := types.CoinPowerStoreKey(account)
	if bz == nil {
		ctx.Logger().Debug("set coins", "account", account, "coin", coin)
		if store.Has(key) {
			store.Delete(key)
		}
		return nil
	}

	store.Set(key, bz)
	return nil
}

func (a AssetKeeper) getCoinsPower(ctx sdk.Context, account types.AccountID) (types.Coins, error) {
	store := store.NewStore(ctx, a.key)
	bz := store.Get(types.CoinPowerStoreKey(account))
	if bz == nil {
		return types.Coins{}, nil
	}

	var coins types.Coins

	if err := a.cdc.UnmarshalBinaryBare(bz, &coins); err != nil {
		return types.Coins{}, sdkerrors.Wrap(err, "get coins unmarshal")
	}

	return coins, nil
}

func (a AssetKeeper) setApprove(ctx sdk.Context, account, spender types.AccountID, data ApproveData) error {
	store := ctx.KVStore(a.key)
	key := types.ApproveStoreKey(account, spender)

	if data.Amount.IsZero() {
		// delete
		if store.Has(key) {
			store.Delete(key)
		}
		return nil
	}

	bz, err := a.cdc.MarshalBinaryBare(data)
	if err != nil {
		return sdkerrors.Wrap(err, "set coins marshal error")
	}

	if bz == nil {
		if store.Has(key) {
			store.Delete(key)
		}
		return nil
	}

	store.Set(key, bz)
	return nil
}

func (a AssetKeeper) updateApproveSum(ctx sdk.Context, account types.AccountID, sum types.Coins) error {
	store := ctx.KVStore(a.key)
	key := types.ApproveSumStoreKey(account)

	bz, err := a.cdc.MarshalBinaryBare(sum)
	if err != nil {
		return sdkerrors.Wrap(err, "set coins marshal error")
	}

	if bz == nil {
		if store.Has(key) {
			store.Delete(key)
		}
		return nil
	}

	store.Set(key, bz)
	return nil
}

func (a AssetKeeper) getApprove(ctx sdk.Context, account, spender types.AccountID) (*ApproveData, error) {
	store := ctx.KVStore(a.key)
	bz := store.Get(types.ApproveStoreKey(account, spender))
	if bz == nil {
		return nil, nil
	}

	var res ApproveData

	if err := a.cdc.UnmarshalBinaryBare(bz, &res); err != nil {
		return nil, sdkerrors.Wrap(err, "get ApproveData unmarshal")
	}

	return &res, nil
}

func (a AssetKeeper) GetApproveSum(ctx sdk.Context, account types.AccountID) (types.Coins, error) {
	store := ctx.KVStore(a.key)
	bz := store.Get(types.ApproveSumStoreKey(account))
	if bz == nil {
		return nil, nil
	}

	var res types.Coins

	if err := a.cdc.UnmarshalBinaryBare(bz, &res); err != nil {
		return nil, sdkerrors.Wrap(err, "get approve sum unmarshal")
	}

	return res, nil
}
