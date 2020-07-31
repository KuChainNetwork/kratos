package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"gopkg.in/yaml.v2"
)

// DefaultStartingProposalID is 1
const DefaultStartingProposalID uint64 = 1

// ProposalStatus is a type alias that represents a proposal status as a byte
type ProposalStatus int32

const (
	// PROPOSAL_STATUS_UNSPECIFIED defines the default propopsal status.
	StatusNil ProposalStatus = 0
	// PROPOSAL_STATUS_DEPOSIT_PERIOD defines a proposal status during the deposit period.
	StatusDepositPeriod ProposalStatus = 1
	// PROPOSAL_STATUS_VOTING_PERIOD defines a proposal status during the voting period.
	StatusVotingPeriod ProposalStatus = 2
	// PROPOSAL_STATUS_PASSED defines a proposal status of a proposal that has passed.
	StatusPassed ProposalStatus = 3
	// PROPOSAL_STATUS_REJECTED defines a proposal status of a proposal that has been rejected.
	StatusRejected ProposalStatus = 4
	// PROPOSAL_STATUS_FAILED defines a proposal status of a proposal that has failed.
	StatusFailed ProposalStatus = 5
)

// ProposalBase defines the core field members of a governance proposal. It includes
// all static fields (i.e fields excluding the dynamic Content). A full proposal
// extends the ProposalBase with Content.
type ProposalBase struct {
	ProposalID       uint64         `json:"id" yaml:"id"`
	Status           ProposalStatus `json:"status,omitempty" yaml:"proposal_status"`
	FinalTallyResult TallyResult    `json:"final_tally_result" yaml:"final_tally_result"`
	SubmitTime       time.Time      `json:"submit_time" yaml:"submit_time"`
	DepositEndTime   time.Time      `json:"deposit_end_time" yaml:"deposit_end_time"`
	TotalDeposit     Coins          `json:"total_deposit" yaml:"total_deposit"`
	VotingStartTime  time.Time      `json:"voting_start_time" yaml:"voting_start_time"`
	VotingEndTime    time.Time      `json:"voting_end_time" yaml:"voting_end_time"`
}

func (p ProposalBase) Equal(other ProposalBase) bool {
	return p.ProposalID == other.ProposalID &&
		p.Status == other.Status &&
		p.FinalTallyResult.Equal(other.FinalTallyResult) &&
		p.SubmitTime.Equal(other.SubmitTime) &&
		p.DepositEndTime.Equal(other.DepositEndTime) &&
		p.TotalDeposit.IsEqual(other.TotalDeposit) &&
		p.VotingEndTime.Equal(other.VotingEndTime) &&
		p.VotingEndTime.Equal(other.VotingEndTime)
}

// Proposal defines a struct used by the governance module to allow for voting
// on network changes.
type Proposal struct {
	Content `json:"content" yaml:"content"` // Proposal content interface
	ProposalBase
}

// NewProposal creates a new Proposal instance
func NewProposal(content Content, id uint64, submitTime, depositEndTime time.Time) Proposal {
	return Proposal{
		Content: content,
		ProposalBase: ProposalBase{
			ProposalID:       id,
			Status:           StatusDepositPeriod,
			FinalTallyResult: EmptyTallyResult(),
			TotalDeposit:     NewCoins(),
			SubmitTime:       submitTime,
			DepositEndTime:   depositEndTime,
		},
	}
}

// Equal returns true if two Proposal types are equal.
func (p Proposal) Equal(other Proposal) bool {
	return p.ProposalBase.Equal(other.ProposalBase) && p.Content.String() == other.Content.String()
}

// String implements stringer interface
func (p Proposal) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Proposals is an array of proposal
type Proposals []Proposal

// Equal returns true if two slices (order-dependant) of proposals are equal.
func (p Proposals) Equal(other Proposals) bool {
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

// String implements stringer interface
func (p Proposals) String() string {
	out := "ID - (Status) [Type] Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) [%s] %s\n",
			prop.ProposalID, prop.Status,
			prop.ProposalType(), prop.GetTitle())
	}
	return strings.TrimSpace(out)
}

type (
	// ProposalQueue defines a queue for proposal ids
	ProposalQueue []uint64
)

