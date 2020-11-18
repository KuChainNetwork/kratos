package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// CommissionRates defines the initial commission rates to be used for creating
// a validator.
type CommissionRates struct {
	Rate          Dec `json:"rate" yaml:"rate"`
	MaxRate       Dec `json:"max_rate" yaml:"max_rate"`
	MaxChangeRate Dec `json:"max_change_rate" yaml:"max_change_rate"`
}

func (cr CommissionRates) Equal(other CommissionRates) bool {
	return cr.Rate.Equal(other.Rate) && cr.MaxRate.Equal(other.MaxRate) && cr.MaxChangeRate.Equal(other.MaxChangeRate)
}

// NewCommissionRates returns an initialized validator commission rates.
func NewCommissionRates(rate, maxRate, maxChangeRate sdk.Dec) CommissionRates {
	return CommissionRates{
		Rate:          rate,
		MaxRate:       maxRate,
		MaxChangeRate: maxChangeRate,
	}
}

// Commission defines a commission parameters for a given validator.
type Commission struct {
	CommissionRates `json:"commission_rates" yaml:"commission_rates"`
	UpdateTime      time.Time `json:"update_time" yaml:"update_time"`
}

// NewCommission returns an initialized validator commission.
func NewCommission(rate, maxRate, maxChangeRate sdk.Dec) Commission {
	return Commission{
		CommissionRates: NewCommissionRates(rate, maxRate, maxChangeRate),
		UpdateTime:      time.Unix(0, 0).UTC(),
	}
}

// NewCommissionWithTime returns an initialized validator commission with a specified
// update time which should be the current block BFT time.
func NewCommissionWithTime(rate, maxRate, maxChangeRate sdk.Dec, updatedAt time.Time) Commission {
	return Commission{
		CommissionRates: NewCommissionRates(rate, maxRate, maxChangeRate),
		UpdateTime:      updatedAt,
	}
}

func (c Commission) Equal(other Commission) bool {
	return c.CommissionRates.Equal(other.CommissionRates) && c.UpdateTime.Equal(other.UpdateTime)
}

// String implements the Stringer interface for a Commission object.
func (c Commission) String() string {
	out, _ := yaml.Marshal(c)
	return string(out)
}

// String implements the Stringer interface for a CommissionRates object.
func (cr CommissionRates) String() string {
	out, _ := yaml.Marshal(cr)
	return string(out)
}

// Validate performs basic sanity validation checks of initial commission
// parameters. If validation fails, an SDK error is returned.
func (cr CommissionRates) Validate() error {
	switch {
	case cr.MaxRate.IsNegative():
		// max rate cannot be negative
		return ErrCommissionNegative

	case cr.MaxRate.GT(sdk.OneDec()):
		// max rate cannot be greater than 1
		return ErrCommissionHuge

	case cr.Rate.IsNegative():
		// rate cannot be negative
		return ErrCommissionNegative

	case cr.Rate.GT(cr.MaxRate):
		// rate cannot be greater than the max rate
		return ErrCommissionGTMaxRate

	case cr.MaxChangeRate.IsNegative():
		// change rate cannot be negative
		return ErrCommissionChangeRateNegative

	case cr.MaxChangeRate.GT(cr.MaxRate):
		// change rate cannot be greater than the max rate
		return ErrCommissionChangeRateGTMaxRate
	}

	return nil
}

// ValidateNewRate performs basic sanity validation checks of a new commission
// rate. If validation fails, an SDK error is returned.
func (c Commission) ValidateNewRate(newRate sdk.Dec, blockTime time.Time) error {
	switch {
	case blockTime.Sub(c.UpdateTime).Hours() < 24:
		// new rate cannot be changed more than once within 24 hours
		return ErrCommissionUpdateTime

	case newRate.IsNegative():
		// new rate cannot be negative
		return ErrCommissionNegative

	case newRate.GT(c.MaxRate):
		// new rate cannot be greater than the max rate
		return ErrCommissionGTMaxRate

	case newRate.Sub(c.Rate).GT(c.MaxChangeRate):
		// new rate % points change cannot be greater than the max change rate
		return ErrCommissionGTMaxChangeRate
	}

	return nil
}
