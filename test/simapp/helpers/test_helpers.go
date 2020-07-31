package helpers

import (
	"math/rand"
	"time"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simulation"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// SimAppChainID hardcoded chainID for simulation
const (
	DefaultGenTxGas = 1000000
	SimAppChainID   = "simulation-app"
)

// GenTx generates a signed mock transaction.
func GenTx(msgs []sdk.Msg, feeAmt types.Coins, gas uint64, payer types.AccountID, chainID string, accNums []uint64, seq []uint64, priv ...crypto.PrivKey) types.StdTx {
	fee := types.StdFee{
		Amount: feeAmt,
		Gas:    gas,
		Payer:  payer,
	}

	sigs := make([]types.StdSignature, len(priv))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	memo := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 0, 100))

	for i, p := range priv {
		// use a empty chainID for ease of testing
		sig, err := p.Sign(types.StdSignBytes(chainID, accNums[i], seq[i], fee, msgs, memo))
		if err != nil {
			panic(err)
		}

		sigs[i] = types.StdSignature{
			PubKey:    p.PubKey(),
			Signature: sig,
		}
	}

	return types.NewStdTx(msgs, fee, sigs, memo)
}
