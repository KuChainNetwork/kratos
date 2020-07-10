package test

import (
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/plugins/test/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// testPlugin for test
type testPlugin struct {
	logger log.Logger
}

func (t *testPlugin) Init(ctx types.Context) error {
	t.logger.Info("plugin init", "name", types.PluginName)

	return nil
}

func (t *testPlugin) Start(ctx types.Context) error {
	t.logger.Info("plugin start", "name", types.PluginName)

	return nil
}

func (t *testPlugin) Stop(ctx types.Context) error {
	t.logger.Info("plugin stop", "name", types.PluginName)

	return nil
}

func (t *testPlugin) OnEvent(ctx types.Context, evt types.Event) {
	t.logger.Info("on event", "type", evt.Type)
}

func (t *testPlugin) OnTx(ctx types.Context, tx chainTypes.StdTx) {
	t.logger.Info("on tx", "tx", tx)
}

func (t *testPlugin) OnMsg(ctx types.Context, msg sdk.Msg) {
	t.logger.Info("on msg", "msg", msg)
}

func (t *testPlugin) MsgHandler() types.PluginMsgHandler {
	return func(ctx types.Context, msg sdk.Msg) {
		t.OnMsg(ctx, msg)
	}
}

func (t *testPlugin) TxHandler() types.PluginTxHandler {
	return func(ctx types.Context, tx chainTypes.StdTx) {
		t.OnTx(ctx, tx)
	}
}

func (t *testPlugin) EvtHandler() types.PluginEvtHandler {
	return func(ctx types.Context, evt types.Event) {
		t.OnEvent(ctx, evt)
	}
}

func (t *testPlugin) Logger() log.Logger {
	return t.logger
}

func (t *testPlugin) Name() string {
	return types.PluginName
}

// NewTestPlugin new test plugin
func NewTestPlugin(ctx types.Context, cfg types.BaseCfg) *testPlugin {
	return &testPlugin{
		logger: ctx.Logger().With("module", types.PluginName),
	}
}
