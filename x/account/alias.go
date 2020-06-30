package account

import (
	"github.com/KuChain-io/kuchain/x/account/keeper"
	"github.com/KuChain-io/kuchain/x/account/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	RouterKey         = types.RouterKey
)

type (
	Keeper       = keeper.AccountKeeper
	GenesisState = types.GenesisState
)

var (
	NewAccountKeeper    = keeper.NewAccountKeeper
	NewQuerier          = keeper.NewQuerier
	NewKuAccount        = types.NewKuAccount
	DefaultGenesisState = types.DefaultGenesisState
	ModuleCdc           = types.ModuleCdc
)
