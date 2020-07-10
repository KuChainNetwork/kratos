package ante

import (
	"github.com/KuChainNetwork/kuchain/plugins"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PluginHandlerDecorator struct {
}

func NewPluginHandlerDecorator() PluginHandlerDecorator {
	return PluginHandlerDecorator{}
}

func (isd PluginHandlerDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	ctx.Logger().Debug("plugin ante handler")

	// no need to increment sequence on CheckTx or RecheckTx
	if ctx.IsCheckTx() && !simulate {
		return next(ctx, tx, simulate)
	}

	if std, ok := tx.(StdTx); ok {
		plugins.HandleTx(ctx, std)
	}

	return next(ctx, tx, simulate)
}
