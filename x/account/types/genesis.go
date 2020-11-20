package types

import (
	"encoding/json"

	"github.com/KuChainNetwork/kuchain/x/account/exported"
)

type GenesisAuth struct {
	Auth
}

// GenesisState genesis state for account module
type GenesisState struct {
	Accounts      exported.GenesisAccounts `json:"accounts"`
	Auths         []exported.GenesisAuth   `json:"auths"`
	AccountNumber uint64                   `json:"account_num"`
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
