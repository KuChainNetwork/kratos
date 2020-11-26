package asset_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants/keys"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/asset"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenesisExport(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test genesis export", t, func() {
		// create assets
		assetsCreated := []struct {
			Creator types.Name
			Symbol  types.Name
			Supply  int64
		}{
			{name4, types.MustName("abc"), defaultCoinSupply},
			{name4, types.MustName("abc1"), defaultCoinSupply},
			{name5, types.MustName("abc2"), defaultCoinSupply},
			{name5, types.MustName("abc3"), defaultCoinSupply},
			{name5, types.MustName("abc4"), defaultCoinSupply},
		}

		for _, ass := range assetsCreated {
			So(createCoin(t, app, true,
				types.NewAccountIDFromName(ass.Creator), ass.Symbol, ass.Supply), ShouldBeNil)
		}

		// issue some assets
		for _, ass := range assetsCreated {
			So(issueCoin(t, app, true,
				types.NewAccountIDFromName(ass.Creator), ass.Symbol,
				types.NewInt64Coin(types.CoinDenom(ass.Creator, ass.Symbol), ass.Supply/2)), ShouldBeNil)
		}

		// export
		genesis := asset.ExportGenesis(app.NewTestContext(), *app.AssetKeeper())
		So(assetTypes.ValidateGenesis(genesis), ShouldBeNil)

		jsonBz, _ := app.Codec().MarshalJSON(genesis)

		// check
		app.Logger().Debug("genesis-data", "data", string(jsonBz))

		coins := make(map[string]asset.GenesisCoin)
		for _, c := range genesis.GenesisCoins {
			coins[types.CoinDenom(c.GetCreator(), c.GetSymbol())] = c
		}

		So(coins[keys.DefaultBondDenom].GetCreator().String(), ShouldEqual, keys.ChainMainNameStr)
		So(coins[keys.DefaultBondDenom].GetSymbol().String(), ShouldEqual, keys.DefaultBondSymbol)
		So(coins[keys.DefaultBondDenom].GetMaxSupply().Amount.Int64(), ShouldEqual, 0)

		for _, cc := range assetsCreated {
			n, ok := coins[types.CoinDenom(cc.Creator, cc.Symbol)]
			So(ok, ShouldBeTrue)

			So(n.GetCreator(), simapp.ShouldEq, cc.Creator)
			So(n.GetSymbol(), simapp.ShouldEq, cc.Symbol)
			So(n.GetSupply().Amount.Int64(), ShouldEqual, cc.Supply/2)
			So(n.GetMaxSupply().Amount.Int64(), ShouldEqual, cc.Supply)
		}
	})
}
