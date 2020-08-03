package account_test

import (
	"errors"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/account/keeper"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
)

var (
	// values for test
	wallet   = simapp.NewWallet()
	addr1    = wallet.NewAccAddressByName(name1)
	addr2    = wallet.NewAccAddressByName(name2)
	addr3    = wallet.NewAccAddress()
	name1    = types.MustName("test01@chain")
	name2    = types.MustName("aaaeeebbbccc")
	account1 = types.NewAccountIDFromName(name1)
	account2 = types.NewAccountIDFromName(name2)
)

func TestCreateNormalAccount(t *testing.T) {
	Convey("TestCreateNormalAccount", t, func() {
		genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth())
		app := simapp.SetupWithGenesisAccounts(genAccs)

		ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

		Convey("check system account", func() {
			sysAccount := app.AccountKeeper().GetAccount(ctxCheck, constants.SystemAccountID)
			So(sysAccount, ShouldNotBeNil)
			So(sysAccount.GetID().Eq(constants.SystemAccountID), ShouldBeTrue)
		})

		origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, wallet.GetRootAuth())
		So(err, ShouldBeNil)

		Convey("create a account by system account", func() {
			msg := accountTypes.NewMsgCreateAccount(
				wallet.GetRootAuth(),
				constants.SystemAccountID,
				name1,
				addr1)
			fee := types.NewInt64Coins(constants.DefaultBondDenom, 100000)

			header := abci.Header{Height: app.LastBlockHeight() + 1}
			_, _, err := simapp.SignCheckDeliver(
				t, app.Codec(), app.BaseApp,
				header, constants.SystemAccountID, fee,
				[]sdk.Msg{&msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
				true, true, wallet.GetRootKey())

			So(err, ShouldBeNil)

			// Check is account is created
			ctxCheck = app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
			accCreated := app.AccountKeeper().GetAccount(ctxCheck, account1)

			So(accCreated, ShouldNotBeNil)
			So(accCreated.GetID().Eq(account1), ShouldBeTrue)
			So(accCreated.GetAuth().Equals(addr1), ShouldBeTrue)
			So(accCreated.GetName().Eq(name1), ShouldBeTrue)
			So(accCreated.GetAccountNumber() == 0, ShouldBeTrue) // we no use this now

			// Check new account auth
			addr1Seq, addr1Num, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addr1)
			So(err, ShouldBeNil)
			So(addr1Num, ShouldEqual, 1) // the 2th account, root account is 0, and this will be 1
			So(addr1Seq, ShouldEqual, 1) // seq will init to 1

			// Check root account auth
			rootCurrSeq, rootCurrNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, wallet.GetRootAuth())
			So(err, ShouldBeNil)
			So(rootCurrNum, ShouldEqual, 0) // root account is 0
			So(rootCurrSeq, ShouldEqual, 2) // seq will init to 1, add is 2
		})

	})
}

func testAccountCreate(
	t *testing.T, app *simapp.SimApp, wallet *simapp.Wallet, shouldBeSuccess bool,
	creator types.AccountID, name types.Name, auth types.AccAddress) error {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
	creatorAuth := app.AccountKeeper().GetAccount(ctxCheck, creator).GetAuth()
	fee := types.NewInt64Coins(constants.DefaultBondDenom, 100000)
	msg := accountTypes.NewMsgCreateAccount(
		creatorAuth,
		creator,
		name,
		auth)

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, creatorAuth)
	So(err, ShouldBeNil)

	header := abci.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(
		t, app.Codec(), app.BaseApp,
		header, creator, fee,
		[]sdk.Msg{&msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		shouldBeSuccess, shouldBeSuccess, wallet.PrivKey(creatorAuth))

	return err
}

