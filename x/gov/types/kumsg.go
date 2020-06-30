package types

import (
	"github.com/KuChain-io/kuchain/chain/msg"
	chaintype "github.com/KuChain-io/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	RouterKeyName = chaintype.MustName(RouterKey)
)

type KuMsgSubmitProposal struct {
	chaintype.KuMsg
	Content Content `json:"content" yaml:"content"`
}

// FIXME: need review of this msg desgin

func NewKuMsgSubmitProposal(auth sdk.AccAddress, content Content, initialDeposit sdk.Coins, proposer chaintype.AccountID) KuMsgSubmitProposal {
	return KuMsgSubmitProposal{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithTransfer(proposer, ModuleAccountID, initialDeposit),
			msg.WithData(Cdc(), &MsgSubmitProposalBase{
				InitialDeposit: initialDeposit,
				Proposer:       proposer,
			}),
		), content,
	}
}

func (msg KuMsgSubmitProposal) ValidateBasic() error {
	if msg.Content == nil {
		return sdkerrors.Wrap(ErrInvalidProposalContent, "missing content")
	}
	msgData := MsgSubmitProposalBase{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return err
	}
	if msgData.Proposer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msgData.Proposer.String())
	}
	if !msgData.InitialDeposit.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msgData.InitialDeposit.String())
	}
	if msgData.InitialDeposit.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msgData.InitialDeposit.String())
	}
	if !IsValidProposalType(msg.Content.ProposalType()) {
		return sdkerrors.Wrap(ErrInvalidProposalType, msg.Content.ProposalType())
	}

	return msg.Content.ValidateBasic()
}

func (msg KuMsgSubmitProposal) GetContent() Content { return msg.Content }
func (msg KuMsgSubmitProposal) GetInitialDeposit() sdk.Coins {
	msgData := MsgSubmitProposalBase{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return sdk.Coins{}
	}
	return msgData.InitialDeposit
}
func (msg KuMsgSubmitProposal) GetProposer() sdk.AccAddress {
	msgData := MsgSubmitProposalBase{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil
	}

	proposerAccAddress, ok := msgData.Proposer.ToAccAddress()
	if ok {
		return proposerAccAddress
	}
	return nil
}
func (msg KuMsgSubmitProposal) GetProposerAccountID() chaintype.AccountID {
	msgData := MsgSubmitProposalBase{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return chaintype.AccountID{}
	}

	return msgData.Proposer
}

type KuMsgDeposit struct {
	chaintype.KuMsg
}

func NewKuMsgDeposit(auth sdk.AccAddress, depositor chaintype.AccountID, proposalID uint64, amount sdk.Coins) KuMsgDeposit {
	return KuMsgDeposit{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithTransfer(depositor, ModuleAccountID, amount),
			msg.WithData(Cdc(), &MsgDeposit{proposalID, depositor, amount}),
		),
	}
}

type KuMsgVote struct {
	chaintype.KuMsg
}

func NewKuMsgVote(auth sdk.AccAddress, voter chaintype.AccountID, proposalID uint64, option VoteOption) KuMsgVote {
	return KuMsgVote{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgVote{proposalID, voter, option}),
		),
	}
}

type MsgGovUnJail struct {
	chaintype.KuMsg
}

func NewMsgGovUnjail(auth sdk.AccAddress, validatoraddr chaintype.AccountID) MsgGovUnJail {
	return MsgGovUnJail{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgGovUnjailBase{ValidatorAccount: validatoraddr}),
		),
	}
}

func (msg MsgGovUnjailBase) ValidateBasic() error {
	if msg.ValidatorAccount.Empty() {
		return ErrBadValidatorAddr //TODO
	}

	return nil
}

// nolint
func (msg MsgGovUnjailBase) GetSigners() []sdk.AccAddress {
	valAccAddress, ok := msg.ValidatorAccount.ToAccAddress()
	if ok {
		return []sdk.AccAddress{valAccAddress}
	}
	return []sdk.AccAddress{}
}
func (msg MsgGovUnjailBase) Route() string        { return RouterKey }
func (msg MsgGovUnjailBase) Type() chaintype.Name { return chaintype.MustName("govunjail") }

func (msg MsgGovUnjailBase) GetUnjailValidator() chaintype.AccountID { return msg.ValidatorAccount }
