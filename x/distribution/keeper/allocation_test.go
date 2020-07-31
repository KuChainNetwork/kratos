package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/constants"
	chainType "github.com/KuChainNetwork/kuchain/chain/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KuChainNetwork/kuchain/x/staking"
	sktypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestAllocateTokensToValidatorWithCommission(t *testing.T) {
	ctx, ak, k, sk, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)

	// create validator with 50% commission
	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))

	description := GetDescription()

	Acc1Name, _ := Acc1.ToName()
	Acc1Auth, _ := ak.GetAuth(ctx, Acc1Name)
	Acc2Name, _ := Acc2.ToName()
	Acc2pubk := AccPubk[Acc2Name.String()]

	//test Acc is validator
	msg := sktypes.NewKuMsgCreateValidator(Acc1Auth, Acc2, Acc2pubk, description, commission.MaxRate, Acc1)
	kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)

	res, err := sh(kuCtx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
	err = ask.Transfer(ctx, Acc1, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
	require.NoError(t, err)

	msg1 := sktypes.NewKuMsgDelegate(Acc1Auth, Acc1, Acc2, chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
	kuCtx = kuCtx.WithTransfMsg(msg1)

	res, err = sh(kuCtx, msg1)
	require.NoError(t, err)
	require.NotNil(t, res)

	val := sk.Validator(ctx, Acc2)

	// allocate tokens
	tokens := chainType.DecCoins{
		{Denom: constants.DefaultBondDenom, Amount: sdk.NewDec(10)},
	}
	k.AllocateTokensToValidator(ctx, val, tokens)

	// check commission
	expected := chainType.DecCoins{
		{Denom: constants.DefaultBondDenom, Amount: sdk.NewDec(5)},
	}

	ru := k.GetValidatorAccumulatedCommission(ctx, val.GetOperatorAccountID())
	require.Equal(t, expected, ru.Commission)

	// check current rewards
	require.Equal(t, expected, k.GetValidatorCurrentRewards(ctx, val.GetOperatorAccountID()).Rewards)
}

func TestAllocateTokensToManyValidators(t *testing.T) {
	ctx, ak, k, sk, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)

	{
		// create validator with 50% commission
		commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))

		description := GetDescription()

		Acc1Name, _ := Acc1.ToName()
		Acc1Auth, _ := ak.GetAuth(ctx, Acc1Name)
		Acc2Name, _ := Acc2.ToName()
		Acc2pubk := AccPubk[Acc2Name.String()]

		//test Acc is validator
		msg := sktypes.NewKuMsgCreateValidator(Acc1Auth, Acc2, Acc2pubk, description, commission.MaxRate, Acc1)

		kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
		res, err := sh(kuCtx, msg)
		require.NoError(t, err)
		require.NotNil(t, res)

		initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
		err = ask.Transfer(ctx, Acc1, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
		require.NoError(t, err)

		msg1 := sktypes.NewKuMsgDelegate(Acc1Auth, Acc1, Acc2, chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
		kuCtx = kuCtx.WithTransfMsg(msg1)

		res, err = sh(kuCtx, msg1)
		require.NoError(t, err)
		require.NotNil(t, res)

	}
	{
		// create second validator with 0% commission
		commission := staking.NewCommissionRates(sdk.NewDec(0), sdk.NewDec(0), sdk.NewDec(0))

		description := GetDescription()

		Acc3Name, _ := Acc3.ToName()
		Acc3Auth, _ := ak.GetAuth(ctx, Acc3Name)
		Acc4Name, _ := Acc4.ToName()
		Acc4pubk := AccPubk[Acc4Name.String()]

		//test Acc is validator
		msg := sktypes.NewKuMsgCreateValidator(Acc3Auth, Acc4, Acc4pubk, description, commission.MaxRate, Acc3)

		kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
		res, err := sh(kuCtx, msg)
		require.NoError(t, err)
		require.NotNil(t, res)

		initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
		err = ask.Transfer(ctx, Acc3, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
		require.NoError(t, err)

		msg1 := sktypes.NewKuMsgDelegate(Acc3Auth, Acc3, Acc4, chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
		kuCtx = kuCtx.WithTransfMsg(msg1)

		res, err = sh(kuCtx, msg1)
		require.NoError(t, err)
		require.NotNil(t, res)

		// assert initial state: zero outstanding rewards, zero community pool, zero commission, zero current rewards
		require.True(t, k.GetValidatorOutstandingRewards(ctx, Acc2).Rewards.IsZero())
		require.True(t, k.GetValidatorOutstandingRewards(ctx, Acc4).Rewards.IsZero())
		require.True(t, k.GetFeePool(ctx).CommunityPool.IsZero())
		require.True(t, k.GetValidatorAccumulatedCommission(ctx, Acc2).Commission.IsZero())
		require.True(t, k.GetValidatorAccumulatedCommission(ctx, Acc4).Commission.IsZero())
		require.True(t, k.GetValidatorCurrentRewards(ctx, Acc2).Rewards.IsZero())
		require.True(t, k.GetValidatorCurrentRewards(ctx, Acc4).Rewards.IsZero())
	}

	{
		Acc2Name, _ := Acc2.ToName()
		Acc2pubk := AccPubk[Acc2Name.String()]
		Acc4Name, _ := Acc4.ToName()
		Acc4pubk := AccPubk[Acc4Name.String()]

		abciValA := abci.Validator{
			Address: Acc2pubk.Address(),
			Power:   100,
		}
		abciValB := abci.Validator{
			Address: Acc4pubk.Address(),
			Power:   100,
		}

		feeCollector := supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName)
		require.NotNil(t, feeCollector)
		{ // allocate tokens as if both had voted and second was proposer
			fees := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, sdk.NewInt(100)))
			_, err := ask.IssueCoinPower(ctx, feeCollector.GetID(), fees)
			require.Nil(t, err)
		}

		votes := []abci.VoteInfo{
			{
				Validator:       abciValA,
				SignedLastBlock: true,
			},
			{
				Validator:       abciValB,
				SignedLastBlock: true,
			},
		}
		addr := sdk.GetConsAddress(Acc2pubk)

		k.AllocateTokens(ctx, 200, 200, addr, votes)

		// 98 outstanding rewards (100 less 2 to community pool)
		r1 := k.GetValidatorOutstandingRewards(ctx, Acc2).Rewards
		require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDecWithPrec(515, 1)}}, r1)
		r2 := k.GetValidatorOutstandingRewards(ctx, Acc4).Rewards
		require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDecWithPrec(465, 1)}}, r2)

		//2 community pool coins
		CommunityPool := k.GetFeePool(ctx).CommunityPool
		require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDec(2)}}, CommunityPool)

		Commission := k.GetValidatorAccumulatedCommission(ctx, Acc2).Commission
		require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDecWithPrec(2575, 2)}}, Commission)
		// zero commission for second proposer
		require.True(t, k.GetValidatorAccumulatedCommission(ctx, Acc4).Commission.IsZero())
		require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDecWithPrec(2575, 2)}}, k.GetValidatorCurrentRewards(ctx, Acc2).Rewards)
		require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDecWithPrec(465, 1)}}, k.GetValidatorCurrentRewards(ctx, Acc4).Rewards)
	}
}

