package eventutil_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/utils/eventutil"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnmarshalEvent(t *testing.T) {
	Convey("test unmarshal event", t, func() {
		var testEvtStruct testEvent
		err := eventutil.UnmarshalEvent(testEvt, &testEvtStruct)

		So(err, ShouldBeNil)

		t.Logf("test evt %v", testEvtStruct)

		Convey("test unmarshal name", func() {
			t.Logf("test evt name %s", testEvtStruct.Name)
			So(testEvtStruct.Name, simapp.ShouldEq, testName)
		})

		Convey("test unmarshal name array", func() {
			So(len(testEvtStruct.Names), ShouldEqual, len(testNames))
			for i, n := range testNames {
				t.Logf("test evt names %d: %s", i, testEvtStruct.Names[i])
				So(testEvtStruct.Names[i], simapp.ShouldEq, n)
			}
		})

		Convey("test unmarshal id", func() {
			So(testEvtStruct.Id, simapp.ShouldEq, testIDName)
			So(testEvtStruct.IdAddr, simapp.ShouldEq, testIDAddr)
		})

		Convey("test unmarshal auth address", func() {
			So(testEvtStruct.Auth, simapp.ShouldEq, testAuth)
		})

		Convey("test unmarshal coin", func() {
			So(testEvtStruct.Coin, simapp.ShouldEq, testCoin)
		})

		Convey("test unmarshal coins", func() {
			So(testEvtStruct.Coins, simapp.ShouldEq, testCoins)
		})

		Convey("test unmarshal base type", func() {
			So(testEvtStruct.Str, ShouldEqual, "test str for event")
			So(testEvtStruct.Int, ShouldEqual, 1234567)
			So(testEvtStruct.Strs, ShouldResemble, []string{"test", "str", "for", "event"})
			So(testEvtStruct.Ints, ShouldResemble, []int64{-12, 34, 56, 7})
			So(testEvtStruct.Bool1, ShouldEqual, true)
			So(testEvtStruct.Bool2, ShouldEqual, false)
		})
	})
}
