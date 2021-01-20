package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/client/utils"
	"github.com/KuChainNetwork/kuchain/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/evidence/{%s}", RestParamEvidenceHash),
		queryEvidenceHandler(cliCtx),
	).Methods(MethodGet)

	r.HandleFunc(
		"/evidence",
		queryAllEvidenceHandler(cliCtx),
	).Methods(MethodGet)

	r.HandleFunc(
		"/evidence/params",
		queryParamsHandler(cliCtx),
	).Methods(MethodGet)
}

func queryEvidenceHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		evidenceHash := vars[RestParamEvidenceHash]

		if strings.TrimSpace(evidenceHash) == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "evidence hash required but not specified")
			return
		}

		cliCtx, ok := utils.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryEvidenceParams(evidenceHash)
		bz, err := cliCtx.Codec().MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to marshal query params: %s", err))
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryEvidence)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		utils.PostProcessResponse(w, cliCtx, res)
	}
}

func queryAllEvidenceHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := utils.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryAllEvidenceParams(page, limit)
		bz, err := cliCtx.Codec().MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to marshal query params: %s", err))
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAllEvidence)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		utils.PostProcessResponse(w, cliCtx, res)
	}
}

func queryParamsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := utils.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		utils.PostProcessResponse(w, cliCtx, res)
	}
}
