package dex_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	dexTypes "github.com/KuChainNetwork/kuchain/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

func SignInMsgForTest(t *testing.T, app *simapp.SimApp, isSuccess bool, id, dex types.AccountID, amount types.Coins) error {
	wallet := app.GetWallet()

	ctx := app.NewTestContext()

	msg := dexTypes.NewMsgDexSigIn(
		wallet.GetAuth(id),
		id,
		dex,
		amount)

	tx := simapp.NewTxForTest(
		id,
		[]sdk.Msg{
			&msg,
		}, wallet.PrivKey(wallet.GetAuth(id)))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return simapp.CheckTxs(t, app, ctx, tx)
}

func SignOutMsgByDexForTest(t *testing.T, app *simapp.SimApp, isSuccess bool, id, dex types.AccountID, amount types.Coins) error {
	wallet := app.GetWallet()

	ctx := app.NewTestContext()

	msg := dexTypes.NewMsgDexSigOut(
		wallet.GetAuth(dex),
		false,
		id,
		dex,
		amount)

	tx := simapp.NewTxForTest(
		dex,
		[]sdk.Msg{
			&msg,
		}, wallet.PrivKey(wallet.GetAuth(dex)))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return simapp.CheckTxs(t, app, ctx, tx)
}

func TestSignInMsg(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test signIn msg", t, func() {
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

		all, locks, err := assetKeeper.GetLockCoins(ctx, account1)
		So(err, ShouldBeNil)
		So(all, simapp.ShouldEq, amt)
		So(len(locks), ShouldEqual, 1)
		So(locks[0].Coins, simapp.ShouldEq, amt)
		So(locks[0].UnlockBlockHeight, ShouldBeLessThan, 0)
	})
}

func TestSignOutMsg(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test signOut msg", t, func() {
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

		all, locks, err := assetKeeper.GetLockCoins(ctx, account1)
		So(err, ShouldBeNil)
		So(all, simapp.ShouldEq, amt)
		So(len(locks), ShouldEqual, 1)
		So(locks[0].Coins, simapp.ShouldEq, amt)
		So(locks[0].UnlockBlockHeight, ShouldBeLessThan, 0)

		out := types.NewInt64CoreCoins(77777)
		left := amt.Sub(out)
		So(SignOutMsgByDexForTest(t, app, true, account1, dexAccount1, out), ShouldBeNil)

		ctx = app.NewTestContext()

		data, err = assetKeeper.GetApproveCoins(ctx, account1, dexAccount1)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		So(data.IsLock, ShouldBeTrue)
		So(data.Amount, simapp.ShouldEq, left)

		all, locks, err = assetKeeper.GetLockCoins(ctx, account1)
		So(err, ShouldBeNil)
		So(all, simapp.ShouldEq, left)
		So(len(locks), ShouldEqual, 1)
		So(locks[0].Coins, simapp.ShouldEq, left)
		So(locks[0].UnlockBlockHeight, ShouldBeLessThan, 0)

		So(SignOutMsgByDexForTest(t, app, true, account1, dexAccount1, left), ShouldBeNil)

		ctx = app.NewTestContext()

		data, err = assetKeeper.GetApproveCoins(ctx, account1, dexAccount1)
		So(err, ShouldBeNil)
		So(data, ShouldBeNil)

		all, locks, err = assetKeeper.GetLockCoins(ctx, account1)
		So(err, ShouldBeNil)
		So(all.IsZero(), ShouldBeTrue)
		So(len(locks), ShouldEqual, 0)
	})
}
