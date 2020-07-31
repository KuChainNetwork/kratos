package genesis

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
)

type GenesisData interface {
	ValidateGenesis(json.RawMessage) error
}

type ModuleBasicBase struct {
	defaultGenesis GenesisData
	cdc            *codec.Codec
}

// DefaultGenesis returns default genesis state as raw bytes for the account module.
func (g ModuleBasicBase) DefaultGenesis() json.RawMessage {
	if g.cdc == nil {
		return nil
	}

	return g.cdc.MustMarshalJSON(g.defaultGenesis)
}

// ValidateGenesis performs genesis state validation for the account module.
func (g ModuleBasicBase) ValidateGenesis(bz json.RawMessage) error {
	if g.cdc == nil {
		return nil
	}

	return g.defaultGenesis.ValidateGenesis(bz)
}

func NewModuleBasicBase(cdc *codec.Codec, defaultGenesis GenesisData) ModuleBasicBase {
	return ModuleBasicBase{
		defaultGenesis: defaultGenesis,
		cdc:            cdc,
	}
}
