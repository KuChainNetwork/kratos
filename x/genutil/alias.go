package genutil

import (
	"github.com/KuChainNetwork/kuchain/x/genutil/types"
)

const (
	ModuleName = types.ModuleName
)

var (
	NewGenesisState             = types.NewGenesisState
	NewGenesisStateFromStdTx    = types.NewGenesisStateFromStdTx
	NewInitConfig               = types.NewInitConfig
	GetGenesisStateFromAppState = types.GetGenesisStateFromAppState
	SetGenesisStateInAppState   = types.SetGenesisStateInAppState
	GenesisStateFromGenDoc      = types.GenesisStateFromGenDoc
	GenesisStateFromGenFile     = types.GenesisStateFromGenFile
	ValidateGenesis             = types.ValidateGenesis

	ModuleCdc = types.ModuleCdc
)

type (
	GenesisState      = types.GenesisState
	AppMap            = types.AppMap
	MigrationCallback = types.MigrationCallback
	MigrationMap      = types.MigrationMap
	InitConfig        = types.InitConfig
)
