package dex_test

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	dexTypes "github.com/KuChainNetwork/kuchain/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

func SignInMsgForTest(t *testing.T, app *simapp.SimApp, isSuccess bool, id, dex types.AccountID, amount types.Coins) error {
	wallet := app.GetWallet()

	ctx := app.NewTestContext()

	msg := dexTypes.NewMsgDexSigIn(
		wallet.GetAuth(id),
		id,
		dex,
		amount)

	tx := simapp.NewTxForTest(
		id,
		[]sdk.Msg{
			&msg,
		}, wallet.PrivKey(wallet.GetAuth(id)))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return simapp.CheckTxs(t, app, ctx, tx)
}

func TestSignInMsg(t *testing.T) {
	app, _ := createAppForTest()

	Convey("test signIn msg", t, func() {
		amt := types.NewInt64CoreCoins(1000000)

		So(SignInMsgForTest(t, app, true, account1, dexAccount1, amt), ShouldBeNil)
	})
}
