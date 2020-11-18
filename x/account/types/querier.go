package types

import (
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
)

// query endpoints supported by the auth Querier
const (
	QueryAccount        = "account"
	QueryAuthByAddress  = "authByAddress"
	QueryAccountsByAuth = "accountsByAuth"
	QueryParams         = "params"
)

// QueryAccountParams defines the params for querying accounts.
type QueryAccountParams struct {
	ID chainTypes.AccountID
}

// NewQueryAccountParams creates a new instance of QueryAccountParams.
func NewQueryAccountParams(id chainTypes.AccountID) QueryAccountParams {
	return QueryAccountParams{ID: id}
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

// QueryAuthSeqParams defines the params for querying accounts.
type QueryAccountsByAuthParams struct {
	Auth chainTypes.AccAddress
}

// NewQueryAccountsByAuthParams creates a new instance of QueryAccountsByAuthParams.
func NewQueryAccountsByAuthParams(auth string) QueryAccountsByAuthParams {
	return QueryAccountsByAuthParams{Auth: chainTypes.MustAccAddressFromBech32(auth)}
}
