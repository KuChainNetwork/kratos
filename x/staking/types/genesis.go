package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	Params               Params                `json:"params" yaml:"params"`
	LastTotalPower       sdk.Int               `json:"last_total_power" yaml:"last_total_power"`
	LastValidatorPowers  []LastValidatorPower  `json:"last_validator_powers" yaml:"last_validator_powers"`
	Validators           Validators            `json:"validators" yaml:"validators"`
	Delegations          Delegations           `json:"delegations" yaml:"delegations"`
	UnbondingDelegations []UnbondingDelegation `json:"unbonding_delegations" yaml:"unbonding_delegations"`
	Redelegations        []Redelegation        `json:"redelegations" yaml:"redelegations"`
	Exported             bool                  `json:"exported" yaml:"exported"`
}

// LastValidatorPower required for validator set update logic
type LastValidatorPower struct {
	ID    AccountID `json:"accountID"`
	Power int64
}

// NewGenesisState creates a new GenesisState instanc e
func NewGenesisState(params Params, validators []Validator, delegations []Delegation) GenesisState {
	return GenesisState{
		Params:      params,
		Validators:  validators,
		Delegations: delegations,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func (g GenesisState) ValidateGenesis(bz json.RawMessage) error {
	gs := DefaultGenesisState()
	if err := Cdc().UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	return ValidateGenesis(gs)
}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	err := validateGenesisStateValidators(data.Validators)
	if err != nil {
		return err
	}
	err = data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}

func validateGenesisStateValidators(validators []Validator) (err error) {
	addrMap := make(map[string]bool, len(validators))
	for i := 0; i < len(validators); i++ {
		val := validators[i]
		strKey := string(val.GetConsPubKey().Bytes())

		if _, ok := addrMap[strKey]; ok {
			return fmt.Errorf("duplicate validator in genesis state: moniker %v, address %v", val.Description.Moniker, val.GetConsAddr())
		}
		if val.Jailed && val.IsBonded() {
			return fmt.Errorf("validator is bonded and jailed in genesis state: moniker %v, address %v", val.Description.Moniker, val.GetConsAddr())
		}
		if val.DelegatorShares.IsZero() && !val.IsUnbonding() {
			return fmt.Errorf("bonded/unbonded genesis validator cannot have zero delegator shares, validator: %v", val)
		}

		addrMap[strKey] = true
	}

	return
}
