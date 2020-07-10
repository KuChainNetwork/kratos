package client

import (
	"github.com/KuChainNetwork/kuchain/x/distribution/client/cli"
	"github.com/KuChainNetwork/kuchain/x/distribution/client/rest"
	"github.com/KuChainNetwork/kuchain/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = client.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
