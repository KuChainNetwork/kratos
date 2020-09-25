package keeper_test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
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

	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(),
		simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1),
	)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	ctxCheck := app.BaseApp.NewContext(true,
		abci.Header{
			Time:   time.Now(),
			Height: app.LastBlockHeight() + 1,
		})
	return app, ctxCheck
}

func createTestAppWithCoins() (*simapp.SimApp, sdk.Context) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)

	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(),
		simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1),
		simapp.NewSimGenesisAccount(account2, addr2).WithAsset(asset1),
		simapp.NewSimGenesisAccount(account3, addr3).WithAsset(asset1),
	)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	ctxCheck := app.BaseApp.NewContext(true,
		abci.Header{
			Time:   time.Now(),
			Height: app.LastBlockHeight() + 1,
		})
	return app, ctxCheck
}

func TestAssetTransfer(t *testing.T) {
	app, ctx := createTestApp()

	Convey("test transfer in keeper", t, func() {
		amt := types.NewInt64Coins(constants.DefaultBondDenom, 100)
		addr := wallet.NewAccAddress()
		err := app.AssetKeeper().Transfer(ctx, account1, types.NewAccountIDFromAccAdd(addr), amt)
		So(err, ShouldBeNil)
	})
}
