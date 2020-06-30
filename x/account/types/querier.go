package types

import (
	chainTypes "github.com/KuChain-io/kuchain/chain/types"
)

// query endpoints supported by the auth Querier
const (
	QueryAccount       = "account"
	QueryAuthByAddress = "authByAddress"
	QueryParams        = "params"
)

// QueryAccountParams defines the params for querying accounts.
type QueryAccountParams struct {
	Id chainTypes.AccountID
}

// NewQueryAccountParams creates a new instance of QueryAccountParams.
func NewQueryAccountParams(id chainTypes.AccountID) QueryAccountParams {
	return QueryAccountParams{Id: id}
}

// QueryAuthSeqParams defines the params for querying accounts.
type QueryAuthSeqParams struct {
	Auth chainTypes.AccAddress
}

// NewQueryAuthSeqParams creates a new instance of QueryAuthSeqParams.
func NewQueryAuthSeqParams(auth chainTypes.AccAddress) QueryAuthSeqParams {
	return QueryAuthSeqParams{Auth: auth}
}

// QueryAuthSeqParams defines the params for querying accounts.
type QueryAuthByAddressParams struct {
	Address chainTypes.AccAddress
}

// NewQueryAddAuthParams creates a new instance of QueryAuthSeqParams.
func NewQueryAddAuthParams(address chainTypes.AccAddress) QueryAuthByAddressParams {
	return QueryAuthByAddressParams{Address: address}
}
