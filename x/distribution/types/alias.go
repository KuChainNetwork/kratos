package types

import (
	"github.com/KuChain-io/kuchain/chain/types"
	GovTypes "github.com/KuChain-io/kuchain/x/gov/types"
	"github.com/KuChain-io/kuchain/x/staking"
	"github.com/KuChain-io/kuchain/x/staking/exported"
	StakingExported "github.com/KuChain-io/kuchain/x/staking/exported"
	StakingKP "github.com/KuChain-io/kuchain/x/staking/keeper"
	StakingTypes "github.com/KuChain-io/kuchain/x/staking/types"
	"github.com/KuChain-io/kuchain/x/supply"
	Sim "github.com/cosmos/cosmos-sdk/x/simulation"
)

type (
	AccountID  = types.AccountID
	Coin       = types.Coin
	Coins      = types.Coins
	KuMsg      = types.KuMsg
	AccAddress = types.AccAddress
	Name       = types.Name
)

var (
	NewAccountIDFromStoreKey = types.NewAccountIDFromStoreKey
	MustName                 = types.MustName
)

const (
	AccIDStoreKeyLen = types.AccIDStoreKeyLen
)

type (
	DelegationI = exported.DelegationI
	ValidatorI  = exported.ValidatorI
)

type (
	StakingExportedValidatorI  = StakingExported.ValidatorI
	StakingExportedDelegationI = StakingExported.DelegationI
)

type (
	StakingDelegation        = staking.Delegation
	StakingDescription       = staking.Description
	StakingKPKeeper          = StakingKP.Keeper
	StakingTypesStakingHooks = StakingTypes.StakingHooks
)

var (
	StakingNewHandler            = staking.NewHandler
	StakingNewCommissionRates    = staking.NewCommissionRates
	StakingNewMsgCreateValidator = staking.NewMsgCreateValidator
	StakingEndBlocker            = staking.EndBlocker
	StakingNewMsgDelegate        = staking.NewMsgDelegate
)

var (
	SupplyRegisterCodec         = supply.RegisterCodec
	SupplyNewModuleAddress      = supply.NewModuleAddress
	SupplyNewEmptyModuleAccount = supply.NewEmptyModuleAccount
)

type (
	SimulationWeightedOperations = Sim.WeightedOperations
	SimulationContentSimulatorFn = Sim.ContentSimulatorFn
	SimulationAccount            = Sim.Account
	SimulationOperation          = Sim.Operation
	SimulationOperationMsg       = Sim.OperationMsg
	SimulationFutureOperation    = Sim.FutureOperation
	SimulationAppParams          = Sim.AppParams
)

var (
	SimulationNoOpMsg              = Sim.NoOpMsg
	SimulationRandomAcc            = Sim.RandomAcc
	SimulationRandPositiveInt      = Sim.RandPositiveInt
	SimulationRandStringOfLength   = Sim.RandStringOfLength
	SimulationNewWeightedOperation = Sim.NewWeightedOperation
	SimulationRandomFees           = Sim.RandomFees
	SimulationNewOperationMsg      = Sim.NewOperationMsg
	SimulationFindAccount          = Sim.FindAccount
	SimulationRandSubsetCoins      = Sim.RandSubsetCoins
)

type (
	GovTypesHandler = GovTypes.Handler
	GovTypesContent = GovTypes.Content
)

var (
	GovTypesRegisterProposalType      = GovTypes.RegisterProposalType
	GovTypesRegisterProposalTypeCodec = GovTypes.RegisterProposalTypeCodec
	GovTypesValidateAbstract          = GovTypes.ValidateAbstract
	GovTypesNewKuMsgSubmitProposal    = GovTypes.NewKuMsgSubmitProposal
)
