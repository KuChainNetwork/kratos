package params

// nolint

import (
	"github.com/KuChainNetwork/kuchain/x/params/keeper"
	"github.com/KuChainNetwork/kuchain/x/params/types"
)

const (
	StoreKey  = types.StoreKey
	TStoreKey = types.TStoreKey
)

var (
	NewKeeper       = keeper.NewKeeper
	NewKeyTable     = types.NewKeyTable
	NewParamSetPair = types.NewParamSetPair
)

type (
	Keeper           = keeper.Keeper
	ParamSetPair     = types.ParamSetPair
	ParamSetPairs    = types.ParamSetPairs
	ParamSet         = types.ParamSet
	Subspace         = types.Subspace
	ReadOnlySubspace = types.ReadOnlySubspace
	KeyTable         = types.KeyTable
)
