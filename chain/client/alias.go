package client

import (
	"github.com/KuChainNetwork/kuchain/chain/transaction"
	"github.com/KuChainNetwork/kuchain/chain/types"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
)

type (
	StdTx     = types.StdTx
	TxBuilder = transaction.TxBuilder
	AccountID = types.AccountID
	Name      = types.Name
)

func NewAccountRetriever(cliCtx Context) accountTypes.AccountRetriever {
	return accountTypes.NewAccountRetriever(cliCtx.Ctx())
}
