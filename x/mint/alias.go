package mint

// nolint

import (
	"github.com/KuChainNetwork/kuchain/x/mint/keeper"
	"github.com/KuChainNetwork/kuchain/x/mint/types"
)

const (
	ModuleName            = types.ModuleName
	DefaultParamspace     = types.DefaultParamspace
	StoreKey              = types.StoreKey
	QuerierRoute          = types.QuerierRoute
	QueryParameters       = types.QueryParameters
	QueryInflation        = types.QueryInflation
	QueryAnnualProvisions = types.QueryAnnualProvisions
)

var (
	NewKeeper            = keeper.NewKeeper
	NewQuerier           = keeper.NewQuerier
	NewGenesisState      = types.NewGenesisState
	DefaultGenesisState  = types.DefaultGenesisState
	ValidateGenesis      = types.ValidateGenesis
	NewMinter            = types.NewMinter
	InitialMinter        = types.InitialMinter
	DefaultInitialMinter = types.DefaultInitialMinter
	ValidateMinter       = types.ValidateMinter
	ParamKeyTable        = types.ParamKeyTable
	NewParams            = types.NewParams
	DefaultParams        = types.DefaultParams

	ModuleCdc              = types.ModuleCdc
	MinterKey              = types.MinterKey
	KeyMintDenom           = types.KeyMintDenom
	KeyInflationRateChange = types.KeyInflationRateChange
	KeyInflationMax        = types.KeyInflationMax
	KeyInflationMin        = types.KeyInflationMin
	KeyGoalBonded          = types.KeyGoalBonded
	KeyBlocksPerYear       = types.KeyBlocksPerYear

	Cdc = types.Cdc
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Minter       = types.Minter
	Params       = types.Params
)
