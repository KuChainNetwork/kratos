package types

import (
	"encoding/json"
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/types"
)

// GenesisState is the supply state that must be provided at genesis.
type GenesisState struct {
	Supply types.Coins `json:"supply" yaml:"supply"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(supply types.Coins) GenesisState {
	return GenesisState{supply}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultSupply().GetTotal())
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func (g GenesisState) ValidateGenesis(bz json.RawMessage) error {
	gs := DefaultGenesisState()
	if err := ModuleCdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	return NewSupply(gs.Supply).ValidateBasic()
}
