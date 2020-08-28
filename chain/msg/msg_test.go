package msg_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/msg"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	routerKeyName = types.MustName(assetTypes.RouterKey)
)

var (
	// values for test
	wallet      = simapp.NewWallet()
	addr1       = wallet.NewAccAddressByName(name1)
	addr2       = wallet.NewAccAddressByName(name2)
	addr3       = wallet.NewAccAddressByName(name3)
	addr4       = wallet.NewAccAddressByName(name4)
	addr5       = wallet.NewAccAddressByName(name5)
	name1       = types.MustName("test01@chain")
	name2       = types.MustName("aaaeeebbbccc")
	name3       = types.MustName("aaaeeebbbcc2")
	name4       = types.MustName("test")
	name5       = types.MustName("foo")
	account1    = types.NewAccountIDFromName(name1)
	account2    = types.NewAccountIDFromName(name2)
	account3    = types.NewAccountIDFromName(name3)
	account4    = types.NewAccountIDFromName(name4)
	account5    = types.NewAccountIDFromName(name5)
	addAccount1 = types.NewAccountIDFromAccAdd(addr1)
	denom       = types.CoinDenom(name5, name4)
)

func createAppForTest() (*simapp.SimApp, sdk.Context) {
	asset1 := types.NewCoins(
		types.NewInt64Coin(denom, 10000000),
		types.NewInt64Coin(constants.DefaultBondDenom, 10000000000))
	asset2 := types.NewCoins(
		types.NewInt64Coin(constants.DefaultBondDenom, 10000000000))
	asset3 := types.NewCoins(
		types.NewInt64Coin(denom, 100000),
		types.NewInt64Coin(constants.DefaultBondDenom, 10000000000))

	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(),
		simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1),
		simapp.NewSimGenesisAccount(account2, addr2).WithAsset(asset2),
		simapp.NewSimGenesisAccount(account3, addr3).WithAsset(asset3),
		simapp.NewSimGenesisAccount(account4, addr4).WithAsset(asset2),
		simapp.NewSimGenesisAccount(account5, addr5).WithAsset(asset2),
	)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
	return app, ctxCheck
}

func TestMsgMultiSend(t *testing.T) {
	msgTest := assetTypes.MsgTransfer{
		*msg.MustNewKuMsg(
			routerKeyName,
			msg.WithAuths([]types.AccAddress{addr1, addr2, addr3}),
			msg.WithTransfer(account2, account4, types.NewInt64Coins(constants.DefaultBondDenom, 2222)),
			msg.WithTransfer(account3, account5, types.NewInt64Coins(denom, 1111)),
		)}

	Convey("test msg with multiple send", t, func() {
		So(len(msgTest.Transfers), ShouldEqual, 2)

		auths := msgTest.GetSigners()

		So(len(auths), ShouldEqual, 3)
		So(auths[0], simapp.ShouldEq, addr1)
		So(auths[1], simapp.ShouldEq, addr2)
		So(auths[2], simapp.ShouldEq, addr3)
	})

	app, ctx := createAppForTest()

	Convey("test msg handler", t, func() {
		var (
			coins1 = app.AssetKeeper().GetAllBalances(ctx, account1)
			coins2 = app.AssetKeeper().GetAllBalances(ctx, account2)
			coins3 = app.AssetKeeper().GetAllBalances(ctx, account3)
			coins4 = app.AssetKeeper().GetAllBalances(ctx, account4)
			coins5 = app.AssetKeeper().GetAllBalances(ctx, account5)
		)

		tx := simapp.NewTxForTest(
			account1,
			[]sdk.Msg{
				&msgTest,
			}, wallet.PrivKey(addr1), wallet.PrivKey(addr2), wallet.PrivKey(addr3))

		So(simapp.CheckTxs(t, app, ctx, tx), ShouldBeNil)

		ctx = app.NewTestContext()
		var (
			coins1curr = app.AssetKeeper().GetAllBalances(ctx, account1)
			coins2curr = app.AssetKeeper().GetAllBalances(ctx, account2)
			coins3curr = app.AssetKeeper().GetAllBalances(ctx, account3)
			coins4curr = app.AssetKeeper().GetAllBalances(ctx, account4)
			coins5curr = app.AssetKeeper().GetAllBalances(ctx, account5)
		)

		So(coins1curr, simapp.ShouldEq, coins1.Sub(simapp.DefaultTestFee))
		So(coins2curr, simapp.ShouldEq, coins2.Sub(types.NewInt64Coins(constants.DefaultBondDenom, 2222)))
		So(coins3curr, simapp.ShouldEq, coins3.Sub(types.NewInt64Coins(denom, 1111)))
		So(coins4curr, simapp.ShouldEq, coins4.Add(types.NewInt64Coin(constants.DefaultBondDenom, 2222)))
		So(coins5curr, simapp.ShouldEq, coins5.Add(types.NewInt64Coin(denom, 1111)))
	})
}
