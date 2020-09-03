package msg

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/plugins"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handler defines the core of the state transition function of an application.
type Handler func(ctx Context, msg sdk.Msg) (*sdk.Result, error)

func WarpHandler(transfer AssetTransfer, auther AccountAuther, h Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		defer func() {
			if r := recover(); r != nil {
				ctx.Logger().Error("handler msg panic", "panic", r)
				panic(r)
			}
		}()

		kuCtx := NewKuMsgCtx(ctx.WithEventManager(sdk.NewEventManager()), auther, msg)
		kuCtx = kuCtx.WithAuths(msg.GetSigners())

		if kuMsg, ok := msg.(KuTransfMsg); ok {
			if err := onHandlerKuMsg(kuCtx, transfer, kuMsg); err != nil {
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

// onHandlerKuMsg handler Ku msg for transfer
func onHandlerKuMsg(ctx Context, k AssetTransfer, msg KuTransfMsg) error {
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

		ctx.RequireAuth(from)

		if err := k.Transfer(ctx.Context(), from, to, amount); err != nil {
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
