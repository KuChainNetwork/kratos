package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

// ValidatorGovInfo used for tallying
type ValidatorGovInfo struct {
	Address             AccountID  // address of the validator operator
	BondedTokens        sdk.Int    // Power of a Validator
	DelegatorShares     sdk.Dec    // Total outstanding delegator shares
	DelegatorDeductions sdk.Dec    // Delegator deductions from validator's delegators voting independently
	Vote                VoteOption // Vote of the validator
}

// NewValidatorGovInfo creates a ValidatorGovInfo instance
func NewValidatorGovInfo(address AccountID, bondedTokens sdk.Int, delegatorShares,
	delegatorDeductions sdk.Dec, vote VoteOption) ValidatorGovInfo {

	return ValidatorGovInfo{
		Address:             address,
		BondedTokens:        bondedTokens,
		DelegatorShares:     delegatorShares,
		DelegatorDeductions: delegatorDeductions,
		Vote:                vote,
	}
}

// TallyResult defines a standard tally for a proposal
type TallyResult struct {
	Yes        sdk.Int `json:"yes" yaml:"yes"`
	Abstain    sdk.Int `json:"abstain" yaml:"abstain"`
	No         sdk.Int `json:"no" yaml:"no"`
	NoWithVeto sdk.Int `json:"no_with_veto" yaml:"no_with_veto"`
}

func (tr TallyResult) Equal(other TallyResult) bool {
	return tr.Yes.Equal(other.Yes) && tr.No.Equal(other.No) && tr.Abstain.Equal(other.Abstain) && tr.NoWithVeto.Equal(other.NoWithVeto)
}

// NewTallyResult creates a new TallyResult instance
func NewTallyResult(yes, abstain, no, noWithVeto sdk.Int) TallyResult {
	return TallyResult{
		Yes:        yes,
		Abstain:    abstain,
		No:         no,
		NoWithVeto: noWithVeto,
	}
}

// NewTallyResultFromMap creates a new TallyResult instance from a Option -> Dec map
func NewTallyResultFromMap(results map[VoteOption]sdk.Dec) TallyResult {
	return NewTallyResult(
		results[OptionYes].TruncateInt(),
		results[OptionAbstain].TruncateInt(),
		results[OptionNo].TruncateInt(),
		results[OptionNoWithVeto].TruncateInt(),
	)
}

// EmptyTallyResult returns an empty TallyResult.
func EmptyTallyResult() TallyResult {
	return NewTallyResult(sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt())
}

// Equals returns if two proposals are equal.
func (tr TallyResult) Equals(comp TallyResult) bool {
	return tr.Yes.Equal(comp.Yes) &&
		tr.Abstain.Equal(comp.Abstain) &&
		tr.No.Equal(comp.No) &&
		tr.NoWithVeto.Equal(comp.NoWithVeto)
}

// String implements stringer interface
func (tr TallyResult) String() string {
	out, _ := yaml.Marshal(tr)
	return string(out)
}
