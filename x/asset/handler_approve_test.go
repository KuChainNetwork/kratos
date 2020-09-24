package asset_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApproveCoins(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test approve coins", t, func() {
		ctx := app.NewTestContext()

		var (
			approveCoins1 = NewInt64CoreCoins(1000000)
			approveCoins2 = NewInt64CoreCoins(200000)
			//approveCoins3 = NewInt64CoreCoins(30000)
		)

		// first coins 1
		msgApprove := assetTypes.NewMsgApprove(
			wallet.GetAuth(account1), account1, account2, approveCoins1)

		tx := simapp.NewTxForTest(
			account1,
			[]sdk.Msg{
				&msgApprove,
			}, wallet.PrivKey(addr1))

		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		ctx = app.NewTestContext()

		app1s, err := app.AssetKeeper().GetApproveCoins(ctx, account1, account2)
		So(err, ShouldBeNil)
		So(app1s, simapp.ShouldEq, approveCoins1)

		// second coins 2
		msgApprove2 := assetTypes.NewMsgApprove(
			wallet.GetAuth(account1), account1, account3, approveCoins2)

		tx = simapp.NewTxForTest(
			account1,
			[]sdk.Msg{
				&msgApprove2,
			}, wallet.PrivKey(addr1))
		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		ctx = app.NewTestContext()

		app1s, err = app.AssetKeeper().GetApproveCoins(ctx, account1, account2)
		So(err, ShouldBeNil)
		So(app1s, simapp.ShouldEq, approveCoins1)

		app2s, err := app.AssetKeeper().GetApproveCoins(ctx, account1, account3)
		So(err, ShouldBeNil)
		So(app2s, simapp.ShouldEq, approveCoins2)

		// no exit
		app4s, err := app.AssetKeeper().GetApproveCoins(ctx, account1, account4)
		So(err, ShouldBeNil)
		So(app4s.IsZero(), ShouldBeTrue)
	})
}

func TestApproveResetCoins(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test approve coins", t, func() {
		ctx := app.NewTestContext()

		var (
			approveCoins1 = NewInt64CoreCoins(1000000)
			approveCoins2 = NewInt64CoreCoins(200000)
		)

		// first coins 1
		msgApprove := assetTypes.NewMsgApprove(
			wallet.GetAuth(account1), account1, account2, approveCoins1)

		tx := simapp.NewTxForTest(
			account1,
			[]sdk.Msg{
				&msgApprove,
			}, wallet.PrivKey(addr1))

		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		ctx = app.NewTestContext()

		app1s, err := app.AssetKeeper().GetApproveCoins(ctx, account1, account2)
		So(err, ShouldBeNil)
		So(app1s, simapp.ShouldEq, approveCoins1)

		// second coins 2
		msgApprove2 := assetTypes.NewMsgApprove(
			wallet.GetAuth(account1), account1, account2, approveCoins2)

		tx = simapp.NewTxForTest(
			account1,
			[]sdk.Msg{
				&msgApprove2,
			}, wallet.PrivKey(addr1))
		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		ctx = app.NewTestContext()

		app2s, err := app.AssetKeeper().GetApproveCoins(ctx, account1, account2)
		So(err, ShouldBeNil)
		So(app2s, simapp.ShouldEq, approveCoins2)

		// second coins zero
		msgApprove3 := assetTypes.NewMsgApprove(
			wallet.GetAuth(account1), account1, account2, NewInt64CoreCoins(0))

		tx = simapp.NewTxForTest(
			account1,
			[]sdk.Msg{
				&msgApprove3,
			}, wallet.PrivKey(addr1))
		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		ctx = app.NewTestContext()

		app3s, err := app.AssetKeeper().GetApproveCoins(ctx, account1, account2)
		So(err, ShouldBeNil)
		So(app3s.IsZero(), ShouldBeTrue)

	})
}
