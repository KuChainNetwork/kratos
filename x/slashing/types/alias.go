package types

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
)

type (
	AccountID = types.AccountID
	KuMsg     = types.KuMsg
	Name      = types.Name
)

var (
	MustName                = types.MustName
	NewAccountIDFromAccAdd  = types.NewAccountIDFromAccAdd
	NewAccountIDFromConsAdd = types.NewAccountIDFromConsAdd
)
