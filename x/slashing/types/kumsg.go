package types

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	RouterKeyName = MustName(RouterKey)
)

type KuMsgUnjail struct {
	KuMsg
}

func NewKuMsgUnjail(auth sdk.AccAddress, validatorAddr AccountID) KuMsgUnjail {
	return KuMsgUnjail{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgUnjail{
				ValidatorAddr: validatorAddr,
			}),
		),
	}
}

func (msg KuMsgUnjail) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	msgData := MsgUnjail{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return err
	}
	return msgData.ValidateBasic()
}
