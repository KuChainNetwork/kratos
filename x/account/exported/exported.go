package exported

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountAuthKeeper is interface for trx auth and account state.
type AccountAuthKeeper interface {
	GetAuthSequence(sdk.Context, types.AccAddress) (uint64, uint64, error)
	IncAuthSequence(sdk.Context, types.AccAddress)
}

// AccountStatKeeper is interface for other modules to get account state.
type AccountStatKeeper interface {
	GetAccount(sdk.Context, types.AccountID) Account // can return nil.
	IterateAccounts(ctx sdk.Context, cb func(account Account) (stop bool))

	GetNextAccountNumber(ctx sdk.Context) uint64
}

type AuthAccountKeeper interface {
	GetAccountsByAuth(sdk.Context, types.AccAddress) []string
	AddAccountByAuth(sdk.Context, types.AccAddress, string)
}

// Account is a interface for kuchain account and address,
// kuchain support both cosmos address and eos-likely account.
type Account interface {
	GetName() types.Name
	SetName(types.Name) error

	GetID() types.AccountID
	SetID(id types.AccountID) error

	GetAuth() types.AccAddress
	SetAuth(types.AccAddress) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	// Ensure that account implements stringer
	String() string
}

// GenesisAccounts defines a slice of GenesisAccount objects
type GenesisAccounts []GenesisAccount

// Contains returns true if the given address exists in a slice of GenesisAccount
// objects.
func (ga GenesisAccounts) Contains(accName types.Name) bool {
	for _, acc := range ga {
		if acc.GetName().Eq(accName) {
			return true
		}
	}

	return false
}

// Append append a account to genesis
func (ga GenesisAccounts) Append(acc GenesisAccount) GenesisAccounts {
	if ga.Contains(acc.GetName()) {
		panic(fmt.Errorf("account %s has put into genesis account", acc.GetName()))
	}

	if err := acc.SetAccountNumber(uint64(len(ga) + 1)); err != nil {
		panic(err)
	}

	return append(ga, acc)
}

// GenesisAccount defines a genesis account that embeds an Account with validation capabilities.
type GenesisAccount interface {
	Account
	Validate() error
}
