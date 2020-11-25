package dex

import (
	"fmt"
	"strconv"
	"time"

	"github.com/KuChainNetwork/kuchain/chain/msg"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
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
		case *types.MsgCreateSymbol:
			return handleMsgCreateSymbol(ctx, k, theMsg)
		case *types.MsgUpdateSymbol:
			return handleMsgUpdateSymbol(ctx, k, theMsg)
		case *types.MsgPauseSymbol:
			return handleMsgPauseSymbol(ctx, k, theMsg)
		case *types.MsgRestoreSymbol:
			return handleMsgRestoreSymbol(ctx, k, theMsg)
		case *types.MsgShutdownSymbol:
			return handleMsgShutdownSymbol(ctx, k, theMsg)
		case *types.MsgDexSigIn:
			return handleMsgDexSigIn(ctx, k, theMsg)
		case *types.MsgDexSigOut:
			return handleMsgDexSigOut(ctx, k, theMsg)
		case *types.MsgDexDeal:
			return handleMsgDexDeal(ctx, k, theMsg)
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
	msg *types.MsgUpdateDexDescription) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData := types.MsgUpdateDexDescriptionData{}
	if err := msg.UnmarshalData(ModuleCdc, &msgData); nil != err {
		return nil, errors.Wrapf(err, "msg dex update description data unmarshal error")
	}

	logger.Debug("handle dex update description",
		"creator", msgData.Creator,
		"desc", string(msgData.Desc))

	ctx.RequireAccount(msgData.Creator)

	if err := keeper.UpdateDexDescription(ctx.Context(),
		msgData.Creator, string(msgData.Desc)); nil != err {
		return nil, errors.Wrapf(err,
			"msg update dex %s description", msgData.Creator)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateDexDescription,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, msgData.Creator.String()),
			sdk.NewAttribute(types.AttributeKeyDescription, string(msgData.Desc)),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgDestroyDex handle Msg destroy dex
