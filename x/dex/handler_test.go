package dex_test

import (
	"testing"
	"time"

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
		var (
			acc     = account4
			accName = name4
			auth    = wallet.GetAuth(acc)
		)

		So(CreateDexForTest(t, app, wallet,
			true,
			acc, types.NewInt64CoreCoins(111),
			[]byte("account4")), ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx := app.NewTestContext()
		dex, ok := app.DexKeeper().GetDex(ctx, accName)

		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, accName)

		msgUpdateDexDescription := dexTypes.NewMsgUpdateDexDescription(auth,
			accName,
			[]byte("xxx.yyy.zzz"))
		So(msgUpdateDexDescription.ValidateBasic(), ShouldBeNil)

		tx := simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgUpdateDexDescription,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(ok, ShouldBeTrue)
		So(dex, ShouldNotBeNil)
		So(dex.Description, ShouldEqual, "xxx.yyy.zzz")

		msgUpdateDexDescription = dexTypes.NewMsgUpdateDexDescription(auth,
			accName,
			[]byte(rand.Str(dexTypes.MaxDexDescriptorLen)))
		So(msgUpdateDexDescription.ValidateBasic(), ShouldNotBeNil)
	})
}

func TestHandleDestroyDex(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test destroy dex", t, func() {
		var (
			acc     = account4
			accName = name4
			auth    = wallet.GetAuth(acc)
		)

		So(CreateDexForTest(t, app, wallet,
			true,
			acc, types.NewInt64CoreCoins(1000),
			[]byte("account4")), ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx := app.NewTestContext()
		dex, ok := app.DexKeeper().GetDex(ctx, accName)
		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, accName)
		So(dex.Number, ShouldEqual, 0)

		coins := app.AssetKeeper().GetCoinPowers(ctx, acc)
		So(coins.IsZero(), ShouldBeTrue)

		msgDestroyDex := dexTypes.NewMsgDestroyDex(
			auth,
			accName)
		So(msgDestroyDex.ValidateBasic(), ShouldBeNil)

		tx := simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgDestroyDex,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(ok, ShouldBeFalse)
		So(dex, ShouldBeNil)
		coins = app.AssetKeeper().GetCoinPowers(ctx, acc)
		So(coins.IsEqual(types.NewInt64CoreCoins(1000)), ShouldBeTrue)
	})
}

