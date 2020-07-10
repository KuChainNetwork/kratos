package external

import (
	"github.com/KuChainNetwork/kuchain/x/gov/client"
	"github.com/KuChainNetwork/kuchain/x/gov/client/rest"
	"github.com/KuChainNetwork/kuchain/x/gov/types"
	costype "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type GovHandler = types.Handler
type GovContent = types.Content

var GovNewProposalHandler = client.NewProposalHandler
var GovNewMsgSubmitProposal = types.NewKuMsgSubmitProposal

type GovProposalRESTHandler = rest.ProposalRESTHandler
type CosGovContent = costype.Content

var GovRegisterProposalType = types.RegisterProposalType
var GovRegisterProposalTypeCodec = types.RegisterProposalTypeCodec
var GovValidateAbstract = types.ValidateAbstract
