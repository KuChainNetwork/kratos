package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"
)

// AppGenesisState type for app 's genesis state, every k-v is for a module
type AppGenesisState map[string]json.RawMessage

// UnmarshalGenesis unmarshal a genesis state for a module in app genesis
func (s AppGenesisState) UnmarshalGenesis(cdc *codec.Codec, key string, val interface{}) error {
	raw, ok := s[key]
	if !ok {
		return fmt.Errorf("key no found for %s", key)
	}

	return cdc.UnmarshalJSON(raw, val)
}

// MarshalGenesis marshals the genesis state for a module in app genesis
func (s AppGenesisState) MarshalGenesis(cdc *codec.Codec, key string, val interface{}) error {
	raw, err := cdc.MarshalJSON(val)
	if err != nil {
		return err
	}

	_, ok := s[key]
	if !ok {
		return fmt.Errorf("key no found for %s", key)
	}

	s[key] = raw
	return nil
}

// LoadGenesisFile reads and unmarshals GenesisDoc from the given file.
func LoadGenesisFile(cdc *codec.Codec, genFile string) (genDoc tmtypes.GenesisDoc, err error) {
	if !tmos.FileExists(genFile) {
		err = fmt.Errorf("%s does not exist", genFile)
		return
	}

	genContents, err := ioutil.ReadFile(genFile)
	if err != nil {
		return genDoc, err
	}

	if err := cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
		return genDoc, err
	}

	return genDoc, err
}

// LoadGenesisStateFromBytes
func LoadGenesisStateFromBytes(cdc *codec.Codec, appState AppGenesisState, key string, val interface{}) error {
	return sdkerrors.Wrap(appState.UnmarshalGenesis(cdc, key, val), "unmarshal genesis error")
}

// LoadGenesisStateFromFile
func LoadGenesisStateFromFile(cdc *codec.Codec, genFile, key string, val interface{}) error {
	doc, err := LoadGenesisFile(cdc, genFile)
	if err != nil {
		return err
	}

	var appState AppGenesisState
	if err = cdc.UnmarshalJSON(doc.AppState, &appState); err != nil {
		return sdkerrors.Wrap(err, "unmarshal app state error")
	}

	return LoadGenesisStateFromBytes(cdc, appState, key, val)
}

// SaveGenesisStateToFile
func SaveGenesisStateToFile(cdc *codec.Codec, genFile, key string, val interface{}) error {
	doc, err := LoadGenesisFile(cdc, genFile)
	if err != nil {
		return err
	}

	var appState AppGenesisState
	if err = cdc.UnmarshalJSON(doc.AppState, &appState); err != nil {
		return err
	}

	if err := appState.MarshalGenesis(cdc, key, val); err != nil {
		return err
	}

	appStateJSON, err := cdc.MarshalJSON(appState)
	if err != nil {
		return err
	}

	doc.AppState = appStateJSON

	if err := doc.ValidateAndComplete(); err != nil {
		return err
	}

	return doc.SaveAs(genFile)
}
