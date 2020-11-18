package constants

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/constants/keys"
	"github.com/KuChainNetwork/kuchain/chain/types"
)

const (
	CommonAccountNameLen = 12
)

var (
	ChainNameStr      = keys.ChainNameStr
	ChainMainNameStr  = keys.ChainMainNameStr
	DefaultBondDenom  = keys.DefaultBondDenom
	DefaultBondSymbol = keys.DefaultBondSymbol
)

var (
	// ChainMainName chain main name, as chain name for all symbol
	ChainMainName         = types.MustName(ChainMainNameStr)
	DefaultBondSymbolName = types.MustName(DefaultBondSymbol)
)

// GetSystemAccount get system account name string
func GetSystemAccount(n string) string {
	return fmt.Sprintf("%s@"+ChainNameStr, n)
}
