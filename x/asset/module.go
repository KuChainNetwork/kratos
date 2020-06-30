package asset

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/simulation"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChain-io/kuchain/chain/msg"
	chainType "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/asset/client/cli"
	"github.com/KuChain-io/kuchain/x/asset/client/rest"
	"github.com/KuChain-io/kuchain/x/asset/keeper"
	"github.com/KuChain-io/kuchain/x/asset/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the asset module.
type AppModuleBasic struct{}

// Name returns the asset module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the asset module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the asset module.
func (AppModuleBasic) DefaultGenesis(codec.JSONMarshaler) json.RawMessage {
	return ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the asset module.
func (AppModuleBasic) ValidateGenesis(_ codec.JSONMarshaler, bz json.RawMessage) error {
	gs := types.DefaultGenesisState()
	if err := ModuleCdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	return gs.ValidateGenesis()
}

// RegisterRESTRoutes registers the REST routes for the asset module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, types.StoreKey)
}

// GetTxCmd returns the root tx command for the asset module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns the root query command for the asset module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the asset module.
type AppModule struct {
	AppModuleBasic

	assetKeeper   keeper.AssetKeeper
	accountAuther chainType.AccountAuther
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountAuther chainType.AccountAuther, assetKeeper keeper.AssetKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		assetKeeper:    assetKeeper,
		accountAuther:  accountAuther,
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
	return msg.WarpHandler(am.assetKeeper, am.accountAuther, NewHandler(am.assetKeeper))
}

// QuerierRoute returns the asset module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the asset module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.assetKeeper)
}

// InitGenesis performs genesis initialization for the asset module. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, marshaler codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.assetKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the asset module.
func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.assetKeeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the asset module.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the asset module. It returns no validator updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

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
	sdr[StoreKey] = simulation.DecodeStore
}

// WeightedOperations doesn't return any asset module operation.
func (AppModule) WeightedOperations(_ module.SimulationState) []sim.WeightedOperation {
	return nil
}
