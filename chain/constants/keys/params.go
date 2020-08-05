package keys

const (
	IssueCoinsWaitBlockNums       int64  = 1000 // how many blocks after coins created that can issue
	DefaultMaxMemoCharacters      int    = 256
	DefaultTxSigLimit             uint64 = 7
	DefaultTxSizeCostPerByte      uint64 = 10
	DefaultSigVerifyCostED25519   uint64 = 590
	DefaultSigVerifyCostSecp256k1 uint64 = 1000
)
