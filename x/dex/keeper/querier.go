package keeper

import (
	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChainNetwork/kuchain/x/dex/types"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(keeper DexKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryDex:
			return queryDex(ctx, req, keeper)
		case types.QuerySigIn:
			return querySigIn(ctx, req, keeper)
		case types.QuerySymbol:
			return querySymbol(ctx, req, keeper)
		default:
			return nil, errors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

// queryDex query dex handler
func queryDex(ctx sdk.Context, req abci.RequestQuery, keeper DexKeeper) ([]byte, error) {
	var params types.QueryDexParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, errors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	dex, ok := keeper.getDex(ctx, params.Creator)
	if dex == nil || !ok {
		return nil, errors.Wrapf(sdkerrors.ErrUnknownAddress, "dex %s does not exist", params.Creator)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, *dex)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// querySigIn query sigIn for a dex
func querySigIn(ctx sdk.Context, req abci.RequestQuery, keeper DexKeeper) ([]byte, error) {
	var params types.QueryDexSigInParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, errors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	coins := keeper.GetSigInForDex(ctx, params.Account, params.Dex)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, coins)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// querySymbol query symbol for dex
func querySymbol(ctx sdk.Context, req abci.RequestQuery, keeper DexKeeper) ([]byte, error) {
	var params types.QuerySymbolParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, errors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	dex, ok := keeper.getDex(ctx, params.Creator)
	if dex == nil || !ok {
		return nil, errors.Wrapf(sdkerrors.ErrUnknownAddress, "dex %s does not exist", params.Creator)
	}

	symbol, ok := dex.Symbol(params.BaseCreator, params.BaseCode, params.QuoteCode, params.QuoteCreator)
	if !ok {
		return nil, errors.Wrapf(types.ErrSymbolNotExists, "%s/%s:%s/%s not exists",
			params.BaseCreator,
			params.BaseCode,
			params.QuoteCreator,
			params.QuoteCode)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, symbol)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
