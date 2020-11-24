package types

import (
	"fmt"
	"time"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/gov/external"
	paramtypes "github.com/KuChainNetwork/kuchain/x/params/types"
	stakingexport "github.com/KuChainNetwork/kuchain/x/staking/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

// Default period for deposits & voting
const (
	DefaultPeriod       time.Duration = time.Hour * 24 * 14 // 14 days
	DefaultPunishPeriod time.Duration = time.Hour * 24 * 7  // 7 days
)

// Default governance params
var (
	DefaultMinDepositTokens = external.TokensFromConsensusPower(500)
	DefaultQuorum           = sdk.NewDecWithPrec(334, 3)
	DefaultThreshold        = sdk.NewDecWithPrec(5, 1)
	DefaultVeto             = sdk.NewDecWithPrec(334, 3)
	DefaultEmergengcy       = sdk.NewDecWithPrec(667, 3)
	DefaultSlashFraction    = types.NewDec(1).Quo(types.NewDec(10000))
)

// Parameter store key
var (
	ParamStoreKeyDepositParams = []byte("depositparams")
	ParamStoreKeyVotingParams  = []byte("votingparams")
	ParamStoreKeyTallyParams   = []byte("tallyparams")
)

// ParamKeyTable - Key declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(ParamStoreKeyDepositParams, DepositParams{}, validateDepositParams),
		paramtypes.NewParamSetPair(ParamStoreKeyVotingParams, VotingParams{}, validateVotingParams),
		paramtypes.NewParamSetPair(ParamStoreKeyTallyParams, TallyParams{}, validateTallyParams),
	)
}

// DepositParams defines the params around deposits for governance
type DepositParams struct {
	// MinDeposit Minimum deposit for a proposal to enter voting period.
	MinDeposit Coins `json:"min_deposit,omitempty" yaml:"min_deposit,omitempty"`
	// MaxDepositPeriod Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
	MaxDepositPeriod time.Duration `json:"max_deposit_period,omitempty" yaml:"max_deposit_period,omitempty"`
}

// NewDepositParams creates a new DepositParams object
func NewDepositParams(minDeposit Coins, maxDepositPeriod time.Duration) DepositParams {
	return DepositParams{
		MinDeposit:       minDeposit,
		MaxDepositPeriod: maxDepositPeriod,
	}
}

// DefaultDepositParams default parameters for deposits
func DefaultDepositParams() DepositParams {
	return NewDepositParams(
		types.NewCoins(types.NewCoin(stakingexport.DefaultBondDenom, DefaultMinDepositTokens)),
		DefaultPeriod,
	)
}

// String implements stringer insterface
func (dp DepositParams) String() string {
	out, _ := yaml.Marshal(dp)
	return string(out)
}

// Equal checks equality of DepositParams
func (dp DepositParams) Equal(dp2 DepositParams) bool {
	return dp.MinDeposit.IsEqual(dp2.MinDeposit) && dp.MaxDepositPeriod == dp2.MaxDepositPeriod
}

