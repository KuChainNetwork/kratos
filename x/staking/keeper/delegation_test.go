package keeper_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/x/staking/exported"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	. "github.com/smartystreets/goconvey/convey"
	"time"
)

func TestDelegation(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("TestDelegation", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
		var validators [3]types.Validator
		for i, amt := range amts {
			validators[i] = types.NewValidator(Accd[i], PKs[i], types.Description{})
			validators[i], _ = validators[i].AddTokensFromDel(amt)
		}
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

		validators[0] = TestingUpdateValidator(app, ctx, validators[0], true)
		validators[1] = TestingUpdateValidator(app, ctx, validators[1], true)
		validators[2] = TestingUpdateValidator(app, ctx, validators[2], true)

		bond1to1 := types.NewDelegation(Accdel[0], Accd[0], sdk.NewDec(9))
		keeper := app.StakeKeeper()

		// check the empty keeper first
		_, found := keeper.GetDelegation(ctx, Accdel[0], Accd[0])
		So(found, ShouldBeFalse)

		// set and retrieve a record
		keeper.SetDelegation(ctx, bond1to1)
		resBond, found := keeper.GetDelegation(ctx, Accdel[0], Accd[0])
		So(found, ShouldBeTrue)
		So(bond1to1.Equal(resBond), ShouldBeTrue)

		// modify a records, save, and retrieve
		bond1to1.Shares = sdk.NewDec(9)
		keeper.SetDelegation(ctx, bond1to1)
		resBond, found = keeper.GetDelegation(ctx, Accdel[0], Accd[0])
		So(found, ShouldBeTrue)
		So(bond1to1.Equal(resBond), ShouldBeTrue)

		// add some more records
		bond1to2 := types.NewDelegation(Accdel[0], Accd[1], sdk.NewDec(9))
		bond1to3 := types.NewDelegation(Accdel[0], Accd[2], sdk.NewDec(9))
		bond2to1 := types.NewDelegation(Accdel[1], Accd[0], sdk.NewDec(9))
		bond2to2 := types.NewDelegation(Accdel[1], Accd[1], sdk.NewDec(9))
		bond2to3 := types.NewDelegation(Accdel[1], Accd[2], sdk.NewDec(9))
		keeper.SetDelegation(ctx, bond1to2)
		keeper.SetDelegation(ctx, bond1to3)
		keeper.SetDelegation(ctx, bond2to1)
		keeper.SetDelegation(ctx, bond2to2)
		keeper.SetDelegation(ctx, bond2to3)

		// test all bond retrieve capabilities
		resBonds := keeper.GetDelegatorDelegations(ctx, Accdel[0], 5)

		So(len(resBonds) == 3, ShouldBeTrue)
		So(bond1to1.Equal(resBonds[0]), ShouldBeTrue)
		So(bond1to2.Equal(resBonds[1]), ShouldBeTrue)
		So(bond1to3.Equal(resBonds[2]), ShouldBeTrue)

		resBonds = keeper.GetAllDelegatorDelegations(ctx, Accdel[0])
		So(len(resBonds) == 3, ShouldBeTrue)

		resBonds = keeper.GetDelegatorDelegations(ctx, Accdel[0], 2)
		So(len(resBonds) == 2, ShouldBeTrue)

		resBonds = keeper.GetDelegatorDelegations(ctx, Accdel[1], 5)
		So(len(resBonds) == 3, ShouldBeTrue)
		So(bond2to1.Equal(resBonds[0]), ShouldBeTrue)
		So(bond2to2.Equal(resBonds[1]), ShouldBeTrue)
		So(bond2to3.Equal(resBonds[2]), ShouldBeTrue)

		allBonds := keeper.GetAllDelegations(ctx)
		So(len(allBonds) == 6, ShouldBeTrue)
		So(bond1to1.Equal(allBonds[0]), ShouldBeTrue)
		So(bond1to2.Equal(allBonds[1]), ShouldBeTrue)
		So(bond1to3.Equal(allBonds[2]), ShouldBeTrue)
		So(bond2to1.Equal(allBonds[3]), ShouldBeTrue)
		So(bond2to2.Equal(allBonds[4]), ShouldBeTrue)
		So(bond2to3.Equal(allBonds[5]), ShouldBeTrue)

		resVals := keeper.GetDelegatorValidators(ctx, Accdel[0], 3)
		So(len(resVals) == 3, ShouldBeTrue)

		resVals = keeper.GetDelegatorValidators(ctx, Accdel[1], 4)
		So(len(resVals) == 3, ShouldBeTrue)

		for i := 0; i < 3; i++ {
			resVal, err := keeper.GetDelegatorValidator(ctx, Accdel[0], Accd[i])
			So(err, ShouldBeNil)
			So(Accd[i].Eq(resVal.GetOperatorAccountID()), ShouldBeTrue)

			resVal, err = keeper.GetDelegatorValidator(ctx, Accdel[1], Accd[i])
			So(err, ShouldBeNil)
			So(Accd[i].Eq(resVal.GetOperatorAccountID()), ShouldBeTrue)

			resDels := keeper.GetValidatorDelegations(ctx, Accd[i])
			So(len(resDels) == 2, ShouldBeTrue)
		}

		// delete a record
		keeper.RemoveDelegation(ctx, bond2to3)
		_, found = keeper.GetDelegation(ctx, Accdel[1], Accd[2])
		So(found, ShouldBeFalse)
		resBonds = keeper.GetDelegatorDelegations(ctx, Accdel[1], 5)
		So(len(resBonds) == 2, ShouldBeTrue)
		So(bond2to1.Equal(resBonds[0]), ShouldBeTrue)
		So(bond2to2.Equal(resBonds[1]), ShouldBeTrue)

		resBonds = keeper.GetAllDelegatorDelegations(ctx, Accdel[1])
		So(len(resBonds) == 2, ShouldBeTrue)

		// delete all the records from delegator 2
		keeper.RemoveDelegation(ctx, bond2to1)
		keeper.RemoveDelegation(ctx, bond2to2)
		_, found = keeper.GetDelegation(ctx, Accdel[1], Accd[0])
		So(found, ShouldBeFalse)
		_, found = keeper.GetDelegation(ctx, Accdel[1], Accd[1])
		So(found, ShouldBeFalse)
		resBonds = keeper.GetDelegatorDelegations(ctx, Accdel[1], 5)
		So(len(resBonds) == 0, ShouldBeTrue)
	})
	Convey("TestUnbondingDelegation", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		ubd := types.NewUnbondingDelegation(Accdel[0], Accd[0], 0,
			time.Unix(0, 0), sdk.NewInt(5))
		keeper := app.StakeKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

		// set and retrieve a record
		keeper.SetUnbondingDelegation(ctx, ubd)
		resUnbond, found := keeper.GetUnbondingDelegation(ctx, Accdel[0], Accd[0])
		So(found, ShouldBeTrue)
		So(ubd.Equal(resUnbond), ShouldBeTrue)

		// modify a records, save, and retrieve
		ubd.Entries[0].Balance = sdk.NewInt(21)
		keeper.SetUnbondingDelegation(ctx, ubd)

		resUnbonds := keeper.GetUnbondingDelegations(ctx, Accdel[0], 5)
		So(len(resUnbonds) == 1, ShouldBeTrue)

		resUnbonds = keeper.GetAllUnbondingDelegations(ctx, Accdel[0])
		So(len(resUnbonds) == 1, ShouldBeTrue)

		resUnbond, found = keeper.GetUnbondingDelegation(ctx, Accdel[0], Accd[0])
		So(found, ShouldBeTrue)
		So(ubd.Equal(resUnbond), ShouldBeTrue)

		// delete a record
		keeper.RemoveUnbondingDelegation(ctx, ubd)
		_, found = keeper.GetUnbondingDelegation(ctx, Accdel[0], Accd[0])
		So(found, ShouldBeFalse)

		resUnbonds = keeper.GetUnbondingDelegations(ctx, Accdel[0], 5)
		So(len(resUnbonds) == 0, ShouldBeTrue)

		resUnbonds = keeper.GetAllUnbondingDelegations(ctx, Accdel[0])
		So(len(resUnbonds) == 0, ShouldBeTrue)
	})
	Convey("TestUnbondDelegation", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		startTokens := exported.TokensFromConsensusPower(10)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), startTokens)))

		// create a validator and a delegator to that validator
		// note this validator starts not-bonded
		validator := types.NewValidator(Accd[0], PKs[0], types.Description{})

		validator, issuedShares := validator.AddTokensFromDel(startTokens)
		So(startTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)

		validator = TestingUpdateValidator(app, ctx, validator, true)

		delegation := types.NewDelegation(Accdel[0], Accd[0], issuedShares)
		keeper.SetDelegation(ctx, delegation)

		bondTokens := exported.TokensFromConsensusPower(6)
		amount, err := keeper.Unbond(ctx, Accdel[0], Accd[0], bondTokens.ToDec()) // no start info
		So(err, ShouldBeNil)
		So(bondTokens.Equal(amount), ShouldBeTrue)

		delegation, found := keeper.GetDelegation(ctx, Accdel[0], Accd[0])
		So(found, ShouldBeTrue)
		validator, found = keeper.GetValidator(ctx, Accd[0])
		So(found, ShouldBeTrue)

		remainingTokens := startTokens.Sub(bondTokens)
		So(remainingTokens.Equal(delegation.Shares.RoundInt()), ShouldBeTrue)
		So(remainingTokens.Equal(validator.BondedTokens()), ShouldBeTrue)

	})
	Convey("TestUnbondingDelegationsMaxEntries", t, func() {
		_, _, _, addAlice, addJack, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		startTokens := exported.TokensFromConsensusPower(10)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		BondedPool := keeper.GetBondedPool(ctx)
		bondDenom := keeper.BondDenom(ctx)

		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), startTokens)))

		// create a validator and a delegator to that validator
		validator := types.NewValidator(addAlice, PKs[0], types.Description{})

		validator, issuedShares := validator.AddTokensFromDel(startTokens)
		So(startTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)

		validator = TestingUpdateValidator(app, ctx, validator, true)
		delegation := types.NewDelegation(addJack, addAlice, issuedShares)
		keeper.SetDelegation(ctx, delegation)

		maxEntries := keeper.MaxEntries(ctx)

		oldBonded := app.AssetKeeper().GetCoinPowerByDenomd(ctx, BondedPool.GetID(), bondDenom).Amount
		oldNotBonded := app.AssetKeeper().GetCoinPowerByDenomd(ctx, notBondedPool.GetID(), bondDenom).Amount

		// should all pass
		var completionTime time.Time
		for i := uint32(0); i < maxEntries; i++ {
			var err error
			completionTime, err = keeper.Undelegate(ctx, addJack, addAlice, sdk.NewDec(1))
			So(err, ShouldBeNil)
		}

		So(app.AssetKeeper().GetCoinPowerByDenomd(ctx, BondedPool.GetID(), bondDenom).Amount.Equal(oldBonded.SubRaw(int64(maxEntries))), ShouldBeTrue)
		So(app.AssetKeeper().GetCoinPowerByDenomd(ctx, notBondedPool.GetID(), bondDenom).Amount.Equal(oldNotBonded.AddRaw(int64(maxEntries))), ShouldBeTrue)

		oldBonded = app.AssetKeeper().GetCoinPowerByDenomd(ctx, BondedPool.GetID(), bondDenom).Amount
		oldNotBonded = app.AssetKeeper().GetCoinPowerByDenomd(ctx, notBondedPool.GetID(), bondDenom).Amount

		// an additional unbond should fail due to max entries
		_, err := keeper.Undelegate(ctx, addJack, addAlice, sdk.NewDec(1))
		So(err, ShouldNotBeNil)

		So(app.AssetKeeper().GetCoinPowerByDenomd(ctx, BondedPool.GetID(), bondDenom).Amount.Equal(oldBonded), ShouldBeTrue)
		So(app.AssetKeeper().GetCoinPowerByDenomd(ctx, notBondedPool.GetID(), bondDenom).Amount.Equal(oldNotBonded), ShouldBeTrue)

		// mature unbonding delegations
		ctx = ctx.WithBlockTime(completionTime)
		_, err = keeper.CompleteUnbonding(ctx, addJack, addAlice)

		So(err, ShouldBeNil)
		So(app.AssetKeeper().GetCoinPowerByDenomd(ctx, BondedPool.GetID(), bondDenom).Amount.Equal(oldBonded), ShouldBeTrue)
		So(app.AssetKeeper().GetCoinPowerByDenomd(ctx, notBondedPool.GetID(), bondDenom).Amount.Equal(oldNotBonded), ShouldBeFalse)

		oldNotBonded = app.AssetKeeper().GetCoinPowerByDenomd(ctx, notBondedPool.GetID(), bondDenom).Amount

		// unbonding  should work again
		_, err = keeper.Undelegate(ctx, addJack, addAlice, sdk.NewDec(1))
		So(err, ShouldBeNil)
		So(app.AssetKeeper().GetCoinPowerByDenomd(ctx, BondedPool.GetID(), bondDenom).Amount.Equal(oldBonded.SubRaw(int64(1))), ShouldBeTrue)
		So(app.AssetKeeper().GetCoinPowerByDenomd(ctx, notBondedPool.GetID(), bondDenom).Amount.Equal(oldNotBonded.AddRaw(int64(1))), ShouldBeTrue)
	})
	Convey("TestUndelegateFromUnbondingValidator", t, func() {
		_, _, _, addAlice, addJack, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		delTokens := exported.TokensFromConsensusPower(10)
		unbondedToken := sdk.NewInt(10000)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		BondedPool := keeper.GetBondedPool(ctx)
		validator := types.NewValidator(addAlice, PKs[0], types.Description{})

		validator, issuedShares := validator.AddTokensFromDel(delTokens)
		So(delTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		// add bonded tokens to pool for delegations
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))
		validator = TestingUpdateValidator(app, ctx, validator, true)
		So(validator.IsBonded(), ShouldBeTrue)

		selfDelegation := types.NewDelegation(addJack, addAlice, issuedShares)
		keeper.SetDelegation(ctx, selfDelegation)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))

		// create a second delegation to this validator
		keeper.DeleteValidatorByPowerIndex(ctx, validator)
		validator, _ = validator.AddTokensFromDel(unbondedToken)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))

		validator = TestingUpdateValidator(app, ctx, validator, true)
		validator, issuedShares = validator.AddTokensFromDel(unbondedToken)
		delegation := types.NewDelegation(addAlice, addAlice, issuedShares)
		keeper.SetDelegation(ctx, delegation)

		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))

		header := ctx.BlockHeader()
		blockHeight := int64(10)
		header.Height = blockHeight
		blockTime := time.Unix(333, 0).UTC()
		header.Time = blockTime
		ctx = ctx.WithBlockHeader(header)

		_, err := keeper.Undelegate(ctx, addJack, addAlice, delTokens.ToDec())
		So(err, ShouldBeNil)
		// end block
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		So(len(updates) == 1, ShouldBeTrue)

		validator, found := keeper.GetValidator(ctx, addAlice)
		So(found, ShouldBeTrue)

		So(blockHeight == validator.UnbondingHeight, ShouldBeTrue)
		params := keeper.GetParams(ctx)

		So(blockTime.Add(params.UnbondingTime).Equal(validator.UnbondingTime), ShouldBeTrue)

		blockHeight2 := int64(20)
		blockTime2 := time.Unix(444, 0).UTC()
		ctx = ctx.WithBlockHeight(blockHeight2)
		ctx = ctx.WithBlockTime(blockTime2)

		validator, found = keeper.GetValidator(ctx, addAlice)
		So(found, ShouldBeTrue)
		// unbond some of the other delegation's shares
		_, err = keeper.Undelegate(ctx, addAlice, addAlice, sdk.NewDec(6))
		So(err, ShouldBeNil)

		// retrieve the unbonding delegation
		ubd, found := keeper.GetUnbondingDelegation(ctx, addAlice, addAlice)
		So(found, ShouldBeTrue)
		So(len(ubd.Entries) == 1, ShouldBeTrue)
		So(ubd.Entries[0].Balance.Equal(sdk.NewInt(6)), ShouldBeTrue)
		So(blockTime2.Add(params.UnbondingTime).Equal(ubd.Entries[0].CompletionTime), ShouldBeTrue)
		So(blockHeight2 == ubd.Entries[0].CreationHeight, ShouldBeTrue)
	})
	Convey("TestUndelegateFromUnbondedValidator", t, func() {
		_, _, _, addAlice, addJack, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		delTokens := exported.TokensFromConsensusPower(10)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		BondedPool := keeper.GetBondedPool(ctx)
		validator := types.NewValidator(addAlice, PKs[0], types.Description{})
		validator, issuedShares := validator.AddTokensFromDel(delTokens)
		So(delTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		// add bonded tokens to pool for delegations
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))
		validator = TestingUpdateValidator(app, ctx, validator, true)
		So(validator.IsBonded(), ShouldBeTrue)

		selfDelegation := types.NewDelegation(addAlice, addAlice, issuedShares)
		keeper.SetDelegation(ctx, selfDelegation)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))

		keeper.DeleteValidatorByPowerIndex(ctx, validator)
		validator, issuedShares = validator.AddTokensFromDel(delTokens)
		validator = TestingUpdateValidator(app, ctx, validator, true)
		delegation := types.NewDelegation(addJack, addAlice, issuedShares)
		keeper.SetDelegation(ctx, delegation)

		ctx = ctx.WithBlockHeight(10)
		ctx = ctx.WithBlockTime(time.Unix(333, 0))

		// unbond the all self-delegation to put validator in unbonding state
		_, err := keeper.Undelegate(ctx, addAlice, addAlice, delTokens.ToDec())
		So(err, ShouldBeNil)

		// end block
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		So(len(updates) == 1, ShouldBeTrue)

		validator, found := keeper.GetValidator(ctx, addAlice)
		So(found, ShouldBeTrue)
		So(ctx.BlockHeight() == validator.UnbondingHeight, ShouldBeTrue)
		params := keeper.GetParams(ctx)
		So(ctx.BlockHeader().Time.Add(params.UnbondingTime).Equal(validator.UnbondingTime), ShouldBeTrue)

		// unbond the validator
		ctx = ctx.WithBlockTime(validator.UnbondingTime)
		keeper.UnbondAllMatureValidatorQueue(ctx)

		// Make sure validator is still in state because there is still an outstanding delegation
		validator, found = keeper.GetValidator(ctx, addAlice)
		So(found, ShouldBeTrue)
		So(validator.Status.Equal(exported.Unbonded), ShouldBeTrue)

		// unbond some of the other delegation's shares
		unbondTokens := sdk.TokensFromConsensusPower(6)
		_, err = keeper.Undelegate(ctx, addJack, addAlice, unbondTokens.ToDec())
		So(err, ShouldBeNil)

		// unbond rest of the other delegation's shares
		remainingTokens := delTokens.Sub(unbondTokens)
		_, err = keeper.Undelegate(ctx, addJack, addAlice, remainingTokens.ToDec())
		So(err, ShouldBeNil)

		//  now validator should now be deleted from state
		validator, found = keeper.GetValidator(ctx, addAlice)
		So(found, ShouldBeFalse)
	})
	Convey("TestUnbondingAllDelegationFromValidator", t, func() {
		_, _, _, addAlice, addJack, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		delTokens := exported.TokensFromConsensusPower(10)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		BondedPool := keeper.GetBondedPool(ctx)
		validator := types.NewValidator(addAlice, PKs[0], types.Description{})
		validator, issuedShares := validator.AddTokensFromDel(delTokens)
		So(delTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		// add bonded tokens to pool for delegations
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))
		validator = TestingUpdateValidator(app, ctx, validator, true)
		So(validator.IsBonded(), ShouldBeTrue)
		selfDelegation := types.NewDelegation(addAlice, addAlice, issuedShares)
		keeper.SetDelegation(ctx, selfDelegation)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))

		keeper.DeleteValidatorByPowerIndex(ctx, validator)
		validator, issuedShares = validator.AddTokensFromDel(delTokens)
		validator = TestingUpdateValidator(app, ctx, validator, true)
		delegation := types.NewDelegation(addJack, addAlice, issuedShares)
		keeper.SetDelegation(ctx, delegation)

		ctx = ctx.WithBlockHeight(10)
		ctx = ctx.WithBlockTime(time.Unix(333, 0).UTC())
		// unbond the all self-delegation to put validator in unbonding state
		_, err := keeper.Undelegate(ctx, addAlice, addAlice, delTokens.ToDec())
		So(err, ShouldBeNil)

		// end block
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		So(len(updates) == 1, ShouldBeTrue)
		//require.Equal(t, 1, len(updates))

		// unbond all the remaining delegation
		_, err = keeper.Undelegate(ctx, addJack, addAlice, delTokens.ToDec())
		So(err, ShouldBeNil)

		// validator should still be in state and still be in unbonding state
		validator, found := keeper.GetValidator(ctx, addAlice)
		So(found, ShouldBeTrue)
		So(validator.Status == exported.Unbonding, ShouldBeTrue)

		// unbond the validator
		ctx = ctx.WithBlockTime(validator.UnbondingTime)
		keeper.UnbondAllMatureValidatorQueue(ctx)

		// validator should now be deleted from state
		_, found = keeper.GetValidator(ctx, Accd[0])
		So(found, ShouldBeFalse)
	})
	Convey("TestGetRedelegationsFromSrcValidator", t, func() {
		_, _, _, addAlice, addJack, addValidator, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		rd := types.NewRedelegation(addJack, addAlice, addValidator, 0,
			time.Unix(0, 0), sdk.NewInt(5),
			sdk.NewDec(5))

		// set and retrieve a record
		keeper.SetRedelegation(ctx, rd)
		resBond, found := keeper.GetRedelegation(ctx, addJack, addAlice, addValidator)
		So(found, ShouldBeTrue)

		// get the redelegations one time
		redelegations := keeper.GetRedelegationsFromSrcValidator(ctx, addAlice)
		So(len(redelegations) == 1, ShouldBeTrue)
		So(redelegations[0].Equal(resBond), ShouldBeTrue)

		// get the redelegations a second time, should be exactly the same
		redelegations = keeper.GetRedelegationsFromSrcValidator(ctx, addAlice)
		So(len(redelegations) == 1, ShouldBeTrue)
		So(redelegations[0].Equal(resBond), ShouldBeTrue)
	})
	Convey("TestRedelegation", t, func() {
		_, _, _, addAlice, addJack, addValidator, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

		rd := types.NewRedelegation(addJack, addAlice, addValidator, 0,
			time.Unix(0, 0), sdk.NewInt(5),
			sdk.NewDec(5))

		// test shouldn't have and redelegations
		has := keeper.HasReceivingRedelegation(ctx, addJack, addValidator)
		So(has, ShouldBeFalse)

		// set and retrieve a record
		keeper.SetRedelegation(ctx, rd)
		resRed, found := keeper.GetRedelegation(ctx, addJack, addAlice, addValidator)
		So(found, ShouldBeTrue)

		redelegations := keeper.GetRedelegationsFromSrcValidator(ctx, addAlice)
		So(len(redelegations) == 1, ShouldBeTrue)
		So(redelegations[0].Equal(resRed), ShouldBeTrue)

		redelegations = keeper.GetRedelegations(ctx, addJack, 5)
		So(len(redelegations) == 1, ShouldBeTrue)
		So(redelegations[0].Equal(resRed), ShouldBeTrue)

		redelegations = keeper.GetAllRedelegations(ctx, addJack, chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID())
		So(len(redelegations) == 1, ShouldBeTrue)
		So(redelegations[0].Equal(resRed), ShouldBeTrue)

		// check if has the redelegation
		has = keeper.HasReceivingRedelegation(ctx, addJack, addValidator)
		So(has, ShouldBeTrue)
		// modify a records, save, and retrieve
		rd.Entries[0].SharesDst = sdk.NewDec(21)
		keeper.SetRedelegation(ctx, rd)

		resRed, found = keeper.GetRedelegation(ctx, addJack, addAlice, addValidator)
		So(found, ShouldBeTrue)
		So(rd.Equal(resRed), ShouldBeTrue)

		redelegations = keeper.GetRedelegationsFromSrcValidator(ctx, addAlice)
		So(len(redelegations) == 1, ShouldBeTrue)
		So(redelegations[0].Equal(resRed), ShouldBeTrue)

		redelegations = keeper.GetRedelegations(ctx, addJack, 5)
		So(len(redelegations) == 1, ShouldBeTrue)
		So(redelegations[0].Equal(resRed), ShouldBeTrue)

		// delete a record
		keeper.RemoveRedelegation(ctx, rd)
		_, found = keeper.GetRedelegation(ctx, addJack, addAlice, addValidator)
		So(found, ShouldBeFalse)

		redelegations = keeper.GetRedelegations(ctx, addJack, 5)
		So(len(redelegations) == 0, ShouldBeTrue)

		redelegations = keeper.GetAllRedelegations(ctx, addJack, chainTypes.EmptyAccountID(), chainTypes.EmptyAccountID())
		So(len(redelegations) == 0, ShouldBeTrue)
	})
	Convey("TestRedelegateToSameValidator", t, func() {
		_, _, _, addAlice, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		delTokens := exported.TokensFromConsensusPower(10)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		BondedPool := keeper.GetBondedPool(ctx)
		validator := types.NewValidator(addAlice, PKs[0], types.Description{})
		validator, issuedShares := validator.AddTokensFromDel(delTokens)
		So(delTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		// add bonded tokens to pool for delegations
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))
		validator = TestingUpdateValidator(app, ctx, validator, true)
		So(validator.IsBonded(), ShouldBeTrue)
		selfDelegation := types.NewDelegation(addAlice, addAlice, issuedShares)
		keeper.SetDelegation(ctx, selfDelegation)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))

		_, err := keeper.BeginRedelegation(ctx, addAlice, addAlice, addAlice, sdk.NewDec(5))
		So(err, ShouldNotBeNil)
	})
	Convey("TestRedelegationMaxEntries", t, func() {
		_, _, _, addAlice, addJack, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		delTokens := exported.TokensFromConsensusPower(20)
		valTokens := exported.TokensFromConsensusPower(10)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		BondedPool := keeper.GetBondedPool(ctx)
		validator := types.NewValidator(addAlice, PKs[0], types.Description{})
		validator, issuedShares := validator.AddTokensFromDel(delTokens)
		So(delTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		// add bonded tokens to pool for delegations
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))
		validator = TestingUpdateValidator(app, ctx, validator, true)
		selfDelegation := types.NewDelegation(addAlice, addAlice, issuedShares)
		keeper.SetDelegation(ctx, selfDelegation)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), delTokens)))
		validator2 := types.NewValidator(addJack, PKs[1], types.Description{})
		validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		So(valTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		validator2 = TestingUpdateValidator(app, ctx, validator2, true)
		So(validator2.Status == exported.Bonded, ShouldBeTrue)
		maxEntries := keeper.MaxEntries(ctx)
		// redelegations should pass
		var completionTime time.Time
		for i := uint32(0); i < maxEntries; i++ {
			var err error
			completionTime, err = keeper.BeginRedelegation(ctx, addAlice, addAlice, addJack, sdk.NewDec(1))
			So(err, ShouldBeNil)
		}
		// an additional redelegation should fail due to max entries
		_, err := keeper.BeginRedelegation(ctx, addAlice, addAlice, addJack, sdk.NewDec(1))
		So(err, ShouldNotBeNil)

		ctx = ctx.WithBlockTime(completionTime)
		_, err = keeper.CompleteRedelegation(ctx, addAlice, addAlice, addJack)
		So(err, ShouldBeNil)
		// redelegation should work again
		_, err = keeper.BeginRedelegation(ctx, addAlice, addAlice, addJack, sdk.NewDec(1))
		So(err, ShouldBeNil)
	})
	Convey("TestRedelegateSelfDelegation", t, func() {
		_, _, _, addAlice, addJack, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		valTokens := exported.TokensFromConsensusPower(10)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		BondedPool := keeper.GetBondedPool(ctx)
		validator := types.NewValidator(addAlice, PKs[0], types.Description{})
		validator, issuedShares := validator.AddTokensFromDel(valTokens)
		So(valTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		// add bonded tokens to pool for delegations
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		validator = TestingUpdateValidator(app, ctx, validator, true)
		selfDelegation := types.NewDelegation(addAlice, addAlice, issuedShares)
		keeper.SetDelegation(ctx, selfDelegation)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		validator2 := types.NewValidator(addJack, PKs[1], types.Description{})
		validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		So(valTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		validator2 = TestingUpdateValidator(app, ctx, validator2, true)
		So(validator2.Status == exported.Bonded, ShouldBeTrue)

		keeper.DeleteValidatorByPowerIndex(ctx, validator)
		validator, _ = validator.AddTokensFromDel(valTokens)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		validator = TestingUpdateValidator(app, ctx, validator, true)

		delegation := types.NewDelegation(addAlice, addAlice, issuedShares)
		keeper.SetDelegation(ctx, delegation)

		_, err := keeper.BeginRedelegation(ctx, addAlice, addAlice, addJack, valTokens.ToDec())
		So(err, ShouldBeNil)
		// end block
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		So(len(updates) == 2, ShouldBeTrue)
		validator, found := keeper.GetValidator(ctx, addAlice)
		So(found, ShouldBeTrue)
		So(valTokens.Equal(validator.Tokens), ShouldBeTrue)
		So(validator.Status == exported.Unbonding, ShouldBeTrue)
	})
	Convey("TestRedelegateFromUnbondingValidator", t, func() {
		_, _, _, addAlice, addJack, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		//	delTokens := exported.TokensFromConsensusPower(5)
		valTokens := exported.TokensFromConsensusPower(10)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		BondedPool := keeper.GetBondedPool(ctx)
		validator := types.NewValidator(addAlice, PKs[0], types.Description{})
		validator, issuedShares := validator.AddTokensFromDel(valTokens)
		So(valTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		validator = TestingUpdateValidator(app, ctx, validator, true)
		selfDelegation := types.NewDelegation(addAlice, addAlice, issuedShares)
		keeper.SetDelegation(ctx, selfDelegation)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))

		keeper.DeleteValidatorByPowerIndex(ctx, validator)
		validator, _ = validator.AddTokensFromDel(valTokens)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		validator = TestingUpdateValidator(app, ctx, validator, true)

		delegation := types.NewDelegation(addJack, addAlice, issuedShares)
		keeper.SetDelegation(ctx, delegation)

		validator2 := types.NewValidator(addJack, PKs[1], types.Description{})
		validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		So(valTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		validator2 = TestingUpdateValidator(app, ctx, validator2, true)
		So(validator2.Status == exported.Bonded, ShouldBeTrue)

		header := ctx.BlockHeader()
		blockHeight := int64(10)
		header.Height = blockHeight
		blockTime := time.Unix(333, 0).UTC()
		header.Time = blockTime
		ctx = ctx.WithBlockHeader(header)

		// unbond the all self-delegation to put validator in unbonding state
		_, err := keeper.Undelegate(ctx, addAlice, addAlice, valTokens.ToDec())
		So(err, ShouldBeNil)

		// end block
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		So(len(updates) == 1, ShouldBeTrue)

		validator, found := keeper.GetValidator(ctx, addAlice)
		So(found, ShouldBeTrue)
		So(blockHeight == validator.UnbondingHeight, ShouldBeTrue)
		params := keeper.GetParams(ctx)
		So(blockTime.Add(params.UnbondingTime).Equal(validator.UnbondingTime), ShouldBeTrue)

		//change the context
		header = ctx.BlockHeader()
		blockHeight2 := int64(20)
		header.Height = blockHeight2
		blockTime2 := time.Unix(444, 0).UTC()
		header.Time = blockTime2
		ctx = ctx.WithBlockHeader(header)

		// unbond some of the other delegation's shares
		redelegateTokens := sdk.TokensFromConsensusPower(6)
		_, err = keeper.BeginRedelegation(ctx, addJack, addAlice, addJack, redelegateTokens.ToDec())
		So(err, ShouldBeNil)

		// retrieve the unbonding delegation
		ubd, found := keeper.GetRedelegation(ctx, addJack, addAlice, addJack)
		So(found, ShouldBeTrue)
		So(len(ubd.Entries) == 1, ShouldBeTrue)
		So(blockHeight == ubd.Entries[0].CreationHeight, ShouldBeTrue)
		So(blockTime.Add(params.UnbondingTime).Equal(ubd.Entries[0].CompletionTime), ShouldBeTrue)
	})
	Convey("TestRedelegateFromUnbondedValidator", t, func() {
		_, _, _, addAlice, addJack, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		//	delTokens := exported.TokensFromConsensusPower(5)
		valTokens := exported.TokensFromConsensusPower(10)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		BondedPool := keeper.GetBondedPool(ctx)
		validator := types.NewValidator(addAlice, PKs[0], types.Description{})
		validator, issuedShares := validator.AddTokensFromDel(valTokens)
		So(valTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		validator = TestingUpdateValidator(app, ctx, validator, true)
		selfDelegation := types.NewDelegation(addAlice, addAlice, issuedShares)
		keeper.SetDelegation(ctx, selfDelegation)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))

		keeper.DeleteValidatorByPowerIndex(ctx, validator)
		validator, _ = validator.AddTokensFromDel(valTokens)
		app.AssetKeeper().IssueCoinPower(ctx, BondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		validator = TestingUpdateValidator(app, ctx, validator, true)

		delegation := types.NewDelegation(addJack, addAlice, issuedShares)
		keeper.SetDelegation(ctx, delegation)

		validator2 := types.NewValidator(addJack, PKs[1], types.Description{})
		validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
		app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), valTokens)))
		So(valTokens.Equal(issuedShares.RoundInt()), ShouldBeTrue)
		validator2 = TestingUpdateValidator(app, ctx, validator2, true)
		So(validator2.Status == exported.Bonded, ShouldBeTrue)

		ctx = ctx.WithBlockHeight(10)
		ctx = ctx.WithBlockTime(time.Unix(333, 0).UTC())

		// unbond the all self-delegation to put validator in unbonding state
		_, err := keeper.Undelegate(ctx, addAlice, addAlice, valTokens.ToDec())
		So(err, ShouldBeNil)

		// end block
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		So(len(updates) == 1, ShouldBeTrue)

		validator, found := keeper.GetValidator(ctx, addAlice)
		So(found, ShouldBeTrue)
		So(ctx.BlockHeight() == validator.UnbondingHeight, ShouldBeTrue)
		params := keeper.GetParams(ctx)
		So(ctx.BlockHeader().Time.Add(params.UnbondingTime).Equal(validator.UnbondingTime), ShouldBeTrue)

		// unbond the validator
		keeper.UnbondingToUnbonded(ctx, validator)

		// redelegate some of the delegation's shares
		redelegationTokens := sdk.TokensFromConsensusPower(6)
		_, err = keeper.BeginRedelegation(ctx, addJack, addAlice, addJack, redelegationTokens.ToDec())
		So(err, ShouldBeNil)

		// no red should have been found
		_, found = keeper.GetRedelegation(ctx, addJack, addAlice, addJack)
		So(found, ShouldBeFalse)
	})
}
