package genutil

import (
	types2 "github.com/KuChainNetwork/kuchain/x/genutil/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis - initialize accounts and deliver genesis transactions
func InitGenesis(
	ctx sdk.Context, cdc *codec.LegacyAmino, stakingKeeper types2.StakingKeeper,
	deliverTx DeliverTxfn, genesisState GenesisState,
) []abci.ValidatorUpdate {
	var validators []abci.ValidatorUpdate
	if len(genesisState.GenTxs) > 0 {
		validators = DeliverGenTxs(ctx, cdc, genesisState.GenTxs, stakingKeeper, deliverTx)
	}

	return validators
}
