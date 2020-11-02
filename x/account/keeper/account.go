package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/store"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/KuChainNetwork/kuchain/x/account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewAccountByName create a account struct by name.
func (ak AccountKeeper) NewAccountByName(ctx sdk.Context, n Name) exported.Account {
	acc := ak.proto()
	if err := acc.SetName(n); err != nil {
		panic(err)
	}
	return ak.NewAccount(ctx, acc)
}

// NewAccount sets the next account number to a given account interface
func (ak AccountKeeper) NewAccount(ctx sdk.Context, acc exported.Account) exported.Account {
	// FIXME: update account number
	return acc
}

// GetAccount get account from keeper
func (ak AccountKeeper) GetAccount(ctx sdk.Context, id AccountID) exported.Account {
	store := store.NewStore(ctx, ak.key)
	bz := store.Get(types.AccountIDStoreKey(id))
	if bz == nil {
		return nil
	}
	acc := ak.decodeAccount(bz)

	return acc
}

// GetAccount get account from keeper
func (ak AccountKeeper) GetAccountByName(ctx sdk.Context, name Name) exported.Account {
	return ak.GetAccount(ctx, chainTypes.NewAccountIDFromName(name))
}

// SetAccount implements sdk.AccountKeeper.
func (ak AccountKeeper) SetAccount(ctx sdk.Context, acc exported.Account) {
	n := acc.GetID()
	store := store.NewStore(ctx, ak.key)

	bz, err := ak.cdc.MarshalBinaryBare(acc)
	if err != nil {
		panic(err)
	}

	store.Set(types.AccountIDStoreKey(n), bz)
}

// EnsureAccount ensure account is exist, if not create a account with init data
func (ak AccountKeeper) EnsureAccount(ctx sdk.Context, id AccountID) error {
	if accAddress, ok := id.ToAccAddress(); ok {
		// init auth
		ak.EnsureAuthInited(ctx, accAddress)
		return nil
	}

	if !ak.isAccountExist(ctx, id) {
		return sdkerrors.Wrapf(types.ErrAccountNoFound, "ensure account no exit %s", id.String())
	}

	return nil
}

// IterateAccounts iterates over all the stored accounts and performs a callback function
func (ak AccountKeeper) IterateAccounts(ctx sdk.Context, cb func(account exported.Account) (stop bool)) {
	store := store.NewStore(ctx, ak.key)
	iterator := sdk.KVStorePrefixIterator(store, types.AddressStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		account := ak.decodeAccount(iterator.Value())

		if cb(account) {
			break
		}
	}
}

// decodeAccount decode account by cdc
func (ak AccountKeeper) decodeAccount(bz []byte) (acc exported.Account) {
	err := ak.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}

// isAccountExist is account exist from keeper
func (ak AccountKeeper) isAccountExist(ctx sdk.Context, id AccountID) bool {
	return store.NewStore(ctx, ak.key).Has(types.AccountIDStoreKey(id))
}
func (ak AccountKeeper) IsAccountExist(ctx sdk.Context, id AccountID) bool {
	return ak.isAccountExist(ctx, id)
}
