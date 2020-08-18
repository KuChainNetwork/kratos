package asset_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/KuChainNetwork/kuchain/chain/types"
)

func TestCreateAssetOpt(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test create asset issue to height", t, func() {
		var (
			se         = types.MustName("abc")
			demon      = types.CoinDenom(name4, se) // creator has @
			maxSupply  = types.Coin{demon, types.NewInt(10000000000000)}
			initSupply = types.Coin{demon, types.NewInt(0)}
			desc       = []byte(fmt.Sprintf("desc for %s", demon))
		)

		So(createCoinExt(t, app, true, account4, se, maxSupply, true, true, true, 1000, initSupply, desc), ShouldBeNil)
	})

}
