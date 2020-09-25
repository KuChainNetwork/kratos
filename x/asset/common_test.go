package asset_test

import (
	"fmt"
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
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
	addAccount2 = types.NewAccountIDFromAccAdd(addr2)
	addAccount3 = types.NewAccountIDFromAccAdd(addr3)
	addAccount4 = types.NewAccountIDFromAccAdd(addr4)
	addAccount5 = types.NewAccountIDFromAccAdd(addr5)
)

var (
	NewInt64Coins     = types.NewInt64Coins
	NewInt64Coin      = types.NewInt64Coin
	MustName          = types.MustName
	CheckTxs          = simapp.CheckTxs
	NewInt64CoreCoin  = types.NewInt64CoreCoin
	NewInt64CoreCoins = types.NewInt64CoreCoins
)

func createAppForTest() (*simapp.SimApp, sdk.Context) {
	asset1 := types.NewCoins(
		types.NewInt64Coin("foo/coin", 10000000),
		types.NewInt64Coin(constants.DefaultBondDenom, 1000000000000))
	asset2 := types.NewCoins(
		types.NewInt64Coin(constants.DefaultBondDenom, 1000000000000))
	asset3 := types.NewCoins(
		types.NewInt64Coin("foo/coin", 100),
		types.NewInt64Coin(constants.DefaultBondDenom, 1000000000000))

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

func transfer(t *testing.T, app *simapp.SimApp, isSuccess bool, from, to types.AccountID, amt types.Coins, payer types.AccountID) error {
	return simapp.CommitTransferTx(t, app, wallet, isSuccess, from, to, amt, payer)
}

func createCoin(t *testing.T, app *simapp.SimApp, isSuccess bool,
	creator types.AccountID, symbol types.Name, maxSupplyAmount int64) error {
	ctx := app.NewTestContext()
	creatorName := creator.MustName()

	var (
		demon      = types.CoinDenom(creatorName, symbol)
		maxSupply  = types.NewCoin(demon, types.NewInt(maxSupplyAmount))
		initSupply = types.NewCoin(demon, types.NewInt(0))
		desc       = []byte(fmt.Sprintf("desc for %s", demon))
	)

	auth := app.AccountKeeper().GetAccount(ctx, creator).GetAuth()

	msg := assetTypes.NewMsgCreate(auth, creatorName, symbol, maxSupply, true, true, true, 0, initSupply, desc)
	tx := simapp.NewTxForTest(
		creator,
		[]sdk.Msg{
			&msg,
		}, wallet.PrivKey(auth))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return simapp.CheckTxs(t, app, ctx, tx)
}

func createCoinExt(t *testing.T, app *simapp.SimApp, isSuccess bool,
	creator types.AccountID, symbol types.Name, maxSupply types.Coin, canIssue, canLock, canBurn bool, issue2Height int64, initSupply types.Coin, desc []byte) error {
	ctx := app.NewTestContext()
	creatorName := creator.MustName()

	auth := app.AccountKeeper().GetAccount(ctx, creator).GetAuth()

	msg := assetTypes.NewMsgCreate(auth, creatorName, symbol, maxSupply, canIssue, canLock, canBurn, issue2Height, initSupply, desc)
	tx := simapp.NewTxForTest(
		creator,
		[]sdk.Msg{
			&msg,
		}, wallet.PrivKey(auth))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return simapp.CheckTxs(t, app, ctx, tx)
}

func issueCoin(t *testing.T, app *simapp.SimApp, isSuccess bool,
	creator types.AccountID, symbol types.Name, amt types.Coin) error {
	ctx := app.NewTestContext()
	creatorName := creator.MustName()

	auth := app.AccountKeeper().GetAccount(ctx, creator).GetAuth()

	msg := assetTypes.NewMsgIssue(auth, creatorName, symbol, amt)
	tx := simapp.NewTxForTest(
		creator,
		[]sdk.Msg{
			&msg,
		}, wallet.PrivKey(auth))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return simapp.CheckTxs(t, app, ctx, tx)
}

// BurnCoinTest burn coins test
func BurnCoinTest(t *testing.T, app *simapp.SimApp, isSuccess bool,
	from types.AccountID, amt types.Coin) error {
	ctx := app.NewTestContext()

	auth := app.AccountKeeper().GetAccount(ctx, from).GetAuth()

	msg := assetTypes.NewMsgBurn(auth, from, amt)
	tx := simapp.NewTxForTest(
		from,
		[]sdk.Msg{
			&msg,
		}, wallet.PrivKey(auth))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return simapp.CheckTxs(t, app, ctx, tx)
}