func TestCreateAccountNameCheck(t *testing.T) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)
	genAcc := simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1)
	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAcc)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	Convey("normal account name length check", t, func() {
		// name length should be 12
		err1 := testAccountCreate(t, app, wallet, false, account1, types.MustName("aaa"), addr1)
		So(errors.Is(err1, accountTypes.ErrAccountNameLenInvalid), ShouldBeTrue)

		err2 := testAccountCreate(t, app, wallet, false, account1, types.MustName("aaabbbcccddde"), addr1)
		So(errors.Is(err2, accountTypes.ErrAccountNameLenInvalid), ShouldBeTrue)

		err3 := testAccountCreate(t, app, wallet, true, account1, types.MustName("aaabbbcccddd"), addr1)
		So(err3, ShouldBeNil)
	})

	Convey("normal system account check", t, func() {
		// name cannot be system account
		err1 := testAccountCreate(t, app, wallet, false, account1, constants.FeeSystemAccount, addr1)
		So(errors.Is(err1, accountTypes.ErrAccountCannotCreateSysAccount), ShouldBeTrue)

		err2 := testAccountCreate(t, app, wallet, false, account1, types.MustName(constants.GetSystemAccount("abc")), addr1)
		So(errors.Is(err2, accountTypes.ErrAccountCannotCreateSysAccount), ShouldBeTrue)
	})

	Convey("account name string error", t, func() {
		err1 := testAccountCreate(t, app, wallet, false, account1, types.MustName("@aabbbcccddd"), addr1)
		So(err1, simapp.ShouldErrIs, accountTypes.ErrAccountNameInvalid)

		err2 := testAccountCreate(t, app, wallet, false, account1, types.MustName("aaabbbcccdd@"), addr1)
		So(err2, simapp.ShouldErrIs, accountTypes.ErrAccountNameInvalid)

		err3 := testAccountCreate(t, app, wallet, false, account1, types.MustName("aaa@bbcc@ddd"), addr1)
		So(err3, simapp.ShouldErrIs, accountTypes.ErrAccountNameInvalid)
	})
}

func TestCreateAccountDoubleTimes(t *testing.T) {
	assets := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)
	genAccs := simapp.NewGenesisAccounts(
		wallet.GetRootAuth(),
		simapp.NewSimGenesisAccount(account1, addr1).WithAsset(assets),
		simapp.NewSimGenesisAccount(account2, addr2).WithAsset(assets))
	app := simapp.SetupWithGenesisAccounts(genAccs)

	Convey("genesis account create double times", t, func() {
		err2 := testAccountCreate(t, app, wallet, false, account1, name2, addr3)
		So(err2, simapp.ShouldErrIs, accountTypes.ErrAccountHasCreated)
	})

	Convey("account double times", t, func() {
		name := types.MustName("aaabbbcccddd")
		err1 := testAccountCreate(t, app, wallet, true, account1, name, addr1)
		So(err1 == nil, ShouldBeTrue)

		err2 := testAccountCreate(t, app, wallet, false, account1, name, addr1)
		So(errors.Is(err2, accountTypes.ErrAccountHasCreated), ShouldBeTrue)
	})
}

func TestUpdateAccountAuth(t *testing.T) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)
	genAcc := simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1)
	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAcc)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	Convey("update a account auth", t, func() {
		ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addr1)
		So(err, ShouldBeNil)

		newAuth := wallet.NewAccAddress()

		msg := accountTypes.NewMsgUpdateAccountAuth(
			addr1,
			name1,
			newAuth)
		fee := types.NewInt64Coins(constants.DefaultBondDenom, 100000)

		header := abci.Header{Height: app.LastBlockHeight() + 1}
		_, _, err = simapp.SignCheckDeliver(
			t, app.Codec(), app.BaseApp,
			header, account1, fee,
			[]sdk.Msg{&msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
			true, true, wallet.PrivKey(addr1))

		So(err, ShouldBeNil)

		// Check is account is created
		ctxCheck = app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
		currAccount := app.AccountKeeper().GetAccount(ctxCheck, account1)

		So(currAccount, ShouldNotBeNil)
		So(currAccount.GetID().Eq(account1), ShouldBeTrue)
		So(currAccount.GetAuth().Equals(newAuth), ShouldBeTrue)
		So(currAccount.GetName().Eq(name1), ShouldBeTrue)

		// Check account old auth
		addr1Seq, addr1Num, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addr1)
		So(err, ShouldBeNil)
		So(addr1Num, ShouldEqual, origAuthNum)
		So(addr1Seq, ShouldEqual, origAuthSeq+1)

		// Check account new auth
		currSeq, currNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, newAuth)
		So(err, ShouldBeNil)
		So(currSeq, ShouldEqual, 1) // this will be init 1
		So(currNum, ShouldEqual, 2) // it is 3rd auth in chain
	})
}

