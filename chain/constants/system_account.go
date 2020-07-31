package constants

import (
	"strings"

	"github.com/KuChainNetwork/kuchain/chain/types"
)

var (
	FeeSystemAccountStr = GetSystemAccount("fee")
	FeeSystemAccount    = types.MustName(FeeSystemAccountStr)
	SystemAccount       = ChainMainName
	SystemAccountID     = types.NewAccountIDFromName(SystemAccount)
)

func IsSystemAccount(name types.Name) bool {
	if name.Eq(FeeSystemAccount) {
		return true
	}

	if name.Eq(SystemAccount) {
		return true
	}

	// TODO: add api to name to fast
	str := name.String()
	splits := strings.Split(str, "@")
	if len(splits) == 2 && splits[1] == ChainNameStr {
		return true
	}

	return false
}

func GetFeeCollector() types.AccountID {
	return types.NewAccountIDFromName(FeeSystemAccount)
}
