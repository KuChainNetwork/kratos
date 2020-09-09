package asset_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
)

func TestSendAddressNotEnoughBalance(t *testing.T) {
	Convey("TestSendAddressNotEnoughBalance", t, func() {
		asset1 := types.NewCoins(
			types.NewInt64Coin("foo/coin", 67),
			types.NewInt64Coin(constants.DefaultBondDenom, 10000000000))
		genAcc := simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1)

		genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAcc)
		app := simapp.SetupWithGenesisAccounts(genAccs)

		ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

		acc1 := app.AccountKeeper().GetAccount(ctxCheck, account1)

		So(acc1, ShouldNotBeNil)
		So(genAcc.GetID().Eq(acc1.GetID()), ShouldBeTrue)
		So(genAcc.GetAuth().Equals(acc1.GetAuth()), ShouldBeTrue)

		origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addr1)
		So(err, ShouldBeNil)

		ctxCheck.Logger().Info("auth nums", "seq", origAuthSeq, "num", origAuthNum)

		msg := assetTypes.NewMsgTransfer(addr1, account1, addAccount1, types.NewInt64Coins("foo/coin", 100))
		fee := types.NewInt64Coins(constants.DefaultBondDenom, 100000)

		header := abci.Header{Height: app.LastBlockHeight() + 1}
		_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
			header, account1, fee,
			[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
			false, false, wallet.PrivKey(addr1))
		So(err, ShouldNotBeNil)
		So(err, simapp.ShouldErrIs, assetTypes.ErrAssetCoinNoEnough)

		simapp.CheckBalance(t, app, account1, asset1.Sub(fee))
		simapp.CheckBalance(t, app, addAccount1, types.Coins{})

		ctxCheck = app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

		So(app.AccountKeeper().GetAccount(ctxCheck, account1), ShouldNotBeNil)

		authSeq, authNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addr1)
		So(err, ShouldBeNil)

		ctxCheck.Logger().Info("account seq", "seq", authSeq, "num", authNum)

		So(authNum, ShouldEqual, origAuthNum)
		So(authSeq, ShouldEqual, origAuthSeq+1)
	})
}

func TestCreateAsset(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test create asset", t, func() {
		So(createCoin(t, app, true, account4, types.MustName("abc"), 10000000000000), ShouldBeNil)
	})

	Convey("test create asset has exists", t, func() {
		// this has created in last
		So(createCoin(t, app, false, account4, types.MustName("abc"), 10000000000000),
			simapp.ShouldErrIs, assetTypes.ErrAssetHasCreated)
	})

	Convey("test create asset symbol no Validate", t, func() {
		symbolErrs := []types.Name{
			types.MustName("abc@aa"), // has @
			types.MustName("1aaaaa"), // begin with 0-9
		}

		for _, se := range symbolErrs {
			var (
				demon      = types.CoinDenom(name4, se)
				maxSupply  = types.Coin{demon, types.NewInt(10000000000000)}
				initSupply = types.Coin{demon, types.NewInt(0)}
				desc       = []byte(fmt.Sprintf("desc for %s", demon))
			)

			So(createCoinExt(t, app, false, account4, se, maxSupply, true, true, true, 0, initSupply, desc),
				simapp.ShouldErrIs, types.ErrCoinDenomInvalid)
		}

		var (
			se         = types.MustName("abc")
			demon      = types.CoinDenom(name1, se) // creator has @
			maxSupply  = types.Coin{demon, types.NewInt(10000000000000)}
			initSupply = types.Coin{demon, types.NewInt(0)}
			desc       = []byte(fmt.Sprintf("desc for %s", demon))
		)

		So(createCoinExt(t, app, false, account1, se, maxSupply, true, true, true, 0, initSupply, desc),
			simapp.ShouldErrIs, types.ErrCoinDenomInvalid)
	})

	Convey("test create asset desc too large", t, func() {
		// create in last test
		var (
			symbol     = types.MustName("abcd")
			demon      = types.CoinDenom(name4, symbol)
			maxSupply  = types.NewCoin(demon, types.NewInt(10000000000000))
			initSupply = types.NewCoin(demon, types.NewInt(0))
			desc       = []byte(make([]byte, assetTypes.CoinDescriptionLen+1))
		)

		So(createCoinExt(t, app, false, account4, symbol, maxSupply, true, true, true, 0, initSupply, desc),
			simapp.ShouldErrIs, assetTypes.ErrAssetDescriptorTooLarge)
	})

	Convey("test create symbol error asset", t, func() {
		var (
			symbol        = types.MustName("abc")
			demon         = types.CoinDenom(name4, symbol)
			demon2        = types.CoinDenom(name5, symbol)                      // wrong creator
			demon3        = types.CoinDenom(name4, types.MustName("aaa"))       // wrong symbol
			maxSupply     = types.NewCoin(demon, types.NewInt(10000000000000))  // use wrong creator
			maxSupplyErr1 = types.NewCoin(demon2, types.NewInt(10000000000000)) // use wrong creator
			maxSupplyErr2 = types.NewCoin(demon3, types.NewInt(10000000000000)) // use wrong symbol
			initSupply    = types.NewCoin(demon, types.NewInt(0))
			initSupplyErr = types.NewCoin(demon3, types.NewInt(0))
			desc          = []byte(fmt.Sprintf("desc for %s", demon))
		)

		So(createCoinExt(t, app, false, account4, symbol, maxSupplyErr1, true, true, true, 0, initSupply, desc),
			simapp.ShouldErrIs, assetTypes.ErrAssetSymbolError)

		So(createCoinExt(t, app, false, account4, symbol, maxSupplyErr2, true, true, true, 0, initSupply, desc),
			simapp.ShouldErrIs, assetTypes.ErrAssetSymbolError)

		So(createCoinExt(t, app, false, account4, symbol, maxSupply, true, true, true, 0, initSupplyErr, desc),
			simapp.ShouldErrIs, assetTypes.ErrAssetSymbolError)
	})
}

// TODO: add coin stat setOpt test

func TestIssueCoins(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test issue coins", t, func() {
		// account4 is create a asset first
		var (
			symbol             = types.MustName("abc")
			denom              = types.CoinDenom(name4, symbol)
			maxSupplyAmt int64 = 10000000000000
		)
		So(createCoin(t, app, true, account4, symbol, maxSupplyAmt), ShouldBeNil)

		ctx := app.NewTestContext()
		amt := app.AssetKeeper().GetAllBalances(ctx, account4)
		amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee
		So(amt.AmountOf(denom).Int64(), ShouldEqual, 0)

		// issue 10000 to self
		issueAmt := types.NewInt64Coin(denom, 10000)
		So(issueCoin(t, app, true, account4, symbol, issueAmt), ShouldBeNil)

		// account4 will add issueAmt
		simapp.CheckBalance(t, app, account4, amt.Add(issueAmt))
	})

	Convey("test issue coins err by amt symbol no equal", t, func() {
		// account4 is create a asset first
		var (
			symbol             = types.MustName("abc1")
			denom              = types.CoinDenom(name4, symbol)
			maxSupplyAmt int64 = 10000000000000
		)
		So(createCoin(t, app, true, account4, symbol, maxSupplyAmt), ShouldBeNil)

		ctx := app.NewTestContext()
		amt := app.AssetKeeper().GetAllBalances(ctx, account4)
		//amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee
		So(amt.AmountOf(denom).Int64(), ShouldEqual, 0)

		// issue 10000 to self
		issueAmt := types.NewInt64Coin(denom, 10000)
		So(issueCoin(t, app, false, account4, types.MustName("abc"), issueAmt),
			simapp.ShouldErrIs, assetTypes.ErrAssetSymbolError)

		// account4 will add issueAmt
		simapp.CheckBalance(t, app, account4, amt)

		// issue 10000 to self
		issueAmtOther := types.NewInt64Coin(types.CoinDenom(name4, types.MustName("aaa")), 10000)
		So(issueCoin(t, app, false, account4, symbol, issueAmtOther),
			simapp.ShouldErrIs, assetTypes.ErrAssetSymbolError)

		//amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee

		// account4 will add issueAmt
		simapp.CheckBalance(t, app, account4, amt)
	})

	Convey("test issue coins err by not creator", t, func() {
		var (
			symbol             = types.MustName("abc")
			denom              = types.CoinDenom(name4, symbol) // create in last
			maxSupplyAmt int64 = 10000000000000
		)
		So(createCoin(t, app, true, account3, symbol, maxSupplyAmt), ShouldBeNil)

		ctx := app.NewTestContext()
		amt := app.AssetKeeper().GetAllBalances(ctx, account3)
		//amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee
		So(amt.AmountOf(denom).Int64(), ShouldEqual, 0)

		issueAmt := types.NewInt64Coin(denom, 10000)
		So(issueCoin(t, app, false, account3, symbol, issueAmt),
			simapp.ShouldErrIs, assetTypes.ErrAssetSymbolError)
		simapp.CheckBalance(t, app, account3, amt) // account3 asset will no change
	})

	Convey("test issue coins err by coins not found", t, func() {
		var (
			symbol = types.MustName("abcdd")
			denom  = types.CoinDenom(name4, symbol) // create in last
		)

		ctx := app.NewTestContext()
		amt := app.AssetKeeper().GetAllBalances(ctx, account4)
		amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee
		So(amt.AmountOf(denom).Int64(), ShouldEqual, 0)

		issueAmt := types.NewInt64Coin(denom, 10000)
		So(issueCoin(t, app, false, account4, symbol, issueAmt),
			simapp.ShouldErrIs, assetTypes.ErrAssetCoinNoExit)
		simapp.CheckBalance(t, app, account4, amt)
	})
}

func TestIssueMaxSupplyCoreCoin(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test core coins no limit by max supply", t, func() {
		ctx := app.NewTestContext()
		stat, err := app.AssetKeeper().GetCoinStat(ctx, constants.ChainMainName, constants.DefaultBondSymbolName)
		So(err, ShouldBeNil)

		So(stat.GetCurrentMaxSupplyLimit(111).IsZero(), ShouldBeTrue)
		So(stat.MaxSupply.IsZero(), ShouldBeTrue)

		simapp.AfterBlockCommitted(app, 1)

		ctx = app.NewTestContext()
		statNew, err := app.AssetKeeper().GetCoinStat(ctx, constants.ChainMainName, constants.DefaultBondSymbolName)
		So(err, ShouldBeNil)

		So(statNew.GetCurrentMaxSupplyLimit(111).IsZero(), ShouldBeTrue)
		So(statNew.MaxSupply.IsZero(), ShouldBeTrue)
		So(statNew.Supply.IsGTE(stat.Supply), ShouldBeTrue)
		So(statNew.Supply.IsZero(), ShouldBeFalse)
	})
}

func TestIssueMaxSupply(t *testing.T) {
	app, _ := createAppForTest()
	Convey("test issue coins cannot > max supply", t, func() {
		// account4 is create a asset first
		var (
			symbol             = types.MustName("abc1")
			denom              = types.CoinDenom(name4, symbol)
			maxSupplyAmt int64 = 100000000000000
		)
		So(createCoin(t, app, true, account4, symbol, maxSupplyAmt), ShouldBeNil)

		amt := app.AssetKeeper().GetAllBalances(app.NewTestContext(), account4)
		So(amt.AmountOf(denom).Int64(), ShouldEqual, 0)
		amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee

		issueAmt := types.NewInt64Coin(denom, maxSupplyAmt-100) // this can ok
		So(issueCoin(t, app, true, account4, symbol, issueAmt), ShouldBeNil)
		amt = amt.Add(issueAmt)
		simapp.CheckBalance(t, app, account4, amt) // check if added

		issueAmtTooLarge := types.NewInt64Coin(denom, 101) // this can ok
		So(issueCoin(t, app, false, account4, symbol, issueAmtTooLarge),
			simapp.ShouldErrIs, assetTypes.ErrAssetIssueGTMaxSupply)
		amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee
		simapp.CheckBalance(t, app, account4, amt)  // check if not added

		issueAmtJustOK := types.NewInt64Coin(denom, 100) // this can ok
		So(issueCoin(t, app, true, account4, symbol, issueAmtJustOK), ShouldBeNil)
		amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee
		amt = amt.Add(issueAmtJustOK)
		simapp.CheckBalance(t, app, account4, amt) // check if added
	})
}

func TestIssueIsCanIssueTag(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test issue coins err coins cannot issue, but can issue in some blocks", t, func() {
		var (
			symbol             = types.MustName("abcd1")
			denom              = types.CoinDenom(name4, symbol) // create in last
			maxSupplyAmt int64 = 10000000000000
		)

		currBlockHeight := app.LastBlockHeight() + 1
		So(createCoinExt(t, app, true,
			account4, symbol,
			types.NewInt64Coin(denom, maxSupplyAmt),
			false, // cannot issue
			true, true, 0,
			types.NewInt64Coin(denom, 0), []byte("cannot issue")), ShouldBeNil)

		ctx := app.NewTestContext()
		amt := app.AssetKeeper().GetAllBalances(ctx, account4)
		amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee
		So(amt.AmountOf(denom).Int64(), ShouldEqual, 0)

		stat, err := app.AssetKeeper().GetCoinStat(ctx, name4, symbol)
		So(err, ShouldBeNil)

		So(stat.CanIssue, ShouldEqual, false)
		So(stat.CanLock, ShouldEqual, true)
		So(stat.CreateHeight, ShouldEqual, currBlockHeight)

		issueAmt := types.NewInt64Coin(denom, 10000)
		So(issueCoin(t, app, true, account4, symbol, issueAmt), ShouldBeNil) // after a block, can issue
		amt = amt.Add(issueAmt)
		simapp.CheckBalance(t, app, account4, amt)

		simapp.AfterBlockCommitted(app, int(constants.IssueCoinsWaitBlockNums-2)) // in this time, next block can issue

		amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee
		amt = amt.Add(issueAmt)
		So(issueCoin(t, app, true, account4, symbol, issueAmt), ShouldBeNil) // after IssueCoinsWaitBlockNums - 1 block, can issue
		simapp.CheckBalance(t, app, account4, amt)

		simapp.AfterBlockCommitted(app, 1) // in this time, next block cannot issue

		amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee
		So(issueCoin(t, app, false, account4, symbol, issueAmt),
			simapp.ShouldErrIs, assetTypes.ErrAssetCoinCannotBeIssue) // after IssueCoinsWaitBlockNums - 1 block, can issue
		simapp.CheckBalance(t, app, account4, amt)

		simapp.AfterBlockCommitted(app, 10) // in this time, next block also cannot issue

		amt, _ = amt.SafeSub(simapp.DefaultTestFee) // issue will cost fee
		So(issueCoin(t, app, false, account4, symbol, issueAmt),
			simapp.ShouldErrIs, assetTypes.ErrAssetCoinCannotBeIssue) // after IssueCoinsWaitBlockNums - 1 block, can issue
		simapp.CheckBalance(t, app, account4, amt)
	})
}

