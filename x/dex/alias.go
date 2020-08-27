package dex

import (
	"github.com/KuChainNetwork/kuchain/x/dex/keeper"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
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

type (
	Keeper = keeper.DexKeeper
)

var (
	NewKeeper = keeper.NewDexKeeper
)
