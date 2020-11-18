package types

import (
	"fmt"
	"time"

	"github.com/KuChainNetwork/kuchain/x/evidence/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	"gopkg.in/yaml.v2"
)

// Evidence type constants
const (
	RouteEquivocation = "equivocation"
	TypeEquivocation  = "equivocation"
)

var _ exported.Evidence = (*Equivocation)(nil)

// Equivocation implements the Evidence interface and defines evidence of double
// signing misbehavior.
type Equivocation struct {
	Height           int64           `json:"height,omitempty" yaml:"height"`
	Time             time.Time       `json:"time" yaml:"time"`
	Power            int64           `json:"power,omitempty" yaml:"power"`
	ConsensusAddress sdk.ConsAddress `json:"consensus_address,omitempty" yaml:"consensus_address"`
}

// Route returns the Evidence Handler route for an Equivocation type.
func (e Equivocation) Route() string { return RouteEquivocation }

// Type returns the Evidence Handler type for an Equivocation type.
func (e Equivocation) Type() string { return TypeEquivocation }

func (e Equivocation) String() string {
	bz, _ := yaml.Marshal(e)
	return string(bz)
}

// Hash returns the hash of an Equivocation object.
func (e Equivocation) Hash() tmbytes.HexBytes {
	return tmhash.Sum(EvidenceCdc.amino.MustMarshalBinaryBare(&e))
}

// ValidateBasic performs basic stateless validation checks on an Equivocation object.
func (e Equivocation) ValidateBasic() error {
	if e.Time.IsZero() {
		return fmt.Errorf("invalid equivocation time: %s", e.Time)
	}
	if e.Height < 1 {
		return fmt.Errorf("invalid equivocation height: %d", e.Height)
	}
	if e.Power < 1 {
		return fmt.Errorf("invalid equivocation validator power: %d", e.Power)
	}
	if e.ConsensusAddress.Empty() {
		return fmt.Errorf("invalid equivocation validator consensus address: %s", e.ConsensusAddress)
	}

	return nil
}

// GetConsensusAddress returns the validator's consensus address at time of the
// Equivocation infraction.
func (e Equivocation) GetConsensusAddress() sdk.ConsAddress {
	return e.ConsensusAddress
}

// GetHeight returns the height at time of the Equivocation infraction.
func (e Equivocation) GetHeight() int64 {
	return e.Height
}

// GetTime returns the time at time of the Equivocation infraction.
func (e Equivocation) GetTime() time.Time {
	return e.Time
}

// GetValidatorPower returns the validator's power at time of the Equivocation
// infraction.
func (e Equivocation) GetValidatorPower() int64 {
	return e.Power
}

// GetTotalPower is a no-op for the Equivocation type.
func (e Equivocation) GetTotalPower() int64 { return 0 }

// ConvertDuplicateVoteEvidence converts a Tendermint concrete Evidence type to
// SDK Evidence using Equivocation as the concrete type.
func ConvertDuplicateVoteEvidence(dupVote abci.Evidence) exported.Evidence {
	return Equivocation{
		Height:           dupVote.Height,
		Power:            dupVote.Validator.Power,
		ConsensusAddress: sdk.ConsAddress(dupVote.Validator.Address),
		Time:             dupVote.Time,
	}
}

// Evidence defines the application-level allowed Evidence to be submitted via a
// MsgSubmitEvidence message.
type Evidence struct {
	Sum interface{} `protobuf_oneof:"sum"`
}

type EvidenceEquivocation struct {
	Equivocation *Equivocation `json:"equivocation,omitempty" yaml:"equivocation"`
}

func (e *Evidence) GetEquivocation() *Equivocation {
	if x, ok := e.Sum.(*EvidenceEquivocation); ok {
		return x.Equivocation
	}
	return nil
}

func (e *Evidence) GetEvidence() exported.Evidence {
	if x := e.GetEquivocation(); x != nil {
		return x
	}
	return nil
}

func (e *Evidence) SetEvidence(value exported.Evidence) error {
	if value == nil {
		e.Sum = nil
		return nil
	}
	switch vt := value.(type) {
	case *Equivocation:
		e.Sum = &EvidenceEquivocation{vt}
		return nil
	case Equivocation:
		e.Sum = &EvidenceEquivocation{&vt}
		return nil
	}
	return fmt.Errorf("can't encode value of type %T as message Evidence", value)
}

// MsgSubmitEvidence defines the application-level message type for handling
// evidence submission.
type MsgSubmitEvidence struct {
	Evidence              *Evidence `json:"evidence,omitempty"`
	MsgSubmitEvidenceBase `json:"base"`
}
