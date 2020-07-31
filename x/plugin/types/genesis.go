package types

import "encoding/json"

type GenesisState struct {
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState()
}

func (gs GenesisState) ValidateGenesis(_ json.RawMessage) error {
	return nil
}
