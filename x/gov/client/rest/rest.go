package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	rest "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/client/context"
)

// REST Variable names
// nolint
const (
	RestParamsType     = "type"
	RestProposalID     = "proposal-id"
	RestDepositor      = "depositor"
	RestVoter          = "voter"
	RestProposalStatus = "status"
	RestNumLimit       = "limit"
)

// ProposalRESTHandler defines a REST handler implemented in another module. The
// sub-route is mounted on the governance REST handler.
type ProposalRESTHandler struct {
	SubRoute string
	Handler  func(http.ResponseWriter, *http.Request)
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, phs []ProposalRESTHandler) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r, phs)
}

// PostProposalReq defines the properties of a proposal request's body.
type PostProposalReq struct {
	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title          string       `json:"title" yaml:"title"`                     // Title of the proposal
	Description    string       `json:"description" yaml:"description"`         // Description of the proposal
	InitialDeposit string       `json:"initial_deposit" yaml:"initial_deposit"` // Coins to add to the proposal's deposit
	ProposerAcc    string       `json:"proposer_acc" yaml:"proposer_acc"`       // account of the proposer
}

// DepositReq defines the properties of a deposit request's body.
type DepositReq struct {
	ProposalID string       `json:"proposal_id" yaml:"proposal_id"`
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	Depositor  string       `json:"depositor" yaml:"depositor"` // Address of the depositor
	Amount     string       `json:"amount" yaml:"amount"`       // Coins to add to the proposal's deposit
}

// VoteReq defines the properties of a vote request's body.
type VoteReq struct {
	ProposalID string       `json:"proposal_id" yaml:"proposal_id"`
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	Voter      string       `json:"voter" yaml:"voter"`
	Option     string       `json:"option" yaml:"option"`
}
