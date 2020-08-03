package ante_test

import (
	"testing"

	. "github.com/KuChainNetwork/kuchain/chain/ante"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTxFee(t *testing.T) {
	Convey("test fee check", t, func() {

	})
}

func TestAnteFeeHandler(t *testing.T) {
	app, _ := createAppForTest()

	ak := *app.AccountKeeper()
	asset := app.AssetKeeper()
	handler := sdk.ChainAnteDecorators(
		NewDeductFeeDecorator(ak, asset),
	)

	Convey("DeductFeeDecorator test fee ante fee handler", t, func() {
		stdTx4Test := testStdTx(app, account4)

		Convey("test fee cost", func() {
			amtOld := asset.GetAllBalances(app.NewTestContext(), account4)

			_, err := handler(app.NewTestContext(), stdTx4Test, true)
			So(err, ShouldBeNil)

			amtAfter := asset.GetAllBalances(app.NewTestContext(), account4)
			So(amtOld.Sub(simapp.DefaultTestFee), simapp.ShouldEq, amtAfter)
		})

		Convey("test fee payer account auth", func() {
			stdTx4Test.Fee.Payer = account2

			amtOld := asset.GetAllBalances(app.NewTestContext(), account4)
			amtOld2 := asset.GetAllBalances(app.NewTestContext(), account2)

			_, err := handler(app.NewTestContext(), stdTx4Test, true)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrUnauthorized)

			amtAfter := asset.GetAllBalances(app.NewTestContext(), account4)
			amtAfter2 := asset.GetAllBalances(app.NewTestContext(), account2)

			So(amtOld, simapp.ShouldEq, amtAfter)
			So(amtOld2, simapp.ShouldEq, amtAfter2)
		})

		Convey("test fee payer use its address auth, but it no have coins", func() {
			stdTx4Test2 := testStdTx(app, account2)
			stdTx4Test2.Fee.Payer = addAccount2

			amtOld := asset.GetAllBalances(app.NewTestContext(), account2)

			_, err := handler(app.NewTestContext(), stdTx4Test2, true)
			So(err, simapp.ShouldErrIs, assetTypes.ErrAssetCoinNoEnough)

			amtAfter := asset.GetAllBalances(app.NewTestContext(), account2)
			So(amtOld, simapp.ShouldEq, amtAfter)
		})

		Convey("test fee payer use its address auth, have coins to sub", func() {
			stdTx4Test.Fee.Payer = types.NewAccountIDFromAccAdd(wallet.GetAuth(account4))

			amtOld := asset.GetAllBalances(app.NewTestContext(), account4)
			addAmtOld := asset.GetAllBalances(app.NewTestContext(), types.NewAccountIDFromAccAdd(addr4))

			_, err := handler(app.NewTestContext(), stdTx4Test, true)
			So(err, ShouldBeNil)

			amtAfter := asset.GetAllBalances(app.NewTestContext(), account4)
			addAmtAfter := asset.GetAllBalances(app.NewTestContext(), types.NewAccountIDFromAccAdd(addr4))

			So(amtOld, simapp.ShouldEq, amtAfter)

			// address coins will cost fee
			So(addAmtOld.Sub(simapp.DefaultTestFee), simapp.ShouldEq, addAmtAfter)
		})

		Convey("test fee payer address auth", func() {
			stdTx4Test.Fee.Payer = types.NewAccountIDFromAccAdd(wallet.GetAuth(account2))

			amtOld := asset.GetAllBalances(app.NewTestContext(), account4)
			amtOld2 := asset.GetAllBalances(app.NewTestContext(), account2)

			_, err := handler(app.NewTestContext(), stdTx4Test, true)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrUnauthorized)

			amtAfter := asset.GetAllBalances(app.NewTestContext(), account4)
			amtAfter2 := asset.GetAllBalances(app.NewTestContext(), account2)

			So(amtOld, simapp.ShouldEq, amtAfter)
			So(amtOld2, simapp.ShouldEq, amtAfter2)
		})

		Convey("test fee payer no existing", func() {
			stdTx4Test.Fee.Payer = types.MustAccountID("adddddd")

			amtOld := asset.GetAllBalances(app.NewTestContext(), account4)

			_, err := handler(app.NewTestContext(), stdTx4Test, true)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrUnknownAddress)

			amtAfter := asset.GetAllBalances(app.NewTestContext(), account4)
			So(amtOld, simapp.ShouldEq, amtAfter)
		})

	})
}

func TestEnsureMempoolFees(t *testing.T) {
	// setup
	app, ctx := createAppForTest()

	Convey("test fee ensure mempool fees", t, func() {
		mfd := NewMempoolFeeDecorator()
		antehandler := sdk.ChainAnteDecorators(mfd)

		tx := testStdTx(app, account4)
		tx.Fee.Amount = types.NewInt64CoreCoins(100)
		tx.Fee.Gas = 100000

		atomPrice := types.NewDecCoinFromDec(constants.DefaultBondDenom,
			types.NewDec(200).Quo(types.NewDec(100000))) // 200 coin
		highGasPrice := types.NewDecCoins(atomPrice)

		atomPrice = types.NewDecCoinFromDec(constants.DefaultBondDenom,
			types.NewDec(0).Quo(types.NewDec(100000)))
		lowGasPrice := types.NewDecCoins(atomPrice)

		Convey("Decorator should have error on too low fee for local gasPrice", func() {
			// Set high gas price so standard test fee fails
			ctx = ctx.WithMinGasPrices(highGasPrice.ToSDK())

			// Set IsCheckTx to true
			ctx := ctx.WithIsCheckTx(true)

			// antehandler errors with insufficient fees
			_, err := antehandler(ctx, tx, false)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrInsufficientFee)
		})

		Convey("MempoolFeeDecorator returned error in DeliverTx", func() {
			// Set high gas price so standard test fee fails
			ctx = ctx.WithMinGasPrices(highGasPrice.ToSDK())

			// Set IsCheckTx to false
			ctx := ctx.WithIsCheckTx(false)

			// antehandler should not error since we do not check minGasPrice in DeliverTx
			_, err := antehandler(ctx, tx, false)
			So(err, ShouldBeNil)

		})

		Convey("Decorator should not have error on fee higher than local gasPrice", func() {

			// Set IsCheckTx back to true for testing sufficient mempool fee
			ctx := ctx.WithIsCheckTx(true)

			// Set low fee price
			ctx = ctx.WithMinGasPrices(lowGasPrice.ToSDK())

			_, err := antehandler(ctx, tx, false)
			So(err, ShouldBeNil)
		})

	})
}
