package types

import (
	err "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	Register = err.Register
)

// x/staking module sentinel errors
//
// TODO: Many of these errors are redundant. They should be removed and replaced
// by sdkerrors.ErrInvalidRequest.
//
// REF: https://github.com/cosmos/cosmos-sdk/issues/5450
var (
	ErrEmptyValidatorAddr              = Register(ModuleName, 2, "empty validator account")
	ErrBadValidatorAddr                = Register(ModuleName, 3, "validator account is invalid")
	ErrNoValidatorFound                = Register(ModuleName, 4, "validator does not exist")
	ErrValidatorOwnerExists            = Register(ModuleName, 5, "validator already exist for this operator account; must use new validator operator account")
	ErrValidatorPubKeyExists           = Register(ModuleName, 6, "validator already exist for this pubkey; must use new validator pubkey")
	ErrValidatorPubKeyTypeNotSupported = Register(ModuleName, 7, "validator pubkey type is not supported")
	ErrValidatorJailed                 = Register(ModuleName, 8, "validator for this account is currently jailed")
	ErrBadRemoveValidator              = Register(ModuleName, 9, "failed to remove validator")
	ErrCommissionNegative              = Register(ModuleName, 10, "commission must be positive")
	ErrCommissionHuge                  = Register(ModuleName, 11, "commission cannot be more than 100%")
	ErrCommissionGTMaxRate             = Register(ModuleName, 12, "commission cannot be more than the max rate")
	ErrCommissionUpdateTime            = Register(ModuleName, 13, "commission cannot be changed more than once in 24h")
	ErrCommissionChangeRateNegative    = Register(ModuleName, 14, "commission change rate must be positive")
	ErrCommissionChangeRateGTMaxRate   = Register(ModuleName, 15, "commission change rate cannot be more than the max rate")
	ErrCommissionGTMaxChangeRate       = Register(ModuleName, 16, "commission cannot be changed more than max change rate")
	ErrSelfDelegationBelowMinimum      = Register(ModuleName, 17, "validator's self delegation must be greater than their minimum self delegation")
	ErrMinSelfDelegationInvalid        = Register(ModuleName, 18, "minimum self delegation must be a positive integer")
	ErrMinSelfDelegationDecreased      = Register(ModuleName, 19, "minimum self delegation cannot be decrease")
	ErrEmptyDelegatorAddr              = Register(ModuleName, 20, "empty delegator account")
	ErrBadDenom                        = Register(ModuleName, 21, "invalid coin denomination")
	ErrBadDelegationAddr               = Register(ModuleName, 22, "invalid account for (account, validator) tuple")
	ErrBadDelegationAmount             = Register(ModuleName, 23, "invalid delegation amount")
	ErrNoDelegation                    = Register(ModuleName, 24, "no delegation for (account, validator) tuple")
	ErrBadDelegatorAddr                = Register(ModuleName, 25, "delegator does not exist with account")
	ErrNoDelegatorForAddress           = Register(ModuleName, 26, "delegator does not contain delegation")
	ErrInsufficientShares              = Register(ModuleName, 27, "insufficient delegation shares")
	ErrDelegationValidatorEmpty        = Register(ModuleName, 28, "cannot delegate to an empty validator")
	ErrNotEnoughDelegationShares       = Register(ModuleName, 29, "not enough delegation shares")
	ErrBadSharesAmount                 = Register(ModuleName, 30, "invalid shares amount")
	ErrBadSharesPercent                = Register(ModuleName, 31, "Invalid shares percent")
	ErrNotMature                       = Register(ModuleName, 32, "entry not mature")
	ErrNoUnbondingDelegation           = Register(ModuleName, 33, "no unbonding delegation found")
	ErrMaxUnbondingDelegationEntries   = Register(ModuleName, 34, "too many unbonding delegation entries for (delegator, validator) tuple")
	ErrBadRedelegationAddr             = Register(ModuleName, 35, "invalid account for (account, src-validator, dst-validator) tuple")
	ErrNoRedelegation                  = Register(ModuleName, 36, "no redelegation found")
	ErrSelfRedelegation                = Register(ModuleName, 37, "cannot redelegate to the same validator")
	ErrTinyRedelegationAmount          = Register(ModuleName, 38, "too few tokens to redelegate (truncates to zero tokens)")
	ErrBadRedelegationDst              = Register(ModuleName, 39, "redelegation destination validator not found")
	ErrTransitiveRedelegation          = Register(ModuleName, 40,
		"redelegation to this validator already in progress; first redelegation to this validator must complete before next redelegation")
	ErrMaxRedelegationEntries      = Register(ModuleName, 41, "too many redelegation entries for (delegator, src-validator, dst-validator) tuple")
	ErrDelegatorShareExRateInvalid = Register(ModuleName, 42, "cannot delegate to validators with invalid (zero) ex-rate")
	ErrBothShareMsgsGiven          = Register(ModuleName, 43, "both shares amount and shares percent provided")
	ErrNeitherShareMsgsGiven       = Register(ModuleName, 44, "neither shares amount nor shares percent provided")
	ErrInvalidHistoricalInfo       = Register(ModuleName, 45, "invalid historical info")
	ErrNoHistoricalInfo            = Register(ModuleName, 46, "no historical info found")
	ErrEmptyValidatorPubKey        = Register(ModuleName, 47, "empty validator public key")
	ErrUnKnowAccount               = Register(ModuleName, 48, "validator operator is not a known account")
)
