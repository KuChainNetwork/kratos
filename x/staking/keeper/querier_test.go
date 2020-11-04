package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/staking/exported"
	stakeKeeprer "github.com/KuChainNetwork/kuchain/x/staking/keeper"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	addrAcc1, addrAcc2 = Accd[0], Accd[1]
	addrVal1, addrVal2 = Accd[0], Accd[1]
	pk1, pk2           = PKs[0], PKs[1]
)

func TestQuerier(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("TestNewQuerier", t, func() {
		cdc := codec.New()
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1}) // Create Validators
		amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8)}
		var validators [2]types.Validator
		for i, amt := range amts {
			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})
			validators[i], _ = validators[i].AddTokensFromDel(amt)
			keeper.SetValidator(ctx, validators[i])
			keeper.SetValidatorByPowerIndex(ctx, validators[i])
		}

		header := abci.Header{
			ChainID: "HelloChain",
			Height:  5,
		}
		hi := types.NewHistoricalInfo(header, validators[:])
		keeper.SetHistoricalInfo(ctx, 5, hi)

		query := abci.RequestQuery{
			Path: "",
			Data: []byte{},
		}

		querier := stakeKeeprer.NewQuerier(*keeper)

		bz, err := querier(ctx, []string{"other"}, query)
		require.Error(t, err)
		require.Nil(t, bz)

		_, err = querier(ctx, []string{"pool"}, query)
		require.NoError(t, err)

		_, err = querier(ctx, []string{"parameters"}, query)
		require.NoError(t, err)

		queryValParams := types.NewQueryValidatorParams(addrVal1)
		bz, errRes := cdc.MarshalJSON(queryValParams)
		require.NoError(t, errRes)

		query.Path = "/custom/kustaking/validator"
		query.Data = bz

		_, err = querier(ctx, []string{"validator"}, query)
		require.NoError(t, err)

		_, err = querier(ctx, []string{"validatorDelegations"}, query)
		require.NoError(t, err)

		_, err = querier(ctx, []string{"validatorUnbondingDelegations"}, query)
		require.NoError(t, err)

		queryDelParams := types.NewQueryDelegatorParams(addrAcc2)
		bz, errRes = cdc.MarshalJSON(queryDelParams)
		require.NoError(t, errRes)

		query.Path = "/custom/kustaking/validator"
		query.Data = bz

		_, err = querier(ctx, []string{"delegatorDelegations"}, query)
		require.NoError(t, err)

		_, err = querier(ctx, []string{"delegatorUnbondingDelegations"}, query)
		require.NoError(t, err)

		_, err = querier(ctx, []string{"delegatorValidators"}, query)
		require.NoError(t, err)

		bz, errRes = cdc.MarshalJSON(types.NewQueryRedelegationParams(chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID()))
		require.NoError(t, errRes)
		query.Data = bz

		_, err = querier(ctx, []string{"redelegations"}, query)
		require.NoError(t, err)

		queryHisParams := types.NewQueryHistoricalInfoParams(5)
		bz, errRes = cdc.MarshalJSON(queryHisParams)
		require.NoError(t, errRes)

		query.Path = "/custom/kustaking/historicalInfo"
		query.Data = bz

		_, err = querier(ctx, []string{"historicalInfo"}, query)
		require.NoError(t, err)
	})
	Convey("TestQueryParametersPool", t, func() {
		cdc := codec.New()
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1}) // Create Validators		bondDenom := keeper.BondDenom(ctx)

		querier := stakeKeeprer.NewQuerier(*keeper)

		query := abci.RequestQuery{
			Path: "",
			Data: []byte{},
		}

		res, err := querier(ctx, []string{"parameters"}, query)
		require.NoError(t, err)

		var params types.Params
		errRes := cdc.UnmarshalJSON(res, &params)
		require.NoError(t, errRes)
		require.Equal(t, keeper.GetParams(ctx), params)

		res, err = querier(ctx, []string{"pool"}, query)
		require.NoError(t, err)

		var pool types.Pool
		oldBondedAmount := keeper.TotalBondedTokens(ctx)
		oldNotBondedAmount := keeper.TotalNotBondedTokens(ctx)

		errRes = cdc.UnmarshalJSON(res, &pool)
		require.NoError(t, errRes)
		require.Equal(t, oldBondedAmount, pool.BondedTokens)
		require.Equal(t, oldNotBondedAmount, pool.NotBondedTokens)
	})
	Convey("TestQueryValidators", t, func() {
		cdc := codec.New()
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1}) // Create Validators		bondDenom := keeper.BondDenom(ctx)
		params := keeper.GetParams(ctx)

		// Create Validators
		amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
		status := []exported.BondStatus{exported.Bonded, exported.Unbonded, exported.Unbonding}
		var validators [3]types.Validator
		for i, amt := range amts {
			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})
			validators[i], _ = validators[i].AddTokensFromDel(amt)
			validators[i] = validators[i].UpdateStatus(status[i])
		}

		keeper.SetValidator(ctx, validators[0])
		keeper.SetValidator(ctx, validators[1])
		keeper.SetValidator(ctx, validators[2])

		// Query Validators
		queriedValidators := keeper.GetValidators(ctx, params.MaxValidators)
		querier := stakeKeeprer.NewQuerier(*keeper)

		for i, s := range status {
			queryValsParams := types.NewQueryValidatorsParams(1, int(params.MaxValidators), s.String())
			bz, err := cdc.MarshalJSON(queryValsParams)
			require.NoError(t, err)

			req := abci.RequestQuery{
				Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryValidators),
				Data: bz,
			}

			res, err := querier(ctx, []string{"validators"}, req)
			require.NoError(t, err)

			var validatorsResp []types.Validator
			err = cdc.UnmarshalJSON(res, &validatorsResp)
			require.NoError(t, err)

			require.Equal(t, 1, len(validatorsResp))
			require.ElementsMatch(t, validators[i].OperatorAccount.Bytes(), validatorsResp[0].OperatorAccount.Bytes())
		}
		// Query each validator
		queryParams := types.NewQueryValidatorParams(addrVal1)
		bz, err := cdc.MarshalJSON(queryParams)
		require.NoError(t, err)

		query := abci.RequestQuery{
			Path: "/custom/kustaking/validator",
			Data: bz,
		}
		res, err := querier(ctx, []string{"validator"}, query)
		require.NoError(t, err)

		var validator types.Validator
		err = cdc.UnmarshalJSON(res, &validator)
		require.NoError(t, err)

		require.Equal(t, queriedValidators[0], validator)
	})
	Convey("TestQueryDelegation", t, func() {
		cdc := codec.New()
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		bondedToken := exported.TokensFromConsensusPower(1234)
		nobondedToken := exported.TokensFromConsensusPower(10000)
		bondedPool := keeper.GetBondedPool(ctx)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, bondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), bondedToken)))
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken)))
		ModuleName, _ := types.ModuleAccountID.ToName()
		app.AssetKeeper().Issue(ctx, ModuleName, ModuleName, chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken))
		params := keeper.GetParams(ctx)

		// Create Validators and Delegation
		val1 := types.NewValidator(addrVal1, pk1, types.Description{})
		keeper.SetValidator(ctx, val1)
		keeper.SetValidatorByPowerIndex(ctx, val1)

		val2 := types.NewValidator(addrVal2, pk2, types.Description{})
		keeper.SetValidator(ctx, val2)
		keeper.SetValidatorByPowerIndex(ctx, val2)
		delTokens := exported.TokensFromConsensusPower(20)
		keeper.Delegate(ctx, addrAcc2, delTokens, exported.Unbonded, val1, true)

		// apply TM updates
		keeper.ApplyAndReturnValidatorSetUpdates(ctx)

		// Query Delegator bonded validators
		queryParams := types.NewQueryDelegatorParams(addrAcc2)
		bz, errRes := cdc.MarshalJSON(queryParams)
		require.NoError(t, errRes)

		query := abci.RequestQuery{
			Path: "/custom/kustaking/delegatorValidators",
			Data: bz,
		}
		querier := stakeKeeprer.NewQuerier(*keeper)

		delValidators := keeper.GetDelegatorValidators(ctx, addrAcc2, params.MaxValidators)

		res, err := querier(ctx, []string{"delegatorValidators"}, query)
		//res, err := queryDelegatorValidators(ctx, query, keeper)
		require.NoError(t, err)

		var validatorsResp []types.Validator
		errRes = cdc.UnmarshalJSON(res, &validatorsResp)
		require.NoError(t, errRes)

		require.Equal(t, len(delValidators), len(validatorsResp))
		require.ElementsMatch(t, delValidators, validatorsResp)

		// error unknown request
		query.Data = bz[:len(bz)-1]

		_, err = querier(ctx, []string{"delegatorValidators"}, query)
		require.Error(t, err)

		// Query bonded validator
		queryBondParams := types.NewQueryBondsParams(addrAcc2, addrVal1)
		bz, errRes = cdc.MarshalJSON(queryBondParams)
		require.NoError(t, errRes)

		query = abci.RequestQuery{
			Path: "/custom/kustaking/delegatorValidator",
			Data: bz,
		}

		res, err = querier(ctx, []string{"delegatorValidator"}, query)
		require.NoError(t, err)

		var validator types.Validator
		errRes = cdc.UnmarshalJSON(res, &validator)
		require.NoError(t, errRes)

		require.Equal(t, delValidators[0], validator)

		// error unknown request
		query.Data = bz[:len(bz)-1]
		_, err = querier(ctx, []string{"delegatorValidator"}, query)
		require.Error(t, err)

		// Query delegation

		query = abci.RequestQuery{
			Path: "/custom/kustaking/delegation",
			Data: bz,
		}

		delegation, found := keeper.GetDelegation(ctx, addrAcc2, addrVal1)
		require.True(t, found)

		res, err = querier(ctx, []string{"delegation"}, query)
		require.NoError(t, err)

		var delegationRes types.DelegationResponse
		errRes = cdc.UnmarshalJSON(res, &delegationRes)
		require.NoError(t, errRes)

		require.Equal(t, delegation.ValidatorAccount, delegationRes.ValidatorAccount)
		require.Equal(t, delegation.DelegatorAccount, delegationRes.DelegatorAccount)
		require.Equal(t, chainTypes.NewCoin(keeper.BondDenom(ctx), delegation.Shares.TruncateInt()), delegationRes.Balance)

		// Query Delegator Delegations
		query = abci.RequestQuery{
			Path: "/custom/kustaking/delegatorDelegations",
			Data: bz,
		}
		res, err = querier(ctx, []string{"delegatorDelegations"}, query)
		require.NoError(t, err)

		var delegatorDelegations types.DelegationResponses
		errRes = cdc.UnmarshalJSON(res, &delegatorDelegations)
		require.NoError(t, errRes)
		require.Len(t, delegatorDelegations, 1)
		require.Equal(t, delegation.ValidatorAccount, delegatorDelegations[0].ValidatorAccount)
		require.Equal(t, delegation.DelegatorAccount, delegatorDelegations[0].DelegatorAccount)
		require.Equal(t, chainTypes.NewCoin(keeper.BondDenom(ctx), delegation.Shares.TruncateInt()), delegatorDelegations[0].Balance)

		// error unknown request
		query.Data = bz[:len(bz)-1]

		_, err = querier(ctx, []string{"delegation"}, query)
		require.Error(t, err)

		// Query validator delegations

		bz, errRes = cdc.MarshalJSON(types.NewQueryValidatorParams(addrVal1))
		require.NoError(t, errRes)

		query = abci.RequestQuery{
			Path: "custom/kustaking/validatorDelegations",
			Data: bz,
		}

		res, err = querier(ctx, []string{"validatorDelegations"}, query)
		require.NoError(t, err)

		var delegationsRes types.DelegationResponses
		errRes = cdc.UnmarshalJSON(res, &delegationsRes)
		require.NoError(t, errRes)
		require.Len(t, delegatorDelegations, 1)
		require.Equal(t, delegation.ValidatorAccount, delegationsRes[0].ValidatorAccount)
		require.Equal(t, delegation.DelegatorAccount, delegationsRes[0].DelegatorAccount)
		require.Equal(t, chainTypes.NewCoin(keeper.BondDenom(ctx), delegation.Shares.TruncateInt()), delegationsRes[0].Balance)

		// Query unbonging delegation
		unbondingTokens := sdk.TokensFromConsensusPower(10)
		_, err = keeper.Undelegate(ctx, addrAcc2, val1.OperatorAccount, unbondingTokens.ToDec())
		require.NoError(t, err)

		queryBondParams = types.NewQueryBondsParams(addrAcc2, addrVal1)
		bz, errRes = cdc.MarshalJSON(queryBondParams)
		require.NoError(t, errRes)

		query = abci.RequestQuery{
			Path: "/custom/kustaking/unbondingDelegation",
			Data: bz,
		}

		unbond, found := keeper.GetUnbondingDelegation(ctx, addrAcc2, addrVal1)
		require.True(t, found)

		res, err = querier(ctx, []string{"unbondingDelegation"}, query)
		require.NoError(t, err)

		var unbondRes types.UnbondingDelegation
		errRes = cdc.UnmarshalJSON(res, &unbondRes)
		require.NoError(t, errRes)

		require.Equal(t, unbond, unbondRes)

		// error unknown request
		query.Data = bz[:len(bz)-1]

		_, err = querier(ctx, []string{"unbondingDelegation"}, query)
		require.Error(t, err)

		// Query Delegator Delegations

		query = abci.RequestQuery{
			Path: "/custom/kustaking/delegatorUnbondingDelegations",
			Data: bz,
		}

		res, err = querier(ctx, []string{"delegatorUnbondingDelegations"}, query)
		require.NoError(t, err)

		var delegatorUbds []types.UnbondingDelegation
		errRes = cdc.UnmarshalJSON(res, &delegatorUbds)
		require.NoError(t, errRes)
		require.Equal(t, unbond, delegatorUbds[0])

		// error unknown request
		query.Data = bz[:len(bz)-1]

		_, err = querier(ctx, []string{"delegatorUnbondingDelegations"}, query)
		require.Error(t, err)

		// Query redelegation
		redelegationTokens := sdk.TokensFromConsensusPower(10)
		_, err = keeper.BeginRedelegation(ctx, addrAcc2, val1.OperatorAccount,
			val2.OperatorAccount, redelegationTokens.ToDec())
		require.NoError(t, err)
		redel, found := keeper.GetRedelegation(ctx, addrAcc2, val1.OperatorAccount, val2.OperatorAccount)
		require.True(t, found)

		bz, errRes = cdc.MarshalJSON(types.NewQueryRedelegationParams(addrAcc2, val1.OperatorAccount, val2.OperatorAccount))
		require.NoError(t, errRes)

		query = abci.RequestQuery{
			Path: "/custom/kustaking/redelegations",
			Data: bz,
		}

		res, err = querier(ctx, []string{"redelegations"}, query)
		require.NoError(t, err)

		var redelRes types.RedelegationResponses
		errRes = cdc.UnmarshalJSON(res, &redelRes)
		require.NoError(t, errRes)
		require.Len(t, redelRes, 1)
		require.Equal(t, redel.DelegatorAccount, redelRes[0].DelegatorAccount)
		require.Equal(t, redel.ValidatorSrcAccount, redelRes[0].ValidatorSrcAccount)
		require.Equal(t, redel.ValidatorDstAccount, redelRes[0].ValidatorDstAccount)
		require.Len(t, redel.Entries, len(redelRes[0].Entries))
	})
	Convey("TestQueryRedelegations", t, func() {
		cdc := codec.New()
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		bondedToken := exported.TokensFromConsensusPower(1234)
		nobondedToken := exported.TokensFromConsensusPower(10000)
		bondedPool := keeper.GetBondedPool(ctx)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, bondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), bondedToken)))
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken)))
		ModuleName, _ := types.ModuleAccountID.ToName()
		app.AssetKeeper().Issue(ctx, ModuleName, ModuleName, chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken))
		// Create Validators and Delegation
		val1 := types.NewValidator(addrVal1, pk1, types.Description{})
		val2 := types.NewValidator(addrVal2, pk2, types.Description{})
		keeper.SetValidator(ctx, val1)
		keeper.SetValidator(ctx, val2)

		delAmount := exported.TokensFromConsensusPower(100)
		keeper.Delegate(ctx, addrAcc2, delAmount, exported.Unbonded, val1, true)
		_ = keeper.ApplyAndReturnValidatorSetUpdates(ctx)

		rdAmount := exported.TokensFromConsensusPower(20)
		keeper.BeginRedelegation(ctx, addrAcc2, val1.GetOperatorAccountID(), val2.GetOperatorAccountID(), rdAmount.ToDec())
		keeper.ApplyAndReturnValidatorSetUpdates(ctx)

		redel, found := keeper.GetRedelegation(ctx, addrAcc2, val1.GetOperatorAccountID(), val2.GetOperatorAccountID())
		require.True(t, found)

		// delegator redelegations
		queryDelegatorParams := types.NewQueryDelegatorParams(addrAcc2)
		bz, errRes := cdc.MarshalJSON(queryDelegatorParams)
		require.NoError(t, errRes)

		query := abci.RequestQuery{
			Path: "/custom/kustaking/redelegations",
			Data: bz,
		}
		querier := stakeKeeprer.NewQuerier(*keeper)

		res, err := querier(ctx, []string{"redelegations"}, query)
		require.NoError(t, err)

		var redelRes types.RedelegationResponses
		errRes = cdc.UnmarshalJSON(res, &redelRes)
		require.NoError(t, errRes)
		require.Len(t, redelRes, 1)
		require.Equal(t, redel.DelegatorAccount, redelRes[0].DelegatorAccount)
		require.Equal(t, redel.ValidatorSrcAccount, redelRes[0].ValidatorSrcAccount)
		require.Equal(t, redel.ValidatorDstAccount, redelRes[0].ValidatorDstAccount)
		require.Len(t, redel.Entries, len(redelRes[0].Entries))

		// validator redelegations
		queryValidatorParams := types.NewQueryValidatorParams(val1.GetOperatorAccountID())
		bz, errRes = cdc.MarshalJSON(queryValidatorParams)
		require.NoError(t, errRes)

		query = abci.RequestQuery{
			Path: "/custom/kustaking/redelegations",
			Data: bz,
		}

		res, err = querier(ctx, []string{"redelegations"}, query)
		require.NoError(t, err)

		errRes = cdc.UnmarshalJSON(res, &redelRes)
		require.NoError(t, errRes)
		require.Len(t, redelRes, 1)
		require.Equal(t, redel.DelegatorAccount, redelRes[0].DelegatorAccount)
		require.Equal(t, redel.ValidatorSrcAccount, redelRes[0].ValidatorSrcAccount)
		require.Equal(t, redel.ValidatorDstAccount, redelRes[0].ValidatorDstAccount)
		require.Len(t, redel.Entries, len(redelRes[0].Entries))
	})
	Convey("TestQueryUnbondingDelegation", t, func() {
		cdc := codec.New()
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		bondedToken := exported.TokensFromConsensusPower(1234)
		nobondedToken := exported.TokensFromConsensusPower(10000)
		bondedPool := keeper.GetBondedPool(ctx)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, bondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), bondedToken)))
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken)))
		ModuleName, _ := types.ModuleAccountID.ToName()
		app.AssetKeeper().Issue(ctx, ModuleName, ModuleName, chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken))
		querier := stakeKeeprer.NewQuerier(*keeper)

		// Create Validators and Delegation
		val1 := types.NewValidator(addrVal1, pk1, types.Description{})
		keeper.SetValidator(ctx, val1)

		// delegate
		delAmount := sdk.TokensFromConsensusPower(100)
		_, err := keeper.Delegate(ctx, addrAcc1, delAmount, exported.Unbonded, val1, true)
		require.NoError(t, err)
		_ = keeper.ApplyAndReturnValidatorSetUpdates(ctx)

		// undelegate
		undelAmount := sdk.TokensFromConsensusPower(20)
		_, err = keeper.Undelegate(ctx, addrAcc1, val1.GetOperatorAccountID(), undelAmount.ToDec())
		require.NoError(t, err)
		keeper.ApplyAndReturnValidatorSetUpdates(ctx)

		_, found := keeper.GetUnbondingDelegation(ctx, addrAcc1, val1.OperatorAccount)
		require.True(t, found)

		//
		// found: query unbonding delegation by delegator and validator
		//
		queryValidatorParams := types.NewQueryBondsParams(addrAcc1, val1.GetOperatorAccountID())
		bz, errRes := cdc.MarshalJSON(queryValidatorParams)
		require.NoError(t, errRes)
		query := abci.RequestQuery{
			Path: "/custom/kustaking/unbondingDelegation",
			Data: bz,
		}
		res, err := querier(ctx, []string{"unbondingDelegation"}, query)
		require.NoError(t, err)
		require.NotNil(t, res)
		var ubDel types.UnbondingDelegation
		require.NoError(t, cdc.UnmarshalJSON(res, &ubDel))
		require.Equal(t, addrAcc1, ubDel.DelegatorAccount)
		require.Equal(t, val1.OperatorAccount, ubDel.ValidatorAccount)
		require.Equal(t, 1, len(ubDel.Entries))

		//
		// not found: query unbonding delegation by delegator and validator
		//
		queryValidatorParams = types.NewQueryBondsParams(addrAcc2, val1.GetOperatorAccountID())
		bz, errRes = cdc.MarshalJSON(queryValidatorParams)
		require.NoError(t, errRes)
		query = abci.RequestQuery{
			Path: "/custom/kustaking/unbondingDelegation",
			Data: bz,
		}
		_, err = querier(ctx, []string{"unbondingDelegation"}, query)
		require.Error(t, err)

		//
		// found: query unbonding delegation by delegator and validator
		//
		queryDelegatorParams := types.NewQueryDelegatorParams(addrAcc1)
		bz, errRes = cdc.MarshalJSON(queryDelegatorParams)
		require.NoError(t, errRes)
		query = abci.RequestQuery{
			Path: "/custom/kustaking/delegatorUnbondingDelegations",
			Data: bz,
		}
		res, err = querier(ctx, []string{"delegatorUnbondingDelegations"}, query)
		require.NoError(t, err)
		require.NotNil(t, res)
		var ubDels []types.UnbondingDelegation
		require.NoError(t, cdc.UnmarshalJSON(res, &ubDels))
		require.Equal(t, 1, len(ubDels))
		require.Equal(t, addrAcc1, ubDels[0].DelegatorAccount)
		require.Equal(t, val1.OperatorAccount, ubDels[0].ValidatorAccount)

		//
		// not found: query unbonding delegation by delegator and validator
		//
		queryDelegatorParams = types.NewQueryDelegatorParams(addrAcc2)
		bz, errRes = cdc.MarshalJSON(queryDelegatorParams)
		require.NoError(t, errRes)
		query = abci.RequestQuery{
			Path: "/custom/kustaking/delegatorUnbondingDelegations",
			Data: bz,
		}
		res, err = querier(ctx, []string{"delegatorUnbondingDelegations"}, query)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NoError(t, cdc.UnmarshalJSON(res, &ubDels))
		require.Equal(t, 0, len(ubDels))
	})
	Convey("TestQueryHistoricalInfo", t, func() {
		cdc := codec.New()
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		bondedToken := exported.TokensFromConsensusPower(1234)
		nobondedToken := exported.TokensFromConsensusPower(10000)
		bondedPool := keeper.GetBondedPool(ctx)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, bondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), bondedToken)))
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken)))
		ModuleName, _ := types.ModuleAccountID.ToName()
		app.AssetKeeper().Issue(ctx, ModuleName, ModuleName, chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken))
		querier := stakeKeeprer.NewQuerier(*keeper)

		// Create Validators and Delegation
		val1 := types.NewValidator(addrVal1, pk1, types.Description{})
		val2 := types.NewValidator(addrVal2, pk2, types.Description{})
		vals := []types.Validator{val1, val2}
		keeper.SetValidator(ctx, val1)
		keeper.SetValidator(ctx, val2)

		header := abci.Header{
			ChainID: "HelloChain",
			Height:  5,
		}
		hi := types.NewHistoricalInfo(header, vals)
		keeper.SetHistoricalInfo(ctx, 5, hi)

		queryHistoricalParams := types.NewQueryHistoricalInfoParams(4)
		bz, errRes := cdc.MarshalJSON(queryHistoricalParams)
		require.NoError(t, errRes)
		query := abci.RequestQuery{
			Path: "/custom/kustaking/historicalInfo",
			Data: bz,
		}
		res, err := querier(ctx, []string{"historicalInfo"}, query)
		require.Error(t, err, "Invalid query passed")
		require.Nil(t, res, "Invalid query returned non-nil result")

		queryHistoricalParams = types.NewQueryHistoricalInfoParams(5)
		bz, errRes = cdc.MarshalJSON(queryHistoricalParams)
		require.NoError(t, errRes)
		query.Data = bz
		res, err = querier(ctx, []string{"historicalInfo"}, query)
		require.NoError(t, err, "Valid query passed")
		require.NotNil(t, res, "Valid query returned nil result")

		var recv types.HistoricalInfo
		require.NoError(t, cdc.UnmarshalJSON(res, &recv))
		require.Equal(t, hi, recv, "HistoricalInfo query returned wrong result")
	})
}
