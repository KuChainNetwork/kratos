package rest

import (
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"io/ioutil"
	"net/http"
)

type TxGasFeeResp struct {
	EstimatedGasConsumed uint64 `json:"estimated_gas_consumed"`
	GasPrice             string `json:"gas_price"`
	EstimatedFee         string `json:"estimated_fee"`
}

func QueryTxFeeAndGasConsumed(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StdSignMsg

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		estimatedGas := constants.GasTxSizePrice*sdk.Gas(len(req.Bytes())) + constants.EstimatedGasConsumed
		gasPrices, err := types.ParseDecCoins(constants.MinGasPriceString)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		estimatedFee := gasPrices.ToSDK()[0].Amount.Mul(types.NewDec(int64(estimatedGas)))

		response := TxGasFeeResp{
			EstimatedGasConsumed: estimatedGas,
			GasPrice:             constants.MinGasPriceString,
			EstimatedFee:         estimatedFee.String(),
		}
		rest.PostProcessResponseBare(w, cliCtx, response)
	}
}
