package appcreator

import (
	"github.com/KuChainNetwork/kuchain/x/account"
	"github.com/KuChainNetwork/kuchain/x/asset"
	"github.com/KuChainNetwork/kuchain/x/dex"
	distr "github.com/KuChainNetwork/kuchain/x/distribution"
	"github.com/KuChainNetwork/kuchain/x/evidence"
	"github.com/KuChainNetwork/kuchain/x/genutil"
	"github.com/KuChainNetwork/kuchain/x/gov"
	"github.com/KuChainNetwork/kuchain/x/mint"
	"github.com/KuChainNetwork/kuchain/x/params"
	"github.com/KuChainNetwork/kuchain/x/plugin"
	"github.com/KuChainNetwork/kuchain/x/slashing"
	"github.com/KuChainNetwork/kuchain/x/staking"
	"github.com/KuChainNetwork/kuchain/x/supply"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

type KVStoreKeys map[string]*sdk.KVStoreKey

func GenStoreKeys() (map[string]*sdk.KVStoreKey, map[string]*sdk.TransientStoreKey) {
	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, staking.StoreKey,
		slashing.StoreKey, evidence.StoreKey, gov.StoreKey,
		account.StoreKey, asset.StoreKey,
		supply.StoreKey, params.StoreKey,
		mint.StoreKey, distr.StoreKey, params.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(params.TStoreKey, staking.TStoreKey, params.TStoreKey)

	return keys, tKeys
}

func GenAppSubspace(cdc *codec.Codec, key sdk.StoreKey, tKey sdk.StoreKey) (map[string]params.Subspace, params.Keeper) {
	paramsKeeper := params.NewKeeper(cdc, key, tKey)

	subspaces := make(map[string]params.Subspace)
	subspaces[account.ModuleName] = paramsKeeper.Subspace(account.DefaultParamspace)
	subspaces[distr.ModuleName] = paramsKeeper.Subspace(distr.DefaultParamspace)
	subspaces[staking.ModuleName] = paramsKeeper.Subspace(staking.DefaultParamspace)
	subspaces[slashing.ModuleName] = paramsKeeper.Subspace(slashing.DefaultParamspace)
	subspaces[evidence.ModuleName] = paramsKeeper.Subspace(evidence.DefaultParamspace)
	subspaces[mint.ModuleName] = paramsKeeper.Subspace(mint.DefaultParamspace)
	subspaces[gov.ModuleName] = paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())

	return subspaces, paramsKeeper
}

func GenAppModules(app KuAppWithKeeper) *module.Manager {
	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	res := module.NewManager(
		account.NewAppModule(app.AccountKeeper(), app.AssetKeeper()),
		genutil.NewAppModule(app.AccountKeeper(), app.StakingKeeper(), app.GetDeliverTx(), app.GetStakingFuncMng()),
		asset.NewAppModule(app.AccountKeeper(), app.AssetKeeper()),
		supply.NewAppModule(app.SupplyKeeper(), app.AssetKeeper(), app.AccountKeeper()),
		distr.NewAppModule(app.DistrKeeper(), app.AccountKeeper(), app.AssetKeeper(), app.SupplyKeeper(), app.StakingKeeper()),
		slashing.NewAppModule(app.SlashingKeeper(), app.AccountKeeper(), app.AssetKeeper(), app.StakingKeeper()),
		staking.NewAppModule(*app.StakingKeeper(), app.AccountKeeper(), app.AssetKeeper(), app.SupplyKeeper()),
		mint.NewAppModule(app.MintKeeper(), app.SupplyKeeper()),
		evidence.NewAppModule(app.EvidenceKeeper(), app.AccountKeeper(), app.AssetKeeper()),
		gov.NewAppModule(app.GovKeeper(), app.AccountKeeper(), app.AssetKeeper(), app.SupplyKeeper()),
		dex.NewAppModule(app.AccountKeeper(), app.AssetKeeper(), app.SupplyKeeper(), app.DexKeeper()),
		plugin.NewAppModule(),
	)

	// plugin.ModuleName MUST be the last
	res.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName, evidence.ModuleName, dex.ModuleName, plugin.ModuleName)
	res.SetOrderEndBlockers(staking.ModuleName, gov.ModuleName, dex.ModuleName, plugin.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	res.SetOrderInitGenesis(
		account.ModuleName, asset.ModuleName, distr.ModuleName,
		staking.ModuleName, slashing.ModuleName, evidence.ModuleName,
		gov.ModuleName, supply.ModuleName, genutil.ModuleName,
		mint.ModuleName, dex.ModuleName,
	)

	res.RegisterRoutes(app.Router(), app.QueryRouter())

	return res
}

func GenAppSimulationMng(app KuAppWithKeeper) *module.SimulationManager {
	// NOTE: This is not required for apps that don't use the simulator for fuzz testing
	// transactions.
	sm := module.NewSimulationManager(
		account.NewAppModule(app.AccountKeeper(), app.AssetKeeper()),
		supply.NewAppModule(app.SupplyKeeper(), app.AssetKeeper(), app.AccountKeeper()),
		distr.NewAppModule(app.DistrKeeper(), app.AccountKeeper(), app.AssetKeeper(), app.SupplyKeeper(), app.StakingKeeper()),
		staking.NewAppModule(*app.StakingKeeper(), app.AccountKeeper(), app.AssetKeeper(), app.SupplyKeeper()),
		slashing.NewAppModule(app.SlashingKeeper(), app.AccountKeeper(), app.AssetKeeper(), app.StakingKeeper()),
		mint.NewAppModule(app.MintKeeper(), app.SupplyKeeper()),
		gov.NewAppModule(app.GovKeeper(), app.AccountKeeper(), app.AssetKeeper(), app.SupplyKeeper()),
	)

	sm.RegisterStoreDecoders()

	return sm
}
