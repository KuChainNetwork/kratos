package keeper

import (
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	assettypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KuChainNetwork/kuchain/x/distribution/types"
)

func TestSetWithdrawAddr(t *testing.T) {
	ctx, _, keeper, _, supplyKeeper, _ := CreateTestInputDefault(t, false, 1000)
	params := keeper.GetParams(ctx)
	params.WithdrawAddrEnabled = false
	keeper.SetParams(ctx, params)

	err := keeper.SetWithdrawAddr(ctx, Acc1, Acc11)
	require.NotNil(t, err)

	params.WithdrawAddrEnabled = true
	keeper.SetParams(ctx, params)

	err = keeper.SetWithdrawAddr(ctx, Acc1, Acc11)
	require.Nil(t, err)

	distrAcc := supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	keeper.blacklistedAddrs[distrAcc.GetID().String()] = true
	require.Error(t, keeper.SetWithdrawAddr(ctx, Acc1, distrAcc.GetID()))
}

func TestWithdrawValidatorCommission1(t *testing.T) {
	ctx, _, keeper, _, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)

	//set module Account coins
	myTokenName, _ := chainType.NewName("mytoken")
	stakeName, _ := chainType.NewName("stake")

	intNum, _ := sdk.NewIntFromString("100000000000000000000")
	intNumMax, _ := sdk.NewIntFromString("300000000000000000000")

	ask.Create(ctx, MasterName, myTokenName, assettypes.NewCoin("kuchain/mytoken", intNum),
		true, true, 0, assettypes.NewCoin("kuchain/mytoken", intNumMax), []byte("mytoken"))

	ask.Create(ctx, MasterName, stakeName, assettypes.NewCoin("kuchain/stake", intNum),
		true, true, 0, assettypes.NewCoin("kuchain/stake", intNumMax), []byte("stake"))

	intNum0, _ := sdk.NewIntFromString("100033333333333333")
	myTokenCoins := assettypes.Coins{assettypes.NewCoin("kuchain/mytoken", intNum0)}
	_, err := ask.IssueCoinPower(ctx, Master, myTokenCoins)
	require.Nil(t, err)

	stakeCoins := assettypes.Coins{assettypes.NewCoin("kuchain/stake", intNum0)}
	_, err = ask.IssueCoinPower(ctx, Master, stakeCoins)
	require.Nil(t, err)

	{
		distrAcc := supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
		//fmt.Println(distrAcc.GetID().String())

		myTokenCoins := assettypes.Coins{{Denom: "kuchain/mytoken", Amount: intNum0}}
		stakeCoins := assettypes.Coins{{Denom: "kuchain/stake", Amount: intNum0}}

		err := supplyKeeper.SendCoinsFromAccountToModule(ctx, Master, distrAcc.GetID().String(), myTokenCoins)
		require.Nil(t, err)

		err = supplyKeeper.SendCoinsFromAccountToModule(ctx, Master, distrAcc.GetID().String(), stakeCoins)
		require.Nil(t, err)

		// check initial balance
		balance := ask.GetCoinPowers(ctx, Acc3)
		expTokens := sdk.TokensFromConsensusPower(0)
		expCoins := chainType.NewCoins(chainType.NewCoin("kuchain/stake", expTokens), chainType.NewCoin("kuchain/mytoken", expTokens))

		//fmt.Println("e",expCoins, "b",balance)
		require.Equal(t, expCoins, balance)
	}

	{
		// set outstanding rewards
		var valCommission1 types.ValidatorOutstandingRewards
		for _, c := range stakeCoins {
			valCommission1.Rewards = valCommission1.Rewards.Add(chainType.NewDecCoin(c.Denom, c.Amount))
		}
		for _, c := range myTokenCoins {
			valCommission1.Rewards = valCommission1.Rewards.Add(chainType.NewDecCoin(c.Denom, c.Amount))
		}

		keeper.SetValidatorOutstandingRewards(ctx, Acc3, valCommission1)

		// set commission
		var valCommission2 types.ValidatorAccumulatedCommission
		for _, c := range stakeCoins {
			valCommission2.Commission = valCommission2.Commission.Add(chainType.NewDecCoin(c.Denom, c.Amount))
		}
		for _, c := range myTokenCoins {
			valCommission2.Commission = valCommission2.Commission.Add(chainType.NewDecCoin(c.Denom, c.Amount))
		}
		keeper.SetValidatorAccumulatedCommission(ctx, Acc3, valCommission2)
	}

	// withdraw commission
	keeper.WithdrawValidatorCommission(ctx, Acc3)

	// check balance increase
	balance := ask.GetCoinPowers(ctx, Acc3)
	//fmt.Println(balance)

	require.Equal(t, chainType.NewCoins(
		chainType.NewCoin("kuchain/mytoken", intNum0),
		chainType.NewCoin("kuchain/stake", intNum0),
	), balance)

	// check remainder
	remainder := keeper.GetValidatorAccumulatedCommission(ctx, Acc3)
	require.True(t, remainder.Commission.IsZero())
	require.True(t, true)
}