// ProposalStatusFromString turns a string into a ProposalStatus
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	switch str {
	case "DepositPeriod":
		return StatusDepositPeriod, nil

	case "VotingPeriod":
		return StatusVotingPeriod, nil

	case "Passed":
		return StatusPassed, nil

	case "Rejected":
		return StatusRejected, nil

	case "Failed":
		return StatusFailed, nil

	case "":
		return StatusNil, nil

	default:
		return ProposalStatus(0xff), fmt.Errorf("'%s' is not a valid proposal status", str)
	}
}

// ValidProposalStatus returns true if the proposal status is valid and false
// otherwise.
func ValidProposalStatus(status ProposalStatus) bool {
	if status == StatusDepositPeriod ||
		status == StatusVotingPeriod ||
		status == StatusPassed ||
		status == StatusRejected ||
		status == StatusFailed {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status ProposalStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *ProposalStatus) Unmarshal(data []byte) error {
	*status = ProposalStatus(data[0])
	return nil
}

// MarshalJSON Marshals to JSON using string representation of the status
func (status ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// UnmarshalJSON Unmarshals from JSON assuming Bech32 encoding
func (status *ProposalStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ProposalStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

// String implements the Stringer interface.
func (status ProposalStatus) String() string {
	switch status {
	case StatusDepositPeriod:
		return "DepositPeriod"

	case StatusVotingPeriod:
		return "VotingPeriod"

	case StatusPassed:
		return "Passed"

	case StatusRejected:
		return "Rejected"

	case StatusFailed:
		return "Failed"

	default:
		return ""
	}
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (status ProposalStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}

// Proposal types
const (
	ProposalTypeText string = "Text"
)

// Implements Content Interface
var _ Content = TextProposal{}

// TextProposal defines a standard text proposal whose changes need to be
// manually updated in case of approval
type TextProposal struct {
	Title       string `json:"title,omitempty" yaml:"title"`
	Description string `json:"description,omitempty" yaml:"description"`
}

// NewTextProposal creates a text proposal Content
func NewTextProposal(title, description string) Content {
	return TextProposal{title, description}
}

// GetTitle returns the proposal title
func (tp TextProposal) GetTitle() string { return tp.Title }

// GetDescription returns the proposal description
func (tp TextProposal) GetDescription() string { return tp.Description }

// ProposalRoute returns the proposal router key
func (tp TextProposal) ProposalRoute() string { return RouterKey }

// ProposalType is "Text"
func (tp TextProposal) ProposalType() string { return ProposalTypeText }

// ValidateBasic validates the content's title and description of the proposal
func (tp TextProposal) ValidateBasic() error { return ValidateAbstract(tp) }

// String implements Stringer interface
func (tp TextProposal) String() string {
	out, _ := yaml.Marshal(tp)
	return string(out)
}

var validProposalTypes = map[string]struct{}{
	ProposalTypeText: {},
}

// RegisterProposalType registers a proposal type. It will panic if the type is
// already registered.
func RegisterProposalType(ty string) {
	if _, ok := validProposalTypes[ty]; ok {
		panic(fmt.Sprintf("already registered proposal type: %s", ty))
	}

	validProposalTypes[ty] = struct{}{}
}

// ContentFromProposalType returns a Content object based on the proposal type.
func ContentFromProposalType(title, desc, ty string) Content {
	switch ty {
	case ProposalTypeText:
		return NewTextProposal(title, desc)

	default:
		return nil
	}
}

// IsValidProposalType returns a boolean determining if the proposal type is
// valid.
//
// NOTE: Modules with their own proposal types must register them.
func IsValidProposalType(ty string) bool {
	_, ok := validProposalTypes[ty]
	return ok
}

// ProposalHandler implements the Handler interface for governance module-based
// proposals (ie. TextProposal ). Since these are
// merely signaling mechanisms at the moment and do not affect state, it
// performs a no-op.
func ProposalHandler(_ sdk.Context, c Content) error {
	switch c.ProposalType() {
	case ProposalTypeText:
		// both proposal types do not change state so this performs a no-op
		return nil

	default:
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized gov proposal type: %s", c.ProposalType())
	}
}
