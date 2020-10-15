package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	KuMsgMaxAuth    = 3
	KuMsgMaxDataLen = 1024
)

const (
	KuMsgMaxLen = (KuMsgMaxAuth+4)*32 + KuMsgMaxDataLen + 64 // TODO: use fix coins imp
)

type KuMsgTransfer struct {
	From   AccountID `json:"from" yaml:"from"`
	To     AccountID `json:"to" yaml:"to"`
	Amount Coins     `json:"amount" yaml:"amount"`
}

// NewKuMsgTransfer create a kuMsgTransfer
func NewKuMsgTransfer(from, to AccountID, amount Coins) KuMsgTransfer {
	return KuMsgTransfer{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

// KuMsg is the base msg for token transfer msg
type KuMsg struct {
	Auth      []sdk.AccAddress `json:"auth,omitempty" yaml:"auth"`
	Transfers []KuMsgTransfer  `json:"transfers" yaml:"transfers"`
	Router    Name             `json:"router" yaml:"router"`
	Action    Name             `json:"action" yaml:"action"`
	Data      []byte           `json:"data,omitempty" yaml:"data"`
}

// Route Implements Msg.
func (msg KuMsg) Route() string { return msg.Router.String() }

// Type Implements Msg
func (msg KuMsg) Type() string {
	if msg.Action.Empty() {
		return "transfer"
	}
	return msg.Action.String()
}

func (msg KuMsg) GetTransfers() []KuMsgTransfer {
	return msg.Transfers
}

// UnmarshalData unmarshal data to a obj
func (msg KuMsg) UnmarshalData(cdc *codec.Codec, obj interface{}) error {
	return cdc.UnmarshalBinaryLengthPrefixed(msg.Data, obj)
}

func (msg KuMsg) GetData() []byte {
	return msg.Data
}

// GetSignBytes Implements Msg.
func (msg KuMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg KuMsg) GetSigners() []sdk.AccAddress {
	res := make([]sdk.AccAddress, 0, KuMsgMaxAuth+1)

	for _, t := range msg.Transfers {
		// need check from acc address if it is
		from, ok := t.From.ToAccAddress()
		if ok {
			res = append(res, from)
		}
	}

	for _, a := range msg.Auth {
		founded := false
		for _, as := range res {
			if as.Equals(a) {
				founded = true
				break
			}
		}

		if !a.Empty() && !founded {
			res = append(res, a)
		}
	}

	return res
}

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg KuMsg) ValidateTransfer() error {
	if msg.Router.Empty() {
		return ErrKuMsgMissingRouter
	}

	if len(msg.Data) > 0 && msg.Action.Empty() {
		return ErrKuMsgMissingType
	}

	if len(msg.GetSigners()) == 0 {
		return ErrKuMsgMissingAuth
	}

	if len(msg.Data) > KuMsgMaxDataLen {
		return ErrKuMsgDataTooLarge
	}

	for _, t := range msg.Transfers {
		if t.Amount.IsAnyNegative() {
			return ErrTransfNoEnough
		}
	}

	return nil
}

func (msg KuMsg) PrettifyJSON(cdc *codec.Codec) ([]byte, error) {
	alias := struct {
		Auth      []AccAddress    `json:"auth,omitempty" yaml:"auth"`
		Transfers []KuMsgTransfer `json:"transfers" yaml:"transfers"`
		Router    Name            `json:"router" yaml:"router"`
		Action    Name            `json:"action" yaml:"action"`
		Data      json.RawMessage `json:"data" yaml:"data"`
		DataRaw   []byte          `json:"data_raw" yaml:"data_raw"`
	}{
		Auth:      msg.Auth,
		Transfers: msg.Transfers,
		Router:    msg.Router,
		Action:    msg.Action,
		DataRaw:   msg.Data,
	}

	if len(msg.Data) > 0 {
		var msgData KuMsgData = nil
		if err := cdc.UnmarshalBinaryLengthPrefixed(msg.Data, &msgData); err != nil {
			return []byte{}, errors.Wrapf(err, "unmarshal msg data")
		}

		if rawData, err := cdc.MarshalJSON(msgData); err != nil {
			return []byte{}, errors.Wrapf(err, "marshal msg data error")
		} else {
			alias.Data = rawData
		}
	}

	return cdc.MarshalJSON(alias)
}

// ValidateTransferTo validate kumsg is transfer from `from` to `to` with amount
func (msg KuMsg) ValidateTransferTo(from, to AccountID, amount Coins) error {
	if amount.IsZero() {
		return nil
	}

	for _, t := range msg.Transfers {
		if t.From.Eq(from) && t.To.Eq(to) && t.Amount.IsEqual(amount) {
			return nil
		}
	}

	return ErrKuMsgFromNotEqual
}

// ValidateTransferRequire validate kumsg is transfer to `to` with amount
func (msg KuMsg) ValidateTransferRequire(to AccountID, amount Coins) error {
	if amount.IsZero() {
		return nil
	}

	for _, t := range msg.Transfers {
		if t.To.Eq(to) && t.Amount.IsEqual(amount) {
			return nil
		}
	}

	return ErrKuMsgFromNotEqual
}
