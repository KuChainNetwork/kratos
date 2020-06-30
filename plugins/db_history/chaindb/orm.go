package chaindb

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func RegOrm(db *pg.DB) error {
	models := []interface{}{
		(*eventInDB)(nil),
		(*txInDB)(nil),
		(*MessageInDB)(nil),
		(*KuTransferInDB)(nil),
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
