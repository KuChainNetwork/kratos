package types

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	"github.com/KuChainNetwork/kuchain/chain/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgCreateDex msg for create dex
type MsgCreateDex struct {
	types.KuMsg
}

const MaxDexDescriptorLen = 512

// MsgCreateDexData msg data for create dex
type MsgCreateDexData struct {
	Creator  Name   `json:"creator" yaml:"creator"` // Creator coin creator account name
	Stakings Coins  `json:"staking" yaml:"staking"` // Staking for dex
	Desc     []byte `json:"desc" yaml:"desc"`       // Description
}

func (MsgCreateDexData) Type() types.Name { return types.MustName("create@dex") }

func (m MsgCreateDexData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgCreateDex new create dex msg
func NewMsgCreateDex(auth types.AccAddress, creator types.Name, stakings types.Coins, desc []byte) MsgCreateDex {
	return MsgCreateDex{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithTransfer(types.NewAccountIDFromName(creator), ModuleAccountID, stakings),
			msg.WithData(ModuleCdc, &MsgCreateDexData{
				Creator:  creator,
				Stakings: stakings,
				Desc:     desc,
			}),
		),
	}
}

func (msg MsgCreateDex) GetData() (MsgCreateDexData, error) {
	res := MsgCreateDexData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgCreateDexData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

func (msg MsgCreateDex) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.Creator.Empty() {
		return sdkerrors.Wrap(types.ErrKuMSgNameEmpty, "creator name not empty")
	}

	if data.Stakings.IsAnyNegative() {
		return sdkerrors.Wrap(types.ErrKuMsgCoinsHasNegative, "stakings should be positive")
	}

	if len(data.Desc) >= MaxDexDescriptorLen {
		return ErrDexDescTooLong
	}

	if err := msg.ValidateTransferTo(types.NewAccountIDFromName(data.Creator), ModuleAccountID, data.Stakings); err != nil {
		return sdkerrors.Wrap(err, "create dex stakings transfer error")
	}

	return nil
}

// MsgUpdateDexDescription msg for update dex description
type MsgUpdateDexDescription struct {
	types.KuMsg
}

func (msg MsgUpdateDexDescription) GetData() (MsgUpdateDexDescriptionData, error) {
	res := MsgUpdateDexDescriptionData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgUpdateDexDescriptionData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgUpdateDexDescriptionData msg data for update dex description
type MsgUpdateDexDescriptionData struct {
	Creator Name   `json:"creator" yaml:"creator"` // Creator coin creator account name
	Desc    []byte `json:"desc" yaml:"desc"`       // Description
}

func (MsgUpdateDexDescriptionData) Type() types.Name { return types.MustName("updatedesc@dex") }

func (m MsgUpdateDexDescriptionData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgUpdateDexDescription new update dex description msg
func NewMsgUpdateDexDescription(auth types.AccAddress,
	creator types.Name,
	desc []byte) MsgUpdateDexDescription {
	return MsgUpdateDexDescription{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgUpdateDexDescriptionData{
				Creator: creator,
				Desc:    desc,
			})),
	}
}

func (msg MsgUpdateDexDescription) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.Creator.Empty() {
		return sdkerrors.Wrap(types.ErrKuMSgNameEmpty, "creator name not empty")
	}

	if len(data.Desc) >= MaxDexDescriptorLen {
		return ErrDexDescTooLong
	}

	return nil
}

// MsgDestroyDex msg for delete dex
type MsgDestroyDex struct {
	types.KuMsg
}

func (msg MsgDestroyDex) GetData() (MsgDestroyDexData, error) {
	res := MsgDestroyDexData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgDestroyDexData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgDestroyDexData msg data for delete dex
type MsgDestroyDexData struct {
	Creator Name `json:"creator" yaml:"creator"`
}

func (MsgDestroyDexData) Type() types.Name { return types.MustName("destroy@dex") }

func (m MsgDestroyDexData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgDestroyDex new destroy dex msg
func NewMsgDestroyDex(auth types.AccAddress, creator types.Name) MsgDestroyDex {
	return MsgDestroyDex{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgDestroyDexData{
				Creator: creator,
			})),
	}
}

func (msg MsgDestroyDex) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.Creator.Empty() {
		return sdkerrors.Wrap(types.ErrKuMSgNameEmpty, "creator name not empty")
	}

	return nil
}
