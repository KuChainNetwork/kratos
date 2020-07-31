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
	"github.com/KuChainNetwork/kuchain/x/mint"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types/time"
)

func TestQueryAccount(t *testing.T) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)
	genAcc1 := simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1)
	genAcc2 := simapp.NewSimGenesisAccount(types.NewAccountIDFromAccAdd(addr2), addr2).WithAsset(asset1)
	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAcc1, genAcc2)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	ctx := app.BaseApp.NewContext(true,
		abci.Header{
			Time:   time.Now(),
			Height: app.LastBlockHeight() + 1,
		})

	Convey("TestQueryAccount", t, func() {
		querier := keeper.NewQuerier(*app.AccountKeeper())
		path := []string{accountTypes.QueryAccount}

		Convey("no exit query", func() {
			req := abci.RequestQuery{
				Path: "",
				Data: []byte{},
			}

			bz, err := querier(ctx, []string{"other"}, req)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrUnknownRequest)
			So(bz, ShouldBeNil)
		})

		Convey("query a account with no params", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAccount),
				Data: []byte{},
			}
			res, err := querier(ctx, path, req)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrJSONUnmarshal)
			So(res, ShouldBeNil)
		})

		Convey("query a account with empty params", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAccount),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAccountParams(types.EmptyAccountID()))

			res, err := querier(ctx, path, req)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrUnknownAddress)
			So(res, ShouldBeNil)
		})

		Convey("query a account", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAccount),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAccountParams(account1))

			res, err := querier(ctx, path, req)
			So(err, ShouldBeNil)

			var accRes exported.Account
			err = app.Codec().UnmarshalJSON(res, &accRes)
			So(err, ShouldBeNil)

			So(accRes.GetID(), simapp.ShouldEq, account1)
			So(accRes.GetName(), simapp.ShouldEq, name1)
			So(accRes.GetAuth(), simapp.ShouldEq, addr1)
			So(accRes.GetAccountNumber(), ShouldEqual, 2) // the 2rd account
		})

		Convey("query a module account", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAccount),
				Data: []byte{},
			}
			id := types.MustAccountID(mint.ModuleName)
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAccountParams(id))

			res, err := querier(ctx, path, req)
			So(err, ShouldBeNil)

			var accRes exported.Account
			err = app.Codec().UnmarshalJSON(res, &accRes)
			So(err, ShouldBeNil)

			So(accRes.GetID(), simapp.ShouldEq, id)
			So(accRes.GetName(), simapp.ShouldEq, types.MustName(mint.ModuleName))
			So(accRes.GetAuth(), simapp.ShouldEq, types.AccAddress{})
		})

		Convey("query a address", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAccount),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAccountParams(types.NewAccountIDFromAccAdd(addr2)))

			res, err := querier(ctx, path, req)
			So(err, ShouldBeNil)

			var accRes exported.Account
			err = app.Codec().UnmarshalJSON(res, &accRes)
			So(err, ShouldBeNil)

			So(accRes.GetID(), simapp.ShouldEq, types.NewAccountIDFromAccAdd(addr2))
			So(accRes.GetName(), simapp.ShouldEq, types.Name{})
			So(accRes.GetAuth(), simapp.ShouldEq, addr2)
		})

		Convey("query a no existing address", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAccount),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(
				accountTypes.NewQueryAccountParams(
					types.NewAccountIDFromAccAdd(wallet.NewAccAddress())))

			res, err := querier(ctx, path, req)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrUnknownAddress)
			So(res, ShouldBeNil)
		})

		Convey("query a no existing account", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAccount),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAccountParams(types.MustAccountID("aabbccdd")))

			res, err := querier(ctx, path, req)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrUnknownAddress)
			So(res, ShouldBeNil)
		})
	})
}

