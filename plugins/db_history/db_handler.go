package dbHistory

import (
	"github.com/KuChainNetwork/kuchain/plugins/db_history/chaindb"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
)

func (db *dbService) OnEvent(evt types.Event) error {
	if err := chaindb.InsertEvent(db.database, db.logger, &evt); err != nil {
		db.logger.Error("insert event error", "err", err)
		return err
	}

	return nil
}
