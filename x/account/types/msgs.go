package types

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	"github.com/KuChainNetwork/kuchain/chain/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RouterKey is they name of the bank module
const RouterKey = ModuleName

var _, _ types.KuMsgData = (*MsgCreateAccountData)(nil), (*MsgUpdateAccountAuthData)(nil)

// MsgCreateAccountData the data struct of MsgCreateAccount
type MsgCreateAccountData struct {
	Creator types.AccountID  `json:"creator" yaml:"creator"`
	Name    types.Name       `json:"name" yaml:"name"`
	Auth    types.AccAddress `json:"auth" yaml:"auth"`
}

func (MsgCreateAccountData) Type() types.Name { return types.MustName("create@account") }

func (msg MsgCreateAccountData) Sender() AccountID {
	return msg.Creator
}

// MsgCreateAccount create account msg
type MsgCreateAccount struct {
	types.KuMsg
}

func NewMsgCreateAccount(
	auth types.AccAddress,
	creator types.AccountID,
	name types.Name,
	accountAuth types.AccAddress) MsgCreateAccount {
	return MsgCreateAccount{
		*msg.MustNewKuMsg(
			types.MustName(RouterKey),
			msg.WithAuth(auth),
			msg.WithTransfer(creator, types.NewAccountIDFromName(name), Coins{}),
			msg.WithData(Cdc(), &MsgCreateAccountData{
				Creator: creator,
				Name:    name,
				Auth:    accountAuth,
			}),
		),
	}
}

func (msg MsgCreateAccount) GetData() (MsgCreateAccountData, error) {
	res := MsgCreateAccountData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgCreateAccountData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

func (msg MsgCreateAccount) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.Creator.Empty() {
		return types.ErrKuMsgAccountIDNil
	}

	if data.Name.Empty() {
		return types.ErrNameNilString
	}

	return nil
}

// MsgUpdateAccountAuthData the data struct of MsgCreateAccount
type MsgUpdateAccountAuthData struct {
	Name types.Name       `json:"name" yaml:"name"`
	Auth types.AccAddress `json:"auth" yaml:"auth"`
}

func (MsgUpdateAccountAuthData) Type() types.Name { return types.MustName("updateauth") }

func (msg MsgUpdateAccountAuthData) Sender() AccountID {
	return NewAccountIDFromName(msg.Name)
}

// MsgUpdateAccountAuth create account msg
type MsgUpdateAccountAuth struct {
	types.KuMsg
}

// NewMsgUpdateAccountAuth create msg to update account auth
func NewMsgUpdateAccountAuth(
	auth types.AccAddress,
	name types.Name,
	accountAuth types.AccAddress) MsgUpdateAccountAuth {
	return MsgUpdateAccountAuth{
		*msg.MustNewKuMsg(
			types.MustName(RouterKey),
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgUpdateAccountAuthData{
				Name: name,
				Auth: accountAuth,
			}),
		),
	}
}

func (msg MsgUpdateAccountAuth) GetData() (MsgUpdateAccountAuthData, error) {
	res := MsgUpdateAccountAuthData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgUpdateAccountAuthData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

func (msg MsgUpdateAccountAuth) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.Name.Empty() {
		return types.ErrNameNilString
	}

	return nil
}
