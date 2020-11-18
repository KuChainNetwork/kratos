package msg

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
)

type (
	KuMsg       = types.KuMsg
	KuMsgData   = types.KuMsgData
	KuTransfMsg = types.KuTransfMsg
)

const (
	KuMsgMaxAuth = types.KuMsgMaxAuth
)

var (
	NewKuMsgCtx = types.NewKuMsgCtx
)

// NewKuMsg create kuMsg by router and opts
func NewKuMsg(router Name, opts ...Option) (*KuMsg, error) {
	res := &KuMsg{
		Auth:   make([]AccAddress, 0, 4),
		Router: router,
	}

	for _, opt := range opts {
		err := opt.Op(res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

// MustNewKuMsg create kuMsg
func MustNewKuMsg(router Name, opts ...Option) *KuMsg {
	res, err := NewKuMsg(router, opts...)
	if err != nil {
		panic(err)
	}
	return res
}
