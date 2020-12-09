package dex_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	dexTypes "github.com/KuChainNetwork/kuchain/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	feeRateDiv int64 = 10
)

func DexDealForTest(t *testing.T, app *simapp.SimApp, isSuccess bool, dex, a1 types.AccountID, c1 types.Coin, a2 types.AccountID, c2 types.Coin) error {
	wallet := app.GetWallet()
	ctx := app.NewTestContext()

	msg := dexTypes.NewMsgDexDeal(
		wallet.GetAuth(dex),
		dex,
		a1, a2,
		c1, c2,
		types.NewInt64Coin(c1.Denom, c1.Amount.Int64()/feeRateDiv),
		types.NewInt64Coin(c1.Denom, c2.Amount.Int64()/feeRateDiv),
		[]byte("deal dex test"))

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

func TestDexDealMsg(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test dex deal", t, func() {
		So(CreateDexForTest(t, app, true, dexAccount1, types.NewInt64CoreCoins(1000000000), []byte("dex for test")), ShouldBeNil)

		amt := types.NewCoins(types.NewInt64CoreCoin(1000000), types.NewInt64Coin(gDenom1, 1000000000000000))

		So(SignInMsgForTest(t, app, true, account1, dexAccount1, amt), ShouldBeNil)
		So(SignInMsgForTest(t, app, true, account2, dexAccount1, amt), ShouldBeNil)

		ctx := app.NewTestContext()

		account1AssetOld := app.AssetKeeper().GetAllBalances(ctx, account1)
		account2AssetOld := app.AssetKeeper().GetAllBalances(ctx, account2)
		dexAssetOld := app.AssetKeeper().GetAllBalances(ctx, dexAccount1)

		a11, a21 := types.NewInt64Coin(gDenom1, 600), types.NewInt64CoreCoin(1000)
		f1 := types.NewInt64Coin(a11.Denom, a11.Amount.Int64()/feeRateDiv)
		f2 := types.NewInt64Coin(a11.Denom, a21.Amount.Int64()/feeRateDiv)

		So(DexDealForTest(t, app, true, dexAccount1, account1, a11, account2, a21), ShouldBeNil)

		ctx = app.NewTestContext()

		account1AssetNew := app.AssetKeeper().GetAllBalances(ctx, account1)
		account2AssetNew := app.AssetKeeper().GetAllBalances(ctx, account2)
		dexAssetNew := app.AssetKeeper().GetAllBalances(ctx, dexAccount1)

		ctx.Logger().Info("amt", "amt", amt)
		ctx.Logger().Info("deal a1", "from", a11, "fee", f1)
		ctx.Logger().Info("deal a2", "from", a21, "fee", f2)

		// coins
		So(account1AssetOld.Add(a21), simapp.ShouldEq, account1AssetNew.Add(a11.Add(f1)))
		So(account2AssetOld.Add(a11.Sub(f2)), simapp.ShouldEq, account2AssetNew.Add(a21))
		So(dexAssetNew, simapp.ShouldEq, (dexAssetOld.Add(f1).Add(f2)).Sub(simapp.DefaultTestFee))

		// approve
		// approve should add to each deal account
		account1Approve, err := app.AssetKeeper().GetApproveCoins(ctx, account1, dexAccount1)
		So(err, ShouldBeNil)
		So(account1Approve, ShouldNotBeNil)
		// now approve should for account1 == (amt+a21-a11-f1)
		So(account1Approve.Amount.Add(a11).Add(f1), simapp.ShouldEq, amt.Add(a21))

		account2Approve, err := app.AssetKeeper().GetApproveCoins(ctx, account2, dexAccount1)
		So(err, ShouldBeNil)
		So(account2Approve, ShouldNotBeNil)
		// now approve should for account2 == (amt+a11-f2-a21)
		So(account2Approve.Amount.Add(a21).Add(f2), simapp.ShouldEq, amt.Add(a11))

		// sigIn
		sigInState1 := app.DexKeeper().GetSigInForDex(ctx, account1, dexAccount1)
		sigInState2 := app.DexKeeper().GetSigInForDex(ctx, account2, dexAccount1)

		// same as approve
		So(sigInState1.Add(a11).Add(f1), simapp.ShouldEq, amt.Add(a21))
		So(sigInState2.Add(a21).Add(f2), simapp.ShouldEq, amt.Add(a11))
	})
}

