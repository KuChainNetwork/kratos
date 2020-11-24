package client

import (
	"github.com/KuChainNetwork/kuchain/x/params/client/cli"
	"github.com/KuChainNetwork/kuchain/x/params/client/rest"
	"github.com/KuChainNetwork/kuchain/x/params/external"
)

// ProposalHandler param change proposal handler
var ProposalHandler = external.GovNewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
