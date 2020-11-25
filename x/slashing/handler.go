package slashing

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k Keeper) msg.Handler {
	return func(ctx chainTypes.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.KuMsgUnjail:
			return handleKuMsgUnjail(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleKuMsgUnjail(ctx chainTypes.Context, k Keeper, msg types.KuMsgUnjail) (*sdk.Result, error) {
	msgData := types.MsgUnjail{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg Unjail  data unmarshal error")
	}

	ctx.RequireAuth(msgData.ValidatorAddr)

	return handleMsgUnjail(ctx.Context(), msgData, k)
}

// Validators must submit a transaction to unjail itself after
// having been jailed (and thus unbonded) for downtime
func handleMsgUnjail(ctx sdk.Context, msg MsgUnjail, k Keeper) (*sdk.Result, error) {
	err := k.Unjail(ctx, msg.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
