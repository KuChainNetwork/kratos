package eventutil_test

import (
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	wallet = simapp.NewWallet()
)

type testEvent struct {
	Name   types.Name       `json:"name"`
	Names  []types.Name     `json:"names"`
	Id     types.AccountID  `json:"id"`
	IdAddr types.AccountID  `json:"idaddr"`
	Auth   types.AccAddress `json:"auth"`
	Coin   types.Coin       `json:"coin"`
	Coins  types.Coins      `json:"coins"`
	Str    string           `json:"str"`
	Int    int64            `json:"int"`
	Ints   []int64          `json:"ints"`
	Strs   []string         `json:"strs"`
	UInt   uint64           `json:"uint"`
	Bool1  bool             `json:"b1"`
	Bool2  bool             `json:"b2"`
}

var (
	testCoin  = types.NewInt64Coin(constants.DefaultBondDenom, 100001)
	testCoins = types.NewCoins(
		types.NewInt64Coin(constants.DefaultBondDenom, 100),
		types.NewInt64Coin("foo/test1", 200),
		types.NewInt64Coin("foo/test2", 300),
		types.NewInt64Coin("foo/test3", 400),
	)

	testName   = types.MustName("adss@sdssd")
	testNames  = []types.Name{types.MustName("abcde"), types.MustName("aa@aaa"), types.MustName("kuc.ha@in")}
	testIDName = types.MustAccountID("adss@sd2s")
	testAuth   = wallet.NewAccAddress()
	testIDAddr = types.NewAccountIDFromAccAdd(wallet.NewAccAddress())

	testEvt sdk.Event = sdk.NewEvent("test",
		sdk.NewAttribute("name", testName.String()),
		sdk.NewAttribute("names", "abcde, aa@aaa,kuc.ha@in"),
		sdk.NewAttribute("id", testIDName.String()),
		sdk.NewAttribute("idaddr", testIDAddr.String()),
		sdk.NewAttribute("auth", testAuth.String()),
		sdk.NewAttribute("coin", testCoin.String()),
		sdk.NewAttribute("coins", testCoins.String()),
		sdk.NewAttribute("str", "test str for event"),
		sdk.NewAttribute("int", "1234567"),
		sdk.NewAttribute("uint", "7654321"),
		sdk.NewAttribute("strs", "test,str, for, event"),
		sdk.NewAttribute("ints", "-12,34,56,7"),
		sdk.NewAttribute("b1", "true"),
		sdk.NewAttribute("b2", "false"),
	)
)
