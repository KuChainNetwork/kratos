package types

import (
	"github.com/KuChainNetwork/kuchain/x/account/exported"
)

// GenesisState genesis state for account module
type GenesisState struct {
	Accounts exported.GenesisAccounts `json:"accounts"`
}

// DefaultGenesisState get default genesis state for account module
func DefaultGenesisState() GenesisState {
	res := GenesisState{
		Accounts: exported.GenesisAccounts{},
	}

	// TODO: add default root account
	// Fix genesis accounts
	/*
		res.Accounts = res.Accounts.Append(
			NewKuAccountByName(types.MustName("kuchain"), types.MustAccAddressFromBech32("kuchain1xmc2z728py4gtwpc7jgytsan0282ww883qtv07"), 1),
		)
	*/

	return res
}
