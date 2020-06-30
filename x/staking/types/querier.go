package types

import (
	//sdk "github.com/cosmos/cosmos-sdk/types"
	chaintype "github.com/KuChain-io/kuchain/chain/types"
)

// query endpoints supported by the staking Querier
const (
	QueryValidators                    = "validators"
	QueryValidator                     = "validator"
	QueryDelegatorDelegations          = "delegatorDelegations"
	QueryDelegatorUnbondingDelegations = "delegatorUnbondingDelegations"
	QueryRedelegations                 = "redelegations"
	QueryValidatorDelegations          = "validatorDelegations"
	QueryValidatorRedelegations        = "validatorRedelegations"
	QueryValidatorUnbondingDelegations = "validatorUnbondingDelegations"
	QueryDelegation                    = "delegation"
	QueryUnbondingDelegation           = "unbondingDelegation"
	QueryDelegatorValidators           = "delegatorValidators"
	QueryDelegatorValidator            = "delegatorValidator"
	QueryPool                          = "pool"
	QueryParameters                    = "parameters"
	QueryHistoricalInfo                = "historicalInfo"
)

// defines the params for the following queries:
// - 'custom/staking/delegatorDelegations'
// - 'custom/staking/delegatorUnbondingDelegations'
// - 'custom/staking/delegatorRedelegations'
// - 'custom/staking/delegatorValidators'
type QueryDelegatorParams struct {
	DelegatorAddr chaintype.AccountID
}

func NewQueryDelegatorParams(delegatorAddr chaintype.AccountID) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/validator'
// - 'custom/staking/validatorDelegations'
// - 'custom/staking/validatorUnbondingDelegations'
// - 'custom/staking/validatorRedelegations'
type QueryValidatorParams struct {
	ValidatorAddr chaintype.AccountID
}

func NewQueryValidatorParams(validatorAddr chaintype.AccountID) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/delegation'
// - 'custom/staking/unbondingDelegation'
// - 'custom/staking/delegatorValidator'
type QueryBondsParams struct {
	DelegatorAddr chaintype.AccountID
	ValidatorAddr chaintype.AccountID
}

func NewQueryBondsParams(delegatorAddr chaintype.AccountID, validatorAddr chaintype.AccountID) QueryBondsParams {
	return QueryBondsParams{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/redelegation'
type QueryRedelegationParams struct {
	DelegatorAddr    chaintype.AccountID
	SrcValidatorAddr chaintype.AccountID
	DstValidatorAddr chaintype.AccountID
}

func NewQueryRedelegationParams(delegatorAddr chaintype.AccountID,
	srcValidatorAddr, dstValidatorAddr chaintype.AccountID) QueryRedelegationParams {

	return QueryRedelegationParams{
		DelegatorAddr:    delegatorAddr,
		SrcValidatorAddr: srcValidatorAddr,
		DstValidatorAddr: dstValidatorAddr,
	}
}

// QueryValidatorsParams defines the params for the following queries:
// - 'custom/staking/validators'
type QueryValidatorsParams struct {
	Page, Limit int
	Status      string
}

func NewQueryValidatorsParams(page, limit int, status string) QueryValidatorsParams {
	return QueryValidatorsParams{page, limit, status}
}

// QueryHistoricalInfoParams defines the params for the following queries:
// - 'custom/staking/historicalInfo'
type QueryHistoricalInfoParams struct {
	Height int64
}

// NewQueryHistoricalInfoParams creates a new QueryHistoricalInfoParams instance
func NewQueryHistoricalInfoParams(height int64) QueryHistoricalInfoParams {
	return QueryHistoricalInfoParams{height}
}
