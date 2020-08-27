package chaindb

import (
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

type MessageInDB struct {
	tableName struct{} `pg:"messages,alias:messages"` // default values are the same

	ID int64 // both "Id" and "ID" are detected as primary key

	sdk.Msg // FIXME: use petty show msg
}

type KuTransferInDB struct {
	tableName struct{} `pg:"transfer,alias:transfer"` // default values are the same

	ID     int64 // both "Id" and "ID" are detected as primary key
	Route  string
	Type   string
	From   string
	To     string
	Amount int64
	Symbol string
}

func newMsgToDB(msg sdk.Msg) *MessageInDB {
	return &MessageInDB{
		Msg: msg,
	}
}

func processMsg(db *pg.DB, msg sdk.Msg) error {
	if err := db.Insert(newMsgToDB(msg)); err != nil {
		return errors.Wrapf(err, "insert msg")
	}

	if msg, ok := msg.(chainTypes.KuTransfMsg); ok {
		transfers := msg.GetTransfers()

		for _, t := range transfers {
			amounts := t.Amount

			in := &KuTransferInDB{
				Route: msg.Route(),
				Type:  msg.Type(),
				From:  t.From.String(),
				To:    t.To.String(),
			}

			for _, amount := range amounts {
				in.Amount = amount.Amount.BigInt().Int64()
				in.Symbol = amount.Denom
				if err := db.Insert(in); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
