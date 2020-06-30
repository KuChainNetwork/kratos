package plugins

import (
	chainTypes "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/plugins/types"
)

type (
	BaseCfg = types.BaseCfg
	Plugin  = types.Plugin
	Context = types.Context
)

var (
	NewContext = types.NewContext
)

type (
	StdTx = chainTypes.StdTx
)
