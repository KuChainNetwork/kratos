package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/pkg/errors"
)

type TxGasFeeResp struct {
	EstimatedGasConsumed uint64 `json:"estimated_gas_consumed"`
	GasPrice             string `json:"gas_price"`
	EstimatedFee         string `json:"estimated_fee"`
}

func QueryTxFeeAndGasConsumed(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StdSignMsg
		var estimatedMsgsGasConsumed uint64

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

		for _, msg := range req.Msg {
			gas, err := CalculateMsgGas(msg)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			estimatedMsgsGasConsumed += gas
		}

		estimatedGas := constants.GasTxSizePrice*sdk.Gas(len(req.Bytes())) + estimatedMsgsGasConsumed
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

func CalculateMsgGas(msg sdk.Msg) (uint64, error) {
	switch msg.Type() {
	case "create@account":
		return constants.EstimatedGasCreateAcc, nil
	case "updateauth":
		return constants.EstimatedGasUpAuth, nil
	case "transfer":
		return constants.EstimatedGasTransfer, nil
	case "create@asset":
		return constants.EstimatedGasCreateCoin, nil
	case "issue":
		return constants.EstimatedGasIssueCoin, nil
	case "lock@coin":
		return constants.EstimatedGasLockCoin, nil
	case "unlock@coin":
		return constants.EstimatedGasUnlockCoin, nil
	case "burn":
		return constants.EstimatedGasBurnCoin, nil
	case "delegate":
		return constants.EstimatedGasDelegate, nil
	case "begin_redelegate":
		return constants.EstimatedGasReDelegate, nil
	case "create@staking":
		return constants.EstimatedGasCreateVal, nil
	case "edit@staking":
		return constants.EstimatedGasEditVal, nil
	case "beginunbonding":
		return constants.EstimatedGasUnBonding, nil
	case "unjail":
		return constants.EstimatedGasUnJail, nil
	case "deposit":
		return constants.EstimatedGasDeposit, nil
	case "vote":
		return constants.EstimatedGasVote, nil
	case "submitproposal":
		return constants.EstimatedGasProposal, nil
	case "withdrawcccid":
		return constants.EstimatedGasSetWithdrawAddr, nil
	case "withdrawdelreward":
	case "withdrawvalcom":
		return constants.EstimatedGasRewards, nil
	}

	return 0, errors.Errorf("Unknown msg type: %s", msg.Type())
}
