package types

import (
	"encoding/json"

	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/pkg/errors"
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
	gs := DefaultGenesisState()
	if err := ModuleCdc.UnmarshalJSON(bz, &gs); err != nil {
		return errors.Errorf("failed to unmarshal %s genesis state: %s", ModuleName, err.Error())
	}

	return ValidateGenesis(gs)
}

func ValidateGenesis(g GenesisState) error {
	// check account numbers
	if int(g.AccountNumber) != len(g.Auths) {
		return errors.Errorf("account number not match %d to %d", g.AccountNumber, len(g.Auths))
	}

	for i, a := range g.Auths {
		if i != int(a.GetNumber()) {
			return errors.Errorf("account %s number sort not match %d", a.GetAddress(), i)
		}
	}

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
