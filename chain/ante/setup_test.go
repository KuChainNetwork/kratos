package ante_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/KuChainNetwork/kuchain/chain/ante"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestSetup(t *testing.T) {
	// setup
	app, ctx := createAppForTest()

	Convey("test setup", t, func() {
		tx := testStdTx(app, account2)

		sud := ante.NewSetUpContextDecorator()
		antehandler := sdk.ChainAnteDecorators(sud)

		Convey("gasMeter set with limit before setup", func() {
			// Set height to non-zero value for GasMeter to be set
			ctx = ctx.WithBlockHeight(1)

			// Context GasMeter Limit not set
			So(uint64(0), ShouldEqual, ctx.GasMeter().Limit())
		})

		Convey("gasMeter set after setup", func() {
			newCtx, err := antehandler(ctx, tx, false)
			So(err, ShouldBeNil)

			// Context GasMeter Limit should be set after SetUpContextDecorator runs
			So(tx.Fee.Gas, ShouldEqual, newCtx.GasMeter().Limit())
		})
	})
}

func TestRecoverPanic(t *testing.T) {
	// setup
	app, ctx := createAppForTest()

	sud := ante.NewSetUpContextDecorator()

	Convey("test panic recover", t, func() {
		tx := testStdTx(app, account2)
		antehandler := sdk.ChainAnteDecorators(sud, OutOfGasDecorator{})

		// Set height to non-zero value for GasMeter to be set
		ctx = ctx.WithBlockHeight(1)

		newCtx, err := antehandler(ctx, tx, false)

		So(err, ShouldNotBeNil)
		So(err, simapp.ShouldErrIs, sdkerrors.ErrOutOfGas)

		So(tx.Fee.Gas, ShouldEqual, newCtx.GasMeter().Limit())

		antehandler = sdk.ChainAnteDecorators(sud, PanicDecorator{})
		So(func() { antehandler(ctx, tx, false) }, ShouldPanic)
	})
}

type OutOfGasDecorator struct{}

// AnteDecorator that will throw OutOfGas panic
func (ogd OutOfGasDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	overLimit := ctx.GasMeter().Limit() + 1

	// Should panic with outofgas error
	ctx.GasMeter().ConsumeGas(overLimit, "test panic")

	// not reached
	return next(ctx, tx, simulate)
}

type PanicDecorator struct{}

func (pd PanicDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	panic("random error")
}
