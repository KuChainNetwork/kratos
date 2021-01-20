package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func OutputPrettifyJSON(cliCtx client.Context, obj interface{}) ([]byte, error) {
	return outputJSON(cliCtx, obj)
}

func outputJSON(cliCtx client.Context, obj interface{}) ([]byte, error) {
	var (
		resp []byte
		err  error
	)

	if cliCtx.Indent() {
		resp, err = cliCtx.Codec().MarshalJSONIndent(obj, "", "  ")
	} else {
		if prettifier, ok := obj.(types.Prettifier); ok {
			resp, err = prettifier.PrettifyJSON(cliCtx.Codec())
		} else {
			resp, err = cliCtx.Codec().MarshalJSON(obj)
		}
	}

	return resp, err
}

// ParseQueryHeightOrReturnBadRequest sets the height to execute a query if set by the http request.
// It returns false if there was an error parsing the height.
func ParseQueryHeightOrReturnBadRequest(w http.ResponseWriter, cliCtx client.Context, r *http.Request) (client.Context, bool) {
	heightStr := r.FormValue("height")
	if heightStr != "" {
		height, err := strconv.ParseInt(heightStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return cliCtx, false
		}

		if height < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "height must be equal or greater than zero")
			return cliCtx, false
		}

		if height > 0 {
			cliCtx = cliCtx.WithHeight(height)
		}
	} else {
		cliCtx = cliCtx.WithHeight(0)
	}

	return cliCtx, true
}

// PostProcessResponseBare post processes a body similar to PostProcessResponse
// except it does not wrap the body and inject the height.
func PostProcessResponseBare(w http.ResponseWriter, cliCtx client.Context, body interface{}) {
	var (
		resp []byte
		err  error
	)

	switch b := body.(type) {
	case []byte:
		resp = b

	default:
		if cliCtx.Indent() {
			resp, err = cliCtx.Codec().MarshalJSONIndent(body, "", "  ")
		} else {
			resp, err = cliCtx.Codec().MarshalJSON(body)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

// PostProcessResponse performs post processing for a REST response. The result
// returned to clients will contain two fields, the height at which the resource
// was queried at and the original result.
func PostProcessResponse(w http.ResponseWriter, cliCtx client.Context, resp interface{}) {
	var result []byte

	if cliCtx.GetHeight() < 0 {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("negative height in response").Error())
		return
	}

	switch res := resp.(type) {
	case []byte:
		result = res

	default:
		var err error
		result, err = outputJSON(cliCtx, resp)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	wrappedResp := rest.NewResponseWithHeight(cliCtx.GetHeight(), result)

	var (
		output []byte
		err    error
	)

	if cliCtx.Indent() {
		output, err = cliCtx.Codec().MarshalJSONIndent(wrappedResp, "", "  ")
	} else {
		output, err = cliCtx.Codec().MarshalJSON(wrappedResp)
	}
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(output)
}
