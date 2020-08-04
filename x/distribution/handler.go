package distribution

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/keeper"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k keeper.Keeper) msg.Handler {
	return func(ctx chainTypes.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgSetWithdrawAccountId:
			return handleMsgModifyWithdrawAccountId(ctx, msg, k)

		case types.MsgWithdrawDelegatorReward:
			return handleMsgWithdrawDelegatorReward(ctx, msg, k)

		case types.MsgWithdrawValidatorCommission:
			return handleMsgWithdrawValidatorCommission(ctx, msg, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distribution message type: %T", msg)
		}
	}
}

// These functions assume everything has been authenticated (ValidateBasic passed, and signatures checked)
func handleMsgModifyWithdrawAccountId(ctx chainTypes.Context, msg types.MsgSetWithdrawAccountId, k keeper.Keeper) (*sdk.Result, error) {
	checkAcc := func(acc AccountID) bool {
		_, ok := acc.ToName()
		if ok {
			return k.AccKeeper.IsAccountExist(ctx.Context(), acc)
		}
		return false
	}

	types.FindAcc = checkAcc

	dataMsg, _ := msg.GetData()
	ctx.RequireAuth(dataMsg.DelegatorAccountid)
	ExistOk := checkAcc(dataMsg.WithdrawAccountid)
	if ExistOk {
		err := k.SetWithdrawAddr(ctx.Context(), dataMsg.DelegatorAccountid, dataMsg.WithdrawAccountid)
		if err != nil {
			return nil, err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, dataMsg.DelegatorAccountid.String()),
			),
		)
		return &sdk.Result{Events: ctx.EventManager().Events()}, nil
	} else {
		ctx.Logger().Error("handleMsgModifyWithdrawAccountId, not found",
			"WithdrawAccountid", dataMsg.WithdrawAccountid, "ExistOk", ExistOk)
	}

	return nil, types.ErrSetWithdrawAddrDisabled
}

func handleMsgWithdrawDelegatorReward(ctx chainTypes.Context, msg types.MsgWithdrawDelegatorReward, k keeper.Keeper) (*sdk.Result, error) {
	dataMsg, _ := msg.GetData()
	ctx.RequireAuth(dataMsg.DelegatorAccountId)
	ctx.Logger().Debug("handleMsgWithdrawDelegatorReward", "valId", dataMsg.ValidatorAccountId, "delId", dataMsg.DelegatorAccountId)

	_, err := k.WithdrawDelegationRewards(ctx.Context(), dataMsg.DelegatorAccountId, dataMsg.ValidatorAccountId)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, dataMsg.DelegatorAccountId.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdrawValidatorCommission(ctx chainTypes.Context, msg types.MsgWithdrawValidatorCommission, k keeper.Keeper) (*sdk.Result, error) {
	dataMsg, _ := msg.GetData()
	ctx.RequireAuth(dataMsg.ValidatorAccountId)

	_, err := k.WithdrawValidatorCommission(ctx.Context(), dataMsg.ValidatorAccountId)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, dataMsg.ValidatorAccountId.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func NewCommunityPoolSpendProposalHandler(k Keeper) types.GovTypesHandler {
	return func(ctx sdk.Context, content types.GovTypesContent) error {
		switch c := content.(type) {
		case types.CommunityPoolSpendProposal:
			return keeper.HandleCommunityPoolSpendProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distr proposal content type: %T", c)
		}
	}
}
