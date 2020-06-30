package plugin

import (
	"github.com/KuChain-io/kuchain/x/plugin/types"
)

const (
	ModuleName   = types.ModuleName
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
)

var (
	ModuleCdc = types.ModuleCdc
)

var (
	NewGenesisState = types.NewGenesisState
	Logger          = types.Logger
)

type (
	GenesisState = types.GenesisState
)
