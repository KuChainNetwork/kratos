package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(keeper AssetViewKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryCoin:
			return queryCoin(ctx, req, keeper)
		case types.QueryCoins:
			return queryCoins(ctx, req, keeper)
		case types.QueryCoinPower:
			return queryCoinPower(ctx, req, keeper)
		case types.QueryCoinPowers:
			return queryCoinPowers(ctx, req, keeper)
		case types.QueryCoinStat:
			return queryCoinStat(ctx, req, keeper)
		case types.QueryCoinDescription:
			return queryCoinDesc(ctx, req, keeper)
		case types.QueryCoinLocked:
			return queryCoinLocked(ctx, req, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

// queryCoin query account coin
func queryCoin(ctx sdk.Context, req abci.RequestQuery, keeper AssetViewKeeper) ([]byte, error) {
	cdc := keeper.Cdc()

	var params types.QueryCoinParams
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	coin, err := keeper.GetCoin(ctx, params.AccountID, params.Creator, params.Symbol)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "get coin from keeper")
	}

	bz, err := codec.MarshalJSONIndent(cdc, coin)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// queryCoins query account coins
func queryCoins(ctx sdk.Context, req abci.RequestQuery, keeper AssetViewKeeper) ([]byte, error) {
	cdc := keeper.Cdc()

	var params types.QueryCoinParams
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	coin, err := keeper.GetCoins(ctx, params.AccountID)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "get coin from keeper")
	}

	bz, err := codec.MarshalJSONIndent(cdc, coin)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// queryCoinPowers query account coin power
func queryCoinPowers(ctx sdk.Context, req abci.RequestQuery, keeper AssetViewKeeper) ([]byte, error) {
	cdc := keeper.Cdc()

	var params types.QueryCoinPowersParams
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	coin := keeper.GetCoinPowers(ctx, params.AccountID)

	bz, err := codec.MarshalJSONIndent(cdc, coin)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// queryCoinPower query account coin powers
func queryCoinPower(ctx sdk.Context, req abci.RequestQuery, keeper AssetViewKeeper) ([]byte, error) {
	cdc := keeper.Cdc()

	var params types.QueryCoinPowerParams
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	coin, err := keeper.GetCoinPower(ctx, params.AccountID, params.Creator, params.Symbol)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "get coin from keeper")
	}

	bz, err := codec.MarshalJSONIndent(cdc, coin)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// queryCoinStat query coin state data
func queryCoinStat(ctx sdk.Context, req abci.RequestQuery, keeper AssetViewKeeper) ([]byte, error) {
	cdc := keeper.Cdc()

	var params types.QueryCoinStatParams
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	stat, err := keeper.GetCoinStat(ctx, params.Creator, params.Symbol)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "get stat from keeper")
	}

	if stat == nil {
		return []byte("{}"), nil
	}

	bz, err := codec.MarshalJSONIndent(cdc, *stat)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// queryCoinDesc query coin desc data
func queryCoinDesc(ctx sdk.Context, req abci.RequestQuery, keeper AssetViewKeeper) ([]byte, error) {
	cdc := keeper.Cdc()

	var params types.QueryCoinDescParams
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	desc, err := keeper.GetCoinDesc(ctx, params.Creator, params.Symbol)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "get desc from keeper")
	}

	if desc == nil {
		return []byte("{}"), nil
	}

	bz, err := codec.MarshalJSONIndent(cdc, *desc)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// queryCoinDesc query coin locked
func queryCoinLocked(ctx sdk.Context, req abci.RequestQuery, keeper AssetViewKeeper) ([]byte, error) {
	cdc := keeper.Cdc()

	var params types.QueryLockedCoinsParams
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	all, stat, err := keeper.GetLockCoins(ctx, params.AccountID)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "get desc from keeper")
	}

	res := types.QueryLockedCoinsResponse{
		LockedCoins: all,
		Locks:       stat,
	}

	bz, err := codec.MarshalJSONIndent(cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
