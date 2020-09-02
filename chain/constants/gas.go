package constants

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
)

var (
	MinGasPriceString                  = "0.01" + DefaultBondDenom
	GasTxSizePrice              uint64 = 5
	EstimatedGasCreateAcc       uint64 = 40000
	EstimatedGasTransfer        uint64 = 40000
	EstimatedGasUpAuth          uint64 = 50000
	EstimatedGasCreateCoin      uint64 = 40000
	EstimatedGasIssueCoin       uint64 = 40000
	EstimatedGasLockCoin        uint64 = 40000
	EstimatedGasUnlockCoin      uint64 = 40000
	EstimatedGasBurnCoin        uint64 = 50000
	EstimatedGasDelegate        uint64 = 100000
	EstimatedGasUnBonding       uint64 = 150000
	EstimatedGasReDelegate      uint64 = 250000
	EstimatedGasCreateVal       uint64 = 70000
	EstimatedGasEditVal         uint64 = 50000
	EstimatedGasUnJail          uint64 = 50000
	EstimatedGasProposal        uint64 = 100000
	EstimatedGasDeposit         uint64 = 80000
	EstimatedGasVote            uint64 = 100000
	EstimatedGasRewards         uint64 = 100000
	EstimatedGasSetWithdrawAddr uint64 = 40000
)

var (
	MinGasPrice types.DecCoins
)

func init() {
	if minGasPrice, err := types.ParseDecCoins(MinGasPriceString); err != nil {
		panic(err)
	} else {
		MinGasPrice = minGasPrice
	}
}
