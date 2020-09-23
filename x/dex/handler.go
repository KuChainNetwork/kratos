package dex

import (
	"time"

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
		case *types.MsgCreateCurrency:
			return handleMsgCreateCurrency(ctx, k, theMsg)
		case *types.MsgUpdateCurrency:
			return handleMsgUpdateCurrency(ctx, k, theMsg)
		case *types.MsgShutdownCurrency:
			return handleMsgShutdownCurrency(ctx, k, theMsg)
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

// handleMsgCreateCurrency handle Msg create currency
func handleMsgCreateCurrency(ctx chainTypes.Context,
	keeper Keeper,
	msg *types.MsgCreateCurrency) (res *sdk.Result, err error) {
	var data types.MsgCreateCurrencyData
	if data, err = msg.GetData(); nil != err {
		return
	}
	if 0 >= len(data.Base.Code) ||
		0 >= len(data.Base.Name) ||
		0 >= len(data.Base.FullName) ||
		0 >= len(data.Base.IconUrl) ||
		0 >= len(data.Base.TxUrl) ||
		0 >= len(data.Quote.Code) ||
		0 >= len(data.Quote.Name) ||
		0 >= len(data.Quote.FullName) ||
		0 >= len(data.Quote.IconUrl) ||
		0 >= len(data.Quote.TxUrl) ||
		0 >= len(data.DomainAddress) {
		err = errors.Wrapf(types.ErrCurrencyIncorrect,
			"msg create currency %s data is incorrect",
			data.Creator.String())
		return
	}
	logger := ctx.Logger()
	logger.Debug("handle create currency",
		"creator", data.Creator)
	if err = keeper.CreateCurrency(ctx.Context(), data.Creator, &types.Currency{
		Base:          data.Base,
		Quote:         data.Quote,
		DomainAddress: data.DomainAddress,
		CreateTime:    time.Now(),
	}); nil != err {
		err = errors.Wrapf(err,
			"msg create currency error, creator %s",
			data.Creator.String())
		return
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateCurrency,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, data.Creator.String()),
			sdk.NewAttribute(types.AttributeKeyCurrencyBaseCode, data.Base.Code),
			sdk.NewAttribute(types.AttributeKeyCurrencyBaseName, data.Base.Name),
			sdk.NewAttribute(types.AttributeKeyCurrencyBaseFullName, data.Base.FullName),
			sdk.NewAttribute(types.AttributeKeyCurrencyBaseIconUrl, data.Base.IconUrl),
			sdk.NewAttribute(types.AttributeKeyCurrencyBaseTxUrl, data.Base.TxUrl),
			sdk.NewAttribute(types.AttributeKeyCurrencyQuoteCode, data.Quote.Code),
			sdk.NewAttribute(types.AttributeKeyCurrencyQuoteName, data.Quote.Name),
			sdk.NewAttribute(types.AttributeKeyCurrencyQuoteFullName, data.Quote.FullName),
			sdk.NewAttribute(types.AttributeKeyCurrencyQuoteIconUrl, data.Quote.IconUrl),
			sdk.NewAttribute(types.AttributeKeyCurrencyQuoteTxUrl, data.Quote.TxUrl),
			sdk.NewAttribute(types.AttributeKeyCurrencyDomainAddress, data.DomainAddress),
		),
	)
	res = &sdk.Result{Events: ctx.EventManager().Events()}
	return
}

// handleMsgUpdateCurrency handle Msg update currency
func handleMsgUpdateCurrency(ctx chainTypes.Context,
	keeper Keeper,
	msg *types.MsgUpdateCurrency) (res *sdk.Result, err error) {
	var data types.MsgUpdateCurrencyData
	if data, err = msg.GetData(); nil != err {
		return
	}
	if 0 >= len(data.Base.Code) ||
		0 >= len(data.Quote.Code) {
		err = errors.Wrapf(types.ErrCurrencyIncorrect,
			"msg create currency %s data is incorrect",
			data.Creator.String())
		return
	}
	upMap := make(map[string]string)
	for _, e := range []struct {
		Key   string
		Value string
	}{
		{types.AttributeKeyCurrencyBaseName, data.Base.Name},
		{types.AttributeKeyCurrencyBaseFullName, data.Base.FullName},
		{types.AttributeKeyCurrencyBaseIconUrl, data.Base.IconUrl},
		{types.AttributeKeyCurrencyBaseTxUrl, data.Base.TxUrl},
		{types.AttributeKeyCurrencyQuoteName, data.Quote.Name},
		{types.AttributeKeyCurrencyQuoteFullName, data.Quote.FullName},
		{types.AttributeKeyCurrencyQuoteIconUrl, data.Quote.IconUrl},
		{types.AttributeKeyCurrencyQuoteTxUrl, data.Quote.TxUrl},
	} {
		if 0 < len(e.Value) {
			upMap[e.Key] = e.Value
		}
	}
	if 0 >= len(upMap) {
		err = errors.Wrapf(types.ErrCurrencyIncorrect,
			"msg create currency %s data is incorrect",
			data.Creator.String())
		return
	}
	logger := ctx.Logger()
	logger.Debug("handle update currency",
		"creator", data.Creator)
	if err = keeper.UpdateCurrencyInfo(ctx.Context(), data.Creator, &types.Currency{
		Base:  data.Base,
		Quote: data.Quote,
	}); nil != err {
		err = errors.Wrapf(err,
			"msg create currency error, creator %s",
			data.Creator.String())
		return
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateCurrency,
			func() (list []sdk.Attribute) {
				list = append(list, sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory))
				for k, v := range upMap {
					list = append(list, sdk.NewAttribute(k, v))
				}
				return
			}()...,
		),
	)
	res = &sdk.Result{Events: ctx.EventManager().Events()}
	return
}

// handleMsgShutdownCurrency handle Msg shutdown currency
func handleMsgShutdownCurrency(ctx chainTypes.Context,
	keeper Keeper,
	msg *types.MsgShutdownCurrency) (res *sdk.Result, err error) {
	var data types.MsgShutdownCurrencyData
	if data, err = msg.GetData(); nil != err {
		return
	}
	if 0 >= len(data.BaseCode) || 0 >= len(data.QuoteCode) {
		err = errors.Wrapf(types.ErrCurrencyIncorrect,
			"msg shutdown currency base code or quote code is empty",
			data.Creator.String())
		return
	}
	logger := ctx.Logger()
	logger.Debug("handle shutdown currency",
		"creator", data.Creator)
	if err = keeper.ShutdownCurrency(ctx.Context(), data.Creator, data.BaseCode, data.QuoteCode); nil != err {
		err = errors.Wrapf(err,
			"msg shutdown currency error, creator %s",
			data.Creator.String())
		return
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateCurrency,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, data.Creator.String()),
			sdk.NewAttribute(types.AttributeKeyCurrencyBaseCode, data.BaseCode),
			sdk.NewAttribute(types.AttributeKeyCurrencyBaseName, data.QuoteCode),
		),
	)
	res = &sdk.Result{Events: ctx.EventManager().Events()}
	return
	return
}
