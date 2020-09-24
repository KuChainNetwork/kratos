package keeper_test

import (
	"fmt"
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/KuChainNetwork/kuchain/x/account/keeper"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
	"github.com/KuChainNetwork/kuchain/x/distribution"
	"github.com/KuChainNetwork/kuchain/x/gov"
	"github.com/KuChainNetwork/kuchain/x/mint"
	"github.com/KuChainNetwork/kuchain/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types/time"
)

func checkAuthSequenceByQuery(app *simapp.SimApp, ctx sdk.Context, addr types.AccAddress, seq, num uint64) {
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAuthByAddress),
		Data: []byte{},
	}
	req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAddAuthParams(addr))
	path := []string{accountTypes.QueryAuthByAddress}

	res, err := keeper.NewQuerier(*app.AccountKeeper())(ctx, path, req)
	So(err, ShouldBeNil)

	var authData accountTypes.Auth
	err = app.Codec().UnmarshalJSON(res, &authData)
	So(err, ShouldBeNil)

	So(authData.GetAddress(), simapp.ShouldEq, addr)
	So(authData.GetSequence(), ShouldEqual, seq)
	So(authData.GetNumber(), ShouldEqual, num)
}

func TestEnsureAccount(t *testing.T) {
	app, ctx := createTestApp()

	Convey("TestEnsureAccount", t, func() {
		// no exist account
		err := app.AccountKeeper().EnsureAccount(ctx, account2)
		So(err, simapp.ShouldErrIs, accountTypes.ErrAccountNoFound)

		// no account exist
		So(app.AccountKeeper().GetAccount(ctx, account2), ShouldBeNil)

		// a new account with auth no inited
		newAddr1 := wallet.NewAccAddress()
		acc3 := app.AccountKeeper().NewAccountByName(ctx, types.MustName("account3test"))
		err = acc3.SetAuth(newAddr1)
		So(err, ShouldBeNil)

		app.AccountKeeper().SetAccount(ctx, acc3)

		err = app.AccountKeeper().EnsureAccount(ctx, acc3.GetID())
		So(err, ShouldBeNil)

		// address account ensure
		newAddr2 := wallet.NewAccAddress()
		err = app.AccountKeeper().EnsureAccount(ctx, types.NewAccountIDFromAccAdd(newAddr2))
		So(err, ShouldBeNil)
		checkAuthSequenceByQuery(app, ctx, newAddr2, 1, 2)
	})
}

func TestIterateAccount(t *testing.T) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)
	genAcc1 := simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1)
	genAcc2 := simapp.NewSimGenesisAccount(account2, addr1).WithAsset(asset1)
	genAcc3 := simapp.NewSimGenesisAccount(types.NewAccountIDFromAccAdd(addr1), addr1).WithAsset(asset1)
	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAcc1, genAcc2, genAcc3)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	ctx := app.BaseApp.NewContext(true,
		abci.Header{
			Time:   time.Now(),
			Height: app.LastBlockHeight() + 1,
		})

	Convey("TestIterateAccount", t, func() {
		// Create accounts for test
		acc3 := app.AccountKeeper().NewAccountByName(ctx, types.MustName("account3test"))
		err := acc3.SetAuth(wallet.NewAccAddress())
		So(err, ShouldBeNil)
		app.AccountKeeper().SetAccount(ctx, acc3)

		names := make(map[string]types.AccountID)
		app.AccountKeeper().IterateAccounts(ctx, func(account exported.Account) bool {
			ctx.Logger().Info("iterate account", "id", account.GetID())
			names[account.GetID().String()] = account.GetID()
			return false
		})

		So(len(names), ShouldEqual, (1 + 7 + 4)) // kuchain, 7 module account, and 4 genesis account
		ids := []string{constants.SystemAccountID.String(),
			"mint", "kugov", "kustaking", "kubondedpool", "kudistribution", "kunotbondedpool",
			account1.String(), account2.String(), addr1.String(), acc3.GetID().String()}

		for _, id := range ids {
			_, ok := names[id]
			So(ok, ShouldBeTrue)
		}
	})

	Convey("TestIterateAccountBreak", t, func() {
		names := make(map[string]types.AccountID)
		app.AccountKeeper().IterateAccounts(ctx, func(account exported.Account) bool {
			ctx.Logger().Info("iterate account", "id", account.GetID())
			names[account.GetID().String()] = account.GetID()
			return true // will return break
		})

		So(len(names), ShouldEqual, 1)
	})
}

func TestAccountExist(t *testing.T) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)
	genAcc1 := simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1)
	genAcc2 := simapp.NewSimGenesisAccount(account2, addr1).WithAsset(asset1)
	genAcc3 := simapp.NewSimGenesisAccount(types.NewAccountIDFromAccAdd(addr1), addr1).WithAsset(asset1)
	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAcc1, genAcc2, genAcc3)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	ctx := app.BaseApp.NewContext(true,
		abci.Header{
			Time:   time.Now(),
			Height: app.LastBlockHeight() + 1,
		})

	Convey("TestAccountExist", t, func() {
		// Create accounts for test
		acc3 := app.AccountKeeper().NewAccountByName(ctx, types.MustName("account3test"))
		err := acc3.SetAuth(wallet.NewAccAddress())
		So(err, ShouldBeNil)
		app.AccountKeeper().SetAccount(ctx, acc3)

		// genesis and account new
		So(app.AccountKeeper().IsAccountExist(ctx, account1), ShouldBeTrue)
		So(app.AccountKeeper().IsAccountExist(ctx, account2), ShouldBeTrue)
		So(app.AccountKeeper().IsAccountExist(ctx, types.NewAccountIDFromAccAdd(addr1)), ShouldBeTrue)
		So(app.AccountKeeper().IsAccountExist(ctx, acc3.GetID()), ShouldBeTrue)

		// root account
		So(app.AccountKeeper().IsAccountExist(ctx, constants.SystemAccountID), ShouldBeTrue)

		// module accounts
		So(app.AccountKeeper().IsAccountExist(ctx, types.MustAccountID(distribution.ModuleName)), ShouldBeTrue)
		So(app.AccountKeeper().IsAccountExist(ctx, types.MustAccountID(mint.ModuleName)), ShouldBeTrue)
		So(app.AccountKeeper().IsAccountExist(ctx, types.MustAccountID(staking.ModuleName)), ShouldBeTrue)
		So(app.AccountKeeper().IsAccountExist(ctx, types.MustAccountID(gov.ModuleName)), ShouldBeTrue)

		// no exit account
		So(app.AccountKeeper().IsAccountExist(ctx, types.MustAccountID("abcde")), ShouldBeFalse)

		// a auth no account
		So(app.AccountKeeper().IsAccountExist(ctx, types.NewAccountIDFromAccAdd(addr2)), ShouldBeFalse)
	})
}