func TestHandleCreateSymbol(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test create symbol", t, func() {
		var (
			acc     = account4
			accName = name4
			auth    = wallet.GetAuth(acc)
		)

		So(CreateDexForTest(t, app, wallet,
			true,
			acc, types.NewInt64CoreCoins(1000),
			[]byte("account4")), ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx := app.NewTestContext()
		dex, ok := app.DexKeeper().GetDex(ctx, accName)
		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, accName)
		So(dex.Number, ShouldEqual, 0)

		symbol := &dexTypes.Symbol{
			Base: dexTypes.BaseCurrency{
				CurrencyBase: dexTypes.CurrencyBase{
					Code:     "1",
					Name:     "BTC",
					FullName: "BTC",
					IconUrl:  "???",
					TxUrl:    "???",
				},
			},
			Quote: dexTypes.QuoteCurrency{
				CurrencyBase: dexTypes.CurrencyBase{
					Code:     "1",
					Name:     "USDT",
					FullName: "USDT",
					IconUrl:  "???",
					TxUrl:    "???",
				},
			},
			CreateTime:    time.Now(),
			DomainAddress: "http:///www.foo.com",
		}

		msgCreateSymbol := dexTypes.NewMsgCreateSymbol(auth,
			accName,
			&symbol.Base,
			&symbol.Quote,
			symbol.DomainAddress,
			symbol.CreateTime)

		So(msgCreateSymbol.ValidateBasic(), ShouldBeNil)

		tx := simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgCreateSymbol,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(ok, ShouldBeTrue)
		savedSymbol, ok := dex.Symbol(symbol.Base.Code, symbol.Quote.Code)
		So(ok, ShouldBeTrue)
		symbol.Height = savedSymbol.Height // ignore height check
		So(savedSymbol.Equal(symbol), ShouldBeTrue)

		simapp.AfterBlockCommitted(app, 1)

		symbol.Quote.Code = "2"
		msgCreateSymbol = dexTypes.NewMsgCreateSymbol(auth,
			accName,
			&symbol.Base,
			&symbol.Quote,
			symbol.DomainAddress,
			symbol.CreateTime)

		So(msgCreateSymbol.ValidateBasic(), ShouldBeNil)

		tx = simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgCreateSymbol,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err = simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(ok, ShouldBeTrue)
		savedSymbol, ok = dex.Symbol(symbol.Base.Code, symbol.Quote.Code)
		So(ok, ShouldBeTrue)
		symbol.Height = savedSymbol.Height // ignore height check
		So(savedSymbol.Equal(symbol), ShouldBeTrue)

		invalidSymbol := symbol
		invalidSymbol.Base.Code = ""
		So(dexTypes.NewMsgCreateSymbol(auth,
			accName,
			&invalidSymbol.Base,
			&invalidSymbol.Quote,
			invalidSymbol.DomainAddress,
			invalidSymbol.CreateTime).ValidateBasic(), ShouldNotBeNil)

		invalidSymbol = symbol
		invalidSymbol.Quote.Code = ""
		So(dexTypes.NewMsgCreateSymbol(auth,
			accName,
			&invalidSymbol.Base,
			&invalidSymbol.Quote,
			invalidSymbol.DomainAddress,
			invalidSymbol.CreateTime).ValidateBasic(), ShouldNotBeNil)

		invalidSymbol = symbol
		invalidSymbol.Base.Code = ""
		invalidSymbol.Quote.Code = ""
		So(dexTypes.NewMsgCreateSymbol(auth,
			accName,
			&invalidSymbol.Base,
			&invalidSymbol.Quote,
			invalidSymbol.DomainAddress,
			invalidSymbol.CreateTime).ValidateBasic(), ShouldNotBeNil)

		invalidSymbol = symbol
		invalidSymbol.Base.Code = "1"
		invalidSymbol.Base.TxUrl = ""
		invalidSymbol.Quote.Code = "3"
		So(dexTypes.NewMsgCreateSymbol(auth,
			accName,
			&invalidSymbol.Base,
			&invalidSymbol.Quote,
			invalidSymbol.DomainAddress,
			invalidSymbol.CreateTime).ValidateBasic(), ShouldNotBeNil)
	})
}

func TestHandleUpdateSymbol(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test update symbol", t, func() {
		var (
			acc     = account4
			accName = name4
			auth    = wallet.GetAuth(acc)
		)

		So(CreateDexForTest(t, app, wallet,
			true,
			acc, types.NewInt64CoreCoins(1000),
			[]byte("account4")), ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx := app.NewTestContext()
		dex, ok := app.DexKeeper().GetDex(ctx, accName)
		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, accName)
		So(dex.Number, ShouldEqual, 0)

		symbol := &dexTypes.Symbol{
			Base: dexTypes.BaseCurrency{
				CurrencyBase: dexTypes.CurrencyBase{
					Code:     "1",
					Name:     "BTC",
					FullName: "BTC",
					IconUrl:  "???",
					TxUrl:    "???",
				},
			},
			Quote: dexTypes.QuoteCurrency{
				CurrencyBase: dexTypes.CurrencyBase{
					Code:     "1",
					Name:     "USDT",
					FullName: "USDT",
					IconUrl:  "???",
					TxUrl:    "???",
				},
			},
			CreateTime:    time.Now(),
			DomainAddress: "http:///www.foo.com",
		}

		msgCreateSymbol := dexTypes.NewMsgCreateSymbol(auth,
			accName,
			&symbol.Base,
			&symbol.Quote,
			symbol.DomainAddress,
			symbol.CreateTime)
		So(msgCreateSymbol.ValidateBasic(), ShouldBeNil)

		tx := simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgCreateSymbol,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(ok, ShouldBeTrue)
		savedSymbol, ok := dex.Symbol(symbol.Base.Code, symbol.Quote.Code)
		So(ok, ShouldBeTrue)
		symbol.Height = savedSymbol.Height // ignore height check
		So(savedSymbol.Equal(symbol), ShouldBeTrue)

		copiedSymbol := *symbol
		copiedSymbol.Base.IconUrl = "base.icon.url"
		copiedSymbol.Quote.IconUrl = "quote.icon.url"
		msgUpdateSymbol := dexTypes.NewMsgUpdateSymbol(auth,
			accName,
			&copiedSymbol.Base,
			&copiedSymbol.Quote)
		So(msgCreateSymbol.ValidateBasic(), ShouldBeNil)

		tx = simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgUpdateSymbol,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err = simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(ok, ShouldBeTrue)
		savedSymbol, ok = dex.Symbol(copiedSymbol.Base.Code, copiedSymbol.Quote.Code)
		So(ok, ShouldBeTrue)
		So(savedSymbol.Equal(&copiedSymbol), ShouldBeTrue)

		copiedSymbol = *symbol
		copiedSymbol.Base.Code = ""
		So(dexTypes.NewMsgUpdateSymbol(auth,
			accName,
			&copiedSymbol.Base,
			&copiedSymbol.Quote).ValidateBasic(), ShouldNotBeNil)

		copiedSymbol = *symbol
		copiedSymbol.Quote.Code = ""
		So(dexTypes.NewMsgUpdateSymbol(auth,
			accName,
			&copiedSymbol.Base,
			&copiedSymbol.Quote).ValidateBasic(), ShouldNotBeNil)

		copiedSymbol = *symbol
		copiedSymbol.Base.Code = ""
		copiedSymbol.Quote.Code = ""
		So(dexTypes.NewMsgUpdateSymbol(auth,
			accName,
			&copiedSymbol.Base,
			&copiedSymbol.Quote).ValidateBasic(), ShouldNotBeNil)

		var emptySymbol dexTypes.Symbol
		emptySymbol.Base.Code = symbol.Base.Code
		emptySymbol.Quote.Code = symbol.Quote.Code
		So(dexTypes.NewMsgUpdateSymbol(auth,
			accName,
			&emptySymbol.Base,
			&emptySymbol.Quote).ValidateBasic(), ShouldNotBeNil)
	})
}

