package transaction

import "github.com/KuChainNetwork/kuchain/chain/types"

type (
	StdTx        = types.StdTx
	StdSignMsg   = types.StdSignMsg
	StdSignature = types.StdSignature
)

var (
	NewStdTx  = types.NewStdTx
	NewStdFee = types.NewStdFee
)
