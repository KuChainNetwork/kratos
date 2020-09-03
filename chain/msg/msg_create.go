package msg

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Option interface {
	Op(msg *KuMsg) error
}

// WithData create kumsg with data and action
type withData struct {
	data  KuMsgData
	bytes []byte
}

func (k withData) Op(msg *KuMsg) error {
	msg.Action = k.data.Type()
	msg.Data = k.bytes

	return nil
}

// WithData create kumsg with data
func WithData(cdc *codec.Codec, data KuMsgData) Option {
	dataByte, err := cdc.MarshalBinaryLengthPrefixed(data)
	if err != nil {
		panic(err)
	}
	return withData{
		data:  data,
		bytes: dataByte,
	}
}

type withTransfer struct {
	From   AccountID
	To     AccountID
	Amount Coins
}

func (k withTransfer) Op(msg *KuMsg) error {
	msg.Transfers = append(msg.Transfers,
		types.KuMsgTransfer{
			From:   k.From,
			To:     k.To,
			Amount: k.Amount,
		})
	return nil
}

// WithTransfer create kumsg with transfer message
func WithTransfer(from, to AccountID, amount Coins) Option {
	return withTransfer{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

type withAuth struct {
	Auth []AccAddress
}

func (k withAuth) Op(msg *KuMsg) error {
	if len(msg.Auth)+len(k.Auth) > KuMsgMaxAuth {
		return types.ErrKuMsgAuthCountTooLarge
	}

	msg.Auth = append(msg.Auth, k.Auth...)
	return nil
}

// WithAuth create kumsg with auth
func WithAuth(auth AccAddress) Option {
	return withAuth{
		Auth: []AccAddress{auth},
	}
}

// WithAuths create kumsg with auth
func WithAuths(auths []AccAddress) Option {
	return withAuth{
		Auth: auths,
	}
}
