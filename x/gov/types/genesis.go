package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	StartingProposalID uint64            `json:"starting_proposal_id" yaml:"starting_proposal_id"`
	Deposits           Deposits          `json:"deposits" yaml:"deposits"`
	Votes              Votes             `json:"votes" yaml:"votes"`
	Proposals          Proposals         `json:"proposals" yaml:"proposals"`
	DepositParams      DepositParams     `json:"deposit_params" yaml:"deposit_params"`
	VotingParams       VotingParams      `json:"voting_params" yaml:"voting_params"`
	TallyParams        TallyParams       `json:"tally_params" yaml:"tally_params"`
	PunishValidators   []PunishValidator `json:"punish_validators" yaml:"punish_validators"`
}

// NewGenesisState creates a new genesis state for the governance module
func NewGenesisState(startingProposalID uint64, dp DepositParams, vp VotingParams, tp TallyParams, pvs []PunishValidator) GenesisState {
	return GenesisState{
		StartingProposalID: startingProposalID,
		DepositParams:      dp,
		VotingParams:       vp,
		TallyParams:        tp,
		PunishValidators:   pvs,
	}
}

// DefaultGenesisState defines the default governance genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(
		DefaultStartingProposalID,
		DefaultDepositParams(),
		DefaultVotingParams(),
		DefaultTallyParams(),
		DefaultPunishValidators(),
	)
}

func (data GenesisState) Equal(other GenesisState) bool {
	return data.StartingProposalID == other.StartingProposalID &&
		data.Deposits.Equal(other.Deposits) &&
		data.Votes.Equal(other.Votes) &&
		data.Proposals.Equal(other.Proposals) &&
		data.DepositParams.Equal(other.DepositParams) &&
		data.TallyParams.Equal(other.TallyParams) &&
		data.VotingParams.Equal(other.VotingParams)
}

// IsEmpty returns true if a GenesisState is empty
func (data GenesisState) IsEmpty() bool {
	return data.Equal(GenesisState{})
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func (data GenesisState) ValidateGenesis(bz json.RawMessage) error {
	gs := DefaultGenesisState()
	if err := ModuleCdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	return ValidateGenesis(gs)
}

// ValidateGenesis checks if parameters are within valid ranges
func ValidateGenesis(data GenesisState) error {
	threshold := data.TallyParams.Threshold
	if threshold.IsNegative() || threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("governance vote threshold should be positive and less or equal to one, is %s",
			threshold.String())
	}

	veto := data.TallyParams.Veto
	if veto.IsNegative() || veto.GT(sdk.OneDec()) {
		return fmt.Errorf("governance vote veto threshold should be positive and less or equal to one, is %s",
			veto.String())
	}

	if !data.DepositParams.MinDeposit.IsValid() {
		return fmt.Errorf("governance deposit amount must be a valid Coins amount, is %s",
			data.DepositParams.MinDeposit.String())
	}

	return nil
}
