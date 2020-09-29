package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Msg = sdk.Msg
)

// AssetTransfer a interface for asset coins transfer
type AssetTransfer interface {
	Transfer(ctx sdk.Context, from, to AccountID, amount Coins) error
	TransferDetail(ctx sdk.Context, from, to AccountID, amount Coins, isApplyApprove bool) error
	ApplyApporve(ctx sdk.Context, from, to AccountID, amount Coins) error
}

// AccountAuther a interface for account auth getter
type AccountAuther interface {
	GetAuth(ctx sdk.Context, account Name) (AccAddress, error)
}

// KuTransfMsg ku Msg
type KuTransfMsg interface {
	Route() string
	Type() string
	GetSignBytes() []byte
	GetSigners() []AccAddress
	GetTransfers() []KuMsgTransfer
	ValidateTransfer() error
}

var _ KuTransfMsg = &KuMsg{}

type KuMsgData interface {
	Type() Name
	Sender() AccountID
}

// Prettifier a type can prettify a byte
type Prettifier interface {
	PrettifyJSON(cdc *codec.Codec) ([]byte, error)
}
