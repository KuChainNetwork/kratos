package keeper

import (
	"encoding/json"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, path[1:], req, k)

		case types.QueryValidatorOutstandingRewards:
			return queryValidatorOutstandingRewards(ctx, path[1:], req, k)

		case types.QueryValidatorCommission:
			return queryValidatorCommission(ctx, path[1:], req, k)

		case types.QueryValidatorSlashes:
			return queryValidatorSlashes(ctx, path[1:], req, k)

		case types.QueryDelegationRewards:
			return queryDelegationRewards(ctx, path[1:], req, k)

		case types.QueryDelegatorTotalRewards:
			return queryDelegatorTotalRewards(ctx, path[1:], req, k)

		case types.QueryDelegatorValidators:
			return queryDelegatorValidators(ctx, path[1:], req, k)

		case types.QueryWithdrawAddr:
			return queryDelegatorWithdrawAddress(ctx, path[1:], req, k)

		case types.QueryCommunityPool:
			return queryCommunityPool(ctx, path[1:], req, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryValidatorOutstandingRewards(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryValidatorOutstandingRewardsParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)

	ctx.Logger().Debug("queryValidatorOutstandingRewards:", params, err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	rewards := k.GetValidatorOutstandingRewards(ctx, params.ValidatorAddress)

	ctx.Logger().Debug("1 queryValidatorOutstandingRewards:", rewards)

	if rewards.Rewards == nil {
		rewards.Rewards = chainTypes.DecCoins{}
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, rewards)
	ctx.Logger().Debug("queryValidatorOutstandingRewards:", bz, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryValidatorCommission(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryValidatorCommissionParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	ctx.Logger().Debug("queryValidatorCommission:", params, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	commission := k.GetValidatorAccumulatedCommission(ctx, params.ValidatorAddress)
	if commission.Commission == nil {
		commission.Commission = chainTypes.DecCoins{}
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, commission)
	ctx.Logger().Debug("queryValidatorCommission:", commission, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryValidatorSlashes(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryValidatorSlashesParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	ctx.Logger().Debug("queryValidatorSlashes:", params, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var events types.ValidatorSlashEvents
	k.IterateValidatorSlashEventsBetween(ctx, params.ValidatorAddress, params.StartingHeight, params.EndingHeight,
		func(height uint64, event types.ValidatorSlashEvent) (stop bool) {
			events.ValidatorSlashEvents = append(events.ValidatorSlashEvents, event)
			return false
		},
	)

	bz, err := codec.MarshalJSONIndent(k.cdc, events)
	ctx.Logger().Debug("queryValidatorSlashes:", events, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDelegationRewards(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDelegationRewardsParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()

	val := k.stakingKeeper.Validator(ctx, params.ValidatorAddress)
	ctx.Logger().Debug("queryDelegationRewards:", val, "err:", err)
	if val == nil {
		return nil, sdkerrors.Wrap(types.ErrNoValidatorExists, params.ValidatorAddress.String())
	}

	del := k.stakingKeeper.Delegation(ctx, params.DelegatorAddress, params.ValidatorAddress)
	ctx.Logger().Debug("queryDelegationRewards:", del)
	if del == nil {
		return nil, types.ErrNoDelegationExists
	}

	endingPeriod := k.IncrementValidatorPeriod(ctx, val)
	rewards := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)
	ctx.Logger().Debug("queryDelegationRewards:", rewards, "err:", err)
	if rewards == nil {
		rewards = chainTypes.DecCoins{}
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, rewards)
	ctx.Logger().Debug("queryDelegationRewards:", bz, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDelegatorTotalRewards(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDelegatorParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	ctx.Logger().Debug("queryDelegatorTotalRewards:", "params:", params, err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()

	total := chainTypes.DecCoins{}
	var delRewards []types.DelegationDelegatorReward

	k.stakingKeeper.IterateDelegations(
		ctx, params.DelegatorAddress,
		func(_ int64, del types.StakingExportedDelegationI) (stop bool) {
			valID := del.GetValidatorAccountID()
			val := k.stakingKeeper.Validator(ctx, valID)
			ctx.Logger().Debug("queryDelegatorTotalRewards", "valId", del.GetValidatorAccountID(), "del", del.GetDelegatorAccountID())
			endingPeriod := k.IncrementValidatorPeriod(ctx, val)
			delReward := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

			delRewards = append(delRewards, types.NewDelegationDelegatorReward(valID, delReward))
			total = total.Add(delReward...)

			ctx.Logger().Debug("total", total)
			return false
		},
	)

	totalRewards := types.NewQueryDelegatorTotalRewardsResponse(delRewards, total)

	bz, err := json.Marshal(totalRewards)
	ctx.Logger().Debug("queryDelegatorTotalRewards:", "totalRewards:", totalRewards, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDelegatorValidators(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDelegatorParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	ctx.Logger().Debug("queryDelegatorValidators:", "params:", params, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()

	var validators []AccountID

	k.stakingKeeper.IterateDelegations(
		ctx, params.DelegatorAddress,
		func(_ int64, del types.StakingExportedDelegationI) (stop bool) { // bugs , staking module interface
			validators = append(validators, del.GetValidatorAccountID())
			return false
		},
	)

	bz, err := codec.MarshalJSONIndent(k.cdc, validators)
	ctx.Logger().Debug("queryDelegatorValidators:", bz, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDelegatorWithdrawAddress(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDelegatorWithdrawAddrParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	ctx.Logger().Debug("queryDelegatorWithdrawAddress:", params, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()
	withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, params.DelegatorAddress)

	bz, err := codec.MarshalJSONIndent(k.cdc, withdrawAddr)
	ctx.Logger().Debug("queryDelegatorWithdrawAddress:", bz, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	info := types.NewWithDrawAddrInfo(withdrawAddr, chainTypes.EmptyAccountID(), params.DelegatorAddress)
	infoBz, err := codec.MarshalJSONIndent(k.cdc, info)
	ctx.Logger().Debug("queryDelegatorWithdrawAddress:", infoBz, "err:", err)

	return infoBz, err
}

func queryCommunityPool(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	pool := k.GetFeePoolCommunityCoins(ctx)
	ctx.Logger().Debug("queryDelegatorWithdrawAddress:", pool)
	if pool == nil {
		pool = chainTypes.DecCoins{}
	}

	bz, err := k.cdc.MarshalJSON(pool)
	ctx.Logger().Debug("queryDelegatorWithdrawAddress:", bz, "err:", err)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
