package plugin

import (
	"github.com/KuChainNetwork/kuchain/x/plugin/types"
)

const (
	ModuleName   = types.ModuleName
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
)

var (
	ModuleCdc           = types.ModuleCdc
	DefaultGenesisState = types.DefaultGenesisState
)

var (
	NewGenesisState = types.NewGenesisState
	Logger          = types.Logger
)

type (
	GenesisState = types.GenesisState
)
