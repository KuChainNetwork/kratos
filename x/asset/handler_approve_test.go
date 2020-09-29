package asset_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

// ApproveForTest approve for test
func ApproveForTest(t *testing.T, w *simapp.Wallet, app *simapp.SimApp, isSuccess bool, id, spender types.AccountID, amt types.Coins) error {
	ctx := app.NewTestContext()
	msgApprove := assetTypes.NewMsgApprove(
		w.GetAuth(id), id, spender, amt)

	tx := simapp.NewTxForTest(
		id,
		[]sdk.Msg{
			&msgApprove,
		}, w.PrivKey(w.GetAuth(id)))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return simapp.CheckTxs(t, app, ctx, tx)
}

func TransferByApproveForTest(t *testing.T, w *simapp.Wallet, app *simapp.SimApp, isSuccess bool, auth types.AccAddress, from, to types.AccountID, amt types.Coins) error {
	ctx := app.NewTestContext()

	// transfer from a to b, with a auth, should success
	transferMsg := assetTypes.NewMsgTransfer(
		auth, from, to, amt)

	tx := simapp.NewTxForTest(
		types.NewAccountIDFromAccAdd(auth),
		[]sdk.Msg{
			&transferMsg,
		}, w.PrivKey(auth))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return simapp.CheckTxs(t, app, ctx, tx)
}

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
		So(app1s.Amount, simapp.ShouldEq, approveCoins1)

		appAll, err := app.AssetKeeper().GetApproveSum(ctx, account1)
		So(err, ShouldBeNil)
		So(appAll, simapp.ShouldEq, approveCoins1)

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
		So(app1s.Amount, simapp.ShouldEq, approveCoins1)

		app2s, err := app.AssetKeeper().GetApproveCoins(ctx, account1, account3)
		So(err, ShouldBeNil)
		So(app2s.Amount, simapp.ShouldEq, approveCoins2)

		appAll, err = app.AssetKeeper().GetApproveSum(ctx, account1)
		So(err, ShouldBeNil)
		So(appAll, simapp.ShouldEq, approveCoins1.Add(approveCoins2...))

		// no exit
		app4s, err := app.AssetKeeper().GetApproveCoins(ctx, account1, account4)
		So(err, ShouldBeNil)
		So(app4s, ShouldBeNil)
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
		So(app1s.Amount, simapp.ShouldEq, approveCoins1)
		So(app1s.IsLock, ShouldEqual, false)

		appAll, err := app.AssetKeeper().GetApproveSum(ctx, account1)
		So(err, ShouldBeNil)
		So(appAll, simapp.ShouldEq, approveCoins1)

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
		So(app2s.Amount, simapp.ShouldEq, approveCoins2)
		So(app2s.IsLock, ShouldEqual, false)

		appAll, err = app.AssetKeeper().GetApproveSum(ctx, account1)
		So(err, ShouldBeNil)
		So(appAll, simapp.ShouldEq, approveCoins2)

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
		So(app3s == nil || app3s.Amount.IsZero(), ShouldBeTrue)

		appAll, err = app.AssetKeeper().GetApproveSum(ctx, account1)
		So(err, ShouldBeNil)
		So(appAll.IsZero(), ShouldBeTrue)
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
		So(apps.Amount, simapp.ShouldEq, apporveCoins.Sub(transferCoins))
	})
}

func TestApproveTransferError(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test approve for test", t, func() {
		So(transfer(t, app, true, account1, addAccount2, NewInt64CoreCoins(2000000000), account1), ShouldBeNil)
		So(transfer(t, app, true, account1, addAccount3, NewInt64CoreCoins(2000000000), account1), ShouldBeNil)
		So(ApproveForTest(t, wallet, app, true, account1, account2, NewInt64CoreCoins(1000)), ShouldBeNil)

		Convey("test transfer error by no approve", func() {
			So(TransferByApproveForTest(t, wallet, app, false, addr3, account1, account3, NewInt64CoreCoins(100)),
				simapp.ShouldErrIs, types.ErrMissingAuth)
		})

		Convey("test transfer error by approve coins no enough", func() {
			So(TransferByApproveForTest(t, wallet, app, false, addr2, account1, account2, NewInt64CoreCoins(1001)),
				simapp.ShouldErrIs, types.ErrMissingAuth)
		})

		Convey("test muit-transfer error by approve coins no enough", func() {
			ctx := app.NewTestContext()

			transferMsg := assetTypes.NewMsgTransfers(
				wallet.GetAuth(account2),
				[]types.KuMsgTransfer{
					types.NewKuMsgTransfer(account1, account2, NewInt64CoreCoins(100)),
					types.NewKuMsgTransfer(account1, account2, NewInt64CoreCoins(900)),
					types.NewKuMsgTransfer(account1, account2, NewInt64CoreCoins(100)),
				})

			tx := simapp.NewTxForTest(
				account2,
				[]sdk.Msg{
					&transferMsg,
				}, wallet.PrivKey(addr2)).WithCannotPass()
			So(simapp.CheckTxs(t, app, ctx, tx), simapp.ShouldErrIs, types.ErrMissingAuth)
		})

		Convey("test transfer error by approve auth error", func() {
			So(TransferByApproveForTest(t, wallet, app, false, addr3, account1, account2, NewInt64CoreCoins(100)),
				simapp.ShouldErrIs, types.ErrMissingAuth)
		})

	})
}

func TestApproveLockMode(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test approve lock mode for test", t, func() {
		So(transfer(t, app, true, account1, addAccount2, NewInt64CoreCoins(2000000000), account1), ShouldBeNil)
		So(transfer(t, app, true, account1, addAccount3, NewInt64CoreCoins(2000000000), account1), ShouldBeNil)

		// use keeper to set
		ctx := app.NewTestContext()
		So(app.AssetKeeper().Approve(ctx, account1, account2, NewInt64CoreCoins(10000), true), ShouldBeNil)

		// add ok
		So(app.AssetKeeper().Approve(ctx, account1, account2, NewInt64CoreCoins(100000), true), ShouldBeNil)

		appData, err := app.AssetKeeper().GetApproveCoins(ctx, account1, account2)
		So(err, ShouldBeNil)
		So(appData.IsLock, ShouldBeTrue)
		So(appData.Amount, simapp.ShouldEq, NewInt64CoreCoins(100000))

		// sub err TODO: this will check in handler
		//So(app.AssetKeeper().Approve(ctx, account1, account2, NewInt64CoreCoins(10000), true),
		//	simapp.ShouldErrIs, assetTypes.ErrAssetApporveCannotChangeLock)

		So(app.AssetKeeper().Approve(ctx, account1, account2, NewInt64CoreCoins(100000), false),
			simapp.ShouldErrIs, assetTypes.ErrAssetApporveCannotChangeLock)
	})
}
