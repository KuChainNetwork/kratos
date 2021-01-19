package client

import (
	"github.com/KuChainNetwork/kuchain/chain/transaction"
	"github.com/KuChainNetwork/kuchain/chain/types"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
	"github.com/cosmos/cosmos-sdk/client"
)

type (
	StdTx     = types.StdTx
	TxBuilder = transaction.TxBuilder
	AccountID = types.AccountID
	Name      = types.Name
)

var (
	ConfigCmd      = client.ConfigCmd
	Paginate       = client.Paginate
	ValidateCmd    = client.ValidateCmd
	RegisterRoutes = client.RegisterRoutes
)

func NewAccountRetriever(cliCtx Context) accountTypes.AccountRetriever {
	return accountTypes.NewAccountRetriever(cliCtx.Ctx())
}
