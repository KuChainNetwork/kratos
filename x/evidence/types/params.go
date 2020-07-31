package types

import (
	"fmt"
	"time"

	"github.com/KuChainNetwork/kuchain/x/evidence/external"
	"gopkg.in/yaml.v2"
)

// DONTCOVER

// Default parameter values
const (
	DefaultParamspace            = ModuleName
	DefaultMaxEvidenceAge        = 60 * 2 * time.Second
	DefaultDoblesignJailDuration = 60 * 60 * 24 * 14 * time.Second
)

// Parameter store keys
var (
	KeyMaxEvidenceAge         = []byte("MaxEvidenceAge")
	KeyDoubleSignJailDuration = []byte("DoubleSignJailDuration")

	// The Double Sign Jail period ends at Max Time supported by Amino
	// (Dec 31, 9999 - 23:59:59 GMT).
	DoubleSignJailEndTime = time.Unix(253402300799, 0)
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() external.ParamsKeyTable {
	return external.ParamNewKeyTable().RegisterParamSet(&Params{})
}

// Params defines the total set of parameters for the evidence module
type Params struct {
	MaxEvidenceAge         time.Duration `json:"max_evidence_age" yaml:"max_evidence_age"`
	DoubleSignJailDuration time.Duration `json:"double_sign_jail_duration" yaml:"double_sign_jail_duration"`
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() external.ParamSetPairs {
	return external.ParamSetPairs{
		external.ParamNewParamSetPair(KeyMaxEvidenceAge, &p.MaxEvidenceAge, validateMaxEvidenceAge),
		external.ParamNewParamSetPair(KeyDoubleSignJailDuration, &p.DoubleSignJailDuration, validateDoubleSignJailDuration),
	}
}

// DefaultParams returns the default parameters for the evidence module.
func DefaultParams() Params {
	return Params{
		MaxEvidenceAge:         DefaultMaxEvidenceAge,
		DoubleSignJailDuration: DefaultDoblesignJailDuration,
	}
}

func validateMaxEvidenceAge(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("max evidence age must be positive: %s", v)
	}

	return nil
}

func validateDoubleSignJailDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("max evidence age must be positive: %s", v)
	}

	return nil
}
