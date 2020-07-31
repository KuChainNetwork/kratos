package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/KuChainNetwork/kuchain/x/evidence/exported"
)

// DONTCOVER

// GenesisState defines the evidence module's genesis state.
type GenesisState struct {
	Params   Params              `json:"params" yaml:"params"`
	Evidence []exported.Evidence `json:"evidence" yaml:"evidence"`
}

func NewGenesisState(p Params, e []exported.Evidence) GenesisState {
	return GenesisState{
		Params:   p,
		Evidence: e,
	}
}

// DefaultGenesisState returns the evidence module's default genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:   DefaultParams(),
		Evidence: []exported.Evidence{},
	}
}

// ValidateGenesis performs basic gensis state validation returning an error upon any
// failure.
func (g GenesisState) ValidateGenesis(bz json.RawMessage) error {
	gs := DefaultGenesisState()

	if err := ModuleCdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	return gs.Validate()
}

func (gs GenesisState) Validate() error {
	for _, e := range gs.Evidence {
		if err := e.ValidateBasic(); err != nil {
			return err
		}
	}

	maxEvidence := gs.Params.MaxEvidenceAge
	if maxEvidence < 1*time.Minute {
		return fmt.Errorf("max evidence age must be at least 1 minute, is %s", maxEvidence.String())
	}

	return nil
}
