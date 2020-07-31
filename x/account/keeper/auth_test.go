package keeper_test

import (
	"fmt"
	"testing"

	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/account/keeper"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAuthGetAccount(t *testing.T) {
	app, ctx := createTestApp()

	Convey("TestAuthGetAccount", t, func() {
		// get
		res := app.AccountKeeper().GetAccountsByAuth(ctx, addr1)
		So(len(res), ShouldEqual, 1)
		So(res[0], ShouldEqual, account1.String())

		// no found
		res = app.AccountKeeper().GetAccountsByAuth(ctx, wallet.NewAccAddress())
		So(len(res), ShouldEqual, 0)
	})

	Convey("TestAuth2AccountAddAndDel", t, func() {
		// add
		app.AccountKeeper().AddAccountByAuth(ctx, addr2, "testacc1")
		app.AccountKeeper().AddAccountByAuth(ctx, addr2, "testacc2")

		res := app.AccountKeeper().GetAccountsByAuth(ctx, addr2)
		So(len(res), ShouldEqual, 2)

		// del
		app.AccountKeeper().DeleteAccountByAuth(ctx, addr2, "testacc1")
		app.AccountKeeper().DeleteAccountByAuth(ctx, addr2, "testacc2")

		res = app.AccountKeeper().GetAccountsByAuth(ctx, addr2)
		So(len(res), ShouldEqual, 0)

		// no exist account to auth
		app.AccountKeeper().DeleteAccountByAuth(ctx, addr2, "testacc3")

		// no exist auth to auth
		app.AccountKeeper().DeleteAccountByAuth(ctx, wallet.NewAccAddress(), "testacc2")
	})
}

func TestAuthSeq(t *testing.T) {
	app, ctx := createTestApp()

	Convey("TestAuthSeq", t, func() {
		ctxCheck0 := app.BaseApp.NewContext(true,
			abci.Header{
				Height: 0,
			})
		// in genesis seq will 0
		seq, num, err := app.AccountKeeper().GetAuthSequence(ctxCheck0, addr1)
		So(seq, ShouldEqual, 0)
		So(num, ShouldEqual, 0)
		So(err, ShouldBeNil)

		// in curr will be init
		seq, _, err = app.AccountKeeper().GetAuthSequence(ctx, addr1)
		So(seq, ShouldEqual, 1)
		So(err, ShouldBeNil)

		// inc seq
		seq, _, err = app.AccountKeeper().GetAuthSequence(ctx, addr2)
		So(seq, ShouldEqual, 1)
		So(err, ShouldBeNil)

		app.AccountKeeper().IncAuthSequence(ctx, addr1)

		seq, _, err = app.AccountKeeper().GetAuthSequence(ctx, addr1)
		So(seq, ShouldEqual, 2) // addr1 will add by 1
		So(err, ShouldBeNil)

		seq, _, err = app.AccountKeeper().GetAuthSequence(ctx, addr2)
		So(seq, ShouldEqual, 1) // addr2 will not change
		So(err, ShouldBeNil)
	})
}

func TestAuthPublicKey(t *testing.T) {
	app, ctx := createTestApp()

	Convey("TestAuthPublicKey", t, func() {
		// set a exit key
		app.AccountKeeper().SetPubKey(ctx, addr1, wallet.PrivKey(addr1).PubKey())

		acc1 := app.AccountKeeper().GetAccount(ctx, account1)
		So(acc1, ShouldNotBeNil)

		req := abci.RequestQuery{
			Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAuthByAddress),
			Data: []byte{},
		}
		req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAddAuthParams(addr1))
		path := []string{accountTypes.QueryAuthByAddress}

		res, err := keeper.NewQuerier(*app.AccountKeeper())(ctx, path, req)
		So(err, ShouldBeNil)

		var authData accountTypes.Auth
		err = app.Codec().UnmarshalJSON(res, &authData)
		So(err, ShouldBeNil)

		So(authData.GetPubKey(), simapp.ShouldEq, wallet.PrivKey(addr1).PubKey())
	})
}