func TestWithdrawValidatorCommission2(t *testing.T) {

	myTokenCoin := types.DecCoin{Denom: "kuchain/mytoken", Amount: sdk.NewDec(120380).Quo(sdk.NewDec(3))}
	stakeCoin := types.DecCoin{Denom: "kuchain/stake", Amount: sdk.NewDec(900380).Quo(sdk.NewDec(7))}

	ctx, _, keeper, _, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)

	//set module Account coins
	myTokenName, _ := chainType.NewName("mytoken")
	stakeName, _ := chainType.NewName("stake")

	intNum, _ := sdk.NewIntFromString("100000000000000000000")
	intNumMax, _ := sdk.NewIntFromString("300000000000000000000")

	ask.Create(ctx, MasterName, myTokenName, assettypes.NewCoin("kuchain/mytoken", intNum),
		true, true, 0, assettypes.NewCoin("kuchain/mytoken", intNumMax), []byte("mytoken"))

	ask.Create(ctx, MasterName, stakeName, assettypes.NewCoin("kuchain/stake", intNum),
		true, true, 0, assettypes.NewCoin("kuchain/stake", intNumMax), []byte("stake"))

	intNum0, _ := sdk.NewIntFromString("100033333333333333")
	TokenCoins := assettypes.Coins{assettypes.NewCoin("kuchain/mytoken", intNum0)}
	_, err := ask.IssueCoinPower(ctx, Master, TokenCoins)
	require.Nil(t, err)

	StakeCoins := assettypes.Coins{assettypes.NewCoin("kuchain/stake", intNum0)}
	_, err = ask.IssueCoinPower(ctx, Master, StakeCoins)
	require.Nil(t, err)

	{
		distrAcc := supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
		//fmt.Println(distrAcc.GetID().String())

		myTokenCoins := assettypes.Coins{{Denom: "kuchain/mytoken", Amount: intNum0}}
		stakeCoins := assettypes.Coins{{Denom: "kuchain/stake", Amount: intNum0}}

		err := supplyKeeper.SendCoinsFromAccountToModule(ctx, Master, distrAcc.GetID().String(), myTokenCoins)
		require.Nil(t, err)

		err = supplyKeeper.SendCoinsFromAccountToModule(ctx, Master, distrAcc.GetID().String(), stakeCoins)
		require.Nil(t, err)

		// check initial balance
		balance := ask.GetCoinPowers(ctx, Acc3)
		expTokens := sdk.TokensFromConsensusPower(0)
		expCoins := chainType.NewCoins(chainType.NewCoin("kuchain/stake", expTokens), chainType.NewCoin("kuchain/mytoken", expTokens))

		//fmt.Println("e",expCoins, "b",balance)
		require.Equal(t, expCoins, balance)
	}

	{
		// set outstanding rewards
		var valCommission1 types.ValidatorOutstandingRewards
		valCommission1.Rewards = valCommission1.Rewards.Add(myTokenCoin)
		valCommission1.Rewards = valCommission1.Rewards.Add(stakeCoin)

		keeper.SetValidatorOutstandingRewards(ctx, Acc3, valCommission1)

		// set commission
		var valCommission2 types.ValidatorAccumulatedCommission
		valCommission2.Commission = valCommission2.Commission.Add(myTokenCoin)
		valCommission2.Commission = valCommission2.Commission.Add(stakeCoin)

		keeper.SetValidatorAccumulatedCommission(ctx, Acc3, valCommission2)
	}

	// withdraw commission
	keeper.WithdrawValidatorCommission(ctx, Acc3)
	{
		// check balance increase
		balance := ask.GetCoinPowers(ctx, Acc3)
		MainExp := chainType.NewCoins(
			chainType.NewCoin("kuchain/mytoken", sdk.NewInt(sdk.NewDec(120380).Quo(sdk.NewDec(3)).TruncateInt().Int64())),
			chainType.NewCoin("kuchain/stake", sdk.NewInt(sdk.NewDec(900380).Quo(sdk.NewDec(7)).TruncateInt().Int64())),
		)

		require.Equal(t, MainExp, balance)
	}

	{
		// check remainder
		remainder := keeper.GetValidatorAccumulatedCommission(ctx, Acc3)

		MainDecExp := chainType.NewDecCoins(
			chainType.NewDecCoin("kuchain/mytoken", sdk.NewInt(sdk.NewDec(120380).Quo(sdk.NewDec(3)).TruncateInt().Int64())),
			chainType.NewDecCoin("kuchain/stake", sdk.NewInt(sdk.NewDec(900380).Quo(sdk.NewDec(7)).TruncateInt().Int64())),
		)
		AllExp := chainType.NewDecCoins()
		AllExp = AllExp.Add(myTokenCoin).Add(stakeCoin)

		require.Equal(t, remainder.Commission, AllExp.Sub(MainDecExp))
		require.True(t, true)
	}

}

