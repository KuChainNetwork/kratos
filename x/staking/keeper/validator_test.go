package keeper_test

import (
	"testing"

	"fmt"
	"time"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/staking/exported"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

//_______________________________________________________
func ValEq(t *testing.T, exp, got types.Validator) (*testing.T, bool, string, types.Validator, types.Validator) {
	return t, exp.MinEqual(got), "expected:\t%v\ngot:\t\t%v", exp, got
}

func TestValidator(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("TestSetValidator", t, func() {
		_, _, _, valAddr, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valPubKey := PKs[0]
		valTokens := exported.TokensFromConsensusPower(10)
		validator := types.NewValidator(valAddr, valPubKey, types.Description{})
		validator, _ = validator.AddTokensFromDel(valTokens)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))

		So(validator.Status == exported.Unbonded, ShouldBeTrue)
		So(valTokens.Equal(validator.Tokens), ShouldBeTrue)
		So(valTokens.Equal(validator.DelegatorShares.RoundInt()), ShouldBeTrue)
		keeper.SetValidator(ctx, validator)
		keeper.SetValidatorByPowerIndex(ctx, validator)

		// ensure update
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		validator, found := keeper.GetValidator(ctx, valAddr)
		So(found, ShouldBeTrue)
		So(len(updates) == 1, ShouldBeTrue)
		So(updates[0].Equal(validator.ABCIValidatorUpdate()), ShouldBeTrue)

		So(validator.Status == exported.Bonded, ShouldBeTrue)
		So(valTokens.Equal(validator.Tokens), ShouldBeTrue)
		So(valTokens.Equal(validator.DelegatorShares.RoundInt()), ShouldBeTrue)

		resVal, found := keeper.GetValidator(ctx, valAddr)
		So(validator.MinEqual(resVal), ShouldBeTrue)
		So(found, ShouldBeTrue)

		resVals := keeper.GetLastValidators(ctx)
		require.Equal(t, 1, len(resVals))
		So(validator.MinEqual(resVals[0]), ShouldBeTrue)
		resVals = keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, 1, len(resVals))
		require.True(ValEq(t, validator, resVals[0]))

		resVals = keeper.GetValidators(ctx, 1)
		require.Equal(t, 1, len(resVals))
		require.True(ValEq(t, validator, resVals[0]))

		resVals = keeper.GetValidators(ctx, 10)
		require.Equal(t, 1, len(resVals))
		require.True(ValEq(t, validator, resVals[0]))

		allVals := keeper.GetAllValidators(ctx)
		require.Equal(t, 1, len(allVals))
	})
	Convey("TestUpdateValidatorByPowerIndex", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		bondedPool := keeper.GetBondedPool(ctx)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		bondedToken := exported.TokensFromConsensusPower(1234)
		nobondedToken := exported.TokensFromConsensusPower(10000)
		app.AssetKeeper().IssueCoinPower(ctx, bondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), bondedToken)))
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken)))
		// add a validator
		validator := types.NewValidator(Accd[0], PKs[0], types.Description{})
		validator, delSharesCreated := validator.AddTokensFromDel(exported.TokensFromConsensusPower(100))
		require.Equal(t, exported.Unbonded, validator.Status)
		require.Equal(t, exported.TokensFromConsensusPower(100), validator.Tokens)
		TestingUpdateValidator(app, ctx, validator, true)
		validator, found := keeper.GetValidator(ctx, Accd[0])
		So(found, ShouldBeTrue)
		require.Equal(t, exported.TokensFromConsensusPower(100), validator.Tokens)
		power := types.GetValidatorsByPowerIndexKey(validator)
		require.True(t, keeper.HasValidatorByPowerIndex(ctx, power))

		// burn half the delegator shares
		keeper.DeleteValidatorByPowerIndex(ctx, validator)
		validator, burned := validator.RemoveDelShares(delSharesCreated.Quo(sdk.NewDec(2)))
		require.Equal(t, exported.TokensFromConsensusPower(50), burned)
		TestingUpdateValidator(app, ctx, validator, true) // update the validator, possibly kicking it out
		require.False(t, keeper.HasValidatorByPowerIndex(ctx, power))

		validator, found = keeper.GetValidator(ctx, Accd[0])
		So(found, ShouldBeTrue)

		power = types.GetValidatorsByPowerIndexKey(validator)
		require.True(t, keeper.HasValidatorByPowerIndex(ctx, power))
	})
	Convey("TestUpdateBondedValidatorsDecreaseCliff", t, func() {
		numVals := 10
		maxVals := 5

		// create context, keeper, and pool for tests
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		bondedPool := keeper.GetBondedPool(ctx)
		notBondedPool := keeper.GetNotBondedPool(ctx)

		// create keeper parameters
		params := keeper.GetParams(ctx)
		params.MaxValidators = uint32(maxVals)
		keeper.SetParams(ctx, params)

		// create a random pool
		bondedToken := exported.TokensFromConsensusPower(1234)
		nobondedToken := exported.TokensFromConsensusPower(10000)
		app.AssetKeeper().IssueCoinPower(ctx, bondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), bondedToken)))
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken)))

		validators := make([]types.Validator, numVals)
		for i := 0; i < len(validators); i++ {
			moniker := fmt.Sprintf("val#%d", int64(i))
			val := types.NewValidator(Accd[i], PKs[i], types.Description{Moniker: moniker})
			delTokens := exported.TokensFromConsensusPower(int64((i + 1) * 10))
			val, _ = val.AddTokensFromDel(delTokens)

			val = TestingUpdateValidator(app, ctx, val, true)
			validators[i] = val
		}

		nextCliffVal := validators[numVals-maxVals+1]

		// remove enough tokens to kick out the validator below the current cliff
		// validator and next in line cliff validator
		keeper.DeleteValidatorByPowerIndex(ctx, nextCliffVal)
		shares := exported.TokensFromConsensusPower(21)
		nextCliffVal, _ = nextCliffVal.RemoveDelShares(shares.ToDec())
		nextCliffVal = TestingUpdateValidator(app, ctx, nextCliffVal, true)

		expectedValStatus := map[int]exported.BondStatus{
			9: exported.Bonded, 8: exported.Bonded, 7: exported.Bonded, 5: exported.Bonded, 4: exported.Bonded,
			0: exported.Unbonding, 1: exported.Unbonding, 2: exported.Unbonding, 3: exported.Unbonding, 6: exported.Unbonding,
		}

		// require all the validators have their respective statuses
		for valIdx, status := range expectedValStatus {
			valAddr := validators[valIdx].OperatorAccount
			val, _ := keeper.GetValidator(ctx, valAddr)

			assert.Equal(
				t, status, val.GetStatus(),
				fmt.Sprintf("expected validator at index %v to have status: %s", valIdx, status),
			)
		}
	})
	Convey("TestSlashToZeroPowerRemoved", t, func() {
		// create context, keeper, and pool for tests
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		bondedPool := keeper.GetBondedPool(ctx)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		// create a random pool
		//	bondedToken := exported.TokensFromConsensusPower(1234)
		nobondedToken := exported.TokensFromConsensusPower(10000)
		valTokens := exported.TokensFromConsensusPower(100)
		app.AssetKeeper().IssueCoinPower(ctx, bondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken)))

		// add a validator
		validator := types.NewValidator(Accd[0], PKs[0], types.Description{})

		validator, _ = validator.AddTokensFromDel(valTokens)
		require.Equal(t, exported.Unbonded, validator.Status)
		require.Equal(t, valTokens, validator.Tokens)
		keeper.SetValidatorByConsAddr(ctx, validator)
		validator = TestingUpdateValidator(app, ctx, validator, true)
		require.Equal(t, valTokens, validator.Tokens, "\nvalidator %v\npool %v", validator, valTokens)

		// slash the validator by 100%
		consAddr0 := sdk.ConsAddress(PKs[0].Address())
		keeper.Slash(ctx, consAddr0, 0, 100, sdk.OneDec())
		// apply TM updates
		keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		// validator should be unbonding
		validator, _ = keeper.GetValidator(ctx, Accd[0])
		require.Equal(t, validator.GetStatus(), exported.Unbonding)
	})
	Convey("TestValidatorBasics", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		bondedPool := keeper.GetBondedPool(ctx)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		// create a random pool
		//	bondedToken := exported.TokensFromConsensusPower(1234)
		nobondedToken := exported.TokensFromConsensusPower(10000)
		valTokens := exported.TokensFromConsensusPower(100)
		app.AssetKeeper().IssueCoinPower(ctx, bondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken)))

		//construct the validators
		var validators [3]types.Validator
		powers := []int64{9, 8, 7}
		for i, power := range powers {
			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})
			validators[i].Status = exported.Unbonded
			validators[i].Tokens = sdk.ZeroInt()
			tokens := exported.TokensFromConsensusPower(power)

			validators[i], _ = validators[i].AddTokensFromDel(tokens)
		}
		assert.Equal(t, exported.TokensFromConsensusPower(9), validators[0].Tokens)
		assert.Equal(t, exported.TokensFromConsensusPower(8), validators[1].Tokens)
		assert.Equal(t, exported.TokensFromConsensusPower(7), validators[2].Tokens)

		// check the empty keeper first
		_, found := keeper.GetValidator(ctx, Accd[0])
		require.False(t, found)
		resVals := keeper.GetLastValidators(ctx)
		require.Zero(t, len(resVals))

		resVals = keeper.GetValidators(ctx, 2)
		require.Zero(t, len(resVals))

		// set and retrieve a record
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], true)
		keeper.SetValidatorByConsAddr(ctx, validators[0])
		resVal, found := keeper.GetValidator(ctx, Accd[0])
		So(found, ShouldBeTrue)
		assert.True(ValEq(t, validators[0], resVal))

		// retrieve from consensus
		resVal, found = keeper.GetValidatorByConsAddr(ctx, sdk.ConsAddress(PKs[0].Address()))
		So(found, ShouldBeTrue)
		assert.True(ValEq(t, validators[0], resVal))
		resVal, found = keeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(PKs[0]))
		So(found, ShouldBeTrue)
		assert.True(ValEq(t, validators[0], resVal))

		resVals = keeper.GetLastValidators(ctx)
		require.Equal(t, 1, len(resVals))
		assert.True(ValEq(t, validators[0], resVals[0]))
		assert.Equal(t, exported.Bonded, validators[0].Status)
		assert.True(sdk.IntEq(t, exported.TokensFromConsensusPower(9), validators[0].BondedTokens()))

		// modify a records, save, and retrieve
		validators[0].Status = exported.Bonded
		validators[0].Tokens = exported.TokensFromConsensusPower(10)
		validators[0].DelegatorShares = validators[0].Tokens.ToDec()
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], true)
		resVal, found = keeper.GetValidator(ctx, Accd[0])
		So(found, ShouldBeTrue)
		assert.True(ValEq(t, validators[0], resVal))

		resVals = keeper.GetLastValidators(ctx)
		require.Equal(t, 1, len(resVals))
		assert.True(ValEq(t, validators[0], resVals[0]))

		// add other validators
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], true)
		validators[2] = TestingUpdateValidator(app, ctx, validators[2], true)
		resVal, found = keeper.GetValidator(ctx, Accd[1])
		So(found, ShouldBeTrue)
		assert.True(ValEq(t, validators[1], resVal))
		resVal, found = keeper.GetValidator(ctx, Accd[2])
		So(found, ShouldBeTrue)
		assert.True(ValEq(t, validators[2], resVal))

		resVals = keeper.GetLastValidators(ctx)
		require.Equal(t, 3, len(resVals))
		assert.True(ValEq(t, validators[0], resVals[0])) // order doesn't matter here
		assert.True(ValEq(t, validators[1], resVals[1]))
		assert.True(ValEq(t, validators[2], resVals[2]))

		// remove a record

		// shouldn't be able to remove if status is not unbonded
		assert.PanicsWithValue(t,
			"cannot call RemoveValidator on bonded or unbonding validators",
			func() { keeper.RemoveValidator(ctx, validators[1].OperatorAccount) })

		// shouldn't be able to remove if there are still tokens left
		validators[1].Status = exported.Unbonded
		keeper.SetValidator(ctx, validators[1])
		assert.PanicsWithValue(t,
			"attempting to remove a validator which still contains tokens",
			func() { keeper.RemoveValidator(ctx, validators[1].OperatorAccount) })

		validators[1].Tokens = sdk.ZeroInt()                       // ...remove all tokens
		keeper.SetValidator(ctx, validators[1])                    // ...set the validator
		keeper.RemoveValidator(ctx, validators[1].OperatorAccount) // Now it can be removed.
		_, found = keeper.GetValidator(ctx, Accd[1])
		require.False(t, found)
	})
	Convey("TestGetValidatorSortingUnmixed", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		bondedPool := keeper.GetBondedPool(ctx)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		// create a random pool
		//	bondedToken := exported.TokensFromConsensusPower(1234)
		nobondedToken := exported.TokensFromConsensusPower(10000)
		valTokens := exported.TokensFromConsensusPower(100)
		app.AssetKeeper().IssueCoinPower(ctx, bondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken)))
		// initialize some validators into the state
		amts := []sdk.Int{
			sdk.NewInt(0),
			exported.PowerReduction.MulRaw(100),
			exported.PowerReduction,
			exported.PowerReduction.MulRaw(400),
			exported.PowerReduction.MulRaw(200)}
		n := len(amts)
		var validators [5]types.Validator
		for i, amt := range amts {
			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})
			validators[i].Status = exported.Bonded
			validators[i].Tokens = amt
			validators[i].DelegatorShares = sdk.NewDecFromInt(amt)
			TestingUpdateValidator(app, ctx, validators[i], true)
		}

		// first make sure everything made it in to the gotValidator group
		resValidators := keeper.GetBondedValidatorsByPower(ctx)
		assert.Equal(t, n, len(resValidators))
		assert.Equal(t, sdk.NewInt(400).Mul(exported.PowerReduction), resValidators[0].BondedTokens(), "%v", resValidators)
		assert.Equal(t, sdk.NewInt(200).Mul(exported.PowerReduction), resValidators[1].BondedTokens(), "%v", resValidators)
		assert.Equal(t, sdk.NewInt(100).Mul(exported.PowerReduction), resValidators[2].BondedTokens(), "%v", resValidators)
		assert.Equal(t, sdk.NewInt(1).Mul(exported.PowerReduction), resValidators[3].BondedTokens(), "%v", resValidators)
		assert.Equal(t, sdk.NewInt(0), resValidators[4].BondedTokens(), "%v", resValidators)
		assert.Equal(t, validators[3].OperatorAccount, resValidators[0].OperatorAccount, "%v", resValidators)
		assert.Equal(t, validators[4].OperatorAccount, resValidators[1].OperatorAccount, "%v", resValidators)
		assert.Equal(t, validators[1].OperatorAccount, resValidators[2].OperatorAccount, "%v", resValidators)
		assert.Equal(t, validators[2].OperatorAccount, resValidators[3].OperatorAccount, "%v", resValidators)
		assert.Equal(t, validators[0].OperatorAccount, resValidators[4].OperatorAccount, "%v", resValidators)

		// test a basic increase in voting power
		keeper.DeleteValidatorByPowerIndex(ctx, validators[3])
		validators[3].Tokens = sdk.NewInt(500).Mul(exported.PowerReduction)
		TestingUpdateValidator(app, ctx, validators[3], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, len(resValidators), n)
		assert.True(ValEq(t, validators[3], resValidators[0]))

		// test a decrease in voting power
		keeper.DeleteValidatorByPowerIndex(ctx, validators[3])
		validators[3].Tokens = sdk.NewInt(300).Mul(exported.PowerReduction)
		TestingUpdateValidator(app, ctx, validators[3], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)

		require.Equal(t, len(resValidators), n)
		assert.True(ValEq(t, validators[3], resValidators[0]))
		assert.True(ValEq(t, validators[4], resValidators[1]))

		// test equal voting power, different age
		keeper.DeleteValidatorByPowerIndex(ctx, validators[3])
		validators[3].Tokens = sdk.NewInt(200).Mul(exported.PowerReduction)
		ctx = ctx.WithBlockHeight(10)
		TestingUpdateValidator(app, ctx, validators[3], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, len(resValidators), n)
		assert.True(ValEq(t, validators[3], resValidators[0]))
		assert.True(ValEq(t, validators[4], resValidators[1]))

		// no change in voting power - no change in sort
		ctx = ctx.WithBlockHeight(20)
		keeper.DeleteValidatorByPowerIndex(ctx, validators[4])
		TestingUpdateValidator(app, ctx, validators[4], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, len(resValidators), n)
		assert.True(ValEq(t, validators[3], resValidators[0]))
		assert.True(ValEq(t, validators[4], resValidators[1]))

		// change in voting power of both validators, both still in v-set, no age change
		keeper.DeleteValidatorByPowerIndex(ctx, validators[3])
		validators[3].Tokens = sdk.NewInt(300).Mul(exported.PowerReduction)
		TestingUpdateValidator(app, ctx, validators[3], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, len(resValidators), n)
		ctx = ctx.WithBlockHeight(30)
		keeper.DeleteValidatorByPowerIndex(ctx, validators[4])
		validators[4].Tokens = sdk.NewInt(300).Mul(exported.PowerReduction)
		TestingUpdateValidator(app, ctx, validators[4], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, len(resValidators), n, "%v", resValidators)
		assert.True(ValEq(t, validators[3], resValidators[0]))
		assert.True(ValEq(t, validators[4], resValidators[1]))
	})
	Convey("TestGetValidatorSortingMixed", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		bondedPool := keeper.GetBondedPool(ctx)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		// create a random pool
		//	bondedToken := exported.TokensFromConsensusPower(1234)
		nobondedToken := exported.TokensFromConsensusPower(0)
		valTokens := exported.TokensFromConsensusPower(10000)
		app.AssetKeeper().IssueCoinPower(ctx, bondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), nobondedToken)))

		// now 2 max resValidators
		params := keeper.GetParams(ctx)
		params.MaxValidators = 2
		keeper.SetParams(ctx, params)

		// initialize some validators into the state
		amts := []sdk.Int{
			sdk.NewInt(0),
			exported.PowerReduction.MulRaw(100),
			exported.PowerReduction,
			exported.PowerReduction.MulRaw(400),
			exported.PowerReduction.MulRaw(200)}

		var validators [5]types.Validator
		for i, amt := range amts {
			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})
			validators[i].DelegatorShares = sdk.NewDecFromInt(amt)
			validators[i].Status = exported.Bonded
			validators[i].Tokens = amt
			TestingUpdateValidator(app, ctx, validators[i], true)
		}

		val0, found := keeper.GetValidator(ctx, Accd[0])
		So(found, ShouldBeTrue)
		val1, found := keeper.GetValidator(ctx, Accd[1])
		So(found, ShouldBeTrue)
		val2, found := keeper.GetValidator(ctx, Accd[2])
		So(found, ShouldBeTrue)
		val3, found := keeper.GetValidator(ctx, Accd[3])
		So(found, ShouldBeTrue)
		val4, found := keeper.GetValidator(ctx, Accd[4])
		So(found, ShouldBeTrue)
		require.Equal(t, exported.Bonded, val0.Status)
		require.Equal(t, exported.Unbonding, val1.Status)
		require.Equal(t, exported.Unbonding, val2.Status)
		require.Equal(t, exported.Bonded, val3.Status)
		require.Equal(t, exported.Bonded, val4.Status)

		// first make sure everything made it in to the gotValidator group
		resValidators := keeper.GetBondedValidatorsByPower(ctx)
		// The validators returned should match the max validators
		assert.Equal(t, 2, len(resValidators))
		assert.Equal(t, sdk.NewInt(400).Mul(exported.PowerReduction), resValidators[0].BondedTokens(), "%v", resValidators)
		assert.Equal(t, sdk.NewInt(200).Mul(exported.PowerReduction), resValidators[1].BondedTokens(), "%v", resValidators)
		assert.Equal(t, validators[3].OperatorAccount.StoreKey(), resValidators[0].OperatorAccount.StoreKey(), "%v", resValidators)
		assert.Equal(t, validators[4].OperatorAccount.StoreKey(), resValidators[1].OperatorAccount.StoreKey(), "%v", resValidators)
	})
	Convey("TestGetValidatorsEdgeCases", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		notBondedPool := keeper.GetNotBondedPool(ctx)

		// set max validators to 2
		params := keeper.GetParams(ctx)
		nMax := uint32(2)
		params.MaxValidators = nMax
		keeper.SetParams(ctx, params)

		// initialize some validators into the state
		powers := []int64{0, 100, 400, 400}
		var validators [4]types.Validator
		for i, power := range powers {
			moniker := fmt.Sprintf("val#%d", int64(i))
			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{Moniker: moniker})
			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)
			app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), tokens)))
			validators[i] = TestingUpdateValidator(app, ctx, validators[i], true)
		}

		// ensure that the first two bonded validators are the largest validators
		resValidators := keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, nMax, uint32(len(resValidators)))
		assert.True(ValEq(t, validators[2], resValidators[0]))
		assert.True(ValEq(t, validators[3], resValidators[1]))

		// delegate 500 tokens to validator 0
		keeper.DeleteValidatorByPowerIndex(ctx, validators[0])
		delTokens := exported.TokensFromConsensusPower(500)
		validators[0], _ = validators[0].AddTokensFromDel(delTokens)

		// test that the two largest validators are
		//   a) validator 0 with 500 tokens
		//   b) validator 2 with 400 tokens (delegated before validator 3)
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, nMax, uint32(len(resValidators)))
		assert.True(ValEq(t, validators[0], resValidators[0]))
		assert.True(ValEq(t, validators[2], resValidators[1]))

		// A validator which leaves the bonded validator set due to a decrease in voting power,
		// then increases to the original voting power, does not get its spot back in the
		// case of a tie.
		//
		// Order of operations for this test:
		//  - validator 3 enter validator set with 1 new token
		//  - validator 3 removed validator set by removing 201 tokens (validator 2 enters)
		//  - validator 3 adds 200 tokens (equal to validator 2 now) and does not get its spot back

		// validator 3 enters bonded validator set
		ctx = ctx.WithBlockHeight(40)

		validators[3], _ = keeper.GetValidator(ctx, validators[3].OperatorAccount)
		keeper.DeleteValidatorByPowerIndex(ctx, validators[3])
		validators[3], _ = validators[3].AddTokensFromDel(exported.TokensFromConsensusPower(1))

		notBondedPool = keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), exported.TokensFromConsensusPower(1))))

		validators[3] = TestingUpdateValidator(app, ctx, validators[3], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, nMax, uint32(len(resValidators)))
		assert.True(ValEq(t, validators[0], resValidators[0]))
		assert.True(ValEq(t, validators[3], resValidators[1]))

		// validator 3 kicked out temporarily
		keeper.DeleteValidatorByPowerIndex(ctx, validators[3])
		rmTokens := validators[3].TokensFromShares(sdk.NewDec(201)).TruncateInt()
		validators[3], _ = validators[3].RemoveDelShares(sdk.NewDec(201))

		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), rmTokens)))

		validators[3] = TestingUpdateValidator(app, ctx, validators[3], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, nMax, uint32(len(resValidators)))
		assert.True(ValEq(t, validators[0], resValidators[0]))
		assert.True(ValEq(t, validators[2], resValidators[1]))

		// validator 3 does not get spot back
		keeper.DeleteValidatorByPowerIndex(ctx, validators[3])
		validators[3], _ = validators[3].AddTokensFromDel(sdk.NewInt(200))

		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(200))))

		validators[3] = TestingUpdateValidator(app, ctx, validators[3], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, nMax, uint32(len(resValidators)))
		assert.True(ValEq(t, validators[0], resValidators[0]))
		assert.True(ValEq(t, validators[2], resValidators[1]))
		_, exists := keeper.GetValidator(ctx, validators[3].OperatorAccount)
		require.True(t, exists)
	})
	Convey("TestValidatorBondHeight", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))

		// now 2 max resValidators
		params := keeper.GetParams(ctx)
		params.MaxValidators = 2
		keeper.SetParams(ctx, params)

		// initialize some validators into the state
		var validators [3]types.Validator
		validators[0] = types.NewValidator(Accd[0], PKs[0], types.Description{})
		validators[1] = types.NewValidator(Accd[1], PKs[1], types.Description{})
		validators[2] = types.NewValidator(Accd[2], PKs[2], types.Description{})

		tokens0 := exported.TokensFromConsensusPower(200)
		tokens1 := exported.TokensFromConsensusPower(100)
		tokens2 := exported.TokensFromConsensusPower(100)
		validators[0], _ = validators[0].AddTokensFromDel(tokens0)
		validators[1], _ = validators[1].AddTokensFromDel(tokens1)
		validators[2], _ = validators[2].AddTokensFromDel(tokens2)

		validators[0] = TestingUpdateValidator(app, ctx, validators[0], true)

		////////////////////////////////////////
		// If two validators both increase to the same voting power in the same block,
		// the one with the first transaction should become bonded
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], true)
		validators[2] = TestingUpdateValidator(app, ctx, validators[2], true)

		resValidators := keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, uint32(len(resValidators)), params.MaxValidators)

		assert.True(ValEq(t, validators[0], resValidators[0]))
		assert.True(ValEq(t, validators[1], resValidators[1]))
		keeper.DeleteValidatorByPowerIndex(ctx, validators[1])
		keeper.DeleteValidatorByPowerIndex(ctx, validators[2])
		delTokens := exported.TokensFromConsensusPower(50)
		validators[1], _ = validators[1].AddTokensFromDel(delTokens)
		validators[2], _ = validators[2].AddTokensFromDel(delTokens)
		validators[2] = TestingUpdateValidator(app, ctx, validators[2], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)
		require.Equal(t, params.MaxValidators, uint32(len(resValidators)))
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], true)
		assert.True(ValEq(t, validators[0], resValidators[0]))
		assert.True(ValEq(t, validators[2], resValidators[1]))
	})
	Convey("TestFullValidatorSetPowerChange", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		params := keeper.GetParams(ctx)
		max := 2
		params.MaxValidators = uint32(2)
		keeper.SetParams(ctx, params)

		// initialize some validators into the state
		powers := []int64{0, 100, 400, 400, 200}
		var validators [5]types.Validator
		for i, power := range powers {
			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})
			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)
			TestingUpdateValidator(app, ctx, validators[i], true)
		}
		for i := range powers {
			var found bool
			validators[i], found = keeper.GetValidator(ctx, validators[i].OperatorAccount)
			So(found, ShouldBeTrue)
		}
		assert.Equal(t, exported.Unbonded, validators[0].Status)
		assert.Equal(t, exported.Unbonding, validators[1].Status)
		assert.Equal(t, exported.Bonded, validators[2].Status)
		assert.Equal(t, exported.Bonded, validators[3].Status)
		assert.Equal(t, exported.Unbonded, validators[4].Status)
		resValidators := keeper.GetBondedValidatorsByPower(ctx)
		assert.Equal(t, max, len(resValidators))
		assert.True(ValEq(t, validators[2], resValidators[0])) // in the order of txs
		assert.True(ValEq(t, validators[3], resValidators[1]))

		// test a swap in voting power

		tokens := exported.TokensFromConsensusPower(600)
		validators[0], _ = validators[0].AddTokensFromDel(tokens)
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], true)
		resValidators = keeper.GetBondedValidatorsByPower(ctx)
		assert.Equal(t, max, len(resValidators))
		assert.True(ValEq(t, validators[0], resValidators[0]))
		assert.True(ValEq(t, validators[2], resValidators[1]))
	})
	Convey("TestApplyAndReturnValidatorSetUpdatesAllNone", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))

		powers := []int64{10, 20}
		var validators [2]types.Validator
		for i, power := range powers {
			valPubKey := PKs[i+1]
			valAddr := Accd[i+1]

			validators[i] = types.NewValidator(valAddr, valPubKey, types.Description{})
			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)
		}

		// test from nothing to something
		//  tendermintUpdate set: {} -> {c1, c3}
		require.Equal(t, 0, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))
		keeper.SetValidator(ctx, validators[0])
		keeper.SetValidatorByPowerIndex(ctx, validators[0])
		keeper.SetValidator(ctx, validators[1])
		keeper.SetValidatorByPowerIndex(ctx, validators[1])

		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		assert.Equal(t, 2, len(updates))
		validators[0], _ = keeper.GetValidator(ctx, validators[0].OperatorAccount)
		validators[1], _ = keeper.GetValidator(ctx, validators[1].OperatorAccount)
		assert.Equal(t, validators[0].ABCIValidatorUpdate(), updates[1])
		assert.Equal(t, validators[1].ABCIValidatorUpdate(), updates[0])
	})
	Convey("TestApplyAndReturnValidatorSetUpdatesIdentical", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		powers := []int64{10, 20}
		var validators [2]types.Validator
		for i, power := range powers {
			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})

			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)

		}
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], false)
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], false)
		require.Equal(t, 2, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))

		// test identical,
		//  tendermintUpdate set: {} -> {}
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], false)
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], false)
		require.Equal(t, 0, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))
	})
	Convey("TestApplyAndReturnValidatorSetUpdatesSingleValueChange", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))

		powers := []int64{10, 20}
		var validators [2]types.Validator
		for i, power := range powers {

			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})

			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)

		}
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], false)
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], false)
		require.Equal(t, 2, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))

		// test single value change
		//  tendermintUpdate set: {} -> {c1'}
		keeper.DeleteValidatorByPowerIndex(ctx, validators[0])
		validators[0].Status = exported.Bonded
		validators[0].Tokens = exported.TokensFromConsensusPower(600)
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], false)

		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)

		So(len(updates) == 1, ShouldBeTrue)
		require.Equal(t, validators[0].ABCIValidatorUpdate(), updates[0])
	})
	Convey("TestApplyAndReturnValidatorSetUpdatesMultipleValueChange", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))

		powers := []int64{10, 20}
		var validators [2]types.Validator
		for i, power := range powers {

			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})

			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)

		}
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], false)
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], false)
		require.Equal(t, 2, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))
		keeper.DeleteValidatorByPowerIndex(ctx, validators[0])
		keeper.DeleteValidatorByPowerIndex(ctx, validators[1])

		// test multiple value change
		//  tendermintUpdate set: {c1, c3} -> {c1', c3'}
		delTokens1 := exported.TokensFromConsensusPower(190)
		delTokens2 := exported.TokensFromConsensusPower(80)
		validators[0], _ = validators[0].AddTokensFromDel(delTokens1)
		validators[1], _ = validators[1].AddTokensFromDel(delTokens2)
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], false)
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], false)

		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)

		require.Equal(t, 2, len(updates))
		require.Equal(t, validators[0].ABCIValidatorUpdate(), updates[0])
		require.Equal(t, validators[1].ABCIValidatorUpdate(), updates[1])
	})
	Convey("TestApplyAndReturnValidatorSetUpdatesInserted", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))

		powers := []int64{10, 20, 5, 15, 25}
		var validators [5]types.Validator
		for i, power := range powers {

			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})

			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)

		}

		validators[0] = TestingUpdateValidator(app, ctx, validators[0], false)
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], false)
		require.Equal(t, 2, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))

		// test validtor added at the beginning
		//  tendermintUpdate set: {} -> {c0}
		keeper.SetValidator(ctx, validators[2])
		keeper.SetValidatorByPowerIndex(ctx, validators[2])
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		validators[2], _ = keeper.GetValidator(ctx, validators[2].OperatorAccount)
		So(len(updates) == 1, ShouldBeTrue)
		require.Equal(t, validators[2].ABCIValidatorUpdate(), updates[0])

		// test validtor added at the beginning
		//  tendermintUpdate set: {} -> {c0}
		keeper.SetValidator(ctx, validators[3])
		keeper.SetValidatorByPowerIndex(ctx, validators[3])
		updates = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		validators[3], _ = keeper.GetValidator(ctx, validators[3].OperatorAccount)
		So(len(updates) == 1, ShouldBeTrue)
		require.Equal(t, validators[3].ABCIValidatorUpdate(), updates[0])

		// test validtor added at the end
		//  tendermintUpdate set: {} -> {c0}
		keeper.SetValidator(ctx, validators[4])
		keeper.SetValidatorByPowerIndex(ctx, validators[4])
		updates = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		validators[4], _ = keeper.GetValidator(ctx, validators[4].OperatorAccount)
		So(len(updates) == 1, ShouldBeTrue)
		require.Equal(t, validators[4].ABCIValidatorUpdate(), updates[0])
	})
	Convey("TestApplyAndReturnValidatorSetUpdatesWithCliffValidator", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		params := types.DefaultParams()
		params.MaxValidators = 2
		keeper.SetParams(ctx, params)

		powers := []int64{10, 20, 5}
		var validators [5]types.Validator
		for i, power := range powers {

			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})

			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)

		}
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], false)
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], false)
		require.Equal(t, 2, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))

		// test validator added at the end but not inserted in the valset
		//  tendermintUpdate set: {} -> {}
		TestingUpdateValidator(app, ctx, validators[2], false)
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		require.Equal(t, 0, len(updates))

		// test validator change its power and become a gotValidator (pushing out an existing)
		//  tendermintUpdate set: {}     -> {c0, c4}
		require.Equal(t, 0, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))

		tokens := exported.TokensFromConsensusPower(10)
		validators[2], _ = validators[2].AddTokensFromDel(tokens)
		keeper.SetValidator(ctx, validators[2])
		keeper.SetValidatorByPowerIndex(ctx, validators[2])
		updates = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		validators[2], _ = keeper.GetValidator(ctx, validators[2].OperatorAccount)
		require.Equal(t, 2, len(updates), "%v", updates)
		require.Equal(t, validators[0].ABCIValidatorUpdateZero(), updates[1])
		require.Equal(t, validators[2].ABCIValidatorUpdate(), updates[0])
	})
	Convey("TestApplyAndReturnValidatorSetUpdatesPowerDecrease", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))

		powers := []int64{100, 100}
		var validators [2]types.Validator
		for i, power := range powers {

			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})

			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)

		}
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], false)
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], false)
		require.Equal(t, 2, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))

		// check initial power
		require.Equal(t, int64(100), validators[0].GetConsensusPower())
		require.Equal(t, int64(100), validators[1].GetConsensusPower())

		// test multiple value change
		//  tendermintUpdate set: {c1, c3} -> {c1', c3'}
		keeper.DeleteValidatorByPowerIndex(ctx, validators[0])
		keeper.DeleteValidatorByPowerIndex(ctx, validators[1])
		delTokens1 := exported.TokensFromConsensusPower(20)
		delTokens2 := exported.TokensFromConsensusPower(30)
		validators[0], _ = validators[0].RemoveDelShares(delTokens1.ToDec())
		validators[1], _ = validators[1].RemoveDelShares(delTokens2.ToDec())
		validators[0] = TestingUpdateValidator(app, ctx, validators[0], false)
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], false)

		// power has changed
		require.Equal(t, int64(80), validators[0].GetConsensusPower())
		require.Equal(t, int64(70), validators[1].GetConsensusPower())

		// Tendermint updates should reflect power change
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		require.Equal(t, 2, len(updates))
		require.Equal(t, validators[0].ABCIValidatorUpdate(), updates[0])
		require.Equal(t, validators[1].ABCIValidatorUpdate(), updates[1])
	})
	Convey("TestApplyAndReturnValidatorSetUpdatesNewValidator", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))

		params := keeper.GetParams(ctx)
		params.MaxValidators = uint32(3)

		keeper.SetParams(ctx, params)

		powers := []int64{100, 100}
		var validators [2]types.Validator

		// initialize some validators into the state
		for i, power := range powers {

			valPubKey := PKs[i+1]
			valAddr := Accd[i+1]

			validators[i] = types.NewValidator(valAddr, valPubKey, types.Description{})
			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)

			keeper.SetValidator(ctx, validators[i])
			keeper.SetValidatorByPowerIndex(ctx, validators[i])
		}

		// verify initial Tendermint updates are correct
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		require.Equal(t, len(validators), len(updates))
		validators[0], _ = keeper.GetValidator(ctx, validators[0].OperatorAccount)
		validators[1], _ = keeper.GetValidator(ctx, validators[1].OperatorAccount)
		require.Equal(t, validators[0].ABCIValidatorUpdate(), updates[0])
		require.Equal(t, validators[1].ABCIValidatorUpdate(), updates[1])

		require.Equal(t, 0, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))

		// update initial validator set
		for i, power := range powers {

			keeper.DeleteValidatorByPowerIndex(ctx, validators[i])
			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)

			keeper.SetValidator(ctx, validators[i])
			keeper.SetValidatorByPowerIndex(ctx, validators[i])
		}

		// add a new validator that goes from zero power, to non-zero power, back to
		// zero power
		valPubKey := PKs[len(validators)+1]
		valAddr := Accd[len(validators)+1]
		amt := sdk.NewInt(100)

		validator := types.NewValidator(valAddr, valPubKey, types.Description{})
		validator, _ = validator.AddTokensFromDel(amt)

		keeper.SetValidator(ctx, validator)

		validator, _ = validator.RemoveDelShares(amt.ToDec())
		keeper.SetValidator(ctx, validator)
		keeper.SetValidatorByPowerIndex(ctx, validator)

		// add a new validator that increases in power
		valPubKey = PKs[len(validators)+2]
		valAddr = Accd[len(validators)+2]

		validator = types.NewValidator(valAddr, valPubKey, types.Description{})
		tokens := exported.TokensFromConsensusPower(500)
		validator, _ = validator.AddTokensFromDel(tokens)
		keeper.SetValidator(ctx, validator)
		keeper.SetValidatorByPowerIndex(ctx, validator)

		// verify initial Tendermint updates are correct
		updates = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		validator, _ = keeper.GetValidator(ctx, validator.OperatorAccount)
		validators[0], _ = keeper.GetValidator(ctx, validators[0].OperatorAccount)
		validators[1], _ = keeper.GetValidator(ctx, validators[1].OperatorAccount)
		require.Equal(t, len(validators)+1, len(updates))
		require.Equal(t, validator.ABCIValidatorUpdate(), updates[0])
		require.Equal(t, validators[0].ABCIValidatorUpdate(), updates[1])
		require.Equal(t, validators[1].ABCIValidatorUpdate(), updates[2])
	})
	Convey("TestApplyAndReturnValidatorSetUpdatesBondTransition", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		params := keeper.GetParams(ctx)
		params.MaxValidators = uint32(2)

		keeper.SetParams(ctx, params)

		powers := []int64{100, 200, 300}
		var validators [3]types.Validator

		// initialize some validators into the state
		for i, power := range powers {
			moniker := fmt.Sprintf("%d", i)
			valPubKey := PKs[i+1]
			valAddr := Accd[i+1]

			validators[i] = types.NewValidator(valAddr, valPubKey, types.Description{Moniker: moniker})
			tokens := exported.TokensFromConsensusPower(power)
			validators[i], _ = validators[i].AddTokensFromDel(tokens)
			keeper.SetValidator(ctx, validators[i])
			keeper.SetValidatorByPowerIndex(ctx, validators[i])
		}

		// verify initial Tendermint updates are correct
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		require.Equal(t, 2, len(updates))
		validators[2], _ = keeper.GetValidator(ctx, validators[2].OperatorAccount)
		validators[1], _ = keeper.GetValidator(ctx, validators[1].OperatorAccount)
		require.Equal(t, validators[2].ABCIValidatorUpdate(), updates[0])
		require.Equal(t, validators[1].ABCIValidatorUpdate(), updates[1])

		require.Equal(t, 0, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))

		// delegate to validator with lowest power but not enough to bond
		ctx = ctx.WithBlockHeight(1)

		var found bool
		validators[0], found = keeper.GetValidator(ctx, validators[0].OperatorAccount)
		So(found, ShouldBeTrue)

		keeper.DeleteValidatorByPowerIndex(ctx, validators[0])
		tokens := exported.TokensFromConsensusPower(1)
		validators[0], _ = validators[0].AddTokensFromDel(tokens)
		keeper.SetValidator(ctx, validators[0])
		keeper.SetValidatorByPowerIndex(ctx, validators[0])

		// verify initial Tendermint updates are correct
		require.Equal(t, 0, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))

		// create a series of events that will bond and unbond the validator with
		// lowest power in a single block context (height)
		ctx = ctx.WithBlockHeight(2)

		validators[1], found = keeper.GetValidator(ctx, validators[1].OperatorAccount)
		So(found, ShouldBeTrue)

		keeper.DeleteValidatorByPowerIndex(ctx, validators[0])
		validators[0], _ = validators[0].RemoveDelShares(validators[0].DelegatorShares)
		keeper.SetValidator(ctx, validators[0])
		keeper.SetValidatorByPowerIndex(ctx, validators[0])
		updates = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		require.Equal(t, 0, len(updates))

		keeper.DeleteValidatorByPowerIndex(ctx, validators[1])
		tokens = exported.TokensFromConsensusPower(250)
		validators[1], _ = validators[1].AddTokensFromDel(tokens)
		keeper.SetValidator(ctx, validators[1])
		keeper.SetValidatorByPowerIndex(ctx, validators[1])

		// verify initial Tendermint updates are correct
		updates = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		So(len(updates) == 1, ShouldBeTrue)
		require.Equal(t, validators[1].ABCIValidatorUpdate(), updates[0])

		require.Equal(t, 0, len(keeper.ApplyAndReturnValidatorSetUpdates(ctx)))
	})
	Convey("TestApplyAndReturnValidatorSetUpdatesBondTransition", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})

		commission1 := types.NewCommissionWithTime(
			sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(3, 1),
			sdk.NewDecWithPrec(1, 1), time.Now().UTC().Add(time.Duration(-1)*time.Hour),
		)
		commission2 := types.NewCommission(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(3, 1), sdk.NewDecWithPrec(1, 1))

		val1 := types.NewValidator(Accd[0], PKs[0], types.Description{})
		val2 := types.NewValidator(Accd[1], PKs[1], types.Description{})

		val1, _ = val1.SetInitialCommission(commission1)
		val2, _ = val2.SetInitialCommission(commission2)

		keeper.SetValidator(ctx, val1)
		keeper.SetValidator(ctx, val2)

		testCases := []struct {
			validator   types.Validator
			newRate     sdk.Dec
			expectedErr bool
		}{
			{val1, sdk.ZeroDec(), true},
			{val2, sdk.NewDecWithPrec(-1, 1), true},
			{val2, sdk.NewDecWithPrec(4, 1), true},
			{val2, sdk.NewDecWithPrec(3, 1), true},
			{val2, sdk.NewDecWithPrec(2, 1), false},
		}

		for i, tc := range testCases {
			commission, err := keeper.UpdateValidatorCommission(ctx, tc.validator, tc.newRate)

			if tc.expectedErr {
				require.Error(t, err, "expected error for test case #%d with rate: %s", i, tc.newRate)
			} else {
				tc.validator.Commission = commission
				keeper.SetValidator(ctx, tc.validator)
				val, found := keeper.GetValidator(ctx, tc.validator.OperatorAccount)

				require.True(t, found,
					"expected to find validator for test case #%d with rate: %s", i, tc.newRate,
				)
				require.NoError(t, err,
					"unexpected error for test case #%d with rate: %s", i, tc.newRate,
				)
				require.Equal(t, tc.newRate, val.Commission.Rate,
					"expected new validator commission rate for test case #%d with rate: %s", i, tc.newRate,
				)
				require.Equal(t, ctx.BlockHeader().Time, val.Commission.UpdateTime,
					"expected new validator commission update time for test case #%d with rate: %s", i, tc.newRate,
				)
			}
		}
	})
}
