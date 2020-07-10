package types

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgEvent event msg for plugin handler
type MsgEvent struct {
	Evt Event
}

// NewMsgEvent new msg event
func NewMsgEvent(evt sdk.Event) *MsgEvent {
	return &MsgEvent{
		Evt: FromSdkEvent(evt),
	}
}

// MsgStdTx stdTx msg for plugin handler
type MsgStdTx struct {
	Tx types.StdTx
}

// NewMsgStdTx creates a new msg
func NewMsgStdTx(tx types.StdTx) *MsgStdTx {
	return &MsgStdTx{
		Tx: tx, // no need deep copy as it will not be changed
	}
}
