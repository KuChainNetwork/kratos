package dex_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
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
)

var (
	dexName1    = types.MustName("dexaccount1")
	dexAccount1 = types.NewAccountIDFromName(dexName1)
	dexAddr1    = wallet.NewAccAddressByName(dexName1)

	dexName2    = types.MustName("dexaccount2")
	dexAccount2 = types.NewAccountIDFromName(dexName2)
	dexAddr2    = wallet.NewAccAddressByName(dexName2)
)

var (
	gDenom1 = types.CoinDenom(name4, types.MustName("coin1"))
	gDenom2 = types.CoinDenom(name4, types.MustName("coin2"))
	gDenom3 = types.CoinDenom(name4, types.MustName("coin3"))
	gDenom4 = types.CoinDenom(name4, types.MustName("coin4"))
)

func createAppForTest() (*simapp.SimApp, sdk.Context) {
	asset1 := types.NewCoins(
		types.NewInt64Coin(gDenom1, 1000000000000000),
		types.NewInt64Coin(gDenom2, 1000000000000000),
		types.NewInt64Coin(gDenom3, 1000000000000000),
		types.NewInt64Coin(gDenom4, 1000000000000000),
		types.NewInt64Coin(constants.DefaultBondDenom, 10000000000))
	asset2 := types.NewCoins(
		types.NewInt64Coin(constants.DefaultBondDenom, 10000000000))

	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(),
		simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1),
		simapp.NewSimGenesisAccount(account2, addr2).WithAsset(asset1),
		simapp.NewSimGenesisAccount(account3, addr3).WithAsset(asset1),
		simapp.NewSimGenesisAccount(account4, addr4).WithAsset(asset1),
		simapp.NewSimGenesisAccount(account5, addr5).WithAsset(asset1),
		simapp.NewSimGenesisAccount(dexAccount1, dexAddr1).WithAsset(asset2),
		simapp.NewSimGenesisAccount(dexAccount2, dexAddr2).WithAsset(asset2),
	)
	app := simapp.SetupWithGenesisAccounts(genAccs)
	app = app.WithWallet(wallet)

	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
	return app, ctxCheck
}

func transfer(t *testing.T, app *simapp.SimApp, isSuccess bool, from, to types.AccountID, amt types.Coins, payer types.AccountID) error {
	return simapp.CommitTransferTx(t, app, wallet, isSuccess, from, to, amt, payer)
}