func TestHandlePauseSymbol(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test pause symbol handler", t, func() {
		var (
			acc     = account4
			accName = name4
			auth    = wallet.GetAuth(acc)
		)

		So(CreateDexForTest(t, app, wallet,
			true,
			acc, types.NewInt64CoreCoins(1000),
			[]byte("account4")), ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx := app.NewTestContext()
		dex, ok := app.DexKeeper().GetDex(ctx, accName)
		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, accName)
		So(dex.Number, ShouldEqual, 0)

		symbol := dexTypes.Symbol{
			Base: dexTypes.BaseCurrency{
				CurrencyBase: dexTypes.CurrencyBase{
					Code:     "1",
					Name:     "BTC",
					FullName: "BTC",
					IconUrl:  "???",
					TxUrl:    "???",
				},
			},
			Quote: dexTypes.QuoteCurrency{
				CurrencyBase: dexTypes.CurrencyBase{
					Code:     "1",
					Name:     "USDT",
					FullName: "USDT",
					IconUrl:  "???",
					TxUrl:    "???",
				},
			},
			CreateTime:    time.Now(),
			DomainAddress: "http:///www.foo.com",
		}

		msgCreateSymbol := dexTypes.NewMsgCreateSymbol(auth,
			accName,
			&symbol.Base,
			&symbol.Quote,
			symbol.DomainAddress,
			symbol.CreateTime)
		So(msgCreateSymbol.ValidateBasic(), ShouldBeNil)

		tx := simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgCreateSymbol,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(ok, ShouldBeTrue)
		savedSymbol, ok := dex.Symbol(symbol.Base.Code, symbol.Quote.Code)
		So(ok, ShouldBeTrue)
		symbol.Height = savedSymbol.Height // ignore height check
		So(savedSymbol.Equal(&symbol), ShouldBeTrue)

		msgPauseSymbol := dexTypes.NewMsgPauseSymbol(auth,
			accName,
			symbol.Base.Code,
			symbol.Quote.Code)
		So(msgPauseSymbol.ValidateBasic(), ShouldBeNil)

		tx = simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgPauseSymbol,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err = simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, accName)
		So(dex.Number, ShouldEqual, 0)

		symbol, ok = dex.Symbol(symbol.Base.Code, symbol.Quote.Code)
		So(ok, ShouldBeTrue)
		So(symbol.Paused(), ShouldBeTrue)

		copiedSymbol := symbol
		copiedSymbol.Base.Code = ""
		So(dexTypes.NewMsgPauseSymbol(auth,
			accName,
			copiedSymbol.Base.Code,
			copiedSymbol.Quote.Code).ValidateBasic(), ShouldNotBeNil)

		copiedSymbol = symbol
		copiedSymbol.Quote.Code = ""
		So(dexTypes.NewMsgPauseSymbol(auth,
			accName,
			copiedSymbol.Base.Code,
			copiedSymbol.Quote.Code).ValidateBasic(), ShouldNotBeNil)

		copiedSymbol = symbol
		copiedSymbol.Base.Code = ""
		copiedSymbol.Quote.Code = ""
		So(dexTypes.NewMsgPauseSymbol(auth,
			accName,
			copiedSymbol.Base.Code,
			copiedSymbol.Quote.Code).ValidateBasic(), ShouldNotBeNil)
	})
}

