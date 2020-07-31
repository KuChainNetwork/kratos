package types

import (
	"encoding/json"

	"github.com/KuChainNetwork/kuchain/x/account/exported"
)

// GenesisState genesis state for account module
type GenesisState struct {
	Accounts exported.GenesisAccounts `json:"accounts"`
}

func (g GenesisState) ValidateGenesis(bz json.RawMessage) error {
	return nil
}

// DefaultGenesisState get default genesis state for account module
func DefaultGenesisState() GenesisState {
	res := GenesisState{
		Accounts: exported.GenesisAccounts{},
	}

	return res
}

// NewGenesisState new genesis state by genesis accounts, for test
func NewGenesisState(accs []exported.GenesisAccount) GenesisState {
	return GenesisState{
		Accounts: accs,
	}
}
