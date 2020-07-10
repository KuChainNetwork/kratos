package plugins

import (
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/plugins/types"
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