func TestUpdateAccountAuthErrReturn(t *testing.T) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)
	genAcc := simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1)
	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAcc)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	Convey("TestUpdateAccountAuthErrReturn", t, func() {
		ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addr1)
		So(err, ShouldBeNil)

		newAuth := wallet.NewAccAddress()

		msg := accountTypes.NewMsgUpdateAccountAuth(
			addr1,
			types.MustName("noexitaccount"),
			newAuth)
		fee := types.NewInt64Coins(constants.DefaultBondDenom, 100000)

		header := abci.Header{Height: app.LastBlockHeight() + 1}
		_, _, err = simapp.SignCheckDeliver(
			t, app.Codec(), app.BaseApp,
			header, account1, fee,
			[]sdk.Msg{&msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
			false, false, wallet.PrivKey(addr1))

		So(err, simapp.ShouldErrIs, accountTypes.ErrAccountNoFound)

		// Check is account is created
		ctxCheck = app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
		currAccount := app.AccountKeeper().GetAccount(ctxCheck, account1)

		So(currAccount, ShouldNotBeNil)
		So(currAccount.GetID().Eq(account1), ShouldBeTrue)
		So(currAccount.GetAuth().Equals(addr1), ShouldBeTrue)
		So(currAccount.GetName().Eq(name1), ShouldBeTrue)

		// Check account old auth
		addr1Seq, addr1Num, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addr1)
		So(err, ShouldBeNil)
		So(addr1Num, ShouldEqual, origAuthNum)
		So(addr1Seq, ShouldEqual, origAuthSeq+1)
	})
}

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

func TestAuthIsNoWithAccount(t *testing.T) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)
	genAcc := simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1)
	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAcc)
	app := simapp.SetupWithGenesisAccounts(genAccs)
	ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	Convey("new account auth should init", t, func() {
		addr := wallet.NewAccAddress()
		err := testAccountCreate(t, app, wallet, true, account1, types.MustName("aaabbbcccdd7"), addr)
		So(err, ShouldBeNil)

		ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		checkAuthSequenceByQuery(app, ctx, addr, 1, 2)
	})

	Convey("update address should be init", t, func() {
		addr := wallet.NewAccAddress()

		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctx, addr1)
		So(err, ShouldBeNil)
		So(origAuthSeq, ShouldEqual, 2) // account1 account create a account, so it is 2

		msg := accountTypes.NewMsgUpdateAccountAuth(
			addr1,
			name1,
			addr)
		fee := types.NewInt64Coins(constants.DefaultBondDenom, 100000)

		header := abci.Header{Height: app.LastBlockHeight() + 1}
		_, _, err = simapp.SignCheckDeliver(
			t, app.Codec(), app.BaseApp,
			header, account1, fee,
			[]sdk.Msg{&msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
			true, true, wallet.PrivKey(addr1))

		So(err, ShouldBeNil)

		ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		checkAuthSequenceByQuery(app, ctx, addr1, 3, 1) // addr1 create a account so it will be 3
		checkAuthSequenceByQuery(app, ctx, addr, 1, 3)
	})
}

func TestAuthForTwoAccount(t *testing.T) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)

	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(),
		simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1),
		simapp.NewSimGenesisAccount(account2, addr1).WithAsset(asset1),
	)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	Convey("two account with one auth should be shared", t, func() {
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

		origAuthSeq, _, err := app.AccountKeeper().GetAuthSequence(ctx, addr1)
		So(err, ShouldBeNil)
		So(origAuthSeq, ShouldEqual, 1)

		err = testAccountCreate(t, app, wallet, true, account1, types.MustName("aaabbbcccdd7"), addr2)
		So(err, ShouldBeNil)

		ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		checkAuthSequenceByQuery(app, ctx, addr1, 2, 1) // addr1 create a account so it will be 2

		err = testAccountCreate(t, app, wallet, true, account2, types.MustName("aaabbbcccdd8"), addr2)
		So(err, ShouldBeNil)

		ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		checkAuthSequenceByQuery(app, ctx, addr1, 3, 1) // addr1 create a account so it will be 3
	})
}

func TestCreateAccountMsg(t *testing.T) {
	Convey("test create account msg", t, func() {
		genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth())
		app := simapp.SetupWithGenesisAccounts(genAccs)

		ctx := app.NewTestContext()
		auth := wallet.GetRootAuth()

		msg := accountTypes.NewMsgCreateAccount(
			auth,
			types.EmptyAccountID(),
			name1,
			addr1)

		So(msg.ValidateBasic(), simapp.ShouldErrIs, types.ErrKuMsgAccountIDNil)

		tx := simapp.NewTxForTest(
			constants.SystemAccountID,
			[]sdk.Msg{
				&msg,
			}, wallet.PrivKey(auth)).WithCannotPass()

		err := simapp.CheckTxs(t, app, ctx, tx)
		So(err, simapp.ShouldErrIs, types.ErrKuMsgAccountIDNil)
	})
}
