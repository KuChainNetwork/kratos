package simulation

import (
	"fmt"
	"math/rand"

	"github.com/KuChainNetwork/kuchain/chain/transaction/helpers"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	kuSim "github.com/KuChainNetwork/kuchain/test/simulation"
	"github.com/KuChainNetwork/kuchain/x/distribution/keeper"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	stakingKeeper "github.com/KuChainNetwork/kuchain/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation" //fuck bugs by cancer
)

// Simulation operation weights constants
const (
	OpWeightMsgSetWithdrawAddress          = "op_weight_msg_set_withdraw_address"
	OpWeightMsgWithdrawDelegationReward    = "op_weight_msg_withdraw_delegation_reward"
	OpWeightMsgWithdrawValidatorCommission = "op_weight_msg_withdraw_validator_commission"
	OpWeightMsgFundCommunityPool           = "op_weight_msg_fund_community_pool"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simulation.AppParams, cdc *codec.Codec, ak types.AccountKeeperAccountID,
	bk types.BankKeeperAccountID, k keeper.Keeper, sk types.StakingKPKeeper,
) types.SimulationWeightedOperations {

	var weightMsgSetWithdrawAddress int
	appParams.GetOrGenerate(cdc, OpWeightMsgSetWithdrawAddress, &weightMsgSetWithdrawAddress, nil,
		func(_ *rand.Rand) {
			weightMsgSetWithdrawAddress = simappparams.DefaultWeightMsgSetWithdrawAddress
		},
	)

	var weightMsgWithdrawDelegationReward int
	appParams.GetOrGenerate(cdc, OpWeightMsgWithdrawDelegationReward, &weightMsgWithdrawDelegationReward, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawDelegationReward = simappparams.DefaultWeightMsgWithdrawDelegationReward
		},
	)

	var weightMsgWithdrawValidatorCommission int
	appParams.GetOrGenerate(cdc, OpWeightMsgWithdrawValidatorCommission, &weightMsgWithdrawValidatorCommission, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawValidatorCommission = simappparams.DefaultWeightMsgWithdrawValidatorCommission
		},
	)

	var weightMsgFundCommunityPool int
	appParams.GetOrGenerate(cdc, OpWeightMsgFundCommunityPool, &weightMsgFundCommunityPool, nil,
		func(_ *rand.Rand) {
			weightMsgFundCommunityPool = simappparams.DefaultWeightMsgFundCommunityPool
		},
	)

	return types.SimulationWeightedOperations{
		types.SimulationNewWeightedOperation(
			weightMsgSetWithdrawAddress,
			SimulateMsgSetWithdrawAddress(ak, bk, k),
		),
		types.SimulationNewWeightedOperation(
			weightMsgWithdrawDelegationReward,
			SimulateMsgWithdrawDelegatorReward(ak, bk, k, sk),
		),
		types.SimulationNewWeightedOperation(
			weightMsgWithdrawValidatorCommission,
			SimulateMsgWithdrawValidatorCommission(ak, bk, k, sk),
		),
	}
}