func TestHandleRestoreSymbol(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test restore symbol handler", t, func() {
		var (
			acc     = account4
			accName = name4
			auth    = wallet.GetAuth(acc)
		)

		So(CreateDexForTest(t, app, wallet,
			true,
			acc, types.NewInt64CoreCoins(1000),
			[]byte("account4")), ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx := app.NewTestContext()
		dex, ok := app.DexKeeper().GetDex(ctx, accName)
		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, accName)
		So(dex.Number, ShouldEqual, 0)

		symbol := dexTypes.Symbol{
			Base: dexTypes.BaseCurrency{
				CurrencyBase: dexTypes.CurrencyBase{
					Code:     "1",
					Name:     "BTC",
					FullName: "BTC",
					IconUrl:  "???",
					TxUrl:    "???",
				},
			},
			Quote: dexTypes.QuoteCurrency{
				CurrencyBase: dexTypes.CurrencyBase{
					Code:     "1",
					Name:     "USDT",
					FullName: "USDT",
					IconUrl:  "???",
					TxUrl:    "???",
				},
			},
			CreateTime:    time.Now(),
			DomainAddress: "http:///www.foo.com",
		}

		msgCreateSymbol := dexTypes.NewMsgCreateSymbol(auth,
			accName,
			&symbol.Base,
			&symbol.Quote,
			symbol.DomainAddress,
			symbol.CreateTime)
		So(msgCreateSymbol.ValidateBasic(), ShouldBeNil)

		tx := simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgCreateSymbol,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(ok, ShouldBeTrue)
		savedSymbol, ok := dex.Symbol(symbol.Base.Code, symbol.Quote.Code)
		So(ok, ShouldBeTrue)
		symbol.Height = savedSymbol.Height // ignore height check
		So(savedSymbol.Equal(&symbol), ShouldBeTrue)

		msgPauseSymbol := dexTypes.NewMsgPauseSymbol(auth,
			accName,
			symbol.Base.Code,
			symbol.Quote.Code)
		So(msgPauseSymbol.ValidateBasic(), ShouldBeNil)

		tx = simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgPauseSymbol,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err = simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, accName)
		So(dex.Number, ShouldEqual, 0)

		symbol, ok = dex.Symbol(symbol.Base.Code, symbol.Quote.Code)
		So(ok, ShouldBeTrue)
		So(symbol.Paused(), ShouldBeTrue)

		msgRestoreSymbol := dexTypes.NewMsgRestoreSymbol(auth,
			accName,
			symbol.Base.Code,
			symbol.Quote.Code)
		So(msgRestoreSymbol.ValidateBasic(), ShouldBeNil)

		tx = simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgRestoreSymbol,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err = simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, accName)
		So(dex.Number, ShouldEqual, 0)

		symbol, ok = dex.Symbol(symbol.Base.Code, symbol.Quote.Code)
		So(ok, ShouldBeTrue)
		So(symbol.Paused(), ShouldBeFalse)

		copiedSymbol := symbol
		copiedSymbol.Base.Code = ""
		So(dexTypes.NewMsgRestoreSymbol(auth,
			accName,
			copiedSymbol.Base.Code,
			copiedSymbol.Quote.Code).ValidateBasic(), ShouldNotBeNil)

		copiedSymbol = symbol
		copiedSymbol.Quote.Code = ""
		So(dexTypes.NewMsgRestoreSymbol(auth,
			accName,
			copiedSymbol.Base.Code,
			copiedSymbol.Quote.Code).ValidateBasic(), ShouldNotBeNil)

		copiedSymbol = symbol
		copiedSymbol.Base.Code = ""
		copiedSymbol.Quote.Code = ""
		So(dexTypes.NewMsgRestoreSymbol(auth,
			accName,
			copiedSymbol.Base.Code,
			copiedSymbol.Quote.Code).ValidateBasic(), ShouldNotBeNil)
	})
}

