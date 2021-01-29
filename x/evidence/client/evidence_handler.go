package client

import (
	"github.com/spf13/cobra"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/x/evidence/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
)

var (
	NewKuCLICtx = client.NewKuCLICtx
)

type (
	// RESTHandlerFn defines a REST service handler for evidence submission
	RESTHandlerFn func(client.Context) rest.EvidenceRESTHandler

	// CLIHandlerFn defines a CLI command handler for evidence submission
	CLIHandlerFn func(*codec.LegacyAmino) *cobra.Command

	// EvidenceHandler defines a type that exposes REST and CLI client handlers for
	// evidence submission.
	EvidenceHandler struct {
		CLIHandler  CLIHandlerFn
		RESTHandler RESTHandlerFn
	}
)

func NewEvidenceHandler(cliHandler CLIHandlerFn, restHandler RESTHandlerFn) EvidenceHandler {
	return EvidenceHandler{
		CLIHandler:  cliHandler,
		RESTHandler: restHandler,
	}
}
