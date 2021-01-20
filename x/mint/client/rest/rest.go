package rest

import (
	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers minting module REST handlers on the provided router.
func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}
