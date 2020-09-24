package asset_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/types"
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

func TestApproveTransfer(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test simple transfer, a->b with auth(a)", t, func() {
		ctx := app.NewTestContext()

		// transfer from a to b, with a auth, should success
		transferMsg := assetTypes.NewMsgTransfer(
			wallet.GetAuth(account1), account1, account2, NewInt64CoreCoins(100))

		tx := simapp.NewTxForTest(
			account1,
			[]sdk.Msg{
				&transferMsg,
			}, wallet.PrivKey(addr1))
		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)
	})

	Convey("test simple transfer, a->b with auth(b), should failed", t, func() {
		ctx := app.NewTestContext()

		transferMsg := assetTypes.NewMsgTransfer(
			wallet.GetAuth(account2), account1, account2, NewInt64CoreCoins(100))

		tx := simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&transferMsg,
			}, wallet.PrivKey(addr2)).WithCannotPass()
		So(simapp.CheckTxs(t, app, ctx, tx), simapp.ShouldErrIs, types.ErrMissingAuth)
	})

	Convey("test simple transfer, a->b with auth(c), should failed", t, func() {
		ctx := app.NewTestContext()

		transferMsg := assetTypes.NewMsgTransfer(
			wallet.GetAuth(account3), account1, account2, NewInt64CoreCoins(100))

		tx := simapp.NewTxForTest(
			account3,
			[]sdk.Msg{
				&transferMsg,
			}, wallet.PrivKey(addr3)).WithCannotPass()
		So(simapp.CheckTxs(t, app, ctx, tx), simapp.ShouldErrIs, types.ErrMissingAuth)
	})

	Convey("test transfer with appover, a->b with auth(b)", t, func() {
		ctx := app.NewTestContext()

		var (
			apporveCoins  = NewInt64CoreCoins(1000)
			transferCoins = NewInt64CoreCoins(100)
		)

		msgApprove := assetTypes.NewMsgApprove(
			wallet.GetAuth(account1), account1, account2, apporveCoins)

		tx := simapp.NewTxForTest(
			account1,
			[]sdk.Msg{
				&msgApprove,
			}, wallet.PrivKey(addr1))

		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		ctx = app.NewTestContext()

		transferMsg := assetTypes.NewMsgTransfer(
			wallet.GetAuth(account2), account1, account2, transferCoins)

		tx = simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&transferMsg,
			}, wallet.PrivKey(addr2))
		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		ctx = app.NewTestContext()

		apps, err := app.AssetKeeper().GetApproveCoins(ctx, account1, account2)
		So(err, ShouldBeNil)
		So(apps, simapp.ShouldEq, apporveCoins.Sub(transferCoins))
	})
}
