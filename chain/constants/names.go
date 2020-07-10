package constants

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/types"
)

const (
	ChainNameStr     = "kts"
	ChainMainNameStr = "kratos"
	DefaultBondDenom = ChainMainNameStr + "/" + ChainNameStr
)

var (
	// ChainMainName chain main name, as chain name for all symbol
	ChainMainName = types.MustName(ChainMainNameStr)
)

// GetSystemAccount get system account name string
func GetSystemAccount(n string) string {
	return fmt.Sprintf("%s@"+ChainNameStr, n)
}
