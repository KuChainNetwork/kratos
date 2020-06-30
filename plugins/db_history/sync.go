package dbHistory

import (
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	ChainIdx = 1
)

// SyncState sync state in pg database
type SyncState struct {
	tableName struct{} `pg:"sync_stat,alias:sync_stat"` // default values are the same

	ID       int // both "Id" and "ID" are detected as primary key
	BlockNum int64
	ChainID  string `pg:",unique"`
}

func syncChainStat(db *pg.DB, logger log.Logger, num int64, chainID string) error {
	stat := &SyncState{
		ID: ChainIdx,
	}
	err := db.Select(stat)
	if err != nil {
		return errors.Wrapf(err, "get sync stat err")
	}

	logger.Info("get sync stat in %d", stat.BlockNum)

	return db.Update(SyncState{
		ID:       ChainIdx,
		BlockNum: num,
		ChainID:  chainID,
	})
}
