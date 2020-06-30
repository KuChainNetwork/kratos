package types

import (
	"github.com/KuChain-io/kuchain/chain/msg"
	"github.com/KuChain-io/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is they name of the asset module
const RouterKey = ModuleName

var (
	RouterKeyName = types.MustName(RouterKey)
)

type (
	KuMsg = types.KuMsg
)

// NewMsgTransfer create msg transfer
func NewMsgTransfer(auth types.AccAddress, from types.AccountID, to types.AccountID, amount sdk.Coins) KuMsg {
	return *msg.MustNewKuMsg(
		RouterKeyName,
		msg.WithAuth(auth),
		msg.WithTransfer(from, to, amount),
	)
}

type MsgCreateCoin struct {
	types.KuMsg
}

func (MsgCreateCoinData) Type() types.Name { return types.MustName("create@coin") }

// NewMsgCreate new create coin msg
func NewMsgCreate(auth types.AccAddress, creator types.Name, symbol types.Name, maxSupply types.Coin, canIssue, canLock bool, issue2Height int64, initSupply types.Coin, desc []byte) MsgCreateCoin {
	return MsgCreateCoin{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgCreateCoinData{
				Creator:       creator,
				Symbol:        symbol,
				MaxSupply:     maxSupply,
				CanIssue:      canIssue,
				CanLock:       canLock,
				IssueToHeight: issue2Height,
				InitSupply:    initSupply,
				Desc:          desc,
			}),
		),
	}
}

type MsgIssueCoin struct {
	types.KuMsg
}

// Type imp for data KuMsgData
func (MsgIssueCoinData) Type() types.Name { return types.MustName("issue") }

// NewMsgIssue new issue msg
func NewMsgIssue(auth types.AccAddress, creator, symbol types.Name, amount types.Coin) MsgIssueCoin {
	return MsgIssueCoin{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgIssueCoinData{
				Creator: creator,
				Symbol:  symbol,
				Amount:  amount,
			}),
		),
	}
}

type MsgBurnCoin struct {
	types.KuMsg
}

// Type imp for data KuMsgData
func (MsgBurnCoinData) Type() types.Name { return types.MustName("burn") }

// NewMsgBurn new issue msg
func NewMsgBurn(auth types.AccAddress, id types.AccountID, amount types.Coin) MsgIssueCoin {
	return MsgIssueCoin{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgBurnCoinData{
				Id:     id,
				Amount: amount,
			}),
		),
	}
}

// MsgLockCoin msg to lock coin
type MsgLockCoin struct {
	types.KuMsg
}

// Type imp for data KuMsgData
func (m *MsgLockCoinData) Type() types.Name { return types.MustName("lock@coin") }

// NewMsgLockCoin create new lock coin msg
func NewMsgLockCoin(auth types.AccAddress, id types.AccountID, amount types.Coins, unlockBlockHeight int64) MsgLockCoin {
	return MsgLockCoin{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgLockCoinData{
				Id:                id,
				Amount:            amount,
				UnlockBlockHeight: unlockBlockHeight,
			}),
		),
	}
}

// MsgUnlockCoin msg to unlock coin
type MsgUnlockCoin struct {
	types.KuMsg
}

// Type imp for data KuMsgData
func (m *MsgUnlockCoinData) Type() types.Name { return types.MustName("unlock@coin") }

// NewMsgUnlockCoin create new lock coin msg
func NewMsgUnlockCoin(auth types.AccAddress, id types.AccountID, amount types.Coins) MsgUnlockCoin {
	return MsgUnlockCoin{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgUnlockCoinData{
				Id:     id,
				Amount: amount,
			}),
		),
	}
}
