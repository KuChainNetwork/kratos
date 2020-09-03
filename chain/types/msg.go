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

// KuMsg is the base msg for token transfer msg
type KuMsg struct {
	Auth   []sdk.AccAddress `json:"auth,omitempty" yaml:"auth"`
	From   AccountID        `json:"from" yaml:"from"`
	To     AccountID        `json:"to" yaml:"to"`
	Amount Coins            `json:"amount" yaml:"amount"`
	Router Name             `json:"router" yaml:"router"`
	Action Name             `json:"action" yaml:"action"`
	Data   []byte           `json:"data,omitempty" yaml:"data"`
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

	// need check from acc address if it is
	from, ok := msg.From.ToAccAddress()
	if ok {
		res = append(res, from)
	}

	for _, a := range msg.Auth {
		if !a.Empty() && !a.Equals(from) {
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

	if msg.Amount.IsAnyNegative() {
		return ErrTransfNoEnough
	}

	return nil
}

// GetFrom get from account
func (msg KuMsg) GetFrom() AccountID { return msg.From }

// GetTo get to account
func (msg KuMsg) GetTo() AccountID { return msg.To }

// GetAmount get amount coin
func (msg KuMsg) GetAmount() Coins { return msg.Amount }

func (msg KuMsg) PrettifyJSON(cdc *codec.Codec) ([]byte, error) {
	alias := struct {
		Auth   []AccAddress    `json:"auth,omitempty" yaml:"auth"`
		From   AccountID       `json:"from" yaml:"from"`
		To     AccountID       `json:"to" yaml:"to"`
		Amount Coins           `json:"amount" yaml:"amount"`
		Router Name            `json:"router" yaml:"router"`
		Action Name            `json:"action" yaml:"action"`
		Data   json.RawMessage `json:"data" yaml:"data"`
	}{
		Auth:   msg.Auth,
		From:   msg.From,
		To:     msg.To,
		Amount: msg.Amount,
		Router: msg.Router,
		Action: msg.Action,
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

	if msg.From.Eq(from) &&
		msg.To.Eq(to) &&
		msg.Amount.IsEqual(amount) {
		return nil
	}

	return ErrKuMsgFromNotEqual
}

// ValidateTransferRequire validate kumsg is transfer to `to` with amount
func (msg KuMsg) ValidateTransferRequire(to AccountID, amount Coins) error {
	if amount.IsZero() {
		return nil
	}

	if msg.To.Eq(to) &&
		msg.Amount.IsEqual(amount) {
		return nil
	}

	return ErrKuMsgFromNotEqual
}
