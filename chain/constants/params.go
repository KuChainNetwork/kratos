package constants

import "github.com/KuChainNetwork/kuchain/chain/constants/keys"

const (
	IssueCoinsWaitBlockNums = keys.IssueCoinsWaitBlockNums // how many blocks after coins created that can issue
)

// Default parameter values
const (
	DefaultMaxMemoCharacters      = keys.DefaultMaxMemoCharacters
	DefaultTxSigLimit             = keys.DefaultTxSigLimit
	DefaultTxSizeCostPerByte      = keys.DefaultTxSizeCostPerByte
	DefaultSigVerifyCostED25519   = keys.DefaultSigVerifyCostED25519
	DefaultSigVerifyCostSecp256k1 = keys.DefaultSigVerifyCostSecp256k1
)