func validateDepositParams(i interface{}) error {
	v, ok := i.(DepositParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.MinDeposit.IsValid() {
		return fmt.Errorf("invalid minimum deposit: %s", v.MinDeposit)
	}
	if v.MaxDepositPeriod <= 0 {
		return fmt.Errorf("maximum deposit period must be positive: %d", v.MaxDepositPeriod)
	}

	return nil
}

// TallyParams defines the params around Tallying votes in governance
type TallyParams struct {
	// Quorum Minimum percentage of total stake needed to vote for a result to be considered valid
	Quorum sdk.Dec `json:"quorum,omitempty" yaml:"quorum,omitempty"`
	// Threshold Minimum proportion of Yes votes for proposal to pass. Initial value: 0.5
	Threshold sdk.Dec `json:"threshold,omitempty" yaml:"threshold,omitempty"`
	// Veto Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	Veto sdk.Dec `json:"veto,omitempty" yaml:"veto,omitempty"`
	// Emergency Minimum proportion of votes for emergency passage.Initial value: 2/3
	Emergency sdk.Dec `json:"emergency,omitempty" yaml:"emergency,omitempty"`
	// MaxPunishPeriod Maximum punish for validator who donot vote for proposal
	MaxPunishPeriod time.Duration `json:"max_punish_period,omitempty" yaml:"max_punish_period,omitempty"`
	// SlashFraction slash fraction for Veto vote to slah validators.Initial value: 1/1000
	SlashFraction sdk.Dec `json:"slash_fraction,omitempty" yaml:"slash_fraction,omitempty"`
}

// NewTallyParams creates a new TallyParams object
func NewTallyParams(quorum, threshold, veto sdk.Dec, emergency sdk.Dec, punishPeriod time.Duration, slashFraction sdk.Dec) TallyParams {
	return TallyParams{
		Quorum:          quorum,
		Threshold:       threshold,
		Veto:            veto,
		Emergency:       emergency,
		MaxPunishPeriod: punishPeriod,
		SlashFraction:   slashFraction,
	}
}

// DefaultTallyParams default parameters for tallying
func DefaultTallyParams() TallyParams {
	return NewTallyParams(DefaultQuorum, DefaultThreshold, DefaultVeto, DefaultEmergengcy, DefaultPunishPeriod, DefaultSlashFraction)
}

// Equal checks equality of TallyParams
func (tp TallyParams) Equal(other TallyParams) bool {
	return tp.Quorum.Equal(other.Quorum) && tp.Threshold.Equal(other.Threshold) && tp.Veto.Equal(other.Veto)
}

// String implements stringer insterface
func (tp TallyParams) String() string {
	out, _ := yaml.Marshal(tp)
	return string(out)
}

func validateTallyParams(i interface{}) error {
	v, ok := i.(TallyParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Quorum.IsNegative() {
		return fmt.Errorf("quorom cannot be negative: %s", v.Quorum)
	}
	if v.Quorum.GT(sdk.OneDec()) {
		return fmt.Errorf("quorom too large: %s", v)
	}
	if !v.Threshold.IsPositive() {
		return fmt.Errorf("vote threshold must be positive: %s", v.Threshold)
	}
	if v.Threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("vote threshold too large: %s", v)
	}
	if !v.Veto.IsPositive() {
		return fmt.Errorf("veto threshold must be positive: %s", v.Threshold)
	}
	if v.Veto.GT(sdk.OneDec()) {
		return fmt.Errorf("veto threshold too large: %s", v)
	}

	return nil
}

// VotingParams defines the params around Voting in governance
type VotingParams struct {
	VotingPeriod time.Duration `json:"voting_period,omitempty" yaml:"voting_period,omitempty"` //  Length of the voting period.
}

// NewVotingParams creates a new VotingParams object
func NewVotingParams(votingPeriod time.Duration) VotingParams {
	return VotingParams{
		VotingPeriod: votingPeriod,
	}
}

// DefaultVotingParams default parameters for voting
func DefaultVotingParams() VotingParams {
	return NewVotingParams(DefaultPeriod)
}

// Equal checks equality of TallyParams
func (vp VotingParams) Equal(other VotingParams) bool {
	return vp.VotingPeriod == other.VotingPeriod
}

// String implements stringer interface
func (vp VotingParams) String() string {
	out, _ := yaml.Marshal(vp)
	return string(out)
}

func validateVotingParams(i interface{}) error {
	v, ok := i.(VotingParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.VotingPeriod <= 0 {
		return fmt.Errorf("voting period must be positive: %s", v.VotingPeriod)
	}

	return nil
}

// Params returns all of the governance params
type Params struct {
	VotingParams  VotingParams  `json:"voting_params" yaml:"voting_params"`
	TallyParams   TallyParams   `json:"tally_params" yaml:"tally_params"`
	DepositParams DepositParams `json:"deposit_params" yaml:"deposit_parmas"`
}

func (gp Params) String() string {
	return gp.VotingParams.String() + "\n" +
		gp.TallyParams.String() + "\n" + gp.DepositParams.String()
}

// NewParams creates a new gov Params instance
func NewParams(vp VotingParams, tp TallyParams, dp DepositParams) Params {
	return Params{
		VotingParams:  vp,
		DepositParams: dp,
		TallyParams:   tp,
	}
}

// DefaultParams default governance params
func DefaultParams() Params {
	return NewParams(DefaultVotingParams(), DefaultTallyParams(), DefaultDepositParams())
}
