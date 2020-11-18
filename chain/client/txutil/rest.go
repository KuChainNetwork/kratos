package txutil

import (
	"log"
	"net/http"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WriteGenerateStdTxResponse writes response for the generate only mode.
func WriteGenerateStdTxResponse(w http.ResponseWriter, cliCtx KuCLIContext, br types.BaseReq, msgs []sdk.Msg) {
	gasAdj, ok := types.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return
	}

	simAndExec, gas, err := flags.ParseGas(br.Gas)
	if err != nil {
		types.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	txBldr := NewTxBuilder(
		GetTxEncoder(cliCtx.Codec), br.AccountNumber, br.Sequence, gas, gasAdj,
		br.Simulate, br.ChainID, br.Memo, br.Fees, br.GasPrices,
	)

	txBldr = txBldr.WithPayer(br.Payer)

	if br.Simulate || simAndExec {
		if gasAdj < 0 {
			types.WriteErrorResponse(w, http.StatusBadRequest, errInvalidGasAdjustment.Error())
			return
		}

		txBldr, err = EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			types.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if br.Simulate {
			types.WriteSimulationResponse(w, cliCtx.Codec, txBldr.Gas())
			return
		}
	}

	stdMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		types.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	output, err := cliCtx.Codec.MarshalJSON(NewStdTx(stdMsg.Msg, stdMsg.Fee, nil, stdMsg.Memo))
	if err != nil {
		types.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(output); err != nil {
		log.Printf("could not write response: %v", err)
	}
}
