package chaindb

import "github.com/KuChainNetwork/kuchain/chain/types"

type txInDB struct {
	tableName struct{} `pg:"tx,alias:tx"` // default values are the same

	ID int64 // both "Id" and "ID" are detected as primary key

	types.StdTx
}

func newTxInDB(tx types.StdTx) *txInDB {
	return &txInDB{
		StdTx: tx,
	}
}
