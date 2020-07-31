package genesis

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	emptyJSONRaw = json.RawMessage("{}")
)

type EmptyGenesisModuleBasicBase struct {
}

// DefaultGenesis returns default genesis state as raw bytes for the account module.
func (EmptyGenesisModuleBasicBase) DefaultGenesis() json.RawMessage {
	return nil
}

// ValidateGenesis performs genesis state validation for the account module.
func (EmptyGenesisModuleBasicBase) ValidateGenesis(json.RawMessage) error {
	return nil
}

type EmptyGenesisModuleBase struct {
}

// InitGenesis performs genesis initialization for the account module. It returns no validator updates.
func (EmptyGenesisModuleBase) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the account module.
func (EmptyGenesisModuleBase) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return emptyJSONRaw
}
