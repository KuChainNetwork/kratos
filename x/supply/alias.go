package supply

// DONTCOVER
// nolint

import (
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
	"github.com/KuChainNetwork/kuchain/x/supply/keeper"
	"github.com/KuChainNetwork/kuchain/x/supply/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	RouterKey    = types.RouterKey
	QuerierRoute = types.QuerierRoute
	Minter       = accountTypes.Minter
	Burner       = accountTypes.Burner
	Staking      = accountTypes.Staking
	BlackHole    = types.BlackHole
)

var (
	// functions aliases
	RegisterInvariants    = keeper.RegisterInvariants
	AllInvariants         = keeper.AllInvariants
	TotalSupply           = keeper.TotalSupply
	NewKeeper             = keeper.NewKeeper
	NewQuerier            = keeper.NewQuerier
	NewModuleAddress      = accountTypes.NewModuleAddress
	NewEmptyModuleAccount = accountTypes.NewEmptyModuleAccount
	NewModuleAccount      = accountTypes.NewModuleAccount
	RegisterCodec         = types.RegisterCodec
	NewGenesisState       = types.NewGenesisState
	DefaultGenesisState   = types.DefaultGenesisState
	NewSupply             = types.NewSupply
	DefaultSupply         = types.DefaultSupply

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper        = keeper.Keeper
	ModuleAccount = types.ModuleAccount
	GenesisState  = types.GenesisState
	Supply        = types.Supply
)
