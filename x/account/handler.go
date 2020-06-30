package account

import (
	"github.com/KuChain-io/kuchain/chain/constants"
	"github.com/KuChain-io/kuchain/chain/msg"
	chainTypes "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) msg.Handler {
	return func(ctx chainTypes.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *types.MsgCreateAccount:
			return handleMsgCreateAccount(ctx, k, msg)
		case *types.MsgUpdateAccountAuth:
			return handleMsgUpdateAccountAuth(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized account message type: %T", msg)
		}
	}
}

// handleMsgCreateAccount handler msg create account
func handleMsgCreateAccount(ctx chainTypes.Context, k Keeper, msg *types.MsgCreateAccount) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData, err := msg.GetData()
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "msg create account data unmarshal error")
	}

	ctx.RequireAuth(msgData.Creator)

	if constants.IsSystemAccount(msgData.Name) {
		return nil, types.ErrAccountCannotCreateSysAccount
	}

	logger.Debug("msg create account", "name", msgData.Name, "creator", msgData.Creator)

	if a := k.GetAccountByName(ctx.Context(), msgData.Name); a != nil {
		logger.Debug("account has already created", "name", msgData.Name)
		return nil, sdkerrors.Wrapf(types.ErrAccountHasCreated, "name %s", msgData.Name)
	}

	newAccount := k.NewAccountByName(ctx.Context(), msgData.Name)
	if err := newAccount.SetAuth(msgData.Auth); err != nil {
		return nil, sdkerrors.Wrapf(err, "set auth to account error")
	}

	// set account
	k.SetAccount(ctx.Context(), newAccount)

	// add auth
	k.EnsureAuthInited(ctx.Context(), msgData.Auth)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateAccount,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.From.String()),
			sdk.NewAttribute(types.AttributeKeyAccount, msgData.Name.String()),
			sdk.NewAttribute(types.AttributeKeyAuth, msgData.Auth.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgUpdateAccountAuth handler msg update account auth
func handleMsgUpdateAccountAuth(ctx chainTypes.Context, k Keeper, msg *types.MsgUpdateAccountAuth) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData := types.MsgUpdateAccountAuthData{}
	if err := msg.UnmarshalData(types.Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg create coin data unmarshal error")
	}

	logger.Debug("msg update account auth", "name", msgData.Name, "auth", msgData.Auth)

	accountStat := k.GetAccountByName(ctx.Context(), msgData.Name)
	if accountStat == nil {
		logger.Debug("account no found", "name", msgData.Name)
		return nil, sdkerrors.Wrapf(types.ErrAccountHasCreated, "name %s", msgData.Name)
	}

	// Auth will Changed
	ctx.RequireAccountAuth(accountStat.GetAuth())

	if err := accountStat.SetAuth(msgData.Auth); err != nil {
		return nil, sdkerrors.Wrapf(err, "set auth to account error")
	}

	// set account
	k.SetAccount(ctx.Context(), accountStat)

	// add auth
	k.EnsureAuthInited(ctx.Context(), msgData.Auth)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateAccountAuth,
			sdk.NewAttribute(types.AttributeKeyAccount, msgData.Name.String()),
			sdk.NewAttribute(types.AttributeKeyAuth, msgData.Auth.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
