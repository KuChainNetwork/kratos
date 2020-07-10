package evidence

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/evidence/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k Keeper) msg.Handler {
	return func(ctx chainTypes.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgSubmitEvidenceBase:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "%T must be extended to support evidence", msg)

		default:
			msgSubEv, ok := msg.(exported.MsgSubmitEvidence)
			if ok {
				return handleMsgSubmitEvidence(ctx.Context(), k, msgSubEv)
			}

			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgSubmitEvidence(ctx sdk.Context, k Keeper, msg exported.MsgSubmitEvidence) (*sdk.Result, error) {
	evidence := msg.GetEvidence()
	if err := k.SubmitEvidence(ctx, evidence); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.GetSubmitter().String()),
		),
	)

	return &sdk.Result{
		Data:   evidence.Hash(),
		Events: ctx.EventManager().Events(),
	}, nil
}
