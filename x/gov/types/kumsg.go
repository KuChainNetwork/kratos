package types

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	RouterKeyName = MustName(RouterKey)
)

type KuMsgSubmitProposal struct {
	KuMsg
	Content Content `json:"content" yaml:"content"`
}

func NewKuMsgSubmitProposal(auth sdk.AccAddress, content Content, initialDeposit Coins, proposer AccountID) KuMsgSubmitProposal {
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
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	if msg.Content == nil {
		return sdkerrors.Wrap(ErrInvalidProposalContent, "missing content")
	}

	msgData := MsgSubmitProposalBase{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return err
	}

	if err := msg.KuMsg.ValidateTransferRequire(ModuleAccountID, msgData.InitialDeposit); err != nil {
		return chainTypes.ErrKuMsgInconsistentAmount
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
func (msg KuMsgSubmitProposal) GetInitialDeposit() Coins {
	msgData := MsgSubmitProposalBase{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return Coins{}
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
func (msg KuMsgSubmitProposal) GetProposerAccountID() AccountID {
	msgData := MsgSubmitProposalBase{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return AccountID{}
	}
	return msgData.Proposer
}

type KuMsgDeposit struct {
	KuMsg
}

func NewKuMsgDeposit(auth sdk.AccAddress, depositor AccountID, proposalID uint64, amount Coins) KuMsgDeposit {
	return KuMsgDeposit{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithTransfer(depositor, ModuleAccountID, amount),
			msg.WithData(Cdc(), &MsgDeposit{proposalID, depositor, amount}),
		),
	}
}

func (msg KuMsgDeposit) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	msgData := MsgDeposit{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return err
	}

	if err := msg.KuMsg.ValidateTransferRequire(ModuleAccountID, msgData.Amount); err != nil {
		return chainTypes.ErrKuMsgInconsistentAmount
	}

	return msgData.ValidateBasic()
}

type KuMsgVote struct {
	KuMsg
}

func NewKuMsgVote(auth sdk.AccAddress, voter AccountID, proposalID uint64, option VoteOption) KuMsgVote {
	return KuMsgVote{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgVote{proposalID, voter, option}),
		),
	}
}

func (msg KuMsgVote) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	msgData := MsgVote{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return err
	}
	return msgData.ValidateBasic()
}
