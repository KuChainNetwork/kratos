package keeper_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types/time"
)

var (
	// values for test
	wallet   = simapp.NewWallet()
	addr1    = wallet.NewAccAddressByName(name1)
	addr2    = wallet.NewAccAddressByName(name2)
	addr3    = wallet.NewAccAddressByName(name3)
	name1    = types.MustName("test01@chain")
	name2    = types.MustName("aaaeeebbbccc")
	name3    = types.MustName("aaaeeebbbcc2")
	account1 = types.NewAccountIDFromName(name1)
	account2 = types.NewAccountIDFromName(name2)
	account3 = types.NewAccountIDFromName(name3)
)

func createTestApp() (*simapp.SimApp, sdk.Context) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)
	genAcc := simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1)
	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAcc)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	ctxCheck := app.BaseApp.NewContext(true,
		abci.Header{
			Time:   time.Now(),
			Height: app.LastBlockHeight() + 1,
		})
	return app, ctxCheck
}

func TestAccountMapperGetSet(t *testing.T) {
	app, ctx := createTestApp()

	Convey("TestAccountMapperGetSet", t, func() {
		// no account before its created
		acc := app.AccountKeeper().GetAccount(ctx, account2)
		So(acc, ShouldBeNil)

		// create account and check default values
		acc2 := app.AccountKeeper().NewAccountByName(ctx, name2)
		So(acc2, ShouldNotBeNil)
		So(acc2.GetName(), simapp.ShouldEq, name2)
		So(acc2.GetID(), simapp.ShouldEq, account2)

		err := acc2.SetAuth(addr2)
		So(err, ShouldBeNil)
		So(acc2.GetAuth(), simapp.ShouldEq, addr2)

		// NewAccount doesn't call Set, so it's still nil
		So(app.AccountKeeper().GetAccount(ctx, account2), ShouldBeNil)
		So(app.AccountKeeper().GetAccountByName(ctx, name2), ShouldBeNil)

		// set some values on the account and save it
		newSequence := uint64(20)
		err = acc2.SetAccountNumber(newSequence)
		So(err, ShouldBeNil)

		app.AccountKeeper().SetAccount(ctx, acc2)

		// check the new values
		acc2 = app.AccountKeeper().GetAccount(ctx, account2)
		So(acc2, ShouldNotBeNil)
		So(newSequence, ShouldEqual, acc2.GetAccountNumber())

		// check get name
		acc2 = app.AccountKeeper().GetAccountByName(ctx, name2)
		So(acc2, ShouldNotBeNil)
		So(acc2.GetName(), simapp.ShouldEq, name2)
	})
}

func TestGetAuthByName(t *testing.T) {
	app, ctx := createTestApp()

	Convey("TestGetAuthByName", t, func() {
		auth, err := app.AccountKeeper().GetAuth(ctx, name1)
		So(err, ShouldBeNil)
		So(auth, simapp.ShouldEq, addr1)

		_, err = app.AccountKeeper().GetAuth(ctx, types.MustName("acccc"))
		So(err, simapp.ShouldErrIs, accountTypes.ErrAccountNoFound)
	})
}
