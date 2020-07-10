package types

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewAccountID(t *testing.T) {
	config.SealChainConfig()

	nameStr1 := "ktsvcdf2322a3"
	nameStr2 := "k"

	accAddressStr := "kts1srwn4v9smgquglhx4rpl5j93ka4j7hs0wxhy8a"
	accAddress, _ := sdk.AccAddressFromBech32(accAddressStr)

	Convey("test new accountID from empty", t, func() {
		id, err := NewAccountIDFromStr("")
		So(err, ShouldBeNil)
		So(id.Empty(), ShouldBeTrue)
	})

	Convey("test new accountID from name 1", t, func() {
		_, err := NewAccountIDFromStr(nameStr1)
		So(err, ShouldBeNil)
	})

	Convey("test new accountID from name 2", t, func() {
		_, err := NewAccountIDFromStr(nameStr2)
		So(err, ShouldBeNil)
	})

	Convey("test new accountID from Address", t, func() {
		acc, err := NewAccountIDFromStr(accAddressStr)
		So(err, ShouldBeNil)
		So(acc.Equals(accAddress), ShouldBeTrue)

		accAdd, ok := acc.ToAccAddress()
		So(ok, ShouldBeTrue)
		So(accAddress.Equals(accAdd), ShouldBeTrue)
	})
}
