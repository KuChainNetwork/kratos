package keeper_test

import (
	"testing"
	"time"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestLockedForever(t *testing.T) {
	app, _ := createTestApp()

	Convey("test locked forever", t, func() {
		ctx := app.NewTestContext()

		err := app.AssetKeeper().LockCoins(ctx, account1, -1, types.NewInt64CoreCoins(1000))
		So(err, ShouldBeNil)

		all, locks, err := app.AssetKeeper().GetLockCoins(ctx, account1)
		So(all, simapp.ShouldEq, types.NewInt64CoreCoins(1000))
		So(len(locks), ShouldEqual, 1)
		So(locks[0].UnlockBlockHeight < 0, ShouldBeTrue)

		err = app.AssetKeeper().LockCoins(ctx, account1, 100, types.NewInt64CoreCoins(2000))
		So(err, ShouldBeNil)

		all, locks, err = app.AssetKeeper().GetLockCoins(ctx, account1)
		So(all, simapp.ShouldEq, types.NewInt64CoreCoins(3000))
		So(len(locks), ShouldEqual, 2)
		So(locks[0].UnlockBlockHeight < 0, ShouldBeTrue)

		ctx = app.BaseApp.NewContext(true,
			abci.Header{
				Height: 101,
				Time:   time.Now(),
			})

		err = app.AssetKeeper().UnLockCoins(ctx, account1, types.NewInt64CoreCoins(2000))
		So(err, ShouldBeNil)

		err = app.AssetKeeper().UnLockCoins(ctx, account1, types.NewInt64CoreCoins(1000))
		So(err, simapp.ShouldErrIs, assetTypes.ErrAssetUnLockCoins)

		all, locks, err = app.AssetKeeper().GetLockCoins(ctx, account1)
		So(all, simapp.ShouldEq, types.NewInt64CoreCoins(1000))
		So(len(locks), ShouldEqual, 1)
		So(locks[0].UnlockBlockHeight < 0, ShouldBeTrue)

		ctx = app.BaseApp.NewContext(true,
			abci.Header{
				Height: 9223372036854775800,
				Time:   time.Now(),
			})

		err = app.AssetKeeper().UnLockCoins(ctx, account1, types.NewInt64CoreCoins(1000))
		So(err, simapp.ShouldErrIs, assetTypes.ErrAssetUnLockCoins)

		all, locks, err = app.AssetKeeper().GetLockCoins(ctx, account1)
		So(all, simapp.ShouldEq, types.NewInt64CoreCoins(1000))
		So(len(locks), ShouldEqual, 1)
		So(locks[0].UnlockBlockHeight < 0, ShouldBeTrue)
	})
}

func TestLockedForeverUnlock(t *testing.T) {
	app, _ := createTestAppWithCoins()

	Convey("test unlocked forever", t, func() {
		ctx := app.NewTestContext()

		Convey("test first one locked forever", func() {
			So(app.AssetKeeper().LockCoins(ctx, account1, -1, types.NewInt64CoreCoins(1000)), ShouldBeNil)
			So(app.AssetKeeper().LockCoins(ctx, account1, -1, types.NewInt64CoreCoins(1000)), ShouldBeNil)

			all, locks, err := app.AssetKeeper().GetLockCoins(ctx, account1)
			So(err, ShouldBeNil)
			So(all, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(len(locks), ShouldEqual, 1)
			So(locks[0].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[0].UnlockBlockHeight, ShouldEqual, -1)

			So(app.AssetKeeper().UnLockFreezedCoins(ctx, account1, types.NewInt64CoreCoins(1000)), ShouldBeNil)

			all, locks, err = app.AssetKeeper().GetLockCoins(ctx, account1)
			So(err, ShouldBeNil)
			So(all, simapp.ShouldEq, types.NewInt64CoreCoins(1000))
			So(len(locks), ShouldEqual, 1)
			So(locks[0].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(1000))
			So(locks[0].UnlockBlockHeight, ShouldEqual, -1)

			So(app.AssetKeeper().UnLockFreezedCoins(ctx, account1, types.NewInt64CoreCoins(1000)), ShouldBeNil)
			all, locks, err = app.AssetKeeper().GetLockCoins(ctx, account1)
			So(err, ShouldBeNil)
			So(all.IsZero(), ShouldBeTrue)
			So(len(locks), ShouldEqual, 0)
		})

		Convey("test last one locked forever", func() {
			So(app.AssetKeeper().LockCoins(ctx, account2, 10, types.NewInt64CoreCoins(2000)), ShouldBeNil)
			So(app.AssetKeeper().LockCoins(ctx, account2, -1, types.NewInt64CoreCoins(1000)), ShouldBeNil)
			So(app.AssetKeeper().LockCoins(ctx, account2, -1, types.NewInt64CoreCoins(1000)), ShouldBeNil)

			all, locks, err := app.AssetKeeper().GetLockCoins(ctx, account2)
			So(err, ShouldBeNil)
			So(all, simapp.ShouldEq, types.NewInt64CoreCoins(4000))
			So(len(locks), ShouldEqual, 2)
			So(locks[0].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[0].UnlockBlockHeight, ShouldEqual, 10)
			So(locks[1].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[1].UnlockBlockHeight, ShouldEqual, -1)

			So(app.AssetKeeper().UnLockFreezedCoins(ctx, account2, types.NewInt64CoreCoins(1000)), ShouldBeNil)

			all, locks, err = app.AssetKeeper().GetLockCoins(ctx, account2)
			So(err, ShouldBeNil)
			So(all, simapp.ShouldEq, types.NewInt64CoreCoins(3000))
			So(len(locks), ShouldEqual, 2)
			So(locks[0].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[0].UnlockBlockHeight, ShouldEqual, 10)
			So(locks[1].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(1000))
			So(locks[1].UnlockBlockHeight, ShouldEqual, -1)

			So(app.AssetKeeper().UnLockFreezedCoins(ctx, account2, types.NewInt64CoreCoins(1000)), ShouldBeNil)

			all, locks, err = app.AssetKeeper().GetLockCoins(ctx, account2)
			So(err, ShouldBeNil)
			So(all, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(len(locks), ShouldEqual, 1)
			So(locks[0].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[0].UnlockBlockHeight, ShouldEqual, 10)
		})

		Convey("test mid one locked forever", func() {
			So(app.AssetKeeper().LockCoins(ctx, account3, 10, types.NewInt64CoreCoins(2000)), ShouldBeNil)
			So(app.AssetKeeper().LockCoins(ctx, account3, -1, types.NewInt64CoreCoins(1000)), ShouldBeNil)
			So(app.AssetKeeper().LockCoins(ctx, account3, -1, types.NewInt64CoreCoins(1000)), ShouldBeNil)
			So(app.AssetKeeper().LockCoins(ctx, account3, 11, types.NewInt64CoreCoins(2000)), ShouldBeNil)

			all, locks, err := app.AssetKeeper().GetLockCoins(ctx, account3)
			So(err, ShouldBeNil)
			So(all, simapp.ShouldEq, types.NewInt64CoreCoins(6000))
			So(len(locks), ShouldEqual, 3)
			So(locks[0].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[0].UnlockBlockHeight, ShouldEqual, 10)
			So(locks[1].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[1].UnlockBlockHeight, ShouldEqual, -1)
			So(locks[2].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[2].UnlockBlockHeight, ShouldEqual, 11)

			So(app.AssetKeeper().UnLockFreezedCoins(ctx, account3, types.NewInt64CoreCoins(1000)), ShouldBeNil)

			all, locks, err = app.AssetKeeper().GetLockCoins(ctx, account3)
			So(err, ShouldBeNil)
			So(all, simapp.ShouldEq, types.NewInt64CoreCoins(5000))
			So(len(locks), ShouldEqual, 3)
			So(locks[0].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[0].UnlockBlockHeight, ShouldEqual, 10)
			So(locks[1].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(1000))
			So(locks[1].UnlockBlockHeight, ShouldEqual, -1)
			So(locks[2].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[2].UnlockBlockHeight, ShouldEqual, 11)

			So(app.AssetKeeper().UnLockFreezedCoins(ctx, account3, types.NewInt64CoreCoins(1000)), ShouldBeNil)

			all, locks, err = app.AssetKeeper().GetLockCoins(ctx, account3)
			So(err, ShouldBeNil)
			So(all, simapp.ShouldEq, types.NewInt64CoreCoins(4000))
			So(len(locks), ShouldEqual, 2)
			So(locks[0].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[0].UnlockBlockHeight, ShouldEqual, 10)
			So(locks[1].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[1].UnlockBlockHeight, ShouldEqual, 11)
		})
	})
}

func TestLockedForeverUnlockErr(t *testing.T) {
	app, _ := createTestAppWithCoins()

	Convey("test unlocked forever", t, func() {
		ctx := app.NewTestContext()

		Convey("test one locked forever no enough", func() {
			So(app.AssetKeeper().LockCoins(ctx, account1, -1, types.NewInt64CoreCoins(1000)), ShouldBeNil)
			So(app.AssetKeeper().LockCoins(ctx, account1, -1, types.NewInt64CoreCoins(1000)), ShouldBeNil)

			all, locks, err := app.AssetKeeper().GetLockCoins(ctx, account1)
			So(err, ShouldBeNil)
			So(all, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(len(locks), ShouldEqual, 1)
			So(locks[0].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(2000))
			So(locks[0].UnlockBlockHeight, ShouldEqual, -1)

			So(app.AssetKeeper().UnLockFreezedCoins(ctx, account1, types.NewInt64CoreCoins(1000)), ShouldBeNil)

			all, locks, err = app.AssetKeeper().GetLockCoins(ctx, account1)
			So(err, ShouldBeNil)
			So(all, simapp.ShouldEq, types.NewInt64CoreCoins(1000))
			So(len(locks), ShouldEqual, 1)
			So(locks[0].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(1000))
			So(locks[0].UnlockBlockHeight, ShouldEqual, -1)

			So(app.AssetKeeper().UnLockFreezedCoins(ctx, account1, types.NewInt64CoreCoins(1001)),
				simapp.ShouldErrIs, assetTypes.ErrAssetUnLockCoins)

			all, locks, err = app.AssetKeeper().GetLockCoins(ctx, account1)
			So(err, ShouldBeNil)
			So(all, simapp.ShouldEq, types.NewInt64CoreCoins(1000))
			So(len(locks), ShouldEqual, 1)
			So(locks[0].Coins, simapp.ShouldEq, types.NewInt64CoreCoins(1000))
			So(locks[0].UnlockBlockHeight, ShouldEqual, -1)
		})
	})
}