func TestGetTotalRewards(t *testing.T) {
	ctx, _, keeper, _, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)

	//set module Account coins
	myTokenName, _ := chainType.NewName("mytoken")
	stakeName, _ := chainType.NewName("stake")

	intNum, _ := sdk.NewIntFromString("100000000000000000000")
	intNumMax, _ := sdk.NewIntFromString("300000000000000000000")

	ask.Create(ctx, MasterName, myTokenName, assettypes.NewCoin("kuchain/mytoken", intNum),
		true, true, 0, assettypes.NewCoin("kuchain/mytoken", intNumMax), []byte("mytoken"))

	ask.Create(ctx, MasterName, stakeName, assettypes.NewCoin("kuchain/stake", intNum),
		true, true, 0, assettypes.NewCoin("kuchain/stake", intNumMax), []byte("stake"))

	{
		intNum0, _ := sdk.NewIntFromString("100033333333333333")
		IssMyTokenCoins := assettypes.Coins{assettypes.NewCoin("kuchain/mytoken", intNum0)}
		_, err := ask.IssueCoinPower(ctx, Master, IssMyTokenCoins)
		require.Nil(t, err)

		IssStakeCoins := assettypes.Coins{assettypes.NewCoin("kuchain/stake", intNum0)}
		_, err = ask.IssueCoinPower(ctx, Master, IssStakeCoins)
		require.Nil(t, err)
	}

	distrAcc := supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	//fmt.Println(distrAcc.GetID().String())

	myTokenCoins := assettypes.Coins{{Denom: "kuchain/mytoken", Amount: sdk.NewInt(int64(800000000))}}
	stakeCoins := assettypes.Coins{{Denom: "kuchain/stake", Amount: sdk.NewInt(int64(600000000))}}

	err := supplyKeeper.SendCoinsFromAccountToModule(ctx, Master, distrAcc.GetID().String(), myTokenCoins)
	require.Nil(t, err)

	err = supplyKeeper.SendCoinsFromAccountToModule(ctx, Master, distrAcc.GetID().String(), stakeCoins)
	require.Nil(t, err)

	////set outstanding rewards
	var valCommission1 types.ValidatorOutstandingRewards
	for _, c := range stakeCoins {
		valCommission1.Rewards = valCommission1.Rewards.Add(chainType.NewDecCoin(c.Denom, c.Amount))
	}
	for _, c := range myTokenCoins {
		valCommission1.Rewards = valCommission1.Rewards.Add(chainType.NewDecCoin(c.Denom, c.Amount))
	}

	keeper.SetValidatorOutstandingRewards(ctx, Acc7, valCommission1)
	keeper.SetValidatorOutstandingRewards(ctx, Acc8, valCommission1)

	expectedRewards := valCommission1.Rewards.MulDec(sdk.NewDec(2))
	totalRewards := keeper.GetTotalRewards(ctx)

	require.Equal(t, expectedRewards, totalRewards)
}
