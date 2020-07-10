package dbHistory

import (
	"encoding/json"
	"fmt"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/config"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// plugin for test
type plugin struct {
	logger log.Logger

	cfg config.Cfg
	db  *dbService
}

func (t *plugin) Init(ctx types.Context) error {
	t.logger.Info("plugin init", "name", types.PluginName)
	t.db = NewDB(t.cfg, ctx.Logger().With("module", "his-database"))
	return nil
}

func (t *plugin) Start(ctx types.Context) error {
	t.logger.Info("plugin start", "name", types.PluginName)

	if err := t.db.Start(); err != nil {
		return err
	}

	return nil
}

func (t *plugin) Stop(ctx types.Context) error {
	t.logger.Info("plugin stop", "name", types.PluginName)

	if err := t.db.Stop(); err != nil {
		return err
	}

	return nil
}

func (t *plugin) MsgHandler() types.PluginMsgHandler {
	return func(ctx types.Context, msg sdk.Msg) {
		t.OnMsg(ctx, msg)
	}
}

func (t *plugin) TxHandler() types.PluginTxHandler {
	return func(ctx types.Context, tx chainTypes.StdTx) {
		t.OnTx(ctx, tx)
	}
}

func (t *plugin) EvtHandler() types.PluginEvtHandler {
	return func(ctx types.Context, evt types.Event) {
		t.OnEvent(ctx, evt)
	}
}

func (t *plugin) Logger() log.Logger {
	return t.logger
}

func (t *plugin) Name() string {
	return types.PluginName
}

// New new plugin
func New(ctx types.Context, cfg types.BaseCfg) *plugin {
	logger := ctx.Logger().With("module", fmt.Sprintf("plugins/%s", types.PluginName))

	res := &plugin{
		logger: logger,
	}

	if err := json.Unmarshal(cfg.CfgRaw, &res.cfg); err != nil {
		panic(err)
	}

	logger.Info("new plugin", "name", types.PluginName, "cfg", res.cfg)

	return res
}
