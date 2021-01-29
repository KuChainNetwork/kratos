package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

// AppGenesisState type for app 's genesis state, every k-v is for a module
type AppGenesisState map[string]json.RawMessage

// UnmarshalGenesis unmarshal a genesis state for a module in app genesis
func (s AppGenesisState) UnmarshalGenesis(cdc *codec.LegacyAmino, key string, val interface{}) error {
	raw, ok := s[key]
	if !ok {
		return fmt.Errorf("key no found for %s", key)
	}

	return cdc.UnmarshalJSON(raw, val)
}

// MarshalGenesis marshals the genesis state for a module in app genesis
func (s AppGenesisState) MarshalGenesis(cdc *codec.LegacyAmino, key string, val interface{}) error {
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
func LoadGenesisFile(cdc *codec.LegacyAmino, genFile string) (genDoc tmtypes.GenesisDoc, err error) {
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
func LoadGenesisStateFromBytes(cdc *codec.LegacyAmino, appState AppGenesisState, key string, val interface{}) error {
	return sdkerrors.Wrap(appState.UnmarshalGenesis(cdc, key, val), "unmarshal genesis error")
}

// LoadGenesisStateFromFile
func LoadGenesisStateFromFile(cdc *codec.LegacyAmino, genFile, key string, val interface{}) error {
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
func SaveGenesisStateToFile(cdc *codec.LegacyAmino, genFile, key string, val interface{}) error {
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

// GenesisValidator is an initial validator.
type GenesisValidator struct {
	ID      AccountID       `json:"account"`
	Address tmtypes.Address `json:"address"`
	PubKey  crypto.PubKey   `json:"pub_key"`
	Power   int64           `json:"power"`
	Name    string          `json:"name"`
}

// GenesisDoc defines the initial conditions for a tendermint blockchain, in particular its validator set.
type GenesisDoc struct {
	GenesisTime     time.Time                `json:"genesis_time"`
	ChainID         string                   `json:"chain_id"`
	ConsensusParams *tmproto.ConsensusParams `json:"consensus_params,omitempty"`
	Validators      []GenesisValidator       `json:"validators,omitempty"`
	AppHash         tmbytes.HexBytes         `json:"app_hash"`
	AppState        json.RawMessage          `json:"app_state,omitempty"`
}

// SaveAs is a utility method for saving GenensisDoc as a JSON file.
func (genDoc *GenesisDoc) SaveAs(file string) error {
	genDocBytes, err := tmjson.MarshalIndent(genDoc, "", "  ")
	if err != nil {
		return err
	}
	return tmos.WriteFile(file, genDocBytes, 0644)
}

// ValidatorHash returns the hash of the validator set contained in the GenesisDoc
func (genDoc *GenesisDoc) ValidatorHash() []byte {
	vals := make([]*tmtypes.Validator, len(genDoc.Validators))
	for i, v := range genDoc.Validators {
		vals[i] = tmtypes.NewValidator(v.PubKey, v.Power)
	}
	vset := tmtypes.NewValidatorSet(vals)
	return vset.Hash()
}

// ValidateAndComplete checks that all necessary fields are present
// and fills in defaults for optional fields left empty
func (genDoc *GenesisDoc) ValidateAndComplete() error {
	if genDoc.ChainID == "" {
		return errors.New("genesis doc must include non-empty chain_id")
	}
	if len(genDoc.ChainID) > tmtypes.MaxChainIDLen {
		return errors.Errorf("chain_id in genesis doc is too long (max: %d)", tmtypes.MaxChainIDLen)
	}

	if genDoc.ConsensusParams == nil {
		genDoc.ConsensusParams = tmtypes.DefaultConsensusParams()
	} else if err := tmtypes.ValidateConsensusParams(*genDoc.ConsensusParams); err != nil {
		return err
	}

	for i, v := range genDoc.Validators {
		if v.Power == 0 {
			return errors.Errorf("the genesis file cannot contain validators with no voting power: %v", v)
		}
		if len(v.Address) > 0 && !bytes.Equal(v.PubKey.Address(), v.Address) {
			return errors.Errorf("incorrect address for validator %v in the genesis file, should be %v", v, v.PubKey.Address())
		}
		if len(v.Address) == 0 {
			genDoc.Validators[i].Address = v.PubKey.Address()
		}
	}

	if genDoc.GenesisTime.IsZero() {
		genDoc.GenesisTime = tmtime.Now()
	}

	return nil
}

//------------------------------------------------------------
// Make genesis state from file

// GenesisDocFromJSON unmarshalls JSON data into a GenesisDoc.
func GenesisDocFromJSON(jsonBlob []byte) (*GenesisDoc, error) {
	genDoc := GenesisDoc{}
	err := tmjson.Unmarshal(jsonBlob, &genDoc)
	if err != nil {
		return nil, err
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return nil, err
	}

	return &genDoc, err
}

// GenesisDocFromFile reads JSON data from a file and unmarshalls it into a GenesisDoc.
func GenesisDocFromFile(genDocFile string) (*GenesisDoc, error) {
	jsonBlob, err := ioutil.ReadFile(genDocFile)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read GenesisDoc file")
	}
	genDoc, err := GenesisDocFromJSON(jsonBlob)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error reading GenesisDoc at %v", genDocFile))
	}
	return genDoc, nil
}
