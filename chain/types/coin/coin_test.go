package coin_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/types/coin"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ld1 = "chain/sys"
	ld2 = "aaaaaaa/cccc"
)

func TestSubLittleCoinDenom(t *testing.T) {
	Convey("test safe sub", t, func() {
		cc1 := coin.Coins{
			coin.NewInt64Coin(ld1, 1000),
			coin.NewInt64Coin(ld2, 2000),
		}

		cc2 := coin.Coins{
			coin.NewInt64Coin(ld1, 100),
		}

		cc1m2, isNeg := cc1.SafeSub(cc2)
		So(isNeg, ShouldBeFalse)

		t.Logf("cc1m2 %s %v", cc1m2.String(), cc1m2.AmountOf(ld1))
		t.Logf("cc1m2 0 %s %v", cc1m2[0].Denom, cc1m2[0].Amount)
		t.Logf("cc1m2 1 %s %v", cc1m2[1].Denom, cc1m2[1].Amount)

		cc1m2 = cc1m2.Sort()
		t.Logf("cc1m2 %s %v", cc1m2.String(), cc1m2.AmountOf(ld1))
		t.Logf("cc1m2 0 %s %v", cc1m2[0].Denom, cc1m2[0].Amount)
		t.Logf("cc1m2 1 %s %v", cc1m2[1].Denom, cc1m2[1].Amount)

		So(len(cc1m2), ShouldEqual, 2)
		So(cc1m2.AmountOf(ld1).Int64(), ShouldEqual, 900)
		So(cc1m2.AmountOf(ld2).Int64(), ShouldEqual, 2000)
	})

	Convey("test coin denom", t, func() {
		cc1 := coin.NewInt64Coin(ld1, 1000)
		cc1s := coin.Coins{cc1}
		So(cc1.Denom, ShouldEqual, ld1)
		So(cc1s.AmountOf(ld1).Int64(), ShouldEqual, 1000)
	})

}