func handleMsgDestroyDex(ctx chainTypes.Context,
	keeper Keeper,
	msg *types.MsgDestroyDex) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData := types.MsgDestroyDexData{}
	if err := msg.UnmarshalData(ModuleCdc, &msgData); nil != err {
		return nil, errors.Wrapf(err, "msg dex destroy unmarshal error")
	}

	logger.Debug("handle dex destroy", "creator", msgData.Creator)

	if err := keeper.DestroyDex(ctx.Context(), msgData.Creator); nil != err {
		return nil, errors.Wrapf(err, "msg destroy dex %s", msgData.Creator)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDestroyDex,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, msgData.Creator.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgCreateSymbol handle Msg create symbol
func handleMsgCreateSymbol(ctx chainTypes.Context,
	keeper Keeper,
	msg *types.MsgCreateSymbol) (*sdk.Result, error) {
	data, err := msg.GetData()
	if err != nil {
		return nil, err
	}

	if !data.Base.Validate() || !data.Quote.Validate() {
		return nil, errors.Wrapf(types.ErrSymbolIncorrect,
			"msg create symbol %s data is incorrect", data.Creator)
	}

	logger := ctx.Logger()
	logger.Debug("handle create symbol", "creator", data.Creator)

	if err := keeper.CreateSymbol(ctx.Context(), data.Creator, &types.Symbol{
		Base:   data.Base,
		Quote:  data.Quote,
		Height: ctx.BlockHeight(),
		CreateTime: func() time.Time {
			if data.CreateTime.IsZero() {
				return ctx.BlockTime()
			}
			return data.CreateTime
		}(),
	}); err != nil {
		return nil, errors.Wrapf(err,
			"msg create symbol error, creator %s",
			data.Creator.String())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateSymbol,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, data.Creator.String()),
			sdk.NewAttribute(types.AttributeKeySymbolCreateHeight, fmt.Sprint(ctx.BlockHeight())),
			sdk.NewAttribute(types.AttributeKeySymbolBaseCreator, data.Base.Creator),
			sdk.NewAttribute(types.AttributeKeySymbolBaseCode, data.Base.Code),
			sdk.NewAttribute(types.AttributeKeySymbolBaseName, data.Base.Name),
			sdk.NewAttribute(types.AttributeKeySymbolBaseFullName, data.Base.FullName),
			sdk.NewAttribute(types.AttributeKeySymbolBaseIconURL, data.Base.IconURL),
			sdk.NewAttribute(types.AttributeKeySymbolBaseTxURL, data.Base.TxURL),
			sdk.NewAttribute(types.AttributeKeySymbolQuoteCreator, data.Quote.Creator),
			sdk.NewAttribute(types.AttributeKeySymbolQuoteCode, data.Quote.Code),
			sdk.NewAttribute(types.AttributeKeySymbolQuoteName, data.Quote.Name),
			sdk.NewAttribute(types.AttributeKeySymbolQuoteFullName, data.Quote.FullName),
			sdk.NewAttribute(types.AttributeKeySymbolQuoteIconURL, data.Quote.IconURL),
			sdk.NewAttribute(types.AttributeKeySymbolQuoteTxURL, data.Quote.TxURL),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgUpdateSymbol handle Msg update symbol
func handleMsgUpdateSymbol(ctx chainTypes.Context,
	keeper Keeper,
	msg *types.MsgUpdateSymbol) (*sdk.Result, error) {
	data, err := msg.GetData()
	if err != nil {
		return nil, err
	}

	if len(data.Base.Creator) == 0 ||
		len(data.Base.Code) == 0 ||
		len(data.Quote.Creator) == 0 ||
		len(data.Quote.Code) == 0 {
		return nil, errors.Wrapf(types.ErrSymbolIncorrect,
			"msg create symbol %s data is incorrect",
			data.Creator.String())
	}

	attributes := make([]sdk.Attribute, 0, 16)
	for _, e := range []struct {
		Key   string
		Value string
	}{
		{types.AttributeKeySymbolBaseName, data.Base.Name},
		{types.AttributeKeySymbolBaseFullName, data.Base.FullName},
		{types.AttributeKeySymbolBaseIconURL, data.Base.IconURL},
		{types.AttributeKeySymbolBaseTxURL, data.Base.TxURL},
		{types.AttributeKeySymbolQuoteName, data.Quote.Name},
		{types.AttributeKeySymbolQuoteFullName, data.Quote.FullName},
		{types.AttributeKeySymbolQuoteIconURL, data.Quote.IconURL},
		{types.AttributeKeySymbolQuoteTxURL, data.Quote.TxURL},
	} {
		if len(e.Value) > 0 {
			attributes = append(attributes, sdk.NewAttribute(e.Key, e.Value))
		}
	}

	if 0 >= len(attributes) {
		return nil, errors.Wrapf(types.ErrSymbolIncorrect,
			"msg create symbol %s data is incorrect",
			data.Creator.String())
	}

	ctx.Logger().Debug("handle update symbol", "creator", data.Creator)

	if err := keeper.UpdateSymbol(ctx.Context(), data.Creator, &types.Symbol{
		Base:  data.Base,
		Quote: data.Quote,
	}); err != nil {
		return nil, errors.Wrapf(err,
			"msg create symbol error, creator %s",
			data.Creator.String())
	}

	attributes = append(attributes, sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateSymbol,
			attributes...,
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgPauseSymbol handle Msg pause symbol
func handleMsgPauseSymbol(ctx chainTypes.Context,
	keeper Keeper,
	msg *types.MsgPauseSymbol) (res *sdk.Result, err error) {
	data, err := msg.GetData()
	if err != nil {
		return nil, err
	}

	if len(data.BaseCreator) == 0 ||
		len(data.BaseCode) == 0 ||
		len(data.QuoteCreator) == 0 ||
		len(data.QuoteCode) == 0 {
		return nil, errors.Wrapf(types.ErrSymbolIncorrect,
			"msg pause symbol base code or quote code is empty, creator %s",
			data.Creator.String())
	}

	ctx.Logger().Debug("handle pause symbol",
		"creator", data.Creator)

	if err := keeper.PauseSymbol(ctx.Context(),
		data.Creator,
		data.BaseCreator,
		data.BaseCode,
		data.QuoteCreator,
		data.QuoteCode); err != nil {
		return nil, errors.Wrapf(err,
			"msg shutdown symbol error, creator %s",
			data.Creator)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePauseSymbol,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, data.Creator.String()),
			sdk.NewAttribute(types.AttributeKeySymbolBaseCreator, data.BaseCreator),
			sdk.NewAttribute(types.AttributeKeySymbolBaseCode, data.BaseCode),
			sdk.NewAttribute(types.AttributeKeySymbolBaseName, data.QuoteCode),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgRestoreSymbol handle Msg restore symbol
func handleMsgRestoreSymbol(ctx chainTypes.Context,
	keeper Keeper,
	msg *types.MsgRestoreSymbol) (*sdk.Result, error) {
	data, err := msg.GetData()
	if err != nil {
		return nil, err
	}

	if len(data.BaseCreator) == 0 ||
		len(data.BaseCode) == 0 ||
		len(data.QuoteCreator) == 0 ||
		len(data.QuoteCode) == 0 {
		return nil, errors.Wrapf(types.ErrSymbolIncorrect,
			"msg restore symbol base code or quote code is empty, creator %s",
			data.Creator)
	}

	ctx.Logger().Debug("handle restore symbol",
		"creator", data.Creator)

	if err := keeper.RestoreSymbol(ctx.Context(),
		data.Creator,
		data.BaseCreator,
		data.BaseCode,
		data.QuoteCreator,
		data.QuoteCode); err != nil {
		return nil, errors.Wrapf(err,
			"msg restore symbol error, creator %s",
			data.Creator)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRestoreSymbol,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, data.Creator.String()),
			sdk.NewAttribute(types.AttributeKeySymbolBaseCreator, data.BaseCreator),
			sdk.NewAttribute(types.AttributeKeySymbolBaseCode, data.BaseCode),
			sdk.NewAttribute(types.AttributeKeySymbolBaseName, data.QuoteCode),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgShutdownSymbol handle Msg shutdown symbol
func handleMsgShutdownSymbol(ctx chainTypes.Context,
	keeper Keeper,
	msg *types.MsgShutdownSymbol) (*sdk.Result, error) {
	data, err := msg.GetData()
	if err != nil {
		return nil, err
	}

	if len(data.BaseCreator) == 0 ||
		len(data.BaseCode) == 0 ||
		len(data.QuoteCreator) == 0 ||
		len(data.QuoteCode) == 0 {
		return nil, errors.Wrapf(types.ErrSymbolIncorrect,
			"msg shutdown symbol base code or quote code is empty, creator %s",
			data.Creator)
	}

	ctx.Logger().Debug("handle shutdown symbol",
		"creator", data.Creator)

	if err = keeper.ShutdownSymbol(ctx.Context(),
		data.Creator,
		data.BaseCreator,
		data.BaseCode,
		data.QuoteCreator,
		data.QuoteCode); err != nil {
		return nil, errors.Wrapf(err,
			"msg shutdown symbol error, creator %s",
			data.Creator)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeShutdownSymbol,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, data.Creator.String()),
			sdk.NewAttribute(types.AttributeKeySymbolBaseCreator, data.BaseCreator),
			sdk.NewAttribute(types.AttributeKeySymbolBaseCode, data.BaseCode),
			sdk.NewAttribute(types.AttributeKeySymbolBaseName, data.QuoteCode),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
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

func handleMsgDexSigOut(ctx chainTypes.Context, k Keeper, msg *types.MsgDexSigOut) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData, err := msg.GetData()
	if err != nil {
		return nil, err
	}

	logger.Debug("handle dex sigout",
		"user", msgData.User, "dex", msgData.Dex, "amount", msgData.Amount, "isTimeout", msgData.IsTimeout)

	// Note: two mode, if just has user's auth, need to wait
	if msgData.IsTimeout {
		ctx.RequireAuth(msgData.User)
	} else {
		ctx.RequireAuth(msgData.Dex)
	}

	if err := k.SigOut(ctx.Context(), msgData.IsTimeout, msgData.User, msgData.Dex, msgData.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDexSigOut,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyUser, msgData.User.String()),
			sdk.NewAttribute(types.AttributeKeyDex, msgData.Dex.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msgData.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyIsTimeout, strconv.FormatBool(msgData.IsTimeout)),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDexDeal(ctx chainTypes.Context, k Keeper, msg *types.MsgDexDeal) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData, err := msg.GetData()
	if err != nil {
		return nil, err
	}

	ctx.RequireAuth(msgData.Dex)

	logger.Debug("handle dex deal", "dex", msgData.Dex)

	// Auto SigIn for each account

	// Update sigIn status
	acc1, ass1, acc2, ass2 := msg.GetDealByDex()
	fee1, fee2 := msg.GetDealFeeByDex()

	if err := k.Deal(ctx.Context(), msgData.Dex, acc1, acc2, ass1, ass2, fee1, fee2); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDexDeal,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyDex, msgData.Dex.String()),
			sdk.NewAttribute(types.AttributeKeyDealRole1, acc1.String()),
			sdk.NewAttribute(types.AttributeKeyDealToken1, ass1.String()),
			sdk.NewAttribute(types.AttributeKeyDealFee1, fee1.String()),
			sdk.NewAttribute(types.AttributeKeyDealRole2, acc2.String()),
			sdk.NewAttribute(types.AttributeKeyDealToken2, ass2.String()),
			sdk.NewAttribute(types.AttributeKeyDealFee2, fee2.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
