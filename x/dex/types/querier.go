package types

import (
	"fmt"
)

// query endpoints supported by the auth Querier
const (
	QueryDex   = "queryDex"
	QuerySigIn = "querySigIn"
)

// QueryDexParams defines the params for querying dex.
type QueryDexParams struct {
	Creator Name
}

// NewQueryDexParams creates a new instance of QueryDexParams.
func NewQueryDexParams(creator Name) QueryDexParams {
	return QueryDexParams{Creator: creator}
}

// QueryDexSigInParams defines the params for querying dex sigIn.
type QueryDexSigInParams struct {
	Account AccountID
	Dex     AccountID
}

// NewQueryDexSigInParams creates a new instance of QueryDexSigInParams.
func NewQueryDexSigInParams(acc, dex AccountID) QueryDexSigInParams {
	return QueryDexSigInParams{
		Account: acc,
		Dex:     dex,
	}
}

// NodeQuerier is an interface that is satisfied by types that provide the QueryWithData method
type NodeQuerier interface {
	// QueryWithData performs a query to a Tendermint node with the provided path
	// and a data payload. It returns the result and height of the query upon success
	// or an error if the query fails.
	QueryWithData(path string, data []byte) ([]byte, int64, error)
}

// DexRetriever defines the properties of a type that can be used to
// retrieve accounts.
type DexRetriever struct {
	querier NodeQuerier
}

// NewDexRetriever init a new DexRetriever instance.
func NewDexRetriever(querier NodeQuerier) DexRetriever {
	return DexRetriever{querier: querier}
}

// GetDexWithHeight queries for a dex
func (ar DexRetriever) GetDexWithHeight(creator Name) (*Dex, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryDexParams(creator))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryDex), bs)
	if err != nil {
		return nil, height, err
	}

	data := Dex{}
	if err := ModuleCdc.UnmarshalJSON(res, &data); err != nil {
		return nil, height, err
	}

	return &data, height, nil
}

// GetSigInWithHeight queries sigIn for a dex
func (ar DexRetriever) GetSigInWithHeight(account, dex AccountID) (Coins, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryDexSigInParams(account, dex))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QuerySigIn), bs)
	if err != nil {
		return nil, height, err
	}

	data := Coins{}
	if err := ModuleCdc.UnmarshalJSON(res, &data); err != nil {
		return nil, height, err
	}

	return data, height, nil
}
