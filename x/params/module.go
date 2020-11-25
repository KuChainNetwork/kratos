package params

import (
	"math/rand"

	"github.com/KuChainNetwork/kuchain/chain/genesis"
	"github.com/KuChainNetwork/kuchain/x/params/simulation"
	"github.com/KuChainNetwork/kuchain/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var (
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the params module.
type AppModuleBasic struct {
	genesis.EmptyGenesisModuleBasicBase
}

// Name returns the params module's name.
func (AppModuleBasic) Name() string {
	return proposal.ModuleName
}

// NewAppModuleBasic new app module basic
func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{}
}

// RegisterCodec registers the params module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	proposal.RegisterCodec(cdc)
}

// RegisterRESTRoutes registers the REST routes for the params module.
func (AppModuleBasic) RegisterRESTRoutes(_ context.CLIContext, _ *mux.Router) {}

// GetTxCmd returns no root tx command for the params module.
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command { return nil }

// GetQueryCmd returns no root query command for the params module.
func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command { return nil }

// AppModule implements an application module for the distribution module.
type AppModule struct {
	AppModuleBasic
}

// NewAppModule creates a new AppModule object
func NewAppModule() AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
	}
}

// AppModuleSimulation functions

// GenerateGenesisState performs a no-op.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
}

// ProposalContents returns all the params content functions used to
// simulate governance proposals.
func (am AppModule) ProposalContents(simState module.SimulationState) []sim.WeightedProposalContent {
	return simulation.ProposalContents(simState.ParamChanges)
}

// RandomizedParams creates randomized distribution param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []sim.ParamChange {
	return nil
}

// RegisterStoreDecoder doesn't register any type.
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(_ module.SimulationState) []sim.WeightedOperation {
	return nil
}
