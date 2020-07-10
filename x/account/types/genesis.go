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

	return res
}
