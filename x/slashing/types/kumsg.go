package types

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	chaintype "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	RouterKeyName = chaintype.MustName(RouterKey)
)

type KuMsgUnjail struct {
	chaintype.KuMsg
}

func NewKuMsgUnjail(auth sdk.AccAddress, validatorAddr chaintype.AccountID) KuMsgUnjail {
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
