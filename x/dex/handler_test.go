package dex_test

import (
	"fmt"
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	dexTypes "github.com/KuChainNetwork/kuchain/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tendermint/tendermint/libs/rand"
)

func CreateDexForTest(t *testing.T, app *simapp.SimApp, wallet *simapp.Wallet, isSuccess bool, account types.AccountID, stakings types.Coins, desc []byte) error {
	ctx := app.NewTestContext()

	var (
		acc     = account
		accName = account.MustName()
		auth    = wallet.GetAuth(acc)
	)

	msg := dexTypes.NewMsgCreateDex(
		auth,
		accName,
		stakings,
		desc)

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	tx := simapp.NewTxForTest(
		acc,
		[]sdk.Msg{
			&msg,
		}, wallet.PrivKey(auth))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return simapp.CheckTxs(t, app, ctx, tx)
}

func TestHandleCreateDex(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test create dex", t, func() {
		ctx := app.NewTestContext()

		var (
			acc     = account5
			accName = name5
			auth    = wallet.GetAuth(acc)
		)

		msg := dexTypes.NewMsgCreateDex(
			auth,
			accName,
			types.NewInt64CoreCoins(1000000),
			[]byte("dex for test"))

		So(msg.ValidateBasic(), ShouldBeNil)

		tx := simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msg,
			}, wallet.PrivKey(auth))

		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)
	})

	Convey("test transfer no enough", t, func() {
		ctx := app.NewTestContext()

		var (
			acc     = account4
			accName = name4
			auth    = wallet.GetAuth(acc)
		)

		msg := dexTypes.NewMsgCreateDex(
			auth,
			accName,
			types.NewInt64CoreCoins(1000000),
			[]byte("dex for test"))

		msg.Transfers[0].Amount[0].Amount = types.NewInt(111)

		So(msg.ValidateBasic(), simapp.ShouldErrIs, types.ErrKuMsgFromNotEqual)

		tx := simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msg,
			}, wallet.PrivKey(auth)).WithCannotPass()

		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, simapp.ShouldErrIs, types.ErrKuMsgFromNotEqual)
	})

	Convey("test transfer no module", t, func() {
		ctx := app.NewTestContext()

		var (
			acc     = account4
			accName = name4
			auth    = wallet.GetAuth(acc)
		)

		msg := dexTypes.NewMsgCreateDex(
			auth,
			accName,
			types.NewInt64CoreCoins(1000000),
			[]byte("dex for test"))

		msg.Transfers[0].To = acc

		So(msg.ValidateBasic(), simapp.ShouldErrIs, types.ErrKuMsgFromNotEqual)

		tx := simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msg,
			}, wallet.PrivKey(auth)).WithCannotPass()

		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, simapp.ShouldErrIs, types.ErrKuMsgFromNotEqual)
	})

	Convey("test transfer from error", t, func() {
		ctx := app.NewTestContext()

		var (
			acc     = account4
			accName = name4
			auth    = wallet.GetAuth(acc)
		)

		msg := dexTypes.NewMsgCreateDex(
			auth,
			accName,
			types.NewInt64CoreCoins(1000000),
			[]byte("dex for test"))

		msg.Transfers[0].From = account5

		So(msg.ValidateBasic(), simapp.ShouldErrIs, types.ErrKuMsgFromNotEqual)

		tx := simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msg,
			}, wallet.PrivKey(auth)).WithCannotPass()

		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, simapp.ShouldErrIs, types.ErrKuMsgFromNotEqual)
	})

	Convey("test transfer desc too long", t, func() {
		So(CreateDexForTest(t, app, wallet,
			false,
			account4, types.NewInt64CoreCoins(111),
			make([]byte, 520)), simapp.ShouldErrIs, dexTypes.ErrDexDescTooLong)
	})

	Convey("test transfer create two times", t, func() {
		So(CreateDexForTest(t, app, wallet,
			false,
			account5, types.NewInt64CoreCoins(111),
			[]byte("hello")), simapp.ShouldErrIs, dexTypes.ErrDexHadCreated)
	})
}

