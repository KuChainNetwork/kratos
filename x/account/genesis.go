package account

import (
	"encoding/json"
	"sort"

	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/KuChainNetwork/kuchain/x/account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis account genesis init
func InitGenesis(ctx sdk.Context, ak Keeper, data json.RawMessage) {
	logger := ak.Logger(ctx)

	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)

	for _, a := range genesisState.Accounts {
		logger.Info("init genesis account", "name", a.GetName(), "auth", a.GetAuth())
		ak.SetAccount(ctx, ak.NewAccount(ctx, a))

		// ensure auth init
		ak.EnsureAuthInited(ctx, a.GetAuth())
		if _, ok := a.GetID().ToName(); ok {
			ak.AddAccountByAuth(ctx, a.GetAuth(), a.GetName().String())
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, ak Keeper) GenesisState {
	var genAccounts exported.GenesisAccounts
	ak.IterateAccounts(ctx, func(account exported.Account) bool {
		genAccounts = append(genAccounts, account.(exported.GenesisAccount))
		return false
	})

	genAuths := exported.GenesisAuths(make([]exported.GenesisAuth, 0, 5120))
	ak.IterateAuths(ctx, func(auth types.Auth) bool {
		genAuths = append(genAuths, types.GenesisAuth{Auth: auth})
		return false
	})
	sort.Stable(genAuths)

	currNextAccountNumber := ak.QueryCurrNextAccountNumber(ctx)

	return GenesisState{
		Accounts:      genAccounts,
		Auths:         genAuths,
		AccountNumber: currNextAccountNumber,
	}
}
