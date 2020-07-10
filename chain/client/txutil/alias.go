package txutil

import (
	"github.com/KuChainNetwork/kuchain/chain/transaction"
	"github.com/KuChainNetwork/kuchain/chain/types"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	StdTx     = types.StdTx
	TxBuilder = transaction.TxBuilder
	AccountID = types.AccountID
	Name      = types.Name
)

var (
	NewStdTx            = types.NewStdTx
	NewTxBuilder        = transaction.NewTxBuilder
	DefaultTxDecoder    = types.DefaultTxDecoder
	DefaultTxEncoder    = types.DefaultTxEncoder
	NewTxBuilderFromCLI = transaction.NewTxBuilderFromCLI
)

// NewAccountRetriever initialises a new AccountRetriever instance.
func NewAccountRetriever(cliCtx KuCLIContext) accountTypes.AccountRetriever {
	return accountTypes.NewAccountRetriever(cliCtx.CLIContext)
}

// GetSignBytes returns the signBytes of the tx for a given signer
func GetSignBytes(ctx sdk.Context, tx *StdTx, accNum, seq uint64) []byte {
	chainID := ctx.ChainID()

	return types.StdSignBytes(
		chainID, accNum, seq, tx.Fee, tx.Msgs, tx.Memo,
	)
}
