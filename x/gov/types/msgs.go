package types

import (
	"gopkg.in/yaml.v2"

	chaintype "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Governance message types and routes
const (
	TypeMsgDeposit        = "deposit"
	TypeMsgVote           = "vote"
	TypeMsgSubmitProposal = "submitproposal"
)

// MsgSubmitProposalI defines the specific interface a concrete message must
// implement in order to process governance proposals. The concrete MsgSubmitProposal
// must be defined at the application-level.
type MsgSubmitProposalI interface {
	sdk.Msg

	GetContent() Content
	GetInitialDeposit() sdk.Coins
	GetProposer() sdk.AccAddress
	GetProposerAccountID() chaintype.AccountID
}

// NewMsgSubmitProposalBase creates a new MsgSubmitProposalBase.
func NewMsgSubmitProposalBase(initialDeposit sdk.Coins, proposer chaintype.AccountID) MsgSubmitProposalBase {
	return MsgSubmitProposalBase{
		InitialDeposit: initialDeposit,
		Proposer:       proposer,
	}
}

// Route implements Msg
func (msg MsgSubmitProposalBase) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgSubmitProposalBase) Type() chaintype.Name {
	return chaintype.MustName(TypeMsgSubmitProposal)
}

// ValidateBasic implements Msg
func (msg MsgSubmitProposalBase) ValidateBasic() error {
	if msg.Proposer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Proposer.String())
	}
	if !msg.InitialDeposit.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.InitialDeposit.String())
	}
	if msg.InitialDeposit.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.InitialDeposit.String())
	}

	return nil
}

// GetSignBytes implements Msg
func (msg MsgSubmitProposalBase) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgSubmitProposalBase) GetSigners() []sdk.AccAddress {
	proposalAccAddress, ok := msg.Proposer.ToAccAddress()
	if ok {
		return []sdk.AccAddress{proposalAccAddress}
	}
	return []sdk.AccAddress{}
}

// String implements the Stringer interface
func (msg MsgSubmitProposalBase) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// NewMsgDeposit creates a new MsgDeposit instance
func NewMsgDeposit(depositor chaintype.AccountID, proposalID uint64, amount sdk.Coins) MsgDeposit {
	return MsgDeposit{proposalID, depositor, amount}
}

// Route implements Msg
func (msg MsgDeposit) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgDeposit) Type() chaintype.Name { return chaintype.MustName(TypeMsgDeposit) }

// ValidateBasic implements Msg
func (msg MsgDeposit) ValidateBasic() error {
	if msg.Depositor.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Depositor.String())
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgDeposit) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	depositorAccAddress, ok := msg.Depositor.ToAccAddress()
	if ok {
		return []sdk.AccAddress{depositorAccAddress}
	}
	return []sdk.AccAddress{}
}

// NewMsgVote creates a message to cast a vote on an active proposal
func NewMsgVote(voter chaintype.AccountID, proposalID uint64, option VoteOption) MsgVote {
	return MsgVote{proposalID, voter, option}
}

// Route implements Msg
func (msg MsgVote) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgVote) Type() chaintype.Name { return chaintype.MustName(TypeMsgVote) }

// ValidateBasic implements Msg
func (msg MsgVote) ValidateBasic() error {
	if msg.Voter.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Voter.String())
	}
	if !ValidVoteOption(msg.Option) {
		return sdkerrors.Wrap(ErrInvalidVote, msg.Option.String())
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgVote) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	voterAccAddress, ok := msg.Voter.ToAccAddress()
	if ok {
		return []sdk.AccAddress{voterAccAddress}
	}
	return []sdk.AccAddress{}
}

// ---------------------------------------------------------------------------
// Deprecated
//
// TODO: Remove once client-side Protobuf migration has been completed.
// ---------------------------------------------------------------------------

// MsgSubmitProposal defines a (deprecated) message to create/submit a governance
// proposal.
//
// TODO: Remove once client-side Protobuf migration has been completed.
type MsgSubmitProposal struct {
	Content        Content             `json:"content" yaml:"content"`
	InitialDeposit sdk.Coins           `json:"initial_deposit" yaml:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive
	Proposer       chaintype.AccountID `json:"proposer" yaml:"proposer"`               //  Address of the proposer
}

// NewMsgSubmitProposal returns a (deprecated) MsgSubmitProposal message.
//
// TODO: Remove once client-side Protobuf migration has been completed.
func NewMsgSubmitProposal(content Content, initialDeposit sdk.Coins, proposer chaintype.AccountID) MsgSubmitProposal {
	return MsgSubmitProposal{content, initialDeposit, proposer}
}

// ValidateBasic implements Msg
func (msg MsgSubmitProposal) ValidateBasic() error {
	if msg.Content == nil {
		return sdkerrors.Wrap(ErrInvalidProposalContent, "missing content")
	}
	if msg.Proposer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Proposer.String())
	}
	if !msg.InitialDeposit.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.InitialDeposit.String())
	}
	if msg.InitialDeposit.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.InitialDeposit.String())
	}
	if !IsValidProposalType(msg.Content.ProposalType()) {
		return sdkerrors.Wrap(ErrInvalidProposalType, msg.Content.ProposalType())
	}

	return msg.Content.ValidateBasic()
}

// GetSignBytes implements Msg
func (msg MsgSubmitProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// nolint
func (msg MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	proposerAccAddress, ok := msg.Proposer.ToAccAddress()
	if ok {
		return []sdk.AccAddress{proposerAccAddress}
	}
	return []sdk.AccAddress{}
}
func (msg MsgSubmitProposal) Route() string                { return RouterKey }
func (msg MsgSubmitProposal) Type() string                 { return TypeMsgSubmitProposal }
func (msg MsgSubmitProposal) GetContent() Content          { return msg.Content }
func (msg MsgSubmitProposal) GetInitialDeposit() sdk.Coins { return msg.InitialDeposit }
func (msg MsgSubmitProposal) GetProposer() sdk.AccAddress {
	proposerAccAddress, ok := msg.Proposer.ToAccAddress()
	if ok {
		return proposerAccAddress
	}
	return nil
}
func (msg MsgSubmitProposal) GetProposerAccountID() chaintype.AccountID { return msg.Proposer }

func (msg MsgSubmitProposal) Marshal() (dAtA []byte, err error) {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz), nil
}

func (msg MsgSubmitProposal) Unmarshal(dAtA []byte) (err error) {
	ModuleCdc.MustUnmarshalJSON(dAtA, &msg)

	return nil
}
