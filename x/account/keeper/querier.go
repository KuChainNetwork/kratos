package keeper

import (
	"github.com/KuChain-io/kuchain/x/account/types"
	abci "github.com/tendermint/tendermint/abci/types"

	chainTypes "github.com/KuChain-io/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(keeper AccountKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryAccount:
			return queryAccount(ctx, req, keeper)
		case types.QueryAuthByAddress:
			return queryAuthByAddress(ctx, req, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

// queryAccount query account handler
func queryAccount(ctx sdk.Context, req abci.RequestQuery, keeper AccountKeeper) ([]byte, error) {
	var params types.QueryAccountParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	account := keeper.GetAccount(ctx, params.Id)
	if account == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", params.Id)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, account)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// queryAuthByAddress query auth by address handler
func queryAuthByAddress(ctx sdk.Context, req abci.RequestQuery, ak AccountKeeper) ([]byte, error) {
	var params types.QueryAuthByAddressParams
	if err := ak.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	auth := params.Address
	store := ctx.KVStore(ak.key)

	bz := store.Get(types.AuthSeqStoreKey(auth))
	if bz == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "auth data no found %s", auth)
	}

	var an types.Auth
	if err := ak.cdc.UnmarshalBinaryBare(bz, &an); err != nil {
		return nil, errors.Wrapf(chainTypes.ErrKuMsgDataUnmarshal, "query auth data unmarshal by %s by %s", auth, err.Error())
	}

	jsonBz, err := codec.MarshalJSONIndent(ak.cdc, an)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return jsonBz, nil
}
