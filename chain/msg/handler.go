package msg

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/plugins"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Handler defines the core of the state transition function of an application.
type Handler func(ctx Context, msg sdk.Msg) (*sdk.Result, error)

func WarpHandler(transfer AssetTransfer, author AccountAuther, h Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		defer func() {
			if r := recover(); r != nil {
				ctx.Logger().Error("handler msg panic", "panic", r)
				panic(r)
			}
		}()

		kuCtx := NewKuMsgCtx(ctx.WithEventManager(sdk.NewEventManager()), author, msg)
		kuCtx = kuCtx.WithAuths(msg.GetSigners())

		kuMsg, ok := msg.(KuTransfMsg)

		if ok {
			if err := onHandlerKuMsg(kuCtx, transfer, author, kuMsg); err != nil {
				return nil, err
			}
		}

		res, err := h(kuCtx, msg)
		if err != nil {
			return nil, err
		}

		if err := kuCtx.CheckAuths(); err != nil {
			return nil, err
		}

		plugins.HandleEvent(ctx, res.Events)

		return res, err
	}
}

func getAuthByAccountID(ctx Context, author AccountAuther, id AccountID) (AccAddress, error) {
	if add, ok := id.ToAccAddress(); ok {
		return add, nil
	}

	if name, ok := id.ToName(); ok {
		authByAccount, err := author.GetAuth(ctx.Context(), name)
		if err != nil {
			return nil, err
		}

		return authByAccount, nil
	}

	return nil, types.ErrMissingAuth
}

func checkTransferAuth(ctx Context, transfer AssetTransfer, author AccountAuther, msg types.KuMsgTransfer) (bool, error) {
	fromAuth, err := getAuthByAccountID(ctx, author, msg.From)
	if err != nil {
		// no found from auth, there must be a error, even no need from auth
		return false, sdkerrors.Wrapf(err, "get from auth error")
	}

	// check is has fromAuth
	if ctx.IsHasAuth(fromAuth) {
		return false, nil // has fromAuth, so it checked ok
	}

	// if is from has approve to with amt
	toAuth, err := getAuthByAccountID(ctx, author, msg.To)
	if err != nil {
		// no found to auth, there must be a error, even no need to auth
		// that means account can not approve to module account
		return false, sdkerrors.Wrapf(err, "get to auth error")
	}

	if ctx.IsHasAuth(toAuth) {
		if err := transfer.ApplyApporve(ctx.Context(), msg.From, msg.To, msg.Amount); err == nil {
			// if apply apporve success, then checked ok
			return true, nil
		}
	}

	return false, types.ErrMissingAuth
}

// onHandlerKuMsg handler Ku msg for transfer
func onHandlerKuMsg(ctx Context, k AssetTransfer, author AccountAuther, msg KuTransfMsg) error {
	transfers := msg.GetTransfers()

	for _, t := range transfers {
		from := t.From
		to := t.To
		amount := t.Amount

		if from.Empty() || to.Empty() || amount.IsZero() || from.Eq(to) {
			continue
		}

		// check validate for safe
		if err := msg.ValidateTransfer(); err != nil {
			return err
		}

		isApplyApprove, err := checkTransferAuth(ctx, k, author, t)
		if err != nil {
			return err
		}

		if err := k.TransferDetail(ctx.Context(), from, to, amount, isApplyApprove); err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				EventTypeTransfer,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.KuCodeSpace),
				sdk.NewAttribute(AttributeKeyFrom, from.String()),
				sdk.NewAttribute(AttributeKeyTo, to.String()),
				sdk.NewAttribute(AttributeKeyAmount, amount.String()),
			),
		)
	}

	return nil
}
