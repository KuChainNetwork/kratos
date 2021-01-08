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

	appcreator "github.com/KuChainNetwork/kuchain/app/app_creator"
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
	"github.com/KuChainNetwork/kuchain/x/genutil/types"
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
var (
	_ simapp.App                 = (*KuchainApp)(nil)
	_ appcreator.KuAppWithKeeper = (*KuchainApp)(nil)
)

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

	keys, tKeys := appcreator.GenStoreKeys()

	app := &KuchainApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	app.subspaces, app.paramsKeeper =
		appcreator.GenAppSubspace(cdc, keys[params.StoreKey], tKeys[params.TStoreKey])

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
	evidenceKeeper.SetRouter(evidence.NewRouter())
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

	app.mintKeeper = mint.NewKeeper(
		cdc, keys[mint.StoreKey], app.subspaces[mint.ModuleName], &app.stakingKeeper,
		app.supplyKeeper, constants.FeeSystemAccountStr,
	)

	app.mm = appcreator.GenAppModules(app)

	// create the simulation manager and define the order of the modules for deterministic simulations
	app.sm = appcreator.GenAppSimulationMng(app)

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

func (app *KuchainApp) AccountKeeper() account.Keeper {
	return app.accountKeeper
}

func (app *KuchainApp) AssetKeeper() asset.Keeper {
	return app.assetKeeper
}

func (app *KuchainApp) SupplyKeeper() supply.Keeper {
	return app.supplyKeeper
}

func (app *KuchainApp) DistrKeeper() distr.Keeper {
	return app.distrKeeper
}

func (app *KuchainApp) MintKeeper() mint.Keeper {
	return app.mintKeeper
}

func (app *KuchainApp) ParamsKeeper() *params.Keeper {
	return &app.paramsKeeper
}

func (app *KuchainApp) StakingKeeper() *staking.Keeper {
	return &app.stakingKeeper
}

func (app *KuchainApp) SlashingKeeper() slashing.Keeper {
	return app.slashingKeeper
}

func (app *KuchainApp) EvidenceKeeper() evidence.Keeper {
	return app.evidenceKeeper
}

func (app *KuchainApp) GovKeeper() gov.Keeper {
	return app.govKeeper
}

func (app *KuchainApp) DexKeeper() dex.Keeper {
	return app.dexKeeper
}

func (app *KuchainApp) GetDeliverTx() appcreator.DeliverTxfn {
	return app.BaseApp.DeliverTx
}

func (app *KuchainApp) GetStakingFuncMng() types.StakingFuncManager {
	return app.stakingFuncManager
}
