package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/account/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChainNetwork/kuchain/chain/store"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		case types.QueryAccountsByAuth:
			return queryAccountsByAuth(ctx, req, keeper)
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

	account := keeper.GetAccount(ctx, params.ID)
	if account == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", params.ID)
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
	store := store.NewStore(ctx, ak.key)

	bz := store.Get(types.AuthSeqStoreKey(auth))
	if bz == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "auth data no found %s", auth)
	}

	var an types.Auth
	if err := ak.cdc.UnmarshalBinaryBare(bz, &an); err != nil {
		return nil, sdkerrors.Wrapf(chainTypes.ErrKuMsgDataUnmarshal,
			"query auth data unmarshal by %s by %s", auth, err.Error())
	}

	jsonBz, err := codec.MarshalJSONIndent(ak.cdc, an)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return jsonBz, nil
}

func queryAccountsByAuth(ctx sdk.Context, req abci.RequestQuery, ak AccountKeeper) ([]byte, error) {
	var params types.QueryAccountsByAuthParams
	if err := ak.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	accounts := ak.GetAccountsByAuth(ctx, params.Auth)

	bz, err := codec.MarshalJSONIndent(ak.cdc, accounts)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
