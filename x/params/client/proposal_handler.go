package client

import (
	"github.com/KuChain-io/kuchain/x/params/client/cli"
	"github.com/KuChain-io/kuchain/x/params/client/rest"
	"github.com/KuChain-io/kuchain/x/params/external"
)

// param change proposal handler
var ProposalHandler = external.GovNewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
