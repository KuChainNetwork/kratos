package asset

import (
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset/keeper"
	"github.com/KuChainNetwork/kuchain/x/asset/types"
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
	NewAssetKeeper      = keeper.NewAssetKeeper
	NewGenesisState     = types.NewGenesisState
	NewGenesisCoin      = types.NewGenesisCoin
	NewGenesisAsset     = types.NewGenesisAsset
	DefaultGenesisState = types.DefaultGenesisState
)

type (
	Keeper      = keeper.AssetKeeper
	KuTransfMsg = chainTypes.KuTransfMsg

	GenesisState = types.GenesisState
	GenesisAsset = types.GenesisAsset
)
