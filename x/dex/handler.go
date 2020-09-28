package dex

import (
	"github.com/pkg/errors"

	"github.com/KuChainNetwork/kuchain/chain/msg"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) msg.Handler {
	return func(ctx chainTypes.Context, msg chainTypes.Msg) (*sdk.Result, error) {
		switch theMsg := msg.(type) {
		case *types.MsgCreateDex:
			return handleMsgCreateDex(ctx, k, theMsg)
		case *types.MsgUpdateDexDescription:
			return handleMsgUpdateDexDescription(ctx, k, theMsg)
		case *types.MsgDestroyDex:
			return handleMsgDestroyDex(ctx, k, theMsg)
		case *types.MsgDexSigIn:
			return handleMsgDexSigIn(ctx, k, theMsg)
		default:
			return nil, errors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized asset message type: %T", msg)
		}
	}
}

// handleMsgCreate Handle Msg create coin type
func handleMsgCreateDex(ctx chainTypes.Context, k Keeper, msg *types.MsgCreateDex) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData := types.MsgCreateDexData{}
	if err := msg.UnmarshalData(ModuleCdc, &msgData); err != nil {
		return nil, errors.Wrapf(err, "msg create coin data unmarshal error")
	}
	logger.Debug("handle dex create",
		"creator", msgData.Creator,
		"stakings", msgData.Stakings,
		"desc", string(msgData.Desc))
	ctx.RequireAccount(msgData.Creator)
	/* no need check, has check by ValiteBasic
	if err := ctx.RequireTransfer(types.ModuleAccountID, msgData.Stakings); err != nil {
		return nil, errors.Wrapf(err, "msg create dex error no transfer")
	}
	*/
	if err := k.CreateDex(ctx.Context(),
		msgData.Creator, msgData.Stakings, string(msgData.Desc)); err != nil {
		return nil, errors.Wrapf(err, "msg create dex %s", msgData.Creator)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateDex,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, msgData.Creator.String()),
			sdk.NewAttribute(types.AttributeKeyStakings, msgData.Stakings.String()),
			sdk.NewAttribute(types.AttributeKeyDescription, string(msgData.Desc)),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgUpdateDexDescription Handle Msg update dex description
func handleMsgUpdateDexDescription(ctx chainTypes.Context,
	keeper Keeper,
	msg *types.MsgUpdateDexDescription) (res *sdk.Result, err error) {
	logger := ctx.Logger()
	msgData := types.MsgUpdateDexDescriptionData{}
	if err = msg.UnmarshalData(ModuleCdc, &msgData); nil != err {
		err = errors.Wrapf(err, "msg dex update description data unmarshal error")
		return
	}
	// check description max length
	if types.MaxDexDescriptorLen < len(msgData.Desc) {
		err = types.ErrDexDescTooLong
		return
	}
	logger.Debug("handle dex update description",
		"creator", msgData.Creator,
		"desc", string(msgData.Desc))
	ctx.RequireAccount(msgData.Creator)
	if err = keeper.UpdateDexDescription(ctx.Context(),
		msgData.Creator,
		string(msgData.Desc)); nil != err {
		err = errors.Wrapf(err, "msg update dex %s description", msgData.Creator)
		return
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateDexDescription,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, msgData.Creator.String()),
			sdk.NewAttribute(types.AttributeKeyDescription, string(msgData.Desc)),
		),
	)
	res = &sdk.Result{Events: ctx.EventManager().Events()}
	return
}

// handleMsgDestroyDex handle Msg destroy dex
func handleMsgDestroyDex(ctx chainTypes.Context,
	keeper Keeper,
	msg *types.MsgDestroyDex) (res *sdk.Result, err error) {
	logger := ctx.Logger()
	msgData := types.MsgDestroyDexData{}
	if err = msg.UnmarshalData(ModuleCdc, &msgData); nil != err {
		err = errors.Wrapf(err, "msg dex destroy unmarshal error")
		return
	}
	logger.Debug("handle dex destroy",
		"creator", msgData.Creator)
	if err = keeper.DestroyDex(ctx.Context(), msgData.Creator); nil != err {
		err = errors.Wrapf(err, "msg destroy dex %s", msgData.Creator)
		return
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDestroyDex,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, msgData.Creator.String()),
		),
	)
	res = &sdk.Result{Events: ctx.EventManager().Events()}
	return
}

func handleMsgDexSigIn(ctx chainTypes.Context, k Keeper, msg *types.MsgDexSigIn) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData, err := msg.GetData()
	if err != nil {
		return nil, err
	}

	logger.Debug("handle dex sigin",
		"user", msgData.User, "dex", msgData.Dex, "amount", msgData.Amount)

	ctx.RequireAuth(msgData.User)

	if err := k.SigIn(ctx.Context(), msgData.User, msgData.Dex, msgData.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDexSigIn,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyUser, msgData.User.String()),
			sdk.NewAttribute(types.AttributeKeyDex, msgData.Dex.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msgData.Amount.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
