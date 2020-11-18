package gov

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/msg"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k Keeper) msg.Handler {
	return func(ctx chainTypes.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.KuMsgSubmitProposal:
			ctx.RequireAuth(msg.GetProposerAccountID())
			return handleMsgSubmitProposal(ctx.Context(), k, msg)
		case types.KuMsgDeposit:
			return handleKuMsgDeposit(ctx, k, msg)
		case types.KuMsgVote:
			return handleKuMsgVote(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleKuMsgDeposit(ctx chainTypes.Context, k Keeper, msg types.KuMsgDeposit) (*sdk.Result, error) {
	msgData := types.MsgDeposit{}
	if err := msg.UnmarshalData(types.Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg MsgDeposit  data unmarshal error")
	}
	ctx.RequireAuth(msgData.Depositor)
	return handleMsgDeposit(ctx.Context(), k, msgData)
}

func handleKuMsgVote(ctx chainTypes.Context, k Keeper, msg types.KuMsgVote) (*sdk.Result, error) {
	msgData := types.MsgVote{}
	if err := msg.UnmarshalData(types.Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg MsgVote  data unmarshal error")
	}
	ctx.RequireAuth(msgData.Voter)
	return handleMsgVote(ctx.Context(), k, msgData)
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitProposalI) (*sdk.Result, error) {
	proposal, err := keeper.SubmitProposal(ctx, msg.GetContent())
	if err != nil {
		return nil, err
	}

	votingStarted, err := keeper.AddDeposit(ctx, proposal.ProposalID, msg.GetProposerAccountID(), msg.GetInitialDeposit())
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.GetProposer().String()),
		),
	)

	submitEvent := sdk.NewEvent(types.EventTypeSubmitProposal, sdk.NewAttribute(types.AttributeKeyProposalType, msg.GetContent().ProposalType()))
	if votingStarted {
		submitEvent = submitEvent.AppendAttributes(
			sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.ProposalID)),
		)
	}
	ctx.EventManager().EmitEvent(submitEvent)

	return &sdk.Result{
		Data:   GetProposalIDBytes(proposal.ProposalID),
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg MsgDeposit) (*sdk.Result, error) {
	votingStarted, err := keeper.AddDeposit(ctx, msg.ProposalID, msg.Depositor, msg.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor.String()),
		),
	)

	if votingStarted {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposalDeposit,
				sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", msg.ProposalID)),
			),
		)
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgVote(ctx sdk.Context, keeper Keeper, msg MsgVote) (*sdk.Result, error) {
	err := keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
