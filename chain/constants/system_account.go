package constants

import (
	"strings"

	"github.com/KuChain-io/kuchain/chain/types"
)

const (
	FeeSystemAccountStr = "fee@kts"
)

var (
	FeeSystemAccount = types.MustName(FeeSystemAccountStr)
	SystemAccount    = types.MustName("kratos")
)

func IsSystemAccount(name types.Name) bool {
	if name.Eq(FeeSystemAccount) {
		return true
	}

	// TODO use constants
	if name.Eq(SystemAccount) {
		return true
	}

	// TODO: add api to name to fast
	str := name.String()
	splits := strings.Split(str, "@")
	if len(splits) == 2 && splits[1] == "kts" {
		return true
	}

	return false
}

func GetFeeCollector() types.AccountID {
	return types.NewAccountIDFromName(FeeSystemAccount)
}
