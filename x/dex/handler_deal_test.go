package dex_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	dexTypes "github.com/KuChainNetwork/kuchain/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

func DexDealForTest(t *testing.T, app *simapp.SimApp, isSuccess bool, dex, a1 types.AccountID, c1 types.Coin, a2 types.AccountID, c2 types.Coin) error {
	var feeRateDiv int64 = 10

	wallet := app.GetWallet()

	ctx := app.NewTestContext()

	msg := dexTypes.NewMsgDexDeal(
		wallet.GetAuth(dex),
		dex,
		a1, a2,
		c1, c2,
		types.NewInt64Coin(c1.Denom, c1.Amount.Int64()/feeRateDiv),
		types.NewInt64Coin(c2.Denom, c2.Amount.Int64()/feeRateDiv),
		[]byte("deal dex test"))

	tx := simapp.NewTxForTest(
		dex,
		[]sdk.Msg{
			&msg,
		}, wallet.PrivKey(wallet.GetAuth(dex)))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return simapp.CheckTxs(t, app, ctx, tx)
}

func TestDexDealMsg(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test dex deal", t, func() {
		So(CreateDexForTest(t, app, true, dexAccount1, types.NewInt64CoreCoins(1000000000), []byte("dex for test")), ShouldBeNil)

		amt := types.NewCoins(types.NewInt64CoreCoin(1000000), types.NewInt64Coin(gDenom1, 1000000000000000))

		So(SignInMsgForTest(t, app, true, account1, dexAccount1, amt), ShouldBeNil)
		So(SignInMsgForTest(t, app, true, account2, dexAccount1, amt), ShouldBeNil)

		a11, a21 := types.NewInt64Coin(gDenom1, 600), types.NewInt64CoreCoin(1000)
		So(DexDealForTest(t, app, true, dexAccount1, account1, a11, account2, a21), ShouldBeNil)
	})
}
