package mint

import (
	"encoding/json"
	"math/rand"

	"github.com/KuChainNetwork/kuchain/chain/genesis"
	"github.com/KuChainNetwork/kuchain/x/mint/client/cli"
	"github.com/KuChainNetwork/kuchain/x/mint/client/rest"
	"github.com/KuChainNetwork/kuchain/x/mint/simulation"
	"github.com/KuChainNetwork/kuchain/x/mint/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the mint module.
type AppModuleBasic struct {
	genesis.ModuleBasicBase
}

func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{
		ModuleBasicBase: genesis.NewModuleBasicBase(Cdc(), DefaultGenesisState()),
	}
}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the mint module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the mint module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {}

// RegisterRESTRoutes registers the REST routes for the mint module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns no root tx command for the mint module.
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command { return nil }

// GetQueryCmd returns the root query command for the mint module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

// AppModule implements an application module for the mint module.
type AppModule struct {
	AppModuleBasic

	keeper       Keeper
	supplyKeeper types.SupplyKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, supplyKeeper types.SupplyKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		supplyKeeper:   supplyKeeper,
	}
}

// Name returns the mint module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the mint module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the mint module.
func (AppModule) Route() string { return "" }

// NewHandler returns an sdk.Handler for the mint module.
func (am AppModule) NewHandler() sdk.Handler { return nil }

// QuerierRoute returns the mint module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the mint module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the mint module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)

	InitGenesis(ctx, am.keeper, am.supplyKeeper, genesisState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the mint
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return ModuleCdc.MustMarshalJSON(ExportGenesis(ctx, am.keeper))
}

// BeginBlock returns the begin blocker for the mint module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock returns the end blocker for the mint module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the mint module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized mint param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []sim.ParamChange {
	return simulation.ParamChanges(r)
}

// RegisterStoreDecoder registers a decoder for mint module's types.
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[StoreKey] = simulation.DecodeStore
}

// WeightedOperations doesn't return any mint module operation.
func (AppModule) WeightedOperations(_ module.SimulationState) []sim.WeightedOperation {
	return nil
}
