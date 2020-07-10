package ante

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	StdTx = types.StdTx
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewHandler(ak keeper.AccountKeeper, asset AssetKeeper) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewSetUpContextDecorator(),
		NewValidateBasicDecorator(),
		NewMempoolFeeDecorator(),
		NewConsumeGasForTxSizeDecorator(),
		NewDeductFeeDecorator(ak, asset),
		NewSetPubKeyDecorator(ak),
		NewSigVerificationDecorator(ak),
		NewIncrementSequenceDecorator(ak),
		NewPluginHandlerDecorator(),
	)
}