func TestAllocateTokensTruncation(t *testing.T) {
	communityTax := sdk.NewDec(0)
	ctx, ak, ask, k, sk, _, supplyKeeper := CreateTestInputAdvanced(t, false, 1000000, communityTax)
	sh := staking.NewHandler(sk)

	{
		// create validator with 10% commission
		commission := staking.NewCommissionRates(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(1, 1), sdk.NewDec(0))

		description := GetDescription()

		Acc7Name, _ := Acc7.ToName()
		Acc7Auth, _ := ak.GetAuth(ctx, Acc7Name)
		Acc8Name, _ := Acc8.ToName()
		Acc8pubk := AccPubk[Acc8Name.String()]

		//test Acc is validator
		msg := sktypes.NewKuMsgCreateValidator(Acc7Auth, Acc8, Acc8pubk, description, commission.MaxRate, Acc7)

		kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
		res, err := sh(kuCtx, msg)
		require.NoError(t, err)
		require.NotNil(t, res)

		initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
		err = ask.Transfer(ctx, Acc7, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
		require.NoError(t, err)

		msg1 := sktypes.NewKuMsgDelegate(Acc7Auth, Acc7, Acc8, chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
		kuCtx = kuCtx.WithTransfMsg(msg1)

		res, err = sh(kuCtx, msg1)
		require.NoError(t, err)
		require.NotNil(t, res)
	}

	{
		// create second validator with 10% commission
		commission := staking.NewCommissionRates(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(1, 1), sdk.NewDec(0))

		description := GetDescription()

		Acc9Name, _ := Acc9.ToName()
		Acc9Auth, _ := ak.GetAuth(ctx, Acc9Name)
		Acc10Name, _ := Acc10.ToName()
		Acc10pubk := AccPubk[Acc10Name.String()]

		//test Acc is validator
		msg := sktypes.NewKuMsgCreateValidator(Acc9Auth, Acc10, Acc10pubk, description, commission.MaxRate, Acc9)

		kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
		res, err := sh(kuCtx, msg)
		require.NoError(t, err)
		require.NotNil(t, res)

		initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
		err = ask.Transfer(ctx, Acc9, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
		require.NoError(t, err)

		msg1 := sktypes.NewKuMsgDelegate(Acc9Auth, Acc9, Acc8, chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
		kuCtx = kuCtx.WithTransfMsg(msg1)

		res, err = sh(kuCtx, msg1)
		require.NoError(t, err)
		require.NotNil(t, res)

	}

	{
		//// create third validator with 10% commission
		commission := staking.NewCommissionRates(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(1, 1), sdk.NewDec(0))

		description := GetDescription()

		Acc5Name, _ := Acc5.ToName()
		Acc5Auth, _ := ak.GetAuth(ctx, Acc5Name)
		Acc6Name, _ := Acc6.ToName()
		Acc6pubk := AccPubk[Acc6Name.String()]

		//test Acc is validator
		msg := sktypes.NewKuMsgCreateValidator(Acc5Auth, Acc6, Acc6pubk, description, commission.MaxRate, Acc5)

		kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
		res, err := sh(kuCtx, msg)
		require.NoError(t, err)
		require.NotNil(t, res)

		initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
		err = ask.Transfer(ctx, Acc5, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
		require.NoError(t, err)

		msg1 := sktypes.NewKuMsgDelegate(Acc5Auth, Acc5, Acc6, chainType.NewCoin(constants.DefaultBondDenom, chainType.NewInt(100)))
		kuCtx = kuCtx.WithTransfMsg(msg1)

		res, err = sh(kuCtx, msg1)
		require.NoError(t, err)
		require.NotNil(t, res)

	}

	Acc8Name, _ := Acc8.ToName()
	Acc8pubk := AccPubk[Acc8Name.String()]
	Acc10Name, _ := Acc10.ToName()
	Acc10pubk := AccPubk[Acc10Name.String()]
	Acc6Name, _ := Acc6.ToName()
	Acc6pubk := AccPubk[Acc6Name.String()]

	abciValA := abci.Validator{
		Address: Acc8pubk.Address(),
		Power:   10,
	}
	abciValB := abci.Validator{
		Address: Acc10pubk.Address(),
		Power:   10,
	}

	abciValС := abci.Validator{
		Address: Acc6pubk.Address(),
		Power:   11,
	}

	// assert initial state: zero outstanding rewards, zero community pool, zero commission, zero current rewards
	require.True(t, k.GetValidatorOutstandingRewards(ctx, Acc8).Rewards.IsZero())
	require.True(t, k.GetValidatorOutstandingRewards(ctx, Acc10).Rewards.IsZero())
	require.True(t, k.GetValidatorOutstandingRewards(ctx, Acc6).Rewards.IsZero())

	require.True(t, k.GetFeePool(ctx).CommunityPool.IsZero())

	require.True(t, k.GetValidatorAccumulatedCommission(ctx, Acc6).Commission.IsZero())
	require.True(t, k.GetValidatorAccumulatedCommission(ctx, Acc8).Commission.IsZero())

	require.True(t, k.GetValidatorCurrentRewards(ctx, Acc6).Rewards.IsZero())
	require.True(t, k.GetValidatorCurrentRewards(ctx, Acc8).Rewards.IsZero())

	feeCollector := supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	require.NotNil(t, feeCollector)
	{ // allocate tokens as if both had voted and second was proposer
		fees := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, sdk.NewInt(100)))
		_, err := ask.IssueCoinPower(ctx, feeCollector.GetID(), fees)
		require.Nil(t, err)
	}

	votes := []abci.VoteInfo{
		{
			Validator:       abciValA,
			SignedLastBlock: true,
		},
		{
			Validator:       abciValB,
			SignedLastBlock: true,
		},
		{
			Validator:       abciValС,
			SignedLastBlock: true,
		},
	}
	addr := sdk.GetConsAddress(Acc6pubk)
	k.AllocateTokens(ctx, 31, 31, addr, votes)

	require.True(t, k.GetValidatorOutstandingRewards(ctx, Acc6).Rewards.IsValid())
	require.True(t, k.GetValidatorOutstandingRewards(ctx, Acc8).Rewards.IsValid())
	require.True(t, k.GetValidatorOutstandingRewards(ctx, Acc10).Rewards.IsValid())
}
