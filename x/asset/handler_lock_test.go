package asset_test

import (
	"fmt"
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

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

func TestUnLockCoins(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test lock core coins", t, func() {
		ctx := app.NewTestContext()

		lockedCoins := types.NewInt64Coins(constants.DefaultBondDenom, 100000000)
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

		// Wait 5 block 4 + 1, 1 for transfer
		simapp.AfterBlockCommitted(app, 4)

		ctx = app.NewTestContext()

		// check is locked before unlock
		currCoins := app.AssetKeeper().GetAllBalances(ctx, account4)
		useableCoins := currCoins.Sub(lockedCoins).Sub(simapp.DefaultTestFee)

		// when locked, coins cannot be transfer
		err = transfer(t, app, false, account4, account1, useableCoins.Add(types.NewInt64Coin(constants.DefaultBondDenom, 100)), account4)
		So(err, simapp.ShouldErrIs, assetTypes.ErrAssetCoinsLocked)

		ctx = app.NewTestContext()

		// unlock coins
		msgUnLock := assetTypes.NewMsgUnlockCoin(addr4, account4, lockedCoins)
		tx = simapp.NewTxForTest(
			account4,
			[]sdk.Msg{
				&msgUnLock,
			}, wallet.PrivKey(addr4))

		// unlock should success
		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		ctx = app.NewTestContext()
		allLockedCoins, locks, err = app.AssetKeeper().GetLockCoins(ctx, account4)
		So(err, ShouldBeNil)
		So(allLockedCoins.IsZero(), ShouldBeTrue)
		So(len(locks), ShouldEqual, 0)

		// now can transfer

	})

	Convey("test no core coins unlock", t, func() {
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
			lockedBlockNum = app.LastBlockHeight() + 1 + 6
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

		simapp.AfterBlockCommitted(app, 5)
		ctx = app.NewTestContext()

		// lock many coins
		var (
			lockedCoins2    = types.NewInt64Coins(constants.DefaultBondDenom, 100)
			lockedBlockNum2 = app.LastBlockHeight() + 1 + 5
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

		// unlock
		ctx = app.NewTestContext()
		// current 17, unlock1 17, unlock2, 22
		ctx.Logger().Info("unlock block num", "curr", app.LastBlockHeight(), "locked", lockedBlockNum, "locked2", lockedBlockNum2)

		// unlock coins
		msgUnLock := assetTypes.NewMsgUnlockCoin(addr2, account2, lockedCoins)
		tx = simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&msgUnLock,
			}, wallet.PrivKey(addr2))

		// unlock should success
		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		allLockedCoins, locks, err = app.AssetKeeper().GetLockCoins(app.NewTestContext(), account2)
		So(err, ShouldBeNil)
		So(allLockedCoins, simapp.ShouldEq, lockedCoins2)
		So(len(locks), ShouldEqual, 1)
		So(locks[0].UnlockBlockHeight, ShouldEqual, lockedBlockNum2)
		So(locks[0].Coins, simapp.ShouldEq, lockedCoins2)

		simapp.AfterBlockCommitted(app, 5)
		ctx = app.NewTestContext()

		// unlock coins
		msgUnLock = assetTypes.NewMsgUnlockCoin(addr2, account2, lockedCoins2)
		tx = simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&msgUnLock,
			}, wallet.PrivKey(addr2))

		// unlock should success
		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		allLockedCoins, locks, err = app.AssetKeeper().GetLockCoins(app.NewTestContext(), account2)
		So(err, ShouldBeNil)
		So(allLockedCoins.IsZero(), ShouldBeTrue)
		So(len(locks), ShouldEqual, 0)
	})
}

func TestUnLockCoinsMultiple(t *testing.T) {
	app, _ := createAppForTest()
	Convey("test no core coins unlock", t, func() {
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
			lockedBlockNum = app.LastBlockHeight() + 1 + 6
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

		simapp.AfterBlockCommitted(app, 5)
		ctx = app.NewTestContext()

		// lock many coins
		var (
			lockedCoins2    = types.NewInt64Coins(constants.DefaultBondDenom, 100)
			lockedBlockNum2 = app.LastBlockHeight() + 1 + 5
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

		simapp.AfterBlockCommitted(app, 10)
		ctx = app.NewTestContext()

		// unlock coins
		msgUnLock := assetTypes.NewMsgUnlockCoin(addr2, account2, lockedCoins.Add(lockedCoins2...))
		tx = simapp.NewTxForTest(
			account2,
			[]sdk.Msg{
				&msgUnLock,
			}, wallet.PrivKey(addr2))

		// unlock should success
		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		allLockedCoins, locks, err = app.AssetKeeper().GetLockCoins(app.NewTestContext(), account2)
		So(err, ShouldBeNil)
		So(allLockedCoins.IsZero(), ShouldBeTrue)
		So(len(locks), ShouldEqual, 0)
	})
}

func TestUnLockCoinsMultipleInSameHeight(t *testing.T) {
	app, _ := createAppForTest()
	Convey("test lock in same height", t, func() {
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
			lockedBlockNum = app.LastBlockHeight() + 1 + 6
		)

		Convey("lock two times in same height", func() {
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

			simapp.AfterBlockCommitted(app, 1)

			ctx = app.NewTestContext()
			msgLock = assetTypes.NewMsgLockCoin(addr2, account2, lockedCoins, lockedBlockNum)
			tx = simapp.NewTxForTest(
				account2,
				[]sdk.Msg{
					&msgLock,
				}, wallet.PrivKey(addr2))

			So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

			lockedCoins = lockedCoins.Add(lockedCoins...)

			allLockedCoins, locks, err = app.AssetKeeper().GetLockCoins(app.NewTestContext(), account2)
			So(err, ShouldBeNil)
			So(allLockedCoins, simapp.ShouldEq, lockedCoins)
			So(len(locks), ShouldEqual, 1)
			So(locks[0].UnlockBlockHeight, ShouldEqual, lockedBlockNum)
			So(locks[0].Coins, simapp.ShouldEq, lockedCoins)

			// lock many coins
			var (
				lockedCoins2    = types.NewInt64Coins(constants.DefaultBondDenom, 100)
				lockedBlockNum2 = app.LastBlockHeight() + 1 + 5
			)

			ctx = app.NewTestContext()

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

			simapp.AfterBlockCommitted(app, 10)
			ctx = app.NewTestContext()

			// unlock coins
			msgUnLock := assetTypes.NewMsgUnlockCoin(addr2, account2, lockedCoins.Add(lockedCoins2...))
			tx = simapp.NewTxForTest(
				account2,
				[]sdk.Msg{
					&msgUnLock,
				}, wallet.PrivKey(addr2))

			// unlock should success
			So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

			allLockedCoins, locks, err = app.AssetKeeper().GetLockCoins(app.NewTestContext(), account2)
			So(err, ShouldBeNil)
			So(allLockedCoins.IsZero(), ShouldBeTrue)
			So(len(locks), ShouldEqual, 0)
		})
	})
}
