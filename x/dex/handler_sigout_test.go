package dex_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/dex/keeper"
	dexTypes "github.com/KuChainNetwork/kuchain/x/dex/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSigoutByUser(t *testing.T) {
	Convey("test signOut msg by user", t, func() {
		// change var
		keeper.SigOutByUserUnlockHeight = 10

		app, _ := createAppForTest()

		So(CreateDexForTest(t, app, true, dexAccount1, types.NewInt64CoreCoins(1000000000), []byte("dex for test")), ShouldBeNil)

		amt := types.NewInt64CoreCoins(1000000)
		So(SignInMsgForTest(t, app, true, account1, dexAccount1, amt), ShouldBeNil)

		ctx := app.NewTestContext()
		assetKeeper := app.AssetKeeper()

		data, err := assetKeeper.GetApproveCoins(ctx, account1, dexAccount1)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		So(data.IsLock, ShouldBeTrue)
		So(data.Amount, simapp.ShouldEq, amt)

		out := types.NewInt64CoreCoins(77777)
		left := amt.Sub(out)
		// sigout req
		So(SignOutMsgByDexExForTest(t, app, ctx, true, account1, account1, dexAccount1, out), ShouldBeNil)

		ctx = app.NewTestContext()
		height, ok := app.DexKeeper().GetSigOutReqHeight(ctx, account1)
		So(ok, ShouldBeTrue)
		So(height, ShouldEqual, app.LastBlockHeight())

		// req two times, should error
		So(SignOutMsgByDexExForTest(t, app, ctx, false, account1, account1, dexAccount1, out),
			simapp.ShouldErrIs, dexTypes.ErrDexSigOutByUserNoUnlock)

		// wait enough blocks
		simapp.AfterBlockCommitted(app, int(keeper.SigOutByUserUnlockHeight))

		ctx = app.NewTestContext()
		app.Logger().Info("current block", "height", ctx.BlockHeight())

		// should ok
		So(SignOutMsgByDexExForTest(t, app, ctx, true, account1, account1, dexAccount1, out), ShouldBeNil)

		ctx = app.NewTestContext()

		data, err = assetKeeper.GetApproveCoins(ctx, account1, dexAccount1)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		So(data.IsLock, ShouldBeTrue)
		So(data.Amount, simapp.ShouldEq, left)

		// should clean sigout
		height, ok = app.DexKeeper().GetSigOutReqHeight(ctx, account1)
		So(ok, ShouldBeFalse)
		So(height, ShouldEqual, 0)

		So(SignOutMsgByDexForTest(t, app, true, account1, dexAccount1, left), ShouldBeNil)

		ctx = app.NewTestContext()

		data, err = assetKeeper.GetApproveCoins(ctx, account1, dexAccount1)
		So(err, ShouldBeNil)
		So(data, ShouldBeNil)
	})
}

func TestSigoutByUserThanDexSigout(t *testing.T) {
	Convey("test signOut msg by user", t, func() {
		app, _ := createAppForTest()

		So(CreateDexForTest(t, app, true, dexAccount1, types.NewInt64CoreCoins(1000000000), []byte("dex for test")), ShouldBeNil)

		amt := types.NewInt64CoreCoins(1000000)
		So(SignInMsgForTest(t, app, true, account1, dexAccount1, amt), ShouldBeNil)

		ctx := app.NewTestContext()
		assetKeeper := app.AssetKeeper()

		data, err := assetKeeper.GetApproveCoins(ctx, account1, dexAccount1)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		So(data.IsLock, ShouldBeTrue)
		So(data.Amount, simapp.ShouldEq, amt)

		out := types.NewInt64CoreCoins(77777)
		left := amt.Sub(out)
		// sigout req
		So(SignOutMsgByDexExForTest(t, app, ctx, true, account1, account1, dexAccount1, out), ShouldBeNil)

		ctx = app.NewTestContext()
		height, ok := app.DexKeeper().GetSigOutReqHeight(ctx, account1)
		So(ok, ShouldBeTrue)
		So(height, ShouldEqual, app.LastBlockHeight())

		// req two times, should error
		So(SignOutMsgByDexExForTest(t, app, ctx, false, account1, account1, dexAccount1, out),
			simapp.ShouldErrIs, dexTypes.ErrDexSigOutByUserNoUnlock)

		// wait no enough blocks
		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		app.Logger().Info("current block", "height", ctx.BlockHeight())

		// should ok
		So(SignOutMsgByDexForTest(t, app, true, account1, dexAccount1, out), ShouldBeNil)

		ctx = app.NewTestContext()

		data, err = assetKeeper.GetApproveCoins(ctx, account1, dexAccount1)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		So(data.IsLock, ShouldBeTrue)
		So(data.Amount, simapp.ShouldEq, left)

		// should clean sigout
		height, ok = app.DexKeeper().GetSigOutReqHeight(ctx, account1)
		So(ok, ShouldBeFalse)
		So(height, ShouldEqual, 0)

		So(SignOutMsgByDexForTest(t, app, true, account1, dexAccount1, left), ShouldBeNil)

		ctx = app.NewTestContext()

		data, err = assetKeeper.GetApproveCoins(ctx, account1, dexAccount1)
		So(err, ShouldBeNil)
		So(data, ShouldBeNil)
	})
}
