package exported

import (
	"github.com/KuChain-io/kuchain/x/account/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ModuleAccountI defines an account interface for modules that hold tokens in
// an escrow.
type ModuleAccountI interface {
	exported.Account

	GetAddress() sdk.AccAddress
	GetPermissions() []string
	HasPermission(string) bool
}

// SupplyI defines an inflationary supply interface for modules that handle
// token supply.
type SupplyI interface {
	GetTotal() sdk.Coins
	SetTotal(total sdk.Coins)

	String() string
	ValidateBasic() error
}
