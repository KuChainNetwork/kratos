package types

import (
	fmt "fmt"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
)

// NodeQuerier is an interface that is satisfied by types that provide the QueryWithData method
type NodeQuerier interface {
	// QueryWithData performs a query to a Tendermint node with the provided path
	// and a data payload. It returns the result and height of the query upon success
	// or an error if the query fails.
	QueryWithData(path string, data []byte) ([]byte, int64, error)
}

// AccountRetriever defines the properties of a type that can be used to
// retrieve accounts.
type AccountRetriever struct {
	querier NodeQuerier
}

// NewAccountRetriever initialises a new AccountRetriever instance.
func NewAccountRetriever(querier NodeQuerier) AccountRetriever {
	return AccountRetriever{querier: querier}
}

// GetAccount queries for an account given an address and a block height. An
// error is returned if the query or decoding fails.
func (ar AccountRetriever) GetAccount(id types.AccountID) (exported.Account, error) {
	account, _, err := ar.GetAccountWithHeight(id)
	return account, err
}

// GetAccountWithHeight queries for an account given an address. Returns the
// height of the query with the account. An error is returned if the query
// or decoding fails.
func (ar AccountRetriever) GetAccountWithHeight(id types.AccountID) (exported.Account, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryAccountParams(id))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryAccount), bs)
	if err != nil {
		return nil, height, err
	}

	var account exported.Account
	if err := ModuleCdc.UnmarshalJSON(res, &account); err != nil {
		return nil, height, err
	}

	return account, height, nil
}

// GetAddAuth queries for an auth state given an address and a block height.
func (ar AccountRetriever) GetAddAuth(add types.AccAddress) (Auth, error) {
	au, _, err := ar.GetAddAuthWithHeight(add)
	return au, err
}

// GetAddAuthWithHeight queries for an auth state given an address.
func (ar AccountRetriever) GetAddAuthWithHeight(add types.AccAddress) (Auth, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryAddAuthParams(add))
	if err != nil {
		return Auth{}, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryAuthByAddress), bs)
	if err != nil {
		return Auth{}, height, err
	}

	var authData Auth
	if err := ModuleCdc.UnmarshalJSON(res, &authData); err != nil {
		return Auth{}, height, err
	}

	return authData, height, nil
}

// EnsureExists returns an error if no account exists for the given address else nil.
func (ar AccountRetriever) EnsureExists(id types.AccountID) error {
	if _, err := ar.GetAccount(id); err != nil {
		return err
	}
	return nil
}

// GetAccountNumberSequence returns sequence and account number for the given address.
// It returns an error if the account couldn't be retrieved from the state.
func (ar AccountRetriever) GetAuthNumberSequence(id types.AccountID) (uint64, uint64, error) {
	auth, ok := id.ToAccAddress()
	if !ok {
		acc, err := ar.GetAccount(id)
		if err != nil {
			return 0, 0, err
		}

		if acc == nil {
			return 0, 0, ErrAccountNoFound
		}
		auth = acc.GetAuth()
	}

	authData, err := ar.GetAddAuth(auth)
	if err != nil {
		return 0, 0, err
	}

	return authData.GetNumber(), authData.GetSequence(), nil
}
