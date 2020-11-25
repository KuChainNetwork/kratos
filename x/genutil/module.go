package genutil

import (
	"encoding/json"

	"github.com/KuChainNetwork/kuchain/chain/genesis"
	"github.com/KuChainNetwork/kuchain/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModuleGenesis = AppModule{}
	_ module.AppModuleBasic   = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the genutil module.
type AppModuleBasic struct {
	genesis.ModuleBasicBase

	stakingFuncManager types.StakingFuncManager
}

func NewAppModuleBasic(stakingFuncManager types.StakingFuncManager) AppModuleBasic {
	return AppModuleBasic{
		ModuleBasicBase:    genesis.NewModuleBasicBase(ModuleCdc, types.DefaultGenesisState(stakingFuncManager)),
		stakingFuncManager: stakingFuncManager,
	}
}

// Name returns the genutil module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the genutil module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {}

// RegisterRESTRoutes registers the REST routes for the genutil module.
func (AppModuleBasic) RegisterRESTRoutes(_ context.CLIContext, _ *mux.Router) {}

// GetTxCmd returns no root tx command for the genutil module.
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command { return nil }

// GetQueryCmd returns no root query command for the genutil module.
func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command { return nil }

// AppModule implements an application module for the genutil module.
type AppModule struct {
	AppModuleBasic

	accountKeeper types.AccountKeeper
	stakingKeeper types.StakingKeeper
	deliverTx     DeliverTxfn
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountKeeper types.AccountKeeper, stakingKeeper types.StakingKeeper,
	deliverTx DeliverTxfn, stakingFuncManager types.StakingFuncManager) module.AppModule {
	return module.NewGenesisOnlyAppModule(AppModule{
		AppModuleBasic: NewAppModuleBasic(stakingFuncManager),
		accountKeeper:  accountKeeper,
		stakingKeeper:  stakingKeeper,
		deliverTx:      deliverTx,
	})
}

// InitGenesis performs genesis initialization for the genutil module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, ModuleCdc, am.stakingKeeper, am.deliverTx, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the genutil
// module.
func (am AppModule) ExportGenesis(_ sdk.Context) json.RawMessage {
	return am.DefaultGenesis()
}
