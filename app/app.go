package app

import (
	"io"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/KuChainNetwork/kuchain/chain/ante"
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/fee"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/account"
	"github.com/KuChainNetwork/kuchain/x/asset"
	"github.com/KuChainNetwork/kuchain/x/dex"
	distr "github.com/KuChainNetwork/kuchain/x/distribution"
	"github.com/KuChainNetwork/kuchain/x/evidence"
	"github.com/KuChainNetwork/kuchain/x/genutil"
	"github.com/KuChainNetwork/kuchain/x/gov"
	"github.com/KuChainNetwork/kuchain/x/mint"
	"github.com/KuChainNetwork/kuchain/x/params"
	paramsclient "github.com/KuChainNetwork/kuchain/x/params/client"
	paramproposal "github.com/KuChainNetwork/kuchain/x/params/types/proposal"
	"github.com/KuChainNetwork/kuchain/x/plugin"
	"github.com/KuChainNetwork/kuchain/x/slashing"
	"github.com/KuChainNetwork/kuchain/x/staking"
	"github.com/KuChainNetwork/kuchain/x/supply"
)

var (
	// ModuleBasics The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		account.NewAppModuleBasic(),
		asset.NewAppModuleBasic(),
		params.AppModuleBasic{},
		distr.NewAppModuleBasic(),
		supply.NewAppModuleBasic(),
		staking.NewAppModuleBasic(),
		slashing.NewAppModuleBasic(),
		evidence.NewAppModuleBasic(),
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
		mint.NewAppModuleBasic(),
		params.NewAppModuleBasic(),
		dex.NewAppModuleBasic(),
		plugin.NewAppModuleBasic(),
	)

	// maccPerms module account permissions
	maccPerms = map[string][]string{
		fee.CollectorName:         nil,
		distr.ModuleName:          nil,
		supply.BlackHole:          nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		mint.ModuleName:           {supply.Minter},
	}
)

// Verify app interface at compile time
var _ simapp.App = (*KuchainApp)(nil)

// KuchainApp extended ABCI application
type KuchainApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	accountKeeper  account.Keeper
	assetKeeper    asset.Keeper
	supplyKeeper   supply.Keeper
	distrKeeper    distr.Keeper
	mintKeeper     mint.Keeper
	paramsKeeper   params.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	evidenceKeeper evidence.Keeper
	govKeeper      gov.Keeper
	dexKeeper      dex.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// function manager
	stakingFuncManager staking.FuncManager
}

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	chainTypes.RegisterCodec(cdc)
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	cdc.Seal()

	return cdc
}

// NewKuchainApp returns a reference to an initialized KuchainApp.
func NewKuchainApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *KuchainApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, txutil.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, staking.StoreKey, slashing.StoreKey, evidence.StoreKey, gov.StoreKey,
		account.StoreKey, asset.StoreKey, supply.StoreKey, params.StoreKey, mint.StoreKey, distr.StoreKey, params.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(params.TStoreKey, staking.TStoreKey, params.TStoreKey)

	app := &KuchainApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	app.subspaces[account.ModuleName] = app.paramsKeeper.Subspace(account.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[evidence.ModuleName] = app.paramsKeeper.Subspace(evidence.DefaultParamspace)
	app.subspaces[mint.ModuleName] = app.paramsKeeper.Subspace(mint.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())

	// add keepers
	app.accountKeeper = account.NewAccountKeeper(cdc, keys[account.StoreKey])
	app.assetKeeper = asset.NewAssetKeeper(cdc, keys[asset.StoreKey], app.accountKeeper)
	app.supplyKeeper = supply.NewKeeper(
		cdc, keys[supply.StoreKey], app.accountKeeper, app.assetKeeper, maccPerms,
	)

	stakingKeeper := staking.NewKeeper(
		app.cdc, keys[staking.StoreKey], app.assetKeeper, app.supplyKeeper, app.subspaces[staking.ModuleName], app.accountKeeper,
	)
	app.stakingFuncManager = staking.NewFuncManager()

	app.distrKeeper = distr.NewKeeper(
		cdc, keys[distr.StoreKey], app.subspaces[distr.ModuleName],
		app.assetKeeper,
		&app.stakingKeeper,
		app.supplyKeeper,
		app.accountKeeper,
		fee.CollectorName,
		app.ModuleAccountAddrs())

	app.slashingKeeper = slashing.NewKeeper(
		cdc, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName],
	)

	// create evidence keeper with evidence router
	evidenceKeeper := evidence.NewKeeper(
		keys[evidence.StoreKey], app.subspaces[evidence.ModuleName], &stakingKeeper, app.slashingKeeper,
	)
	evidenceRouter := evidence.NewRouter()

	app.evidenceKeeper = *evidenceKeeper

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper))
	app.govKeeper = gov.NewKeeper(cdc,
		keys[gov.StoreKey], app.subspaces[gov.ModuleName],
		app.supplyKeeper, &stakingKeeper, app.distrKeeper, govRouter,
	)

	app.dexKeeper = dex.NewKeeper(cdc, keys[gov.StoreKey], app.assetKeeper)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	// TODO: register evidence routes
	evidenceKeeper.SetRouter(evidenceRouter)
	app.mintKeeper = mint.NewKeeper(
		cdc, keys[mint.StoreKey], app.subspaces[mint.ModuleName], &app.stakingKeeper,
		app.supplyKeeper, constants.FeeSystemAccountStr,
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		account.NewAppModule(app.accountKeeper, app.assetKeeper),
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx, app.stakingFuncManager),
		asset.NewAppModule(app.accountKeeper, app.assetKeeper),
		supply.NewAppModule(app.supplyKeeper, app.assetKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.assetKeeper, app.supplyKeeper, app.stakingKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.assetKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.assetKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper, app.supplyKeeper),
		evidence.NewAppModule(app.evidenceKeeper, app.accountKeeper, app.assetKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.assetKeeper, app.supplyKeeper),
		dex.NewAppModule(app.accountKeeper, app.assetKeeper, app.supplyKeeper, app.dexKeeper),
		plugin.NewAppModule(),
	)

	// plugin.ModuleName MUST be the last
	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName, evidence.ModuleName, dex.ModuleName, plugin.ModuleName)
	app.mm.SetOrderEndBlockers(staking.ModuleName, gov.ModuleName, dex.ModuleName, plugin.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		account.ModuleName,
		asset.ModuleName,
		distr.ModuleName,
		staking.ModuleName,
		slashing.ModuleName, evidence.ModuleName, gov.ModuleName,
		supply.ModuleName,
		genutil.ModuleName,
		mint.ModuleName,
		dex.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: This is not required for apps that don't use the simulator for fuzz testing
	// transactions.
	app.sm = module.NewSimulationManager(
		account.NewAppModule(app.accountKeeper, app.assetKeeper),
		supply.NewAppModule(app.supplyKeeper, app.assetKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.assetKeeper, app.supplyKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.assetKeeper, app.supplyKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.assetKeeper, app.stakingKeeper),
		mint.NewAppModule(app.mintKeeper, app.supplyKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.assetKeeper, app.supplyKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	app.SetAnteHandler(ante.NewHandler(app.accountKeeper, app.assetKeeper))

	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	constants.LogVersion(app.Logger())

	return app
}

// Name returns the name of the App
func (app *KuchainApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *KuchainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *KuchainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *KuchainApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

// LoadHeight loads a particular height
func (app *KuchainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *KuchainApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns the application's sealed codec.
func (app *KuchainApp) Codec() *codec.Codec {
	return app.cdc
}

// SimulationManager implements the SimulationApp interface
func (app *KuchainApp) SimulationManager() *module.SimulationManager {
	return app.sm
}
