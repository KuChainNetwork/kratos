package types

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/supply/exported"
	"gopkg.in/yaml.v2"
)

// Implements Delegation interface
var _ exported.SupplyI = (*Supply)(nil)

// Supply represents a struct that passively keeps track of the total supply
// amounts in the network.
type Supply struct {
	Total types.Coins `json:"total" yaml:"total"`
}

// NewSupply creates a new Supply instance
func NewSupply(total types.Coins) *Supply {
	return &Supply{total}
}

// DefaultSupply creates an empty Supply
func DefaultSupply() *Supply {
	return NewSupply(types.NewCoins())
}

// SetTotal sets the total supply.
func (supply *Supply) SetTotal(total types.Coins) {
	supply.Total = total
}

// GetTotal returns the supply total.
func (supply Supply) GetTotal() types.Coins {
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
