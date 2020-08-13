package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers the auth module REST routes.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(
		"/assets/coins/{account}",
		getCoinsHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/assets/coin_powers/{account}",
		getCoinPowersHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/assets/coins_locked/{account}",
		getCoinsLockedHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/assets/coin_stat/{creator}/{symbol}",
		getCoinStatHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/assets/transfer",
		TransferRequestHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/assets/create",
		CreateRequestHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/assets/issue",
		IssueRequestHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/assets/burn",
		BurnRequestHandlerFn(cliCtx),
	)
	r.HandleFunc(
		"/assets/lock",
		LockRequestHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/assets/unlock",
		UnlockRequestHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/assets/exercise",
		ExerciseRequestHandlerFn(cliCtx),
	).Methods("POST")
}
