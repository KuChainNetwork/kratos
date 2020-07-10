package constants

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/constants/keys"
	"github.com/KuChainNetwork/kuchain/chain/types"
)

const (
	ChainNameStr     = keys.ChainNameStr
	ChainMainNameStr = keys.ChainMainNameStr
	DefaultBondDenom = keys.DefaultBondDenom
)

var (
	// ChainMainName chain main name, as chain name for all symbol
	ChainMainName = types.MustName(ChainMainNameStr)
)

// GetSystemAccount get system account name string
func GetSystemAccount(n string) string {
	return fmt.Sprintf("%s@"+ChainNameStr, n)
}
