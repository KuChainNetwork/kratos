package chaindb

import (
	"reflect"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/plugins/test/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg/v10"
	"github.com/tendermint/tendermint/libs/log"
)

func Process(db *pg.DB, logger log.Logger, msg interface{}) error {
	logger.Debug("process msg", "typ", reflect.TypeOf(msg), "msg", msg)
	switch msg := msg.(type) {
	case types.Event:
		return InsertEvent(db, logger, &msg)
	case chainTypes.StdTx:
		return insert(db, newTxInDB(msg))
	}

	if msg, ok := msg.(sdk.Msg); ok {
		return processMsg(db, msg)
	}

	return nil
}

func insert(db *pg.DB, obj interface{}) error {
	return db.Insert(obj)
}