func TestHandleCreateDexNumber(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test create dex number", t, func() {
		So(CreateDexForTest(t, app, wallet,
			true,
			account4, types.NewInt64CoreCoins(111),
			[]byte("account4")), ShouldBeNil)

		So(CreateDexForTest(t, app, wallet,
			true,
			account5, types.NewInt64CoreCoins(111),
			[]byte("account5")), ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx := app.NewTestContext()
		dex1, ok := app.DexKeeper().GetDex(ctx, name4)

		So(dex1, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex1.Creator, simapp.ShouldEq, name4)
		So(dex1.Number, ShouldEqual, 0)

		dex2, ok := app.DexKeeper().GetDex(ctx, name5)

		So(dex2, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex2.Creator, simapp.ShouldEq, name5)
		So(dex2.Number, ShouldEqual, 1)
	})
}

func TestHandleUpdateDexDescription(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test update dex description", t, func() {
		So(CreateDexForTest(t, app, wallet,
			true,
			account4, types.NewInt64CoreCoins(111),
			[]byte("account4")), ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx := app.NewTestContext()
		dex, ok := app.DexKeeper().GetDex(ctx, name4)

		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, name4)
		So(dex.Number, ShouldEqual, 0)

		for i := 0; i < 2; i++ {
			desc := fmt.Sprintf("test%d", i)
			err := app.DexKeeper().UpdateDexDescription(ctx, name4, desc)
			So(err, ShouldBeNil)
			theDex, ok := app.DexKeeper().GetDex(ctx, name4)
			So(ok, ShouldBeTrue)
			So(theDex.Description, ShouldEqual, desc)
		}
		// check max length
		err := app.DexKeeper().UpdateDexDescription(ctx, name4, "test")
		So(err, ShouldBeNil)
		theDex, ok := app.DexKeeper().GetDex(ctx, name4)
		So(ok, ShouldBeTrue)
		So(theDex.Description, ShouldEqual, "test")

		desc := rand.Str(dexTypes.MaxDexDescriptorLen)
		err = app.DexKeeper().UpdateDexDescription(ctx, name4, desc)
		So(err, ShouldEqual, dexTypes.ErrDexDescTooLong)
		dex, ok = app.DexKeeper().GetDex(ctx, name4)
		So(ok, ShouldBeTrue)
		So(dex, ShouldNotBeNil)
		So(dex.Description, ShouldEqual, "test")
	})
}

func TestHandleDestroyDex(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test destroy dex", t, func() {
		So(CreateDexForTest(t, app, wallet,
			true,
			account4, types.NewInt64CoreCoins(1000),
			[]byte("account4")), ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx := app.NewTestContext()
		dex, ok := app.DexKeeper().GetDex(ctx, name4)

		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, name4)
		So(dex.Number, ShouldEqual, 0)

		theDex, ok := app.DexKeeper().GetDex(ctx, name4)
		So(ok, ShouldBeTrue)
		So(theDex.DestroyFlag, ShouldBeFalse)

		err := app.DexKeeper().UpdateDexDescription(ctx, name4, "test0")
		So(err, ShouldBeNil)

		coins := app.AssetKeeper().GetCoinPowers(ctx, account4)
		So(len(coins), ShouldEqual, 0)

		err = app.DexKeeper().DestroyDex(ctx, name4)
		So(err, ShouldBeNil)

		coins = app.AssetKeeper().GetCoinPowers(ctx, account4)
		So(coins.IsEqual(types.NewInt64CoreCoins(1000)), ShouldBeTrue)

		theDex, ok = app.DexKeeper().GetDex(ctx, name4)
		So(ok, ShouldBeTrue)
		So(theDex.DestroyFlag, ShouldBeTrue)

		err = app.DexKeeper().DestroyDex(ctx, name4)
		So(err, ShouldNotBeNil)

		err = app.DexKeeper().UpdateDexDescription(ctx, name4, "test1")
		So(err, ShouldNotBeNil)
	})
}