// SimulateMsgSetWithdrawAddress generates a MsgSetWithdrawAddress with random values.
func SimulateMsgSetWithdrawAddress(ak types.AccountKeeperAccountID, bk types.BankKeeperAccountID, k keeper.Keeper) types.SimulationOperation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []types.SimulationAccount, chainID string,
	) (types.SimulationOperationMsg, []types.SimulationFutureOperation, error) {
		if !k.GetWithdrawAddrEnabled(ctx) {
			return types.SimulationNoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, _ := types.SimulationRandomAcc(r, accs)
		simToAccount, _ := types.SimulationRandomAcc(r, accs)

		simAId := chainTypes.NewAccountIDFromAccAdd(simAccount.Address)
		account := ak.GetAccount(ctx, simAId)

		aId := chainTypes.NewAccountIDFromAccAdd(account.GetAuth())
		spendable := bk.SpendableCoins(ctx, aId)

		fees, err := kuSim.RandomFees(r, ctx, spendable)
		if err != nil {
			return types.SimulationNoOpMsg(types.ModuleName), nil, err
		}

		simToAccountId := chainTypes.NewAccountIDFromAccAdd(simToAccount.Address)
		msg := types.NewMsgSetWithdrawAccountId(account.GetAuth(), simAId, simToAccountId)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{0}, // TODO: sim support new seq []uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return types.SimulationNoOpMsg(types.ModuleName), nil, err
		}

		return types.SimulationNewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgWithdrawDelegatorReward generates a MsgWithdrawDelegatorReward with random values.
func SimulateMsgWithdrawDelegatorReward(ak types.AccountKeeperAccountID, bk types.BankKeeperAccountID, k keeper.Keeper, sk types.StakingKPKeeper) types.SimulationOperation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []types.SimulationAccount, chainID string,
	) (types.SimulationOperationMsg, []types.SimulationFutureOperation, error) {
		simAccount, _ := types.SimulationRandomAcc(r, accs)
		delAccAddress := chainTypes.NewAccountIDFromAccAdd(simAccount.Address)
		delegations := sk.GetAllDelegatorDelegations(ctx, delAccAddress)
		if len(delegations) == 0 {
			return types.SimulationNoOpMsg(types.ModuleName), nil, nil
		}

		delegation := delegations[r.Intn(len(delegations))]

		validator := sk.Validator(ctx, delegation.GetValidatorAccountID())

		if validator == nil {
			return types.SimulationNoOpMsg(types.ModuleName), nil, fmt.Errorf("validator %s not found", delegation.GetValidatorAddr())
		}

		simAccountId := chainTypes.NewAccountIDFromAccAdd(simAccount.Address)
		account := ak.GetAccount(ctx, simAccountId)

		accountId := chainTypes.NewAccountIDFromAccAdd(account.GetAuth())
		spendable := bk.SpendableCoins(ctx, accountId)

		fees, err := kuSim.RandomFees(r, ctx, spendable)
		if err != nil {
			return types.SimulationNoOpMsg(types.ModuleName), nil, err
		}

		valId := validator.GetOperatorAccountID()
		msg := types.NewMsgWithdrawDelegatorReward(account.GetAuth(), simAccountId, valId)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{0}, // TODO: sim support new seq []uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return types.SimulationNoOpMsg(types.ModuleName), nil, err
		}

		return types.SimulationNewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgWithdrawValidatorCommission generates a MsgWithdrawValidatorCommission with random values.
func SimulateMsgWithdrawValidatorCommission(ak types.AccountKeeperAccountID, bk types.BankKeeperAccountID, k keeper.Keeper, sk types.StakingKPKeeper) types.SimulationOperation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []types.SimulationAccount, chainID string,
	) (types.SimulationOperationMsg, []types.SimulationFutureOperation, error) {

		validator, ok := stakingKeeper.RandomValidator(r, sk, ctx)
		if !ok {
			return types.SimulationNoOpMsg(types.ModuleName), nil, nil
		}

		commission := k.GetValidatorAccumulatedCommission(ctx, validator.GetOperatorAccountID())
		if commission.Commission.IsZero() {
			return types.SimulationNoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, found := types.SimulationFindAccount(accs, sdk.AccAddress(validator.GetOperator()))
		if !found {
			return types.SimulationNoOpMsg(types.ModuleName), nil, fmt.Errorf("validator %s not found", validator.GetOperator())
		}

		simAccountId := chainTypes.NewAccountIDFromAccAdd(simAccount.Address)
		account := ak.GetAccount(ctx, simAccountId)

		accountId := chainTypes.NewAccountIDFromAccAdd(account.GetAuth())
		spendable := bk.SpendableCoins(ctx, accountId)

		fees, err := kuSim.RandomFees(r, ctx, spendable)
		if err != nil {
			return types.SimulationNoOpMsg(types.ModuleName), nil, err
		}

		valId, _ := chainTypes.NewAccountIDFromStr(string(validator.GetOperator())) //bugs, staking interface
		msg := types.NewMsgWithdrawValidatorCommission(account.GetAuth(), valId)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{0}, // TODO: sim support new seq []uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return types.SimulationNoOpMsg(types.ModuleName), nil, err
		}

		return types.SimulationNewOperationMsg(msg, true, ""), nil, nil
	}
}
