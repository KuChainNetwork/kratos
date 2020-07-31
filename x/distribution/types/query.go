package types

import (
	"fmt"
	"strings"
)

// QueryDelegatorTotalRewardsResponse defines the properties of
// QueryDelegatorTotalRewards query's response.
type QueryDelegatorTotalRewardsResponse struct {
	Rewards []DelegationDelegatorReward `json:"rewards" yaml:"rewards"`
	Total   DecCoins                    `json:"total" yaml:"total"`
}

// NewQueryDelegatorTotalRewardsResponse constructs a QueryDelegatorTotalRewardsResponse
func NewQueryDelegatorTotalRewardsResponse(rewards []DelegationDelegatorReward,
	total DecCoins) QueryDelegatorTotalRewardsResponse {
	return QueryDelegatorTotalRewardsResponse{Rewards: rewards, Total: total}
}

func (res QueryDelegatorTotalRewardsResponse) String() string {
	out := "Delegator Total Rewards:\n"
	out += "  Rewards:"
	for _, reward := range res.Rewards {
		out += fmt.Sprintf(`  
	ValidatorAddress: %s
	Reward: %s`, reward.ValidatorAddress, reward.Reward)
	}
	out += fmt.Sprintf("\n  Total: %s\n", res.Total)
	return strings.TrimSpace(out)
}

// DelegationDelegatorReward defines the properties
// of a delegator's delegation reward.
type DelegationDelegatorReward struct {
	ValidatorAddress AccountID `json:"validator_account" yaml:"validator_account"`
	Reward           DecCoins  `json:"reward" yaml:"reward"`
}

// NewDelegationDelegatorReward constructs a DelegationDelegatorReward.
func NewDelegationDelegatorReward(valAddr AccountID, reward DecCoins) DelegationDelegatorReward {
	return DelegationDelegatorReward{ValidatorAddress: valAddr, Reward: reward}
}

type WithDrawAddrInfo struct {
	WithDrawAddress  AccountID `json:"withdraw_account" yaml:"withdraw_account"`
	ValidatorAddress AccountID `json:"validator_account" yaml:"validator_account"`
	DelegatorAddress AccountID `json:"delegator_account" yaml:"delegator_account"`
}

func NewWithDrawAddrInfo(withdrawAddr, valAddr, delAddr AccountID) WithDrawAddrInfo {
	return WithDrawAddrInfo{WithDrawAddress: withdrawAddr, ValidatorAddress: valAddr, DelegatorAddress: delAddr}
}

func (res WithDrawAddrInfo) String() string {
	out := "withdraw info:\n"
	out += fmt.Sprintf("WithDrawAddress:%s ,DelegatorAddress:%s", res.WithDrawAddress, res.DelegatorAddress)
	out += "\n"

	return strings.TrimSpace(out)
}
