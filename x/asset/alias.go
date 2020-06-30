package asset

import (
	chainTypes "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/asset/keeper"
	"github.com/KuChain-io/kuchain/x/asset/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
)

var (
	ModuleCdc = types.ModuleCdc
	Cdc       = types.Cdc
)

var (
	NewAssetKeeper  = keeper.NewAssetKeeper
	NewGenesisState = types.NewGenesisState
	NewGenesisCoin  = types.NewGenesisCoin
	NewGenesisAsset = types.NewGenesisAsset
)

type (
	Keeper      = keeper.AssetKeeper
	KuTransfMsg = chainTypes.KuTransfMsg

	GenesisState = types.GenesisState
	GenesisAsset = types.GenesisAsset
)
