package plugins

import (
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	dbHistory "github.com/KuChainNetwork/kuchain/plugins/db_history"
	"github.com/KuChainNetwork/kuchain/plugins/test"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: use a goroutine
var (
	plugins *Plugins
)

func InitPlugins(ctx Context, cfgs []BaseCfg) error {
	plugins = NewPlugins(ctx.Logger().With("module", "plugins"))
	for _, cfg := range cfgs {
		initPlugin(ctx, cfg, plugins)
	}

	plugins.Start()

	return nil
}

func StopPlugins(ctx Context) {
	if plugins != nil {
		plugins.Stop(ctx)
	}
}

func initPlugin(ctx Context, cfg BaseCfg, plugins *Plugins) {
	switch cfg.Name {
	case test.PluginName:
		plugins.RegPlugin(ctx, test.NewTestPlugin(ctx, cfg))
	case dbHistory.PluginName:
		plugins.RegPlugin(ctx, dbHistory.New(ctx, cfg))
	}
}

// HandleEvent plugins handler Events
func HandleEvent(ctx sdk.Context, evts sdk.Events) {
	if plugins == nil {
		return
	}

	for _, evt := range evts {
		plugins.EmitEvent(evt)
	}
}

// HandleTx handler tx for each plugins
func HandleTx(ctxSdk sdk.Context, tx chainTypes.StdTx) {
	if plugins == nil {
		return
	}

	plugins.EmitTx(tx)
}
