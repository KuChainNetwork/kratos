package asset

import (
	"encoding/json"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis account genesis init
func InitGenesis(ctx sdk.Context, ak Keeper, bz json.RawMessage) {
	logger := ak.Logger(ctx)

	var data GenesisState
	ModuleCdc.MustUnmarshalJSON(bz, &data)

	logger.Debug("init genesis", "module", ModuleName, "data", data)

	for _, a := range data.GenesisCoins {
		logger.Info("init genesis asset coin", "accountID", a.GetCreator(), "coins", a.GetSymbol(), "maxsupply:", a.GetMaxSupply())

		initSupply := types.NewCoin(a.GetMaxSupply().Denom, sdk.ZeroInt())

		err := ak.Create(ctx, a.GetCreator(), a.GetSymbol(), a.GetMaxSupply(), true, true, true, 0, initSupply, []byte{}) // TODO: genesis coins support opt
		if err != nil {
			panic(err)
		}
	}

	for _, a := range data.GenesisAssets {
		logger.Info("init genesis account asset", "accountID", a.GetID(), "coins", a.GetCoins())
		err := ak.GenesisCoins(ctx, a.GetID(), a.GetCoins())
		if err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, ak Keeper) GenesisState {
	return GenesisState{}
}

// GenesisBalancesIterator implements genesis account iteration.
type GenesisBalancesIterator struct{}

// IterateGenesisBalances iterates over all the genesis accounts found in
// appGenesis and invokes a callback on each genesis account. If any call
// returns true, iteration stops.
func (GenesisBalancesIterator) IterateGenesisBalances(
	cdc *codec.Codec, appState types.AppGenesisState, cb func(GenesisAsset) (stop bool),
) {
	var gs GenesisState
	err := types.LoadGenesisStateFromBytes(cdc, appState, ModuleName, &gs)
	if err != nil {
		panic(err)
	}

	for _, a := range gs.GenesisAssets {
		if cb(a) {
			break
		}
	}
}
