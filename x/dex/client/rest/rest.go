package rest

import (
	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers the dex module REST routes.
func RegisterRoutes(cliCtx client.Context, r *mux.Router, storeName string) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
