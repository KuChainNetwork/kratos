package types

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/plugins/types"
	"github.com/tendermint/tendermint/libs/log"
)

type (
	Context          = types.Context
	Event            = types.Event
	BaseCfg          = types.BaseCfg
	PluginMsgHandler = types.PluginMsgHandler
	PluginTxHandler  = types.PluginTxHandler
	PluginEvtHandler = types.PluginEvtHandler
)

func Logger(ctx Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("plugins/%s", PluginName))
}
