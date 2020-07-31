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
