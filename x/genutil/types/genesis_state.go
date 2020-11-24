package types

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	"github.com/cosmos/cosmos-sdk/codec"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"
)

// GenesisState defines the raw genesis transaction in JSON
type GenesisState struct {
	GenTxs []json.RawMessage `json:"gentxs" yaml:"gentxs"`

	stakingFuncManager StakingFuncManager
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(genTxs []json.RawMessage) GenesisState {
	// Ensure genTxs is never nil
	if len(genTxs) == 0 {
		genTxs = make([]json.RawMessage, 0)
	}
	return GenesisState{
		GenTxs: genTxs,
	}
}

// DefaultGenesisState returns the genutil module's default genesis state.
func DefaultGenesisState(stakingFuncManager StakingFuncManager) GenesisState {
	return GenesisState{
		GenTxs:             []json.RawMessage{},
		stakingFuncManager: stakingFuncManager,
	}
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func (g GenesisState) ValidateGenesis(bz json.RawMessage) error {
	gs := DefaultGenesisState(g.stakingFuncManager)
	if err := ModuleCdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	return ValidateGenesis(gs, g.stakingFuncManager)
}

// NewGenesisStateFromStdTx creates a new GenesisState object
// from auth transactions
func NewGenesisStateFromStdTx(genTxs []txutil.StdTx) GenesisState {
	genTxsBz := make([]json.RawMessage, len(genTxs))
	for i, genTx := range genTxs {
		genTxsBz[i] = ModuleCdc.MustMarshalJSON(genTx)
	}
	return NewGenesisState(genTxsBz)
}

// GetGenesisStateFromAppState gets the genutil genesis state from the expected app state
func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}

// SetGenesisStateInAppState sets the genutil genesis state within the expected app state
func SetGenesisStateInAppState(cdc *codec.Codec,
	appState map[string]json.RawMessage, genesisState GenesisState) map[string]json.RawMessage {
	genesisStateBz := cdc.MustMarshalJSON(genesisState)
	appState[ModuleName] = genesisStateBz
	return appState
}

// GenesisStateFromGenDoc creates the core parameters for genesis initialization
// for the application.
//
// NOTE: The pubkey input is this machines pubkey.
func GenesisStateFromGenDoc(
	cdc *codec.Codec, genDoc tmtypes.GenesisDoc) (genesisState map[string]json.RawMessage, err error) {
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return genesisState, err
	}
	return genesisState, nil
}

// GenesisStateFromGenFile creates the core parameters for genesis initialization
// for the application.
//
// NOTE: The pubkey input is this machines pubkey.
func GenesisStateFromGenFile(
	cdc *codec.Codec, genFile string) (genesisState map[string]json.RawMessage, genDoc *tmtypes.GenesisDoc, err error) {
	if !tmos.FileExists(genFile) {
		return genesisState, genDoc,
			fmt.Errorf("%s does not exist, run `init` first", genFile)
	}
	genDoc, err = tmtypes.GenesisDocFromFile(genFile)
	if err != nil {
		return genesisState, genDoc, err
	}

	genesisState, err = GenesisStateFromGenDoc(cdc, *genDoc)
	return genesisState, genDoc, err
}

// ValidateGenesis validates GenTx transactions
func ValidateGenesis(genesisState GenesisState, stakingFuncManager StakingFuncManager) error {
	for i, genTx := range genesisState.GenTxs {
		var tx txutil.StdTx
		if err := ModuleCdc.UnmarshalJSON(genTx, &tx); err != nil {
			return err
		}

		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return errors.New(
				"must provide genesis StdTx with exactly 1 CreateValidator message")
		}

		if stakingFuncManager.MsgCreateValidatorVal(msgs[0]) {
			return fmt.Errorf(
				"genesis transaction %v does not contain a MsgCreateValidator", i)
		}
	}
	return nil
}
