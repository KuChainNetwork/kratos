package types

// DelegatorStartingInfo starting info for a delegator reward period
// tracks the previous validator period, the delegation's amount
// of staking token, and the creation height (to check later on
// if any slashes have occurred)
// NOTE that even though validators are slashed to whole staking tokens, the
// delegators within the validator may be left with less than a full token,
// thus sdk.Dec is used
type DelegatorStartingInfo struct {
	PreviousPeriod uint64 `json:"previous_period,omitempty" yaml:"previous_period"`
	Stake          Dec    `json:"stake" yaml:"stake"`
	Height         uint64 `json:"creation_height" yaml:"creation_height"`
}

// create a new DelegatorStartingInfo
func NewDelegatorStartingInfo(previousPeriod uint64, stake Dec, height uint64) DelegatorStartingInfo {
	return DelegatorStartingInfo{
		PreviousPeriod: previousPeriod,
		Stake:          stake,
		Height:         height,
	}
}