func TestDexDealErr(t *testing.T) {
	Convey("test dex deal err by no dex", t, func() {
		app, _ := createAppForTest()
		a11, a21 := types.NewInt64Coin(gDenom1, 600), types.NewInt64CoreCoin(1000)
		So(DexDealForTest(t, app, false, dexAccount1, account1, a11, account2, a21), ShouldNotBeNil)
	})

	Convey("test dex deal err by no from sigIn", t, func() {
		app, _ := createAppForTest()
		So(CreateDexForTest(t, app, true, dexAccount1, types.NewInt64CoreCoins(1000000000), []byte("dex for test")), ShouldBeNil)
		amt := types.NewCoins(types.NewInt64CoreCoin(1000000), types.NewInt64Coin(gDenom1, 1000000000000000))

		// So(SignInMsgForTest(t, app, true, account1, dexAccount1, amt), ShouldBeNil)
		So(SignInMsgForTest(t, app, true, account2, dexAccount1, amt), ShouldBeNil)

		a11, a21 := types.NewInt64Coin(gDenom1, 600), types.NewInt64CoreCoin(1000)
		So(DexDealForTest(t, app, false, dexAccount1, account1, a11, account2, a21),
			simapp.ShouldErrIs, types.ErrMissingAuth)
	})

	Convey("test dex deal err by no to sigIn", t, func() {
		app, _ := createAppForTest()
		So(CreateDexForTest(t, app, true, dexAccount1, types.NewInt64CoreCoins(1000000000), []byte("dex for test")), ShouldBeNil)
		amt := types.NewCoins(types.NewInt64CoreCoin(1000000), types.NewInt64Coin(gDenom1, 1000000000000000))

		So(SignInMsgForTest(t, app, true, account1, dexAccount1, amt), ShouldBeNil)
		// So(SignInMsgForTest(t, app, true, account2, dexAccount1, amt), ShouldBeNil)

		a11, a21 := types.NewInt64Coin(gDenom1, 600), types.NewInt64CoreCoin(1000)
		So(DexDealForTest(t, app, false, dexAccount1, account1, a11, account2, a21),
			simapp.ShouldErrIs, types.ErrMissingAuth)
	})

	Convey("test dex deal err by approve", t, func() {
		app, _ := createAppForTest1()

		So(CreateDexForTest(t, app, true, dexAccount1, types.NewInt64CoreCoins(1000000000), []byte("dex for test")), ShouldBeNil)

		sigInAmt1 := types.NewCoins(types.NewInt64Coin(gDenom1, 100))
		sigInAmt2 := types.NewCoins(types.NewInt64Coin(gDenom2, 100))

		So(SignInMsgForTest(t, app, true, account1, dexAccount1, sigInAmt1), ShouldBeNil)
		So(SignInMsgForTest(t, app, true, account1, dexAccount1, sigInAmt2), ShouldBeNil)

		a11, a21 := types.NewInt64Coin(gDenom1, 10), types.NewInt64Coin(gDenom2, 10)

		So(DexDealForTest(t, app, true, dexAccount1, account1, a11, account1, a21), ShouldBeNil)

		sigAmt := types.NewCoins(types.NewInt64Coin(gDenom1, 1900))
		So(SignInMsgForTest(t, app, true, account2, dexAccount1, sigAmt), ShouldBeNil)
	})
}

func TestDexDealErr1(t *testing.T) {
	Convey("test dex deal err by approve", t, func() {
		app, _ := createAppForTest1()

		So(CreateDexForTest(t, app, true, dexAccount1, types.NewInt64CoreCoins(1000000000), []byte("dex for test")), ShouldBeNil)

		sigInAmt1 := types.NewCoins(types.NewInt64Coin(gDenom1, 100))
		sigInAmt2 := types.NewCoins(types.NewInt64Coin(gDenom2, 100))

		So(SignInMsgForTest(t, app, true, account1, dexAccount1, sigInAmt1), ShouldBeNil)
		So(SignInMsgForTest(t, app, true, account1, dexAccount1, sigInAmt2), ShouldBeNil)

		a11, a21 := types.NewInt64Coin(gDenom1, 10), types.NewInt64Coin(gDenom2, 10)

		So(DexDealForTest(t, app, true, dexAccount1, account1, a11, account1, a21), ShouldBeNil)

		sigAmt := types.NewCoins(types.NewInt64Coin(gDenom1, 1900))
		So(SignInMsgForTest(t, app, true, account1, dexAccount1, sigAmt), ShouldBeNil)

		sigAmt2 := types.NewCoins(types.NewInt64Coin(gDenom2, 1900))
		So(SignInMsgForTest(t, app, true, account1, dexAccount1, sigAmt2), ShouldBeNil)
	})
}
