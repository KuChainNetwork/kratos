package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/store"
	chaintypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (ak AccountKeeper) GetAccountsByAuth(ctx sdk.Context, auth chaintypes.AccAddress) []string {
	var authAccounts types.AuthAccounts
	store := store.NewStore(ctx, ak.key)

	bz := store.Get(types.AuthAccountsStoreKey(auth))
	if bz == nil {
		return nil
	}

	if err := ak.cdc.UnmarshalBinaryBare(bz, &authAccounts); err != nil {
		panic(err)
	}

	return authAccounts.GetAccounts()
}

func (ak AccountKeeper) AddAccountByAuth(ctx sdk.Context, auth chaintypes.AccAddress, acc string) {
	var authAccounts types.AuthAccounts
	store := store.NewStore(ctx, ak.key)

	bz := store.Get(types.AuthAccountsStoreKey(auth))
	if bz == nil {
		authAccounts = types.NewAuthAccount(auth.String(), acc)
	} else {
		if err := ak.cdc.UnmarshalBinaryBare(bz, &authAccounts); err != nil {
			panic(err)
		}

		authAccounts.AddAccount(acc)
	}

	bz, err := ak.cdc.MarshalBinaryBare(authAccounts)
	if err != nil {
		panic(err)
	}

	store.Set(types.AuthAccountsStoreKey(auth), bz)
}

func (ak AccountKeeper) DeleteAccountByAuth(ctx sdk.Context, auth chaintypes.AccAddress, acc string) {
	var authAccounts types.AuthAccounts
	store := store.NewStore(ctx, ak.key)

	bz := store.Get(types.AuthAccountsStoreKey(auth))
	if bz == nil {
		return
	}

	if err := ak.cdc.UnmarshalBinaryBare(bz, &authAccounts); err != nil {
		panic(err)
	}

	authAccounts.DeleteAccount(acc)

	bz, err := ak.cdc.MarshalBinaryBare(authAccounts)
	if err != nil {
		panic(err)
	}

	store.Set(types.AuthAccountsStoreKey(auth), bz)
}
