package server

import "github.com/cosmos/cosmos-sdk/server"

type (
	Context = server.Context
)

var (
	ShowNodeIDCmd     = server.ShowNodeIDCmd
	ShowValidatorCmd  = server.ShowValidatorCmd
	ShowAddressCmd    = server.ShowAddressCmd
	VersionCmd        = server.VersionCmd
	UnsafeResetAllCmd = server.UnsafeResetAllCmd
)

var (
	TrapSignal = server.TrapSignal
)

var (
	GetPruningOptionsFromFlags = server.GetPruningOptionsFromFlags
)

// Tendermint full-node start flags
const (
	flagWithTendermint     = "with-tendermint"
	flagAddress            = "address"
	flagTraceStore         = "trace-store"
	flagCPUProfile         = "cpu-profile"
	FlagMinGasPrices       = "minimum-gas-prices"
	FlagHaltHeight         = "halt-height"
	FlagHaltTime           = "halt-time"
	FlagInterBlockCache    = "inter-block-cache"
	FlagUnsafeSkipUpgrades = "unsafe-skip-upgrades"
	FlagTrace              = "trace"

	FlagPruning           = "pruning"
	FlagPruningKeepRecent = "pruning-keep-recent"
	FlagPruningKeepEvery  = "pruning-keep-every"
	FlagPruningInterval   = "pruning-interval"
)
