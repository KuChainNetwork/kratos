package types

import (
	"fmt"
	chaintype "github.com/KuChain-io/kuchain/chain/types"
	"time"
)

func NewPunishValidator(validatorAccount chaintype.AccountID, height int64, untilTime time.Time, proposalID uint64) PunishValidator {
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
			validator.GetValidatorAccount().String(),
			validator.GetStartHeight(),
			validator.GetJailedUntil().String(),
			validator.GetMissedProposalID())
	}
	return out
}

func (p Punishvalidators) Equal(other Punishvalidators) bool {
	if len(p) != len(other) {
		return false
	}

	for i, proposal := range p {
		if !proposal.Equal(other[i]) {
			return false
		}
	}

	return true
}
