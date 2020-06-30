package account

import (
	"github.com/KuChain-io/kuchain/x/account/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis account genesis init
func InitGenesis(ctx sdk.Context, ak Keeper, data GenesisState) {
	logger := ak.Logger(ctx)

	for _, a := range data.Accounts {
		logger.Info("init genesis account", "name", a.GetName(), "auth", a.GetAuth())
		ak.SetAccount(ctx, ak.NewAccount(ctx, a))

		// ensure auth init
		ak.EnsureAuthInited(ctx, a.GetAuth())
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, ak Keeper) GenesisState {
	var genAccounts exported.GenesisAccounts
	ak.IterateAccounts(ctx, func(account exported.Account) bool {
		genAccounts = append(genAccounts, account.(exported.GenesisAccount))
		return false
	})

	return GenesisState{
		Accounts: genAccounts,
	}
}
