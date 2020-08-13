package keeper_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoinPower(t *testing.T) {
	app, ctx := createTestApp()

	amt := types.NewCoins(types.NewInt64Coin(constants.DefaultBondDenom, 100000))

	Convey("test coins to power by self", t, func() {
		err := app.AssetKeeper().CoinsToPower(ctx, account1, account1, amt)
		So(err, ShouldBeNil)

		coins := app.AssetKeeper().GetCoinPowers(ctx, account1)
		So(coins, simapp.ShouldEq, amt)
	})

}

func TestCoinPowerExercise(t *testing.T) {
	app, ctx := createTestApp()

	amt := types.NewCoins(types.NewInt64Coin(constants.DefaultBondDenom, 100000))

	Convey("test coin power exercise", t, func() {
		Convey("step1 : test coin to coin power", func() {
			err := app.AssetKeeper().CoinsToPower(ctx, account1, account1, amt)
			So(err, ShouldBeNil)

			coins := app.AssetKeeper().GetCoinPowers(ctx, account1)
			So(coins, simapp.ShouldEq, amt)
		})

		Convey("step2 : test coin power exercise", func() {
			ctx := app.NewTestContext()
			amtAll := app.AssetKeeper().GetAllBalances(ctx, account1)
			coinPowerAll := app.AssetKeeper().GetCoinPowers(ctx, account1)

			amte := types.NewInt64Coin(constants.DefaultBondDenom, 1000)

			err := app.AssetKeeper().ExerciseCoinPower(ctx, account1, amte)
			So(err, ShouldBeNil)
			So(amtAll.Add(amte), simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account1))
			So(coinPowerAll.Sub(types.NewCoins(amte)), simapp.ShouldEq, app.AssetKeeper().GetCoinPowers(ctx, account1))
		})

		Convey("step3 : test coin power exercise zero", func() {
			ctx := app.NewTestContext()
			amtAll := app.AssetKeeper().GetAllBalances(ctx, account1)
			coinPowerAll := app.AssetKeeper().GetCoinPowers(ctx, account1)

			err := app.AssetKeeper().ExerciseCoinPower(ctx, account1, types.NewInt64Coin(constants.DefaultBondDenom, 0))
			So(err, ShouldBeNil)
			So(amtAll, simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account1))
			So(coinPowerAll, simapp.ShouldEq, app.AssetKeeper().GetCoinPowers(ctx, account1))
		})

		Convey("step4 : test coin power exercise all", func() {
			ctx := app.NewTestContext()
			amtAll := app.AssetKeeper().GetAllBalances(ctx, account1)
			coinPowerAll := app.AssetKeeper().GetCoinPowers(ctx, account1)

			err := app.AssetKeeper().ExerciseCoinPower(ctx, account1, coinPowerAll[0])
			So(err, ShouldBeNil)
			So(amtAll.Add(coinPowerAll...), simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account1))
			So(types.NewCoins(), simapp.ShouldEq, app.AssetKeeper().GetCoinPowers(ctx, account1))
		})
	})

}

func TestCoinPowerExerciseErrors(t *testing.T) {
	app, ctx := createTestApp()

	amt := types.NewCoins(types.NewInt64Coin(constants.DefaultBondDenom, 100000))

	Convey("test coin power exercise errors", t, func() {
		Convey("step1 : test coin to coin power", func() {
			err := app.AssetKeeper().CoinsToPower(ctx, account1, account1, amt)
			So(err, ShouldBeNil)

			coins := app.AssetKeeper().GetCoinPowers(ctx, account1)
			So(coins, simapp.ShouldEq, amt)
		})

		Convey("step2 : test coin power exercise error by no enough", func() {
			ctx := app.NewTestContext()
			amtAll := app.AssetKeeper().GetAllBalances(ctx, account1)
			coinPowerAll := app.AssetKeeper().GetCoinPowers(ctx, account1)

			amte := types.NewInt64Coin(constants.DefaultBondDenom, 10000000000)

			err := app.AssetKeeper().ExerciseCoinPower(ctx, account1, amte)
			So(err, simapp.ShouldErrIs, assetTypes.ErrAssetCoinNoEnough)
			So(amtAll, simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account1))
			So(coinPowerAll, simapp.ShouldEq, app.AssetKeeper().GetCoinPowers(ctx, account1))
		})
	})

}
