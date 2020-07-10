package dbHistory

import (
	"sync"

	"github.com/KuChainNetwork/kuchain/plugins/db_history/chaindb"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/config"
	"github.com/go-pg/pg/v10"
	"github.com/tendermint/tendermint/libs/log"
)

type dbWork struct {
	msg interface{}
}

func (w dbWork) IsStopped() bool {
	return w.msg == nil
}

type dbService struct {
	logger   log.Logger
	database *pg.DB

	dbChan chan dbWork
	wg     sync.WaitGroup
}

// NewDB create a connection commit event to db
func NewDB(cfg config.Cfg, logger log.Logger) *dbService {
	return &dbService{
		database: pg.Connect(&pg.Options{
			Addr:     cfg.DB.Address,
			User:     cfg.DB.User,
			Password: cfg.DB.Password,
			Database: cfg.DB.Database,
		}),
		logger: logger,
		dbChan: make(chan dbWork, 512),
	}
}

func (db *dbService) Start() error {
	db.logger.Info("Starting database service")

	db.wg.Add(1)
	go func() {
		defer db.wg.Done()

		if err := createSchema(db.database); err != nil {
			panic(err)
		}

		for {
			work, ok := <-db.dbChan
			if !ok {
				db.logger.Info("msg channel closed")
				return
			}

			if work.IsStopped() {
				db.logger.Info("db service stopped")
				return
			}

			if err := db.Process(&work); err != nil {
				db.logger.Error("db process error", "err", err)
			}
		}
	}()
	return nil
}

func (db *dbService) Process(work *dbWork) error {
	if err := chaindb.Process(db.database, db.logger, work.msg); err != nil {
		return err
	}

	return nil
}

func (db *dbService) Emit(work dbWork) {
	db.dbChan <- work
}

func (db *dbService) Stop() error {
	db.logger.Info("Stopping database service")

	db.dbChan <- dbWork{}
	db.wg.Wait()

	db.logger.Info("Database service stopped")

	db.database.Close()

	db.logger.Info("Database connection closed")
	return nil
}
