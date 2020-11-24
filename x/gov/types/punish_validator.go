package types

import (
	"fmt"
	"time"
)

type PunishValidator struct {
	ValidatorAccount AccountID `json:"account_id" yaml:"account_id"`
	StartHeight      int64     `json:"start_height,omitempty" yaml:"start_height"` // height at which validator was first a candidate OR was unjailed
	JailedUntil      time.Time `json:"jailed_until" yaml:"jailed_until"`
	MissedProposalID uint64    `json:"missed_proposal_id" yaml:"missed_proposal_id"`
}

func (p PunishValidator) Equal(other PunishValidator) bool {
	return p.ValidatorAccount.Eq(other.ValidatorAccount) &&
		p.StartHeight == other.StartHeight &&
		p.JailedUntil.Equal(other.JailedUntil) &&
		p.MissedProposalID == other.MissedProposalID
}

func NewPunishValidator(validatorAccount AccountID, height int64, untilTime time.Time, proposalID uint64) PunishValidator {
	return PunishValidator{
		ValidatorAccount: validatorAccount,
		StartHeight:      height,
		JailedUntil:      untilTime,
		MissedProposalID: proposalID,
	}
}

type Punishvalidators []PunishValidator

func (v Punishvalidators) String() string {
	if len(v) == 0 {
		return "[]"
	}
	out := "Punished Validators"
	for _, validator := range v {
		out += fmt.Sprintf("\n  %s %d %s %d",
			validator.ValidatorAccount.String(),
			validator.StartHeight,
			validator.JailedUntil.String(),
			validator.MissedProposalID)
	}
	return out
}

func (v Punishvalidators) Equal(other Punishvalidators) bool {
	if len(v) != len(other) {
		return false
	}

	for i, proposal := range v {
		if !proposal.Equal(other[i]) {
			return false
		}
	}

	return true
}
