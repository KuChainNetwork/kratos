package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/KuChainNetwork/kuchain/chain/store"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/KuChainNetwork/kuchain/x/account/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ exported.AccountAuthKeeper = (*AccountKeeper)(nil)
var _ exported.AccountStatKeeper = (*AccountKeeper)(nil)
var _ exported.AuthAccountKeeper = (*AccountKeeper)(nil)

// AccountKeeper keeper for account module
type AccountKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	// The prototypical Account constructor.
	proto func() exported.Account
}

// NewAccountKeeper new account keeper
func NewAccountKeeper(cdc *codec.Codec, key sdk.StoreKey) AccountKeeper {
	return AccountKeeper{
		key:   key,
		proto: types.NewProtoKuAccount,
		cdc:   cdc,
	}
}

// Logger returns a module-specific logger.
func (ak AccountKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetNextAccountNumber returns and increments the global account number counter.
// If the global account number is not set, it initializes it with value 0.
func (ak AccountKeeper) GetNextAccountNumber(ctx sdk.Context) uint64 {
	var accNumber uint64
	store := store.NewStore(ctx, ak.key)
	bz := store.Get(types.GlobalAccountNumberKey)
	if bz == nil {
		// initialize the account numbers
		accNumber = 0
	} else {
		err := ak.cdc.UnmarshalBinaryLengthPrefixed(bz, &accNumber)
		if err != nil {
			panic(err)
		}
	}

	bz = ak.cdc.MustMarshalBinaryLengthPrefixed(accNumber + 1)
	store.Set(types.GlobalAccountNumberKey, bz)
	return accNumber
}

// GetAuth get auth for a account
func (ak AccountKeeper) GetAuth(ctx sdk.Context, account Name) (sdk.AccAddress, error) {
	acc := ak.GetAccount(ctx, chainTypes.NewAccountIDFromName(account))
	if acc == nil {
		return sdk.AccAddress{}, types.ErrAccountNoFound
	}

	return acc.GetAuth(), nil
}

func (ak AccountKeeper) GetStoreKey() sdk.StoreKey {
	return ak.key
}
