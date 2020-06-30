package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoinDemond(t *testing.T) {
	Convey("test coin denomd", t, func() {
		account := MustName("kuchainasset")
		symbol := MustName("etheosbtc")

		denomd := CoinDenom(account, symbol)

		So(sdk.ValidateDenom(denomd), ShouldEqual, nil)

		Printf("denomd %s", denomd)
	})
}

func TestCoinDemondCoin(t *testing.T) {
	Convey("test coin denomd", t, func() {
		account := MustName("kuchainasset")
		symbol := MustName("etheosbtc")

		denomd := CoinDenom(account, symbol)

		So(sdk.ValidateDenom(denomd), ShouldBeNil)

		//Printf("denomd %s\n", denomd)

		acoin := sdk.NewCoin(denomd, sdk.NewInt(111111))

		//Printf("coin %s\n", acoin.String())

		_, err := sdk.ParseCoin(acoin.String())

		So(err, ShouldBeNil)

		//Printf("coin %s\n", bcoin.String())
	})
}
