package types

// DONTCOVER

// query endpoints supported by the governance Querier
const (
	QueryParams           = "params"
	QueryProposals        = "proposals"
	QueryProposal         = "proposal"
	QueryDeposits         = "deposits"
	QueryDeposit          = "deposit"
	QueryVotes            = "votes"
	QueryVote             = "vote"
	QueryTally            = "tally"
	QueryPunishValidators = "punishvalidators"
	QueryPunishValidator  = "punishvalidator"

	ParamDeposit  = "deposit"
	ParamVoting   = "voting"
	ParamTallying = "tallying"
)

// QueryProposalParams Params for queries:
// - 'custom/gov/proposal'
// - 'custom/gov/deposits'
// - 'custom/gov/tally'
type QueryProposalParams struct {
	ProposalID uint64
}

// NewQueryProposalParams creates a new instance of QueryProposalParams
func NewQueryProposalParams(proposalID uint64) QueryProposalParams {
	return QueryProposalParams{
		ProposalID: proposalID,
	}
}

// QueryProposalVotesParams used for queries to 'custom/gov/votes'.
type QueryProposalVotesParams struct {
	ProposalID uint64
	Page       int
	Limit      int
}

// NewQueryProposalVotesParams creates new instance of the QueryProposalVotesParams.
func NewQueryProposalVotesParams(proposalID uint64, page, limit int) QueryProposalVotesParams {
	return QueryProposalVotesParams{
		ProposalID: proposalID,
		Page:       page,
		Limit:      limit,
	}
}

// QueryDepositParams params for query 'custom/gov/deposit'
type QueryDepositParams struct {
	ProposalID uint64
	Depositor  AccountID
}

// NewQueryDepositParams creates a new instance of QueryDepositParams
func NewQueryDepositParams(proposalID uint64, depositor AccountID) QueryDepositParams {
	return QueryDepositParams{
		ProposalID: proposalID,
		Depositor:  depositor,
	}
}

// QueryVoteParams Params for query 'custom/gov/vote'
type QueryVoteParams struct {
	ProposalID uint64
	Voter      AccountID
}

// NewQueryVoteParams creates a new instance of QueryVoteParams
func NewQueryVoteParams(proposalID uint64, voter AccountID) QueryVoteParams {
	return QueryVoteParams{
		ProposalID: proposalID,
		Voter:      voter,
	}
}

// QueryProposalsParams Params for query 'custom/gov/proposals'
type QueryProposalsParams struct {
	Page           int
	Limit          int
	Voter          AccountID
	Depositor      AccountID
	ProposalStatus ProposalStatus
}

// NewQueryProposalsParams creates a new instance of QueryProposalsParams
func NewQueryProposalsParams(page, limit int, status ProposalStatus, voter AccountID, depositor AccountID) QueryProposalsParams {
	return QueryProposalsParams{
		Page:           page,
		Limit:          limit,
		Voter:          voter,
		Depositor:      depositor,
		ProposalStatus: status,
	}
}

type QueryPunishValidatorParams struct {
	ValidatorAccount AccountID
}

func NewQueryPunishValidatorParams(validatorAccount AccountID) QueryPunishValidatorParams {
	return QueryPunishValidatorParams{
		ValidatorAccount: validatorAccount,
	}
}
