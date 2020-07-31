package keeper_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoinPower(t *testing.T) {
	app, ctx := createTestApp()

	amt := types.NewCoins(types.NewInt64Coin(constants.DefaultBondDenom, 100000))

	Convey("test coins to power by self", t, func() {
		err := app.AssetKeeper().CoinsToPower(ctx, account1, account1, amt)
		So(err, ShouldBeNil)

		coins := app.AssetKeeper().GetCoinPowers(ctx, account1)
		So(coins, simapp.ShouldEq, amt)
	})

}
