package plugins

import (
	"errors"
	"sync"

	"github.com/KuChainNetwork/kuchain/plugins/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

type pluginMsg interface{}

// Plugins a handler for all plugins to reg
type Plugins struct {
	plugins     []Plugin
	txHandlers  []types.PluginTxHandler
	msgHandlers []types.PluginMsgHandler
	evtHandlers []types.PluginEvtHandler

	msgChan chan pluginMsg
	closed  bool
	logger  log.Logger
	wg      sync.WaitGroup
}

func NewPlugins(logger log.Logger) *Plugins {
	return &Plugins{
		msgChan: make(chan pluginMsg, 512),
		closed:  false,
		logger:  logger,
	}
}

func (p *Plugins) RegPlugin(ctx Context, plugin Plugin) {
	plugin.Logger().Info("init plugin", "name", plugin.Name())

	for _, p := range p.plugins {
		if p.Name() == plugin.Name() {
			panic(errors.New("plugin reg two times"))
		}
	}

	if err := plugin.Init(ctx); err != nil {
		panic(err)
	}

	p.plugins = append(p.plugins, plugin)

	if tx := plugin.TxHandler(); tx != nil {
		p.txHandlers = append(p.txHandlers, tx)
	}

	if msg := plugin.MsgHandler(); msg != nil {
		p.msgHandlers = append(p.msgHandlers, msg)
	}

	if evt := plugin.EvtHandler(); evt != nil {
		p.evtHandlers = append(p.evtHandlers, evt)
	}
}

func (p *Plugins) onTx(ctx types.Context, tx StdTx) {
	for _, h := range p.txHandlers {
		h(ctx, tx)
	}

	for _, h := range p.msgHandlers {
		for _, msg := range tx.Msgs {
			h(ctx, msg)
		}
	}
}

func (p *Plugins) onEvent(ctx types.Context, evt types.Event) {
	for _, h := range p.evtHandlers {
		h(ctx, evt)
	}
}

func (p *Plugins) Start() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		ctx := types.NewContext(p.logger)
		for _, p := range p.plugins {
			if err := p.Start(ctx); err != nil {
				panic(err)
			}
		}

		for {
			msg, ok := <-p.msgChan
			if !ok {
				p.logger.Info("msg channel closed")
				return
			}

			if msg == nil {
				p.logger.Info("stop channel")
				return
			}

			ctx := NewContext(p.logger)

			switch msg := msg.(type) {
			case *types.MsgEvent:
				p.onEvent(ctx, msg.Evt)
			case *types.MsgStdTx:
				p.onTx(ctx, msg.Tx)
			}
		}
	}()
}

func (p *Plugins) EmitEvent(evt sdk.Event) {
	p.msgChan <- types.NewMsgEvent(evt)
}

func (p *Plugins) EmitTx(tx StdTx) {
	p.msgChan <- types.NewMsgStdTx(tx)
}

func (p *Plugins) Stop(ctx types.Context) {
	p.msgChan <- nil
	p.wg.Wait()

	for _, plg := range p.plugins {
		plg.Stop(ctx)
	}
}