func TestLockCoins(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test lock core coins", t, func() {
		ctx := app.NewTestContext()

		lockedCoins := types.NewInt64Coins(constants.DefaultBondDenom, 1000000000)
		lockedBlockNum := app.LastBlockHeight() + 1 + 5

		msgLock := assetTypes.NewMsgLockCoin(addr4, account4, lockedCoins, lockedBlockNum)
		tx := simapp.NewTxForTest(
			account4,
			[]sdk.Msg{
				&msgLock,
			}, wallet.PrivKey(addr4))

		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		ctx = app.NewTestContext()
		allLockedCoins, locks, err := app.AssetKeeper().GetLockCoins(ctx, account4)
		So(err, ShouldBeNil)
		So(allLockedCoins, simapp.ShouldEq, lockedCoins)
		So(len(locks), ShouldEqual, 1)
		So(locks[0].UnlockBlockHeight, ShouldEqual, lockedBlockNum)
		So(locks[0].Coins, simapp.ShouldEq, lockedCoins)

		currCoins := app.AssetKeeper().GetAllBalances(ctx, account4)
		useableCoins := currCoins.Sub(lockedCoins).Sub(simapp.DefaultTestFee)

		// when locked, coins can be transfer in useful coins
		err = transfer(t, app, true, account4, account1, useableCoins, account4)
		So(err, ShouldBeNil)

		currCoinsBeforeFailedTransferCoins := app.AssetKeeper().GetAllBalances(app.NewTestContext(), account4)

		// when locked, coins cannot be transfer
		err = transfer(t, app, false, account4, account1, types.NewInt64Coins(constants.DefaultBondDenom, 100), account4)
		So(err, simapp.ShouldErrIs, assetTypes.ErrAssetCoinsLocked)

		currCoins = app.AssetKeeper().GetAllBalances(app.NewTestContext(), account4)
		So(currCoinsBeforeFailedTransferCoins, simapp.ShouldEq, currCoins)
	})

	Convey("test no core coins lock", t, func() {
		var (
			symbol             = types.MustName("abc")
			denom              = types.CoinDenom(name2, symbol) // create in last
			maxSupplyAmt int64 = 10000000000000
		)

		So(createCoin(t, app, true, account2, symbol, maxSupplyAmt), ShouldBeNil)

		issueAmt := types.NewInt64Coin(denom, 1000)
		So(issueCoin(t, app, true, account2, symbol, issueAmt), ShouldBeNil)

		var (
			ctx            = app.NewTestContext()
			lockedCoins    = types.NewInt64Coins(denom, 100)
			lockedBlockNum = app.LastBlockHeight() + 1 + 5
		)

		msgLock := assetTypes.NewMsgLockCoin(addr2, account2, lockedCoins, lockedBlockNum)
		tx := simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&msgLock,
			}, wallet.PrivKey(addr2))

		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		allLockedCoins, locks, err := app.AssetKeeper().GetLockCoins(app.NewTestContext(), account2)
		So(err, ShouldBeNil)
		So(allLockedCoins, simapp.ShouldEq, lockedCoins)
		So(len(locks), ShouldEqual, 1)
		So(locks[0].UnlockBlockHeight, ShouldEqual, lockedBlockNum)
		So(locks[0].Coins, simapp.ShouldEq, lockedCoins)

		// lock many coins
		var (
			lockedCoins2    = types.NewInt64Coins(constants.DefaultBondDenom, 100)
			lockedBlockNum2 = app.LastBlockHeight() + 100
		)

		msgLock2 := assetTypes.NewMsgLockCoin(
			addr2, account2, lockedCoins2, lockedBlockNum2)
		tx2 := simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&msgLock2,
			}, wallet.PrivKey(addr2))

		So(simapp.CheckTxs(t, app, app.NewTestContext(), tx2), ShouldBeNil)

		allLockedCoins, locks, err = app.AssetKeeper().GetLockCoins(app.NewTestContext(), account2)
		So(err, ShouldBeNil)
		So(allLockedCoins, simapp.ShouldEq, lockedCoins.Add(lockedCoins2...))
		So(len(locks), ShouldEqual, 2)
		So(locks[0].UnlockBlockHeight, ShouldEqual, lockedBlockNum)
		So(locks[0].Coins, simapp.ShouldEq, lockedCoins)
		So(locks[1].UnlockBlockHeight, ShouldEqual, lockedBlockNum2)
		So(locks[1].Coins, simapp.ShouldEq, lockedCoins2)
	})
}

