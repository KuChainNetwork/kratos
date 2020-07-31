package types

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// Deposit defines an amount deposited by an account address to an active proposal
type Deposit struct {
	ProposalID uint64    `json:"proposal_id,omitempty" yaml:"proposal_id"`
	Depositor  AccountID `json:"depositor" yaml:"depositor"`
	Amount     Coins     `json:"amount" yaml:"amount"`
}

func (d Deposit) Equal(other Deposit) bool {
	return d.ProposalID == other.ProposalID &&
		d.Depositor.Eq(other.Depositor) &&
		d.Amount.IsEqual(other.Amount)
}

// NewDeposit creates a new Deposit instance
func NewDeposit(proposalID uint64, depositor AccountID, amount Coins) Deposit {
	return Deposit{proposalID, depositor, amount}
}

func (d Deposit) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

// Deposits is a collection of Deposit objects
type Deposits []Deposit

// Equal returns true if two slices (order-dependant) of deposits are equal.
func (d Deposits) Equal(other Deposits) bool {
	if len(d) != len(other) {
		return false
	}

	for i, deposit := range d {
		if !deposit.Equal(other[i]) {
			return false
		}
	}

	return true
}

func (d Deposits) String() string {
	if len(d) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Deposits for Proposal %d:", d[0].ProposalID)
	for _, dep := range d {
		out += fmt.Sprintf("\n  %s: %s", dep.Depositor, dep.Amount)
	}
	return out
}

// Empty returns whether a deposit is empty.
func (d Deposit) Empty() bool {
	return d.Equal(Deposit{})
}
