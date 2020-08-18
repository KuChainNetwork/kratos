package asset_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	. "github.com/smartystreets/goconvey/convey"
)

func checkBurnState(t *testing.T, app *simapp.SimApp, isSuccess bool, from types.AccountID, amt types.Coin) error {
	ctx := app.NewTestContext()

	creator, symbol, err := types.CoinAccountsFromDenom(amt.Denom)
	So(err, ShouldBeNil)

	amtOld := app.AssetKeeper().GetAllBalances(ctx, from)
	statOld, err := app.AssetKeeper().GetCoinStat(ctx, creator, symbol)
	So(err, ShouldBeNil)

	burnErr := BurnCoinTest(t, app, isSuccess, from, amt)

	currAmt := amtOld.Sub(simapp.DefaultTestFee)
	if isSuccess {
		currAmt = currAmt.Sub(types.NewCoins(amt))
	}

	ctx = app.NewTestContext()
	So(app.AssetKeeper().GetAllBalances(ctx, from), simapp.ShouldEq, currAmt)

	statNew, err := app.AssetKeeper().GetCoinStat(ctx, creator, symbol)
	So(err, ShouldBeNil)

	// should no core coin as it can issue per block
	if isSuccess {
		So(statOld.Supply.Sub(amt), simapp.ShouldEq, statNew.Supply)
	} else {
		So(statOld.Supply, simapp.ShouldEq, statNew.Supply)
	}

	return burnErr
}

func TestBurnCoreCoins(t *testing.T) {
	app, _ := createAppForTest()
	Convey("test core coins", t, func() {
		ctx := app.NewTestContext()
		acc1Coins := app.AssetKeeper().GetAllBalances(ctx, account1)
		burnCoin := types.NewInt64Coin(constants.DefaultBondDenom, 10000)

		So(BurnCoinTest(t, app, true, account1, burnCoin), ShouldBeNil)
		amt := acc1Coins.Sub(simapp.DefaultTestFee)
		ctx = app.NewTestContext()
		So(amt.Sub(types.NewCoins(burnCoin)), simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account1))
	})

	Convey("test burn all core coins error by no fee", t, func() {
		ctx := app.NewTestContext()
		acc2Coins := app.AssetKeeper().GetAllBalances(ctx, account2)
		burnCoin := types.NewInt64CoreCoin(acc2Coins.AmountOf(constants.DefaultBondDenom).Int64())

		So(BurnCoinTest(t, app, false, account2, burnCoin), simapp.ShouldErrIs, assetTypes.ErrAssetCoinNoEnough)
		amt := acc2Coins.Sub(simapp.DefaultTestFee)
		ctx = app.NewTestContext()
		So(amt, simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account2))
	})

	Convey("test burn all core coins", t, func() {
		ctx := app.NewTestContext()
		acc2Coins := app.AssetKeeper().GetAllBalances(ctx, account2)
		burnCoin := types.NewInt64CoreCoin(acc2Coins.AmountOf(constants.DefaultBondDenom).Int64())
		burnCoin = burnCoin.Sub(simapp.DefaultTestFee[0])

		So(BurnCoinTest(t, app, true, account2, burnCoin), ShouldBeNil)
		ctx = app.NewTestContext()
		So(types.NewCoins(), simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account2))
	})
}

func TestBurnOtherCoins(t *testing.T) {
	app, _ := createAppForTest()

	var (
		symbol             = types.MustName("abc")
		denom              = types.CoinDenom(name4, symbol)
		maxSupplyAmt int64 = 10000000000000
		issueAmts          = types.NewInt64Coin(denom, maxSupplyAmt)
	)

	Convey("create coins to test at first", t, func() {
		So(createCoin(t, app, true, account4, symbol, maxSupplyAmt), ShouldBeNil)
		So(issueCoin(t, app, true, account4, symbol, issueAmts), ShouldBeNil)
		So(transfer(t, app, true, account4, account1, types.NewInt64Coins(denom, 1000000), account4), ShouldBeNil)
	})

	Convey("test other coins burn", t, func() {
		burnCoin := types.NewInt64Coin(denom, 10000)

		So(checkBurnState(t, app, true, account1, burnCoin), ShouldBeNil)
	})

	Convey("test burn no enough coins", t, func() {
		burnCoin := types.NewInt64Coin(denom, 1000000)

		So(checkBurnState(t, app, false, account1, burnCoin), simapp.ShouldErrIs, assetTypes.ErrAssetCoinNoEnough)
	})

	Convey("test burn all core coins", t, func() {
		ctx := app.NewTestContext()
		acc1Coins := app.AssetKeeper().GetAllBalances(ctx, account1)
		burnCoin := types.NewInt64Coin(denom, acc1Coins.AmountOf(denom).Int64())

		So(checkBurnState(t, app, true, account1, burnCoin), ShouldBeNil)
	})
}

func TestBurnNoCanBurnCoins(t *testing.T) {
	app, _ := createAppForTest()

	var (
		symbol             = types.MustName("abcd1")
		denom              = types.CoinDenom(name4, symbol) // create in last
		maxSupplyAmt int64 = 10000000000000
		issueAmts          = types.NewInt64Coin(denom, 10000000000000)
	)

	Convey("create coins to test at first", t, func() {
		So(createCoinExt(t, app, true,
			account4, symbol,
			types.NewInt64Coin(denom, maxSupplyAmt),
			true, true,
			false, // cannot burn
			0,
			types.NewInt64Coin(denom, 0), []byte("cannot issue")), ShouldBeNil)
		So(issueCoin(t, app, true, account4, symbol, issueAmts), ShouldBeNil)
		So(transfer(t, app, true, account4, account1, types.NewInt64Coins(denom, 1000000), account4), ShouldBeNil)
	})

	Convey("test other coins burn", t, func() {
		burnCoin := types.NewInt64Coin(denom, 10000)

		So(checkBurnState(t, app, false, account1, burnCoin), simapp.ShouldErrIs, assetTypes.ErrAssetCoinCannotBeBurn)
	})
}
