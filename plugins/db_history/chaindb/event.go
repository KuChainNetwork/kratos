package chaindb

import (
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	"github.com/go-pg/pg/v10"
	"github.com/tendermint/tendermint/libs/log"
)

type eventInDB struct {
	tableName struct{} `pg:"events,alias:events"` // default values are the same

	ID   int // both "Id" and "ID" are detected as primary key
	Type string

	Attributes map[string]string
}

func InsertEvent(db *pg.DB, logger log.Logger, evt *types.Event) error {
	return db.Insert(&eventInDB{
		Type:       evt.Type,
		Attributes: evt.Attributes,
	})
}
