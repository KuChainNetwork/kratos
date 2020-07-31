package staking

import (
	"time"

	"github.com/KuChainNetwork/kuchain/chain/msg"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	stakingexport "github.com/KuChainNetwork/kuchain/x/staking/exported"
	"github.com/KuChainNetwork/kuchain/x/staking/keeper"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/gogo/protobuf/types"
	tmstrings "github.com/tendermint/tendermint/libs/strings"
	tmtypes "github.com/tendermint/tendermint/types"
)

func NewHandler(k keeper.Keeper) msg.Handler {
	return func(ctx chainTypes.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.KuMsgCreateValidator:
			return handleKuMsgCreateValidator(ctx, k, msg)
		case types.KuMsgDelegate:
			return handleKuMsgDelegate(ctx, k, msg)
		case types.KuMsgEditValidator:
			return handleKuMsgEditValidator(ctx, k, msg)
		case types.KuMsgRedelegate:
			return handleKuMsgRedelegate(ctx, k, msg)
		case types.KuMsgUnbond:
			return handleKuMsgUnbond(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleKuMsgCreateValidator(ctx chainTypes.Context, k keeper.Keeper, msg types.KuMsgCreateValidator) (*sdk.Result, error) {
	msgData := types.MsgCreateValidator{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg CreateValidator data unmarshal error")
	}

	ctx.RequireAuth(msgData.ValidatorAccount)
	return handleMsgCreateValidator(ctx.Context(), msgData, k)
}

func handleKuMsgDelegate(ctx chainTypes.Context, k keeper.Keeper, msg types.KuMsgDelegate) (*sdk.Result, error) {
	msgData := types.MsgDelegate{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg Delegate data unmarshal error")
	}
	ctx.RequireAuth(msgData.DelegatorAccount)
	return handleMsgDelegate(ctx, msgData, k)
}

func handleKuMsgEditValidator(ctx chainTypes.Context, k keeper.Keeper, msg types.KuMsgEditValidator) (*sdk.Result, error) {
	msgData := types.MsgEditValidator{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg EditValidator  data unmarshal error")
	}
	ctx.RequireAuth(msgData.ValidatorAccount)
	return handleMsgEditValidator(ctx.Context(), msgData, k)
}

func handleKuMsgRedelegate(ctx chainTypes.Context, k keeper.Keeper, msg types.KuMsgRedelegate) (*sdk.Result, error) {
	msgData := types.MsgBeginRedelegate{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg Redelegate  data unmarshal error")
	}
	ctx.RequireAuth(msgData.DelegatorAccount)
	return handleMsgBeginRedelegate(ctx.Context(), msgData, k)
}

func handleKuMsgUnbond(ctx chainTypes.Context, k keeper.Keeper, msg types.KuMsgUnbond) (*sdk.Result, error) {
	msgData := types.MsgUndelegate{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg Unbond  data unmarshal error")
	}
	ctx.RequireAuth(msgData.DelegatorAccount)
	return handleMsgUndelegate(ctx.Context(), msgData, k)
}

// These functions assume everything has been authenticated,
// now we just perform action and save

func handleMsgCreateValidator(ctx sdk.Context, msg types.MsgCreateValidator, k keeper.Keeper) (*sdk.Result, error) {
	logger := k.Logger(ctx)

	logger.Debug("handle msg create validator", "msg", msg)

	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetValidator(ctx, msg.ValidatorAccount); found {
		return nil, ErrValidatorOwnerExists
	}

	if !k.ValidatorAccount(ctx, msg.ValidatorAccount) {
		return nil, ErrUnKnowAccount
	}

	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, msg.Pubkey)
	if err != nil {
		return nil, err
	}

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return nil, ErrValidatorPubKeyExists
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}

	if ctx.ConsensusParams() != nil {
		tmPubKey := tmtypes.TM2PB.PubKey(pk)
		if !tmstrings.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
			return nil, sdkerrors.Wrapf(
				ErrValidatorPubKeyTypeNotSupported,
				"got: %s, expected: %s", tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes,
			)
		}
	}

	validator := NewValidator(msg.ValidatorAccount, pk, msg.Description)
	commission := NewCommissionWithTime(
		msg.CommissionRates, sdk.OneDec(),
		sdk.OneDec(), ctx.BlockHeader().Time,
	)

	validator, err = validator.SetInitialCommission(commission)
	if err != nil {
		return nil, err
	}

	k.SetValidator(ctx, validator)
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)

	// call the after-creation hook
	k.AfterValidatorCreated(ctx, validator.OperatorAccount)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateValidator,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAccount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAccount.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgEditValidator(ctx sdk.Context, msg types.MsgEditValidator, k keeper.Keeper) (*sdk.Result, error) {
	// validator must already be registered
	validator, found := k.GetValidator(ctx, msg.ValidatorAccount)
	if !found {
		return nil, ErrNoValidatorFound
	}

	// replace all editable fields (clients should autofill existing values)
	description, err := validator.Description.UpdateDescription(msg.Description)
	if err != nil {
		return nil, err
	}

	validator.Description = description
	if msg.CommissionRate != nil {
		commission, err := k.UpdateValidatorCommission(ctx, validator, *msg.CommissionRate)
		if err != nil {
			return nil, err
		}

		// call the before-modification hook since we're about to update the commission
		k.BeforeValidatorModified(ctx, msg.ValidatorAccount)

		validator.Commission = commission
	}

	k.SetValidator(ctx, validator)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditValidator,
			sdk.NewAttribute(types.AttributeKeyCommissionRate, validator.Commission.String()),
			sdk.NewAttribute(types.AttributeKeyMinSelfDelegation, validator.MinSelfDelegation.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAccount.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDelegate(ctx chainTypes.Context, msg types.MsgDelegate, k keeper.Keeper) (*sdk.Result, error) {
	validator, found := k.GetValidator(ctx.Context(), msg.ValidatorAccount)
	if !found {
		return nil, ErrNoValidatorFound
	}

	if err := ctx.RequireTransfer(types.ModuleAccountID, chainTypes.Coins{msg.Amount}); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg delegate required transfer no enough")
	}

	if msg.Amount.Denom != k.BondDenom(ctx.Context()) {
		return nil, ErrBadDenom
	}

	// NOTE: source funds are always unbonded
	_, err := k.Delegate(ctx.Context(), msg.DelegatorAccount, msg.Amount.Amount, stakingexport.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAccount.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAccount.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUndelegate(ctx sdk.Context, msg types.MsgUndelegate, k keeper.Keeper) (*sdk.Result, error) {
	shares, err := k.ValidateUnbondAmount(
		ctx, msg.DelegatorAccount, msg.ValidatorAccount, msg.Amount.Amount,
	)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Denom != k.BondDenom(ctx) {
		return nil, ErrBadDenom
	}

	completionTime, err := k.Undelegate(ctx, msg.DelegatorAccount, msg.ValidatorAccount, shares)
	if err != nil {
		return nil, err
	}

	ts, err := gogotypes.TimestampProto(completionTime)
	if err != nil {
		return nil, ErrBadRedelegationAddr
	}

	completionTimeBz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(ts)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbond,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAccount.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAccount.String()),
		),
	})

	return &sdk.Result{Data: completionTimeBz, Events: ctx.EventManager().Events()}, nil
}

func handleMsgBeginRedelegate(ctx sdk.Context, msg types.MsgBeginRedelegate, k keeper.Keeper) (*sdk.Result, error) {
	shares, err := k.ValidateUnbondAmount(
		ctx, msg.DelegatorAccount, msg.ValidatorSrcAccount, msg.Amount.Amount,
	)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Denom != k.BondDenom(ctx) {
		return nil, ErrBadDenom
	}

	completionTime, err := k.BeginRedelegation(
		ctx, msg.DelegatorAccount, msg.ValidatorSrcAccount, msg.ValidatorDstAccount, shares,
	)
	if err != nil {
		return nil, err
	}

	ts, err := gogotypes.TimestampProto(completionTime)
	if err != nil {
		return nil, ErrBadRedelegationAddr
	}

	completionTimeBz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(ts)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRedelegate,
			sdk.NewAttribute(types.AttributeKeySrcValidator, msg.ValidatorSrcAccount.String()),
			sdk.NewAttribute(types.AttributeKeyDstValidator, msg.ValidatorDstAccount.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAccount.String()),
		),
	})

	return &sdk.Result{Data: completionTimeBz, Events: ctx.EventManager().Events()}, nil
}