func TestLockCoinsError(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test lock err by denom coins no exist", t, func() {
		var (
			symbol             = types.MustName("abc")
			denom              = types.CoinDenom(name2, symbol) // create in last
			maxSupplyAmt int64 = 10000000000000
		)

		So(createCoin(t, app, true, account2, symbol, maxSupplyAmt), ShouldBeNil)

		issueAmt := types.NewInt64Coin(denom, 1000)
		So(issueCoin(t, app, true, account2, symbol, issueAmt), ShouldBeNil)

		msgLock := assetTypes.NewMsgLockCoin(
			addr2, account2, types.NewInt64Coins("aaa/coin", 100), 10)
		tx := simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&msgLock,
			}, wallet.PrivKey(addr2)).WithCannotPass()

		So(simapp.CheckTxs(t, app, app.NewTestContext(), tx),
			simapp.ShouldErrIs, assetTypes.ErrAssetCoinNoExit)

	})

	Convey("test lock err by no coins", t, func() {
		var (
			symbol             = types.MustName("aaa")
			denom              = types.CoinDenom(name2, symbol) // create in last
			maxSupplyAmt int64 = 10000000000000
		)

		So(createCoin(t, app, true, account2, symbol, maxSupplyAmt), ShouldBeNil)
		// no issue, so account2 no have this coin

		msgLock := assetTypes.NewMsgLockCoin(
			addr2, account2, types.NewInt64Coins(denom, 100), 10)
		tx := simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&msgLock,
			}, wallet.PrivKey(addr2)).WithCannotPass()

		So(simapp.CheckTxs(t, app, app.NewTestContext(), tx),
			simapp.ShouldErrIs, assetTypes.ErrAssetLockCoinsNoEnough)
	})

	Convey("test lock err by block height err", t, func() {
		var (
			symbol             = types.MustName("aab")
			denom              = types.CoinDenom(name2, symbol) // create in last
			maxSupplyAmt int64 = 10000000000000
		)

		So(createCoin(t, app, true, account2, symbol, maxSupplyAmt), ShouldBeNil)
		// no issue, so account2 no have this coin
		issueAmt := types.NewInt64Coin(denom, 1000)
		So(issueCoin(t, app, true, account2, symbol, issueAmt), ShouldBeNil)

		simapp.AfterBlockCommitted(app, 2)

		msgLock := assetTypes.NewMsgLockCoin(
			addr2, account2, types.NewInt64Coins(denom, 100), 2) // 2 has already passed
		tx := simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&msgLock,
			}, wallet.PrivKey(addr2)).WithCannotPass()

		So(simapp.CheckTxs(t, app, app.NewTestContext(), tx),
			simapp.ShouldErrIs, assetTypes.ErrAssetLockUnlockBlockHeightErr)
	})

	Convey("test lock err by coins type cannot be locked", t, func() {
		var (
			symbol             = types.MustName("aac")
			denom              = types.CoinDenom(name2, symbol) // create in last
			maxSupplyAmt int64 = 10000000000000
			maxSupply          = types.NewCoin(denom, types.NewInt(maxSupplyAmt))
			initSupply         = types.NewCoin(denom, types.NewInt(0))
			desc               = []byte(fmt.Sprintf("desc for %s", denom))
		)

		So(createCoinExt(t, app, true, account2, symbol, maxSupply, true, false, true, 0, initSupply, desc), ShouldBeNil)
		// no issue, so account2 no have this coin
		lockAmt := types.NewInt64Coin(denom, 1000)
		So(issueCoin(t, app, true, account2, symbol, lockAmt), ShouldBeNil)

		msgLock := assetTypes.NewMsgLockCoin(
			addr2, account2, types.NewInt64Coins(denom, 100), 100) // 2 has already passed
		tx := simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&msgLock,
			}, wallet.PrivKey(addr2)).WithCannotPass()

		So(simapp.CheckTxs(t, app, app.NewTestContext(), tx),
			simapp.ShouldErrIs, assetTypes.ErrAssetCoinCannotBeLock)

		// for multiple coins test
		locks := types.NewInt64Coins(constants.DefaultBondDenom, 10).Add(lockAmt)
		msgLock2 := assetTypes.NewMsgLockCoin(addr2, account2, locks, 100)
		tx2 := simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&msgLock2,
			}, wallet.PrivKey(addr2)).WithCannotPass()

		So(simapp.CheckTxs(t, app, app.NewTestContext(), tx2),
			simapp.ShouldErrIs, assetTypes.ErrAssetCoinCannotBeLock)
	})
}
