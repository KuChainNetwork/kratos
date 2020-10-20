package utils

import (
	"fmt"
	"net/http"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func OutputPrettifyJSON(cliCtx context.CLIContext, obj interface{}) ([]byte, error) {
	return outputJSON(cliCtx, obj)
}

func outputJSON(cliCtx context.CLIContext, obj interface{}) ([]byte, error) {
	var (
		resp []byte
		err  error
	)

	if cliCtx.Indent {
		resp, err = cliCtx.Codec.MarshalJSONIndent(obj, "", "  ")
	} else {
		if prettifier, ok := obj.(types.Prettifier); ok {
			resp, err = prettifier.PrettifyJSON(cliCtx.Codec)
		} else {
			resp, err = cliCtx.Codec.MarshalJSON(obj)
		}
	}

	return resp, err
}

// PostProcessResponseBare post processes a body similar to PostProcessResponse
// except it does not wrap the body and inject the height.
func PostProcessResponseBare(w http.ResponseWriter, cliCtx context.CLIContext, body interface{}) {
	var (
		resp []byte
		err  error
	)

	switch b := body.(type) {
	case []byte:
		resp = b

	default:
		resp, err = outputJSON(cliCtx, body)

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
func PostProcessResponse(w http.ResponseWriter, cliCtx context.CLIContext, resp interface{}) {
	var result []byte

	if cliCtx.Height < 0 {
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

	wrappedResp := rest.NewResponseWithHeight(cliCtx.Height, result)

	var (
		output []byte
		err    error
	)

	if cliCtx.Indent {
		output, err = cliCtx.Codec.MarshalJSONIndent(wrappedResp, "", "  ")
	} else {
		output, err = cliCtx.Codec.MarshalJSON(wrappedResp)
	}
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(output)
}
