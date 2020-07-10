package types

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/x/supply/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

// Implements Delegation interface
var _ exported.SupplyI = (*Supply)(nil)

// NewSupply creates a new Supply instance
func NewSupply(total sdk.Coins) *Supply {
	return &Supply{total}
}

// DefaultSupply creates an empty Supply
func DefaultSupply() *Supply {
	return NewSupply(sdk.NewCoins())
}

// SetTotal sets the total supply.
func (supply *Supply) SetTotal(total sdk.Coins) {
	supply.Total = total
}

// GetTotal returns the supply total.
func (supply Supply) GetTotal() sdk.Coins {
	return supply.Total
}

// String returns a human readable string representation of a supplier.
func (supply Supply) String() string {
	bz, _ := yaml.Marshal(supply)
	return string(bz)
}

// ValidateBasic validates the Supply coins and returns error if invalid
func (supply Supply) ValidateBasic() error {
	if !supply.Total.IsValid() {
		return fmt.Errorf("invalid total supply: %s", supply.Total.String())
	}
	return nil
}
