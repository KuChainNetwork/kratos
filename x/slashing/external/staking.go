package external

import (
	"github.com/KuChainNetwork/kuchain/x/staking"
	"github.com/KuChainNetwork/kuchain/x/staking/exported"
	"github.com/KuChainNetwork/kuchain/x/staking/keeper"
)

type StakingValidatorl = exported.ValidatorI
type StakingDelegatel = exported.DelegationI

type StakingKeeper = keeper.Keeper

var RandomValidator = keeper.RandomValidator

type StakingMsgCreateValidator = staking.MsgCreateValidator
type StakingDescription = staking.Description
type StakingMsgDelegate = staking.MsgDelegate

var StakingNewMsgCreateValidator = staking.NewMsgCreateValidator
var StakingNewMsgDelegate = staking.NewMsgDelegate

var TokensFromConsensusPower = exported.TokensFromConsensusPower
var DefaultBondDenom = exported.DefaultBondDenom
