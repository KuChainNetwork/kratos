package genutil

import (
	"encoding/json"

	"github.com/KuChain-io/kuchain/chain/client/txutil"
	"github.com/KuChain-io/kuchain/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// SetGenTxsInAppGenesisState - sets the genesis transactions in the app genesis state
func SetGenTxsInAppGenesisState(
	cdc *codec.Codec, appGenesisState map[string]json.RawMessage, genTxs []txutil.StdTx,
) (map[string]json.RawMessage, error) {

	genesisState := GetGenesisStateFromAppState(cdc, appGenesisState)
	genTxsBz := make([]json.RawMessage, 0, len(genTxs))

	for _, genTx := range genTxs {
		txBz, err := cdc.MarshalJSON(genTx)
		if err != nil {
			return appGenesisState, err
		}

		genTxsBz = append(genTxsBz, txBz)
	}

	genesisState.GenTxs = genTxsBz
	return SetGenesisStateInAppState(cdc, appGenesisState, genesisState), nil
}

// ValidateAccountInGenesis checks that the provided account has a sufficient
// balance in the set of genesis accounts.
func ValidateAccountInGenesis(
	appGenesisState map[string]json.RawMessage, genBalIterator types.GenesisBalancesIterator,
	addr sdk.Address, coins sdk.Coins, cdc *codec.Codec, stakingFuncManager types.StakingFuncManager,
) error {
	// FIXME: Just allow so that will use account module, until staking module use account module
	return nil
}

type deliverTxfn func(abci.RequestDeliverTx) abci.ResponseDeliverTx

// DeliverGenTxs iterates over all genesis txs, decodes each into a StdTx and
// invokes the provided deliverTxfn with the decoded StdTx. It returns the result
// of the staking module's ApplyAndReturnValidatorSetUpdates.
func DeliverGenTxs(
	ctx sdk.Context, cdc *codec.Codec, genTxs []json.RawMessage,
	stakingKeeper types.StakingKeeper, deliverTx deliverTxfn,
) []abci.ValidatorUpdate {

	logger := ctx.Logger()

	for _, genTx := range genTxs {
		var tx txutil.StdTx
		cdc.MustUnmarshalJSON(genTx, &tx)
		bz := cdc.MustMarshalBinaryLengthPrefixed(tx)
		logger.Debug("deliverTx", "tx", tx)
		res := deliverTx(abci.RequestDeliverTx{Tx: bz})
		if !res.IsOK() {
			panic(res.Log)
		}
	}

	return stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
}
