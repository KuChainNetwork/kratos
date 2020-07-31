package types

import (
	"fmt"
)

// zero fee pool
func InitialFeePool() FeePool {
	return FeePool{
		CommunityPool: DecCoins{},
	}
}

// global fee pool for distribution
type FeePool struct {
	CommunityPool DecCoins `json:"community_pool" yaml:"community_pool"`
}

// ValidateGenesis validates the fee pool for a genesis state
func (f FeePool) ValidateGenesis() error {
	if f.CommunityPool.IsAnyNegative() {
		return fmt.Errorf("negative CommunityPool in distribution fee pool, is %v",
			f.CommunityPool)
	}

	return nil
}
