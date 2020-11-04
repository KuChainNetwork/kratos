package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/store"
	"github.com/KuChainNetwork/kuchain/x/account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"
)

// InitAuthData
func (ak AccountKeeper) InitAuthData(ctx sdk.Context, auth AccAddress) {
	authData := types.NewAuth(auth)
	ak.initAuthData(ctx, &authData)
	ak.setAuthData(ctx, auth, authData)
}

func (ak AccountKeeper) EnsureAuthInited(ctx sdk.Context, auth AccAddress) {
	if !ak.isAuthExist(ctx, auth) {
		ak.InitAuthData(ctx, auth)
	}
}

// GetAuthSequence
func (ak AccountKeeper) GetAuthSequence(ctx sdk.Context, auth AccAddress) (uint64, uint64, error) {
	if ctx.BlockHeight() == 0 {
		// for genesis
		return 0, 0, nil
	}

	authData := ak.getAuthData(ctx, auth)
	return authData.GetSequence(), authData.GetNumber(), nil
}

// IncAuthSequence
func (ak AccountKeeper) IncAuthSequence(ctx sdk.Context, auth AccAddress) {
	authData := ak.getAuthData(ctx, auth)
	authData.SetSequence(authData.GetSequence() + 1)
	ak.setAuthData(ctx, auth, authData)
}

// SetPubKey
func (ak AccountKeeper) SetPubKey(ctx sdk.Context, auth AccAddress, pubKey crypto.PubKey) {
	authData := ak.getAuthData(ctx, auth)
	authData.SetPubKey(pubKey)
	ak.setAuthData(ctx, auth, authData)
}

func (ak AccountKeeper) initAuthData(ctx sdk.Context, auth *types.Auth) {
	auth.SetAccountNum(ak.GetNextAccountNumber(ctx))
}

// getAuthData
func (ak AccountKeeper) getAuthData(ctx sdk.Context, auth AccAddress) types.Auth {
	store := store.NewStore(ctx, ak.key)

	bz := store.Get(types.AuthSeqStoreKey(auth))
	if bz == nil {
		authData := types.NewAuth(auth)
		ak.initAuthData(ctx, &authData)
		ak.setAuthData(ctx, auth, authData)
		return authData
	}

	var an types.Auth
	if err := ak.cdc.UnmarshalBinaryBare(bz, &an); err != nil {
		panic(errors.Wrap(err, "get auth data unmarshal"))
	}

	return an
}

// setAuthData
func (ak AccountKeeper) setAuthData(ctx sdk.Context, auth AccAddress, data types.Auth) {
	store := store.NewStore(ctx, ak.key)

	bz, err := ak.cdc.MarshalBinaryBare(data)
	if err != nil {
		panic(err)
	}

	store.Set(types.AuthSeqStoreKey(auth), bz)
}

// isAuthExist
func (ak AccountKeeper) isAuthExist(ctx sdk.Context, auth AccAddress) bool {
	return store.NewStore(ctx, ak.key).Has(types.AuthSeqStoreKey(auth))
}
