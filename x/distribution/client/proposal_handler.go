package client

import (
	"github.com/KuChain-io/kuchain/x/distribution/client/cli"
	"github.com/KuChain-io/kuchain/x/distribution/client/rest"
	"github.com/KuChain-io/kuchain/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = client.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
