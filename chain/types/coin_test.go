package types

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/types/coin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoinDemond(t *testing.T) {
	Convey("test coin denomd", t, func() {
		account := MustName("kuchainasset")
		symbol := MustName("etheosbtc")

		denomd := CoinDenom(account, symbol)

		So(coin.ValidateDenom(denomd), ShouldEqual, nil)

		Printf("denomd %s", denomd)
	})
}

func TestCoinDemondCoin(t *testing.T) {
	Convey("test coin denomd", t, func() {
		account := MustName("kuchainasset")
		symbol := MustName("etheosbtc")

		denomd := CoinDenom(account, symbol)

		So(coin.ValidateDenom(denomd), ShouldBeNil)

		//Printf("denomd %s\n", denomd)

		acoin := NewCoin(denomd, NewInt(111111))

		//Printf("coin %s\n", acoin.String())

		_, err := ParseCoin(acoin.String())

		So(err, ShouldBeNil)

		//Printf("coin %s\n", bcoin.String())
	})
}
