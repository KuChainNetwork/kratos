package external

import (
	"github.com/KuChain-io/kuchain/x/staking/exported"
)

type StakingValidatorI = exported.ValidatorI
type StakingDelegationI = exported.DelegationI

var TokensFromConsensusPower = exported.TokensFromConsensusPower
