package dbHistory

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/KuChainNetwork/kuchain/plugins/db_history/chaindb"
)

// createSchema creates database schema for User and Story models.
func createSchema(db *pg.DB) error {
	if err := chaindb.RegOrm(db); err != nil {
		return err
	}

	models := []interface{}{
		(*SyncState)(nil),
	}

	for _, model := range models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