func TestShutdownSymbol(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test shutdown symbol handler", t, func() {
		var (
			acc     = account4
			accName = name4
			auth    = wallet.GetAuth(acc)
		)

		So(CreateDexForTest(t, app, wallet,
			true,
			acc, types.NewInt64CoreCoins(1000),
			[]byte("account4")), ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx := app.NewTestContext()
		dex, ok := app.DexKeeper().GetDex(ctx, accName)
		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, accName)
		So(dex.Number, ShouldEqual, 0)

		symbol := dexTypes.Symbol{
			Base: dexTypes.BaseCurrency{
				CurrencyBase: dexTypes.CurrencyBase{
					Code:     "1",
					Name:     "BTC",
					FullName: "BTC",
					IconUrl:  "???",
					TxUrl:    "???",
				},
			},
			Quote: dexTypes.QuoteCurrency{
				CurrencyBase: dexTypes.CurrencyBase{
					Code:     "1",
					Name:     "USDT",
					FullName: "USDT",
					IconUrl:  "???",
					TxUrl:    "???",
				},
			},
			CreateTime:    time.Now(),
			DomainAddress: "http:///www.foo.com",
		}

		msgCreateSymbol := dexTypes.NewMsgCreateSymbol(auth,
			accName,
			&symbol.Base,
			&symbol.Quote,
			symbol.DomainAddress,
			symbol.CreateTime)
		So(msgCreateSymbol.ValidateBasic(), ShouldBeNil)

		tx := simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgCreateSymbol,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(ok, ShouldBeTrue)
		savedSymbol, ok := dex.Symbol(symbol.Base.Code, symbol.Quote.Code)
		So(ok, ShouldBeTrue)
		symbol.Height = savedSymbol.Height // ignore height check
		So(savedSymbol.Equal(&symbol), ShouldBeTrue)

		msgShutdownSymbol := dexTypes.NewMsgShutdownSymbol(auth,
			accName,
			symbol.Base.Code,
			symbol.Quote.Code)
		So(msgShutdownSymbol.ValidateBasic(), ShouldBeNil)

		tx = simapp.NewTxForTest(
			acc,
			[]sdk.Msg{
				&msgShutdownSymbol,
			}, wallet.PrivKey(auth))
		ctx = app.NewTestContext()
		err = simapp.CheckTxs(t, app, ctx, tx)
		So(err, ShouldBeNil)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		dex, ok = app.DexKeeper().GetDex(ctx, accName)
		So(dex, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
		So(dex.Creator, simapp.ShouldEq, accName)
		So(dex.Number, ShouldEqual, 0)

		symbol, ok = dex.Symbol(symbol.Base.Code, symbol.Quote.Code)
		So(ok, ShouldBeFalse)

		copiedSymbol := symbol
		copiedSymbol.Base.Code = ""
		So(dexTypes.NewMsgShutdownSymbol(auth,
			accName,
			copiedSymbol.Base.Code,
			copiedSymbol.Quote.Code).ValidateBasic(), ShouldNotBeNil)

		copiedSymbol = symbol
		copiedSymbol.Quote.Code = ""
		So(dexTypes.NewMsgShutdownSymbol(auth,
			accName,
			copiedSymbol.Base.Code,
			copiedSymbol.Quote.Code).ValidateBasic(), ShouldNotBeNil)

		copiedSymbol = symbol
		copiedSymbol.Base.Code = ""
		copiedSymbol.Quote.Code = ""
		So(dexTypes.NewMsgShutdownSymbol(auth,
			accName,
			copiedSymbol.Base.Code,
			copiedSymbol.Quote.Code).ValidateBasic(), ShouldNotBeNil)
	})
}
