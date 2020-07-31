package ante_test

import (
	"encoding/json"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/KuChainNetwork/kuchain/chain/ante"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestValidateBasic(t *testing.T) {
	app, ctx := createAppForTest()

	Convey("test validate check", t, func() {
		invalidTx := testStdTx(app, account2)
		invalidTx.Memo = strings.Repeat("aaaaaaaa", 100)

		vbd := ante.NewValidateBasicDecorator()
		antehandler := sdk.ChainAnteDecorators(vbd)
		_, err := antehandler(ctx, invalidTx, false)
		So(err, simapp.ShouldErrIs, sdkerrors.ErrMemoTooLarge)

		validTx := testStdTx(app, account2)

		_, err = antehandler(ctx, validTx, false)
		So(err, ShouldBeNil)

		// test decorator skips on recheck
		ctx = ctx.WithIsReCheckTx(true)

		// decorator should skip processing invalidTx on recheck and thus return nil-error
		_, err = antehandler(ctx, invalidTx, false)
		So(err, ShouldBeNil)
	})
}

func TestConsumeGasForTxSize(t *testing.T) {
	app, ctx := createAppForTest()

	Convey("test consume gas for tx size check", t, func() {

		tx := testStdTx(app, account2)
		tx.Memo = strings.Repeat("01234567890", 10) // 100 len
		txBytes, err := json.Marshal(tx)
		So(err, ShouldBeNil)

		antehandler := sdk.ChainAnteDecorators(ante.NewConsumeGasForTxSizeDecorator())
		expectedGas := sdk.Gas(len(txBytes)) * constants.GasTxSizePrice

		Convey("decorator consume should be the correct amount of gas", func() {
			// Set ctx with TxBytes manually
			ctx = ctx.WithTxBytes(txBytes)

			ctx, err = antehandler(ctx, tx, false)
			So(err, ShouldBeNil)

			// require that decorator consumes expected amount of gas
			consumedGas := ctx.GasMeter().GasConsumed()
			So(expectedGas, ShouldEqual, consumedGas)
		})

		Convey("require that simulated tx is smaller than tx with signatures", func() {
			// simulation must not underestimate gas of this decorator even with nil signatures
			sigTx := tx
			sigTx.Signatures = []types.StdSignature{{}}

			simTxBytes, err := json.Marshal(sigTx)
			So(err, ShouldBeNil)
			So(len(simTxBytes), ShouldBeLessThan, len(txBytes))

			// Set ctx with smaller simulated TxBytes manually
			ctx = ctx.WithTxBytes(txBytes)

			beforeSimGas := ctx.GasMeter().GasConsumed()

			// run antehandler with simulate=true
			ctx, err = antehandler(ctx, sigTx, true)
			consumedSimGas := ctx.GasMeter().GasConsumed() - beforeSimGas

			// require that antehandler passes and does not underestimate decorator cost
			So(err, ShouldBeNil)
			So(expectedGas, ShouldBeLessThanOrEqualTo, consumedSimGas)
		})
	})
}