func TestQueryAuth(t *testing.T) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)
	genAcc1 := simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1)
	genAcc2 := simapp.NewSimGenesisAccount(types.NewAccountIDFromAccAdd(addr2), addr2).WithAsset(asset1)
	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAcc1, genAcc2)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	ctx := app.BaseApp.NewContext(true,
		abci.Header{
			Time:   time.Now(),
			Height: app.LastBlockHeight() + 1,
		})

	Convey("TestQueryAuth", t, func() {
		querier := keeper.NewQuerier(*app.AccountKeeper())
		path := []string{accountTypes.QueryAuthByAddress}

		Convey("query a auth no params", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAuthByAddress),
				Data: []byte{},
			}
			res, err := querier(ctx, path, req)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrJSONUnmarshal)
			So(res, ShouldBeNil)
		})

		Convey("query a auth with empty params", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAuthByAddress),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAddAuthParams(types.AccAddress{}))

			res, err := querier(ctx, path, req)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrUnknownAddress)
			So(res, ShouldBeNil)
		})

		Convey("query a auth", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAuthByAddress),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAddAuthParams(addr1))

			res, err := querier(ctx, path, req)
			So(err, ShouldBeNil)

			var authData accountTypes.Auth
			err = app.Codec().UnmarshalJSON(res, &authData)
			So(err, ShouldBeNil)

			So(authData.GetAddress(), simapp.ShouldEq, addr1)
			So(authData.GetNumber(), ShouldEqual, 1) // first genesis auth
			So(authData.GetSequence(), ShouldEqual, 1)
			So(authData.GetPubKey(), ShouldBeNil) // no init pubkey
		})

		Convey("query a address auth", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAuthByAddress),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAddAuthParams(addr2))

			res, err := querier(ctx, path, req)
			So(err, ShouldBeNil)

			var authData accountTypes.Auth
			err = app.Codec().UnmarshalJSON(res, &authData)
			So(err, ShouldBeNil)

			So(authData.GetAddress(), simapp.ShouldEq, addr2)
			So(authData.GetNumber(), ShouldEqual, 2)
			So(authData.GetSequence(), ShouldEqual, 1)
			So(authData.GetPubKey(), ShouldBeNil) // no init pubkey
		})

		Convey("query a no existing auth", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAuthByAddress),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAddAuthParams(wallet.NewAccAddress()))

			res, err := querier(ctx, path, req)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrUnknownAddress)
			So(res, ShouldBeNil)
		})
	})
}

func TestQueryAccountsByAuth(t *testing.T) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)

	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(),
		simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1),
		simapp.NewSimGenesisAccount(types.MustAccountID("aabbcc1"), addr1).WithAsset(asset1),
		simapp.NewSimGenesisAccount(types.MustAccountID("aabbcc2"), addr1).WithAsset(asset1),
		simapp.NewSimGenesisAccount(types.MustAccountID("aabbcc3"), addr1).WithAsset(asset1),
		simapp.NewSimGenesisAccount(types.NewAccountIDFromAccAdd(addr1), addr1).WithAsset(asset1),
		simapp.NewSimGenesisAccount(types.NewAccountIDFromAccAdd(addr2), addr2).WithAsset(asset1),
		simapp.NewSimGenesisAccount(account3, addr3).WithAsset(asset1),
	)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	ctx := app.BaseApp.NewContext(true,
		abci.Header{
			Time:   time.Now(),
			Height: app.LastBlockHeight() + 1,
		})

	Convey("TestQueryAccountsByAuth", t, func() {
		querier := keeper.NewQuerier(*app.AccountKeeper())
		path := []string{accountTypes.QueryAccountsByAuth}

		Convey("query accounts no params", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAccountsByAuth),
				Data: []byte{},
			}
			res, err := querier(ctx, path, req)
			So(err, simapp.ShouldErrIs, sdkerrors.ErrJSONUnmarshal)
			So(res, ShouldBeNil)
		})

		Convey("query no existing auth", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAccountsByAuth),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAccountsByAuthParams(wallet.NewAccAddress().String()))

			res, err := querier(ctx, path, req)
			So(err, ShouldBeNil)

			var accounts []string
			err = app.Codec().UnmarshalJSON(res, &accounts)
			So(err, ShouldBeNil)
			So(len(accounts), ShouldEqual, 0)
		})

		Convey("query auth", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAccountsByAuth),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAccountsByAuthParams(addr1.String()))

			res, err := querier(ctx, path, req)
			So(err, ShouldBeNil)

			var accounts []string
			err = app.Codec().UnmarshalJSON(res, &accounts)
			So(err, ShouldBeNil)
			So(len(accounts), ShouldEqual, 4) // no address, just 4 account
		})

		Convey("query auth address", func() {
			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAccountsByAuth),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAccountsByAuthParams(addr2.String()))

			res, err := querier(ctx, path, req)
			So(err, ShouldBeNil)

			var accounts []string
			err = app.Codec().UnmarshalJSON(res, &accounts)
			So(err, ShouldBeNil)
			So(len(accounts), ShouldEqual, 0) // address no in accounts
		})
	})
}
