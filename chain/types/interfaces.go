package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AssetTransfer a interface for asset coins transfer
type AssetTransfer interface {
	Transfer(ctx sdk.Context, from, to AccountID, amount Coins) error
}

// AccountAuther a interface for account auth getter
type AccountAuther interface {
	GetAuth(ctx sdk.Context, account Name) (AccAddress, error)
}

// KuTransfMsg ku Msg
type KuTransfMsg interface {
	sdk.Msg

	GetFrom() AccountID
	GetTo() AccountID
	GetAmount() Coins
	Type() string
	GetData() []byte
}

var _ KuTransfMsg = &KuMsg{}

type KuMsgData interface {
	Type() Name
}

// Prettifier a type can prettify a byte
type Prettifier interface {
	PrettifyJSON(cdc *codec.Codec) ([]byte, error)
}
