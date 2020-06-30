package types

import (
	"github.com/KuChain-io/kuchain/chain/msg"
	"github.com/KuChain-io/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RouterKey is they name of the bank module
const RouterKey = ModuleName

// MsgCreateAccountData the data struct of MsgCreateAccount
type MsgCreateAccountData struct {
	Creator types.AccountID  `json:"creator" yaml:"creator"`
	Name    types.Name       `json:"name" yaml:"name"`
	Auth    types.AccAddress `json:"auth" yaml:"auth"`
}

func (MsgCreateAccountData) Type() types.Name { return types.MustName("create@account") }

func (m MsgCreateAccountData) Marshal() ([]byte, error) {
	return ModuleCdc.MarshalJSON(m)
}

func (m *MsgCreateAccountData) Unmarshal(b []byte) error {
	return ModuleCdc.UnmarshalJSON(b, m)
}

// MsgCreateAccount create account msg
type MsgCreateAccount struct {
	types.KuMsg
}

func NewMsgCreateAccount(auth types.AccAddress, creator types.AccountID, name types.Name, accountAuth types.AccAddress) MsgCreateAccount {
	return MsgCreateAccount{
		*msg.MustNewKuMsg(
			types.MustName(RouterKey),
			msg.WithAuth(auth),
			msg.WithTransfer(creator, types.NewAccountIDFromName(name), sdk.Coins{}), // TODO: with first coin
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

// MsgUpdateAccountAuthData the data struct of MsgCreateAccount
type MsgUpdateAccountAuthData struct {
	Name types.Name       `json:"name" yaml:"name"`
	Auth types.AccAddress `json:"auth" yaml:"auth"`
}

func (MsgUpdateAccountAuthData) Type() types.Name { return types.MustName("updateauth") }

// MsgUpdateAccountAuth create account msg
type MsgUpdateAccountAuth struct {
	types.KuMsg
}

// NewMsgUpdateAccountAuth create msg to update account auth
func NewMsgUpdateAccountAuth(auth types.AccAddress, name types.Name, accountAuth types.AccAddress) MsgUpdateAccountAuth {
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
