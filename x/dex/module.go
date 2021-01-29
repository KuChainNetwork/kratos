package dex

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	clientSDK "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/simulation"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/genesis"
	"github.com/KuChainNetwork/kuchain/chain/msg"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	assetKeeper "github.com/KuChainNetwork/kuchain/x/asset/keeper"
	"github.com/KuChainNetwork/kuchain/x/dex/client/cli"
	"github.com/KuChainNetwork/kuchain/x/dex/client/rest"
	"github.com/KuChainNetwork/kuchain/x/dex/keeper"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the asset module.
type AppModuleBasic struct {
	genesis.ModuleBasicBase
}

func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{
		ModuleBasicBase: genesis.NewModuleBasicBase(ModuleCdc, DefaultGenesisState()),
	}
}

// Name returns the asset module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the asset module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// RegisterRESTRoutes registers the REST routes for the asset module.
func (AppModuleBasic) RegisterRESTRoutes(ctx clientSDK.Context, rtr *mux.Router) {
	rest.RegisterRoutes(client.NewKuCLICtx(ctx), rtr, types.StoreKey)
}

// GetTxCmd returns the root tx command for the asset module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns the root query command for the asset module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

// AppModule implements an application module for the asset module.
type AppModule struct {
	AppModuleBasic

	dexKeeper     Keeper
	assetKeeper   assetKeeper.AssetKeeper
	accountAuther chainTypes.AccountAuther
	supplyKeeper  types.SupplyKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountAuther chainTypes.AccountAuther,
	assetKeeper assetKeeper.AssetKeeper,
	supplyKeeper types.SupplyKeeper,
	dexKeeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(),
		assetKeeper:    assetKeeper,
		accountAuther:  accountAuther,
		dexKeeper:      dexKeeper,
		supplyKeeper:   supplyKeeper,
	}
}

// Name returns the asset module's name.
func (AppModule) Name() string { return types.ModuleName }

// RegisterInvariants performs a no-op.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the asset module.
func (AppModule) Route() string { return RouterKey }

// NewHandler returns an sdk.Handler for the asset module.
func (am AppModule) NewHandler() sdk.Handler {
	return msg.WarpHandler(am.assetKeeper, am.accountAuther, NewHandler(am.dexKeeper))
}

// QuerierRoute returns the asset module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the asset module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.dexKeeper)
}

// InitGenesis performs genesis initialization for the asset module. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	// check if the module account exists
	ctx.Logger().Info("genesis module account", "name", types.ModuleName)
	if err := am.supplyKeeper.InitModuleAccount(ctx, types.ModuleName); err != nil {
		panic(err)
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the asset module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(GenesisState{})
}

// BeginBlock returns the begin blocker for the asset module.
func (AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req)
}

// EndBlock returns the end blocker for the asset module. It returns no validator updates.
func (AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, req)
}

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the asset module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized auth param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []sim.ParamChange {
	return simulation.ParamChanges(r)
}

// RegisterStoreDecoder registers a decoder for asset module's types
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
}

// WeightedOperations doesn't return any asset module operation.
func (AppModule) WeightedOperations(_ module.SimulationState) []sim.WeightedOperation {
	return nil
}
