package rest

import (
	"fmt"
	"net/http"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// QueryGenesisTxs writes the genesis transactions to the response if no error occurs.
func QueryGenesisTxs(cliCtx client.Context, w http.ResponseWriter) {
	_, err := cliCtx.GetClient().Genesis()
	if err != nil {
		rest.WriteErrorResponse(
			w, http.StatusInternalServerError,
			fmt.Sprintf("failed to retrieve genesis from client: %s", err),
		)
		return
	}

	genTxs := make([]sdk.Tx, 0, 64)

	// FIXME: support zero block txs
	/*
		appState, err := types.GenesisStateFromGenDoc(cliCtx.Codec, *resultGenesis.Genesis)
		if err != nil {
			rest.WriteErrorResponse(
				w, http.StatusInternalServerError,
				fmt.Sprintf("failed to decode genesis doc: %s", err),
			)
			return
		}

		genState := types.GetGenesisStateFromAppState(cliCtx.Codec, appState)

		for i, tx := range genState.GenTxs {
			err := cliCtx.Codec.UnmarshalJSON(tx, &genTxs[i])
			if err != nil {
				rest.WriteErrorResponse(
					w, http.StatusInternalServerError,
					fmt.Sprintf("failed to decode genesis transaction: %s", err),
				)
				return
			}
		}

	*/

	utils.PostProcessResponseBare(w, cliCtx, genTxs)
}
