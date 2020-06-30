package external

import (
	"github.com/KuChain-io/kuchain/x/staking"
	"github.com/KuChain-io/kuchain/x/staking/exported"
	"github.com/KuChain-io/kuchain/x/staking/keeper"
)

type StakingValidatorl = exported.ValidatorI
type StakingDelegatel = exported.DelegationI

type StakingKeeper = keeper.Keeper

var StakingRandomValidator = keeper.RandomValidator

type StakingMsgCreateValidator = staking.MsgCreateValidator
type StakingDescription = staking.Description
type StakingMsgDelegate = staking.MsgDelegate

var StakingNewMsgCreateValidator = staking.NewMsgCreateValidator
var StakingNewMsgDelegate = staking.NewMsgDelegate

var TokensFromConsensusPower = exported.TokensFromConsensusPower
var DefaultBondDenom = exported.DefaultBondDenom
