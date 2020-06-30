package account

import (
	"encoding/json"
	"math/rand"

	"github.com/KuChain-io/kuchain/chain/client/txutil"
	"github.com/KuChain-io/kuchain/chain/msg"
	chainTypes "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/account/client/cli"
	"github.com/KuChain-io/kuchain/x/account/client/rest"
	"github.com/KuChain-io/kuchain/x/account/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	// TODO: account module sim
	"github.com/cosmos/cosmos-sdk/x/auth/simulation"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the account module.
type AppModuleBasic struct{}

// Name returns the account module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the account module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	txutil.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the account module.
func (AppModuleBasic) DefaultGenesis(codec.JSONMarshaler) json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the account module.
func (AppModuleBasic) ValidateGenesis(codec.JSONMarshaler, json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes registers the REST routes for the account module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, types.StoreKey)
}

// GetTxCmd returns the root tx command for the account module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns the root query command for the account module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the account module.
type AppModule struct {
	AppModuleBasic

	accountKeeper Keeper
	assetTransfer chainTypes.AssetTransfer
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountKeeper Keeper, assetTransfer chainTypes.AssetTransfer) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		accountKeeper:  accountKeeper,
		assetTransfer:  assetTransfer,
	}
}

// Name returns the account module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants performs a no-op.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the account module.
func (AppModule) Route() string { return RouterKey }

// NewHandler returns an sdk.Handler for the account module.
func (am AppModule) NewHandler() sdk.Handler {
	return msg.WarpHandler(am.assetTransfer, am.accountKeeper, NewHandler(am.accountKeeper))
}

// QuerierRoute returns the account module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the account module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.accountKeeper)
}

// InitGenesis performs genesis initialization for the account module. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, marshaler codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.accountKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the account module.
func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.accountKeeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the account module.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the account module. It returns no validator updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the account module
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

// RegisterStoreDecoder registers a decoder for account module's types
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[StoreKey] = simulation.DecodeStore
}

// WeightedOperations doesn't return any account module operation.
func (AppModule) WeightedOperations(_ module.SimulationState) []sim.WeightedOperation {
	return nil
}
