package common

import (
	"fmt"

	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueryDelegationRewards queries a delegation rewards between a delegator and a
// validator.
func QueryDelegationRewards(cliCtx context.CLIContext, queryRoute, delAddr, valAddr string) ([]byte, int64, error) {
	delegatorAddr, err := chainType.NewAccountIDFromStr(delAddr)
	if err != nil {
		return nil, 0, err
	}

	validatorAddr, err := chainType.NewAccountIDFromStr(valAddr)
	if err != nil {
		return nil, 0, err
	}

	fmt.Println("delegatorAddr:", delegatorAddr, "validatorAddr", validatorAddr)
	params := types.NewQueryDelegationRewardsParams(delegatorAddr, validatorAddr)
	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal params: %w", err)
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDelegationRewards)
	return cliCtx.QueryWithData(route, bz)
}

// QueryDelegatorValidators returns delegator's list of validators
// it submitted delegations to.
func QueryDelegatorValidators(cliCtx context.CLIContext, queryRoute string, delegatorID chainType.AccountID) ([]byte, error) {
	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDelegatorValidators),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryDelegatorParams(delegatorID)),
	)
	return res, err
}

// QueryValidatorCommission returns a validator's commission.
func QueryValidatorCommission(cliCtx context.CLIContext, queryRoute string, validatorID chainType.AccountID) ([]byte, error) {
	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryValidatorCommission),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryValidatorCommissionParams(validatorID)),
	)
	return res, err
}

// WithdrawAllDelegatorRewards builds a multi-message slice to be used
// to withdraw all delegations rewards for the given delegator.
func WithdrawAllDelegatorRewards(cliCtx context.CLIContext, auth sdk.AccAddress, queryRoute string, delegatorID chainType.AccountID) ([]sdk.Msg, error) {
	// retrieve the comprehensive list of all validators which the
	// delegator had submitted delegations to
	bz, err := QueryDelegatorValidators(cliCtx, queryRoute, delegatorID)
	if err != nil {
		return nil, err
	}

	var validators []chainType.AccountID //bugs ,x
	if err := cliCtx.Codec.UnmarshalJSON(bz, &validators); err != nil {
		return nil, err
	}

	// build multi-message transaction
	msgs := make([]sdk.Msg, 0, len(validators))

	for _, valAcc := range validators {
		fmt.Println("validator:", valAcc.String())
		msg := types.NewMsgWithdrawDelegatorReward(auth, delegatorID, valAcc)
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

// WithdrawValidatorRewardsAndCommission builds a two-message message slice to be
// used to withdraw both validation's commission and self-delegation reward.
func WithdrawValidatorRewardsAndCommission(validatorAcc chainType.AccountID) ([]sdk.Msg, error) {
	commissionMsg := types.NewMsgWithdrawValidatorCommission(validatorAcc.MustAccAddress(), validatorAcc)
	if err := commissionMsg.ValidateBasic(); err != nil {
		return nil, err
	}

	// build and validate MsgWithdrawDelegatorReward

	rewardMsg := types.NewMsgWithdrawDelegatorReward(validatorAcc.MustAccAddress(), validatorAcc, validatorAcc) //bugs ??
	if err := rewardMsg.ValidateBasic(); err != nil {
		return nil, err
	}

	return []sdk.Msg{commissionMsg, rewardMsg}, nil
}
