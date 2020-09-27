package asset_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTransferCoins(t *testing.T) {
	app, _ := createAppForTest()
	Convey("test normal transfer", t, func() {
		ctx := app.NewTestContext()
		acc1Coins := app.AssetKeeper().GetAllBalances(ctx, account1)
		acc2Coins := app.AssetKeeper().GetAllBalances(ctx, account2)
		coins2Transfer := types.NewInt64Coins(constants.DefaultBondDenom, 1000)

		So(transfer(t, app, true, account1, account2, coins2Transfer, account1), ShouldBeNil)
		amt := acc1Coins.Sub(simapp.DefaultTestFee)
		ctx = app.NewTestContext()
		So(amt.Sub(coins2Transfer), simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account1))
		So(acc2Coins.Add(coins2Transfer...), simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account2))
	})
}

func TestTransferCoinsErr(t *testing.T) {
	app, _ := createAppForTest()
	Convey("test transfer negative coins", t, func() {
		ctx := app.NewTestContext()
		acc1Coins := app.AssetKeeper().GetAllBalances(ctx, account1)
		acc2Coins := app.AssetKeeper().GetAllBalances(ctx, account2)
		coins2Transfer := types.Coins{types.Coin{constants.DefaultBondDenom, types.NewInt(-1000)}}

		So(transfer(t, app, false, account1, account2, coins2Transfer, account1),
			simapp.ShouldErrIs, types.ErrTransfNoEnough)

		amt := acc1Coins
		//amt := acc1Coins.Sub(simapp.DefaultTestFee) will err by basic check

		ctx = app.NewTestContext()
		So(amt, simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account1))
		So(acc2Coins, simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account2))
	})

	Convey("test transfer no enough coins", t, func() {
		ctx := app.NewTestContext()
		acc1Coins := app.AssetKeeper().GetAllBalances(ctx, account1)
		acc2Coins := app.AssetKeeper().GetAllBalances(ctx, account2)
		coins2Transfer := NewInt64CoreCoins(100000000000001)

		So(transfer(t, app, false, account1, account2, coins2Transfer, account1),
			simapp.ShouldErrIs, assetTypes.ErrAssetCoinNoEnough)

		amt := acc1Coins.Sub(simapp.DefaultTestFee)

		ctx = app.NewTestContext()
		So(amt, simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account1))
		So(acc2Coins, simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account2))
	})

	Convey("test transfer empty account", t, func() {
		ctx := app.NewTestContext()
		acc1Coins := app.AssetKeeper().GetAllBalances(ctx, account1)
		acc2Coins := app.AssetKeeper().GetAllBalances(ctx, account2)
		coins2Transfer := NewInt64CoreCoins(10000000001)

		So(transfer(t, app, true, account1, types.AccountID{}, coins2Transfer, account1), ShouldBeNil)

		amt := acc1Coins.Sub(simapp.DefaultTestFee)

		ctx = app.NewTestContext()
		So(amt, simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account1))
		So(acc2Coins, simapp.ShouldEq, app.AssetKeeper().GetAllBalances(ctx, account2))
	})
}
