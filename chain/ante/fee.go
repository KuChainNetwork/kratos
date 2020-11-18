package ante

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ FeeTx = (*types.StdTx)(nil) // assert StdTx implements FeeTx
)

// FeeTx defines the interface to be implemented by Tx to use the FeeDecorators
type FeeTx interface {
	sdk.Tx
	GetGas() uint64
	GetFee() Coins
	FeePayer() AccountID
}

// MempoolFeeDecorator will check if the transaction's fee is at least as large
// as the local validator's minimum gasFee (defined in validator config).
// If fee is too low, decorator returns error and tx is rejected from mempool.
// Note this only applies when ctx.CheckTx = true
// If fee is high enough or not CheckTx, then call next AnteHandler
// CONTRACT: Tx must implement FeeTx to use MempoolFeeDecorator
type MempoolFeeDecorator struct{}

func NewMempoolFeeDecorator() MempoolFeeDecorator {
	return MempoolFeeDecorator{}
}

func (mfd MempoolFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}
	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() && !simulate {
		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredFees := make(Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := types.NewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = types.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			if !feeCoins.IsAnyGTE(requiredFees) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}
	}

	return next(ctx, tx, simulate)
}

// DeductFeeDecorator deducts fees from payer or the first signer of the tx
type DeductFeeDecorator struct {
	ak      AssetKeeper
	account AccountKeeper
}

func NewDeductFeeDecorator(acc AccountKeeper, ak AssetKeeper) DeductFeeDecorator {
	return DeductFeeDecorator{
		ak:      ak,
		account: acc,
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()
	ctx.Logger().Debug("fee deduct", "feePayer", feePayer, "fee", feeTx.GetFee(), "gas", feeTx.GetGas())

	// deduct the fees
	if !feeTx.GetFee().IsZero() {
		if err := checkPayerAuth(ctx, dfd.account, tx, simulate, feePayer); err != nil {
			return ctx, err
		}

		err = DeductFees(ctx, dfd.ak, feePayer, feeTx.GetFee())
		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

// DeductFees deducts fees from the given account.
func DeductFees(ctx sdk.Context, assetKeeper AssetKeeper, payer AccountID, fees Coins) error {
	return assetKeeper.PayFee(ctx, payer, fees)
}

func checkPayerAuth(ctx sdk.Context, ak AccountKeeper, tx sdk.Tx, simulate bool, payer AccountID) error {
	// if !ctx.IsCheckTx() {
	// 	return nil
	// }

	sigTx, ok := tx.(SigVerifiableTx)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type for SigVerifiableTx")
	}

	auths := sigTx.GetSigners()

	if add, ok := payer.ToAccAddress(); ok {
		if !isHasAuth(auths, add) {
			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "signature verification failed; fee payer address no found")
		}
		return nil
	}

	acc := ak.GetAccount(ctx, payer)
	if acc == nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "payer not found")
	}

	accAuth := acc.GetAuth()

	if !isHasAuth(auths, accAuth) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "signature verification failed; fee payer account auth no found")
	}

	return nil
}

func isHasAuth(auths []types.AccAddress, auth types.AccAddress) bool {
	for _, a := range auths {
		if a.Equals(auth) {
			return true
		}
	}

	return false
}
