package supply

import (
	"encoding/json"
	"math/rand"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/genesis"
	"github.com/KuChainNetwork/kuchain/x/supply/client/cli"
	"github.com/KuChainNetwork/kuchain/x/supply/client/rest"
	"github.com/KuChainNetwork/kuchain/x/supply/simulation"
	"github.com/KuChainNetwork/kuchain/x/supply/types"
	clientSDK "github.com/cosmos/cosmos-sdk/client"
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

// AppModuleBasic defines the basic application module used by the supply module.
type AppModuleBasic struct {
	genesis.ModuleBasicBase
}

// Name returns the supply module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{
		ModuleBasicBase: genesis.NewModuleBasicBase(ModuleCdc, DefaultGenesisState()),
	}
}

// RegisterCodec registers the supply module's typesxx for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {
	RegisterCodec(cdc)
}

// RegisterRESTRoutes registers the REST routes for the supply module.
func (AppModuleBasic) RegisterRESTRoutes(ctx clientSDK.Context, rtr *mux.Router) {
	rest.RegisterRoutes(client.NewKuCLICtx(ctx), rtr)
}

// GetTxCmd returns the root tx command for the supply module.
func (AppModuleBasic) GetTxCmd(_ *codec.LegacyAmino) *cobra.Command { return nil }

// GetQueryCmd returns no root query command for the supply module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.LegacyAmino) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

// AppModule implements an application module for the supply module.
type AppModule struct {
	AppModuleBasic

	keeper Keeper
	bk     types.BankKeeper
	ak     types.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, bk types.BankKeeper, ak types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		bk:             bk,
		ak:             ak,
	}
}

// Name returns the supply module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the supply module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	RegisterInvariants(ir, am.keeper)
}

// Route returns the message routing key for the supply module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the supply module.
func (am AppModule) NewHandler() sdk.Handler { return nil }

// QuerierRoute returns the supply module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the supply module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the supply module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.bk, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the supply
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return ModuleCdc.MustMarshalJSON(ExportGenesis(ctx, am.keeper))
}

// BeginBlock returns the begin blocker for the supply module.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the supply module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the supply module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalContent {
	return nil
}

// RandomizedParams doesn't create any randomized supply param changes for the simulator.
func (AppModule) RandomizedParams(_ *rand.Rand) []sim.ParamChange {
	return nil
}

// RegisterStoreDecoder registers a decoder for supply module's typesxx
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[StoreKey] = simulation.DecodeStore
}

// WeightedOperations doesn't return any operation for the supply module.
func (AppModule) WeightedOperations(_ module.SimulationState) []sim.WeightedOperation {
	return nil
}
