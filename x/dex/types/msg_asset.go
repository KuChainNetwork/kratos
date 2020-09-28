package types

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	"github.com/KuChainNetwork/kuchain/chain/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgDexSigIn struct {
	types.KuMsg
}

type MsgDexSigInData struct {
	User   AccountID `json:"user" yaml:"user"`
	Dex    AccountID `json:"dex" yaml:"dex"`
	Amount Coins     `json:"amount" yaml:"amount"`
}

// Type imp for data KuMsgData
func (MsgDexSigInData) Type() types.Name { return types.MustName("sigin") }

func (msg MsgDexSigInData) Sender() AccountID {
	return msg.User
}

// NewMsgDexSigIn new dex sig in msg
func NewMsgDexSigIn(auth types.AccAddress, user, dex AccountID, amount Coins) MsgDexSigIn {
	return MsgDexSigIn{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgDexSigInData{
				User:   user,
				Dex:    dex,
				Amount: amount,
			}),
		),
	}
}

func (msg MsgDexSigIn) GetData() (MsgDexSigInData, error) {
	res := MsgDexSigInData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgDexSigInData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

func (msg MsgDexSigIn) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.User.Empty() {
		return sdkerrors.Wrap(types.ErrKuMsgAccountIDNil, "user accountID empty")
	}

	if data.Dex.Empty() {
		return sdkerrors.Wrap(types.ErrKuMsgAccountIDNil, "dex accountID empty")
	}

	if data.User.Eq(data.Dex) {
		return sdkerrors.Wrap(types.ErrKuMsgSpenderShouldNotEqual, "dex should not be equal to user")
	}

	if data.Amount.IsAnyNegative() {
		return types.ErrKuMsgCoinsHasNegative
	}

	return nil
}