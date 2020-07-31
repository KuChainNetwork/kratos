package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/constants"
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	assettypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	sktypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	"testing"

	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	"github.com/stretchr/testify/require"

	"github.com/KuChainNetwork/kuchain/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestCalculateRewardsBasic(t *testing.T) {
	ctx, ak, k, sk, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)

	// create validator with 50% commission
	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))

	description := GetDescription()

	Acc1Name, _ := Acc1.ToName()
	Acc1Auth, _ := ak.GetAuth(ctx, Acc1Name)
	Acc2Name, _ := Acc2.ToName()
	Acc2pubk := AccPubk[Acc2Name.String()]

	kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
	msg := sktypes.NewKuMsgCreateValidator(Acc1Auth, Acc2, Acc2pubk, description, commission.MaxRate, Acc1)
	kuCtx = kuCtx.WithTransfMsg(msg)
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

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator and delegation

	val := sk.Validator(ctx, Acc2)

	del := sk.Delegation(ctx, Acc1, Acc2)

	// historical count should be 2 (once for validator init, once for delegation init)
	require.Equal(t, uint64(2), k.GetValidatorHistoricalReferenceCount(ctx))

	// end period
	endingPeriod := k.IncrementValidatorPeriod(ctx, val)

	// historical count should be 2 still
	require.Equal(t, uint64(2), k.GetValidatorHistoricalReferenceCount(ctx))

	// calculate delegation rewards
	rewards := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

	// rewards should be zero
	require.True(t, rewards.IsZero())

	// allocate some rewards
	initial := int64(10)
	tokens := chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDec(initial)}}
	k.AllocateTokensToValidator(ctx, val, tokens)

	// end period
	endingPeriod = k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards
	rewards = k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

	// rewards should be half the tokens
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDec(initial / 2)}}, rewards)

	// commission should be the other half
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDec(initial / 2)}}, k.GetValidatorAccumulatedCommission(ctx, Acc2).Commission)
}

func TestCalculateRewardsAfterSlash(t *testing.T) {
	ctx, ak, k, sk, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)

	// create validator with 50% commission
	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	description := GetDescription()

	Acc3Name, _ := Acc3.ToName()
	Acc3Auth, _ := ak.GetAuth(ctx, Acc3Name)
	Acc4Name, _ := Acc4.ToName()
	Acc4pubk := AccPubk[Acc4Name.String()]

	kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
	msg := sktypes.NewKuMsgCreateValidator(Acc3Auth, Acc4, Acc4pubk, description, commission.MaxRate, Acc3)
	kuCtx = kuCtx.WithTransfMsg(msg)
	res, err := sh(kuCtx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	valPower := int64(100)

	intNum, _ := sdk.NewIntFromString("100000000000000000000")
	initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, intNum))
	err = ask.Transfer(ctx, Acc3, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
	require.NoError(t, err)

	msg1 := sktypes.NewKuMsgDelegate(Acc3Auth, Acc3, Acc4, chainType.NewCoin(constants.DefaultBondDenom, intNum))
	kuCtx = kuCtx.WithTransfMsg(msg1)

	res, err = sh(kuCtx, msg1)
	require.NoError(t, err)
	require.NotNil(t, res)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator and delegation
	val := sk.Validator(ctx, Acc4)
	del := sk.Delegation(ctx, Acc3, Acc4)

	// end period
	endingPeriod := k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards
	rewards := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

	// rewards should be zero
	require.True(t, rewards.IsZero())

	// start out block height
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)

	// slash the validator by 50%
	addr := sdk.GetConsAddress(Acc4pubk)

	sk.Slash(ctx, addr, ctx.BlockHeight(), valPower, sdk.NewDecWithPrec(5, 1))

	// retrieve validator
	val = sk.Validator(ctx, Acc4)

	// increase block height
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)

	// allocate some rewards
	initial := sdk.TokensFromConsensusPower(10)
	tokens := chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial.ToDec()}}
	k.AllocateTokensToValidator(ctx, val, tokens)

	// end period
	endingPeriod = k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards
	rewards = k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

	// rewards should be half the tokens
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial.QuoRaw(2).ToDec()}}, rewards)

	// commission should be the other half
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial.QuoRaw(2).ToDec()}},
		k.GetValidatorAccumulatedCommission(ctx, Acc4).Commission)
}

func TestCalculateRewardsAfterManySlashes(t *testing.T) {
	ctx, ak, k, sk, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)

	// create validator with 50% commission
	power := int64(100)

	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))

	description := GetDescription()

	Acc5Name, _ := Acc5.ToName()
	Acc5Auth, _ := ak.GetAuth(ctx, Acc5Name)
	Acc6Name, _ := Acc4.ToName()
	Acc6pubk := AccPubk[Acc6Name.String()]

	kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
	msg := sktypes.NewKuMsgCreateValidator(Acc5Auth, Acc6, Acc6pubk, description, commission.MaxRate, Acc5)
	kuCtx = kuCtx.WithTransfMsg(msg)
	res, err := sh(kuCtx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	intNum, _ := sdk.NewIntFromString("100000000000000000000")
	initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, intNum))
	err = ask.Transfer(ctx, Acc5, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
	require.NoError(t, err)

	msg1 := sktypes.NewKuMsgDelegate(Acc5Auth, Acc5, Acc6, chainType.NewCoin(constants.DefaultBondDenom, intNum))
	kuCtx = kuCtx.WithTransfMsg(msg1)

	res, err = sh(kuCtx, msg1)
	require.NoError(t, err)
	require.NotNil(t, res)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator and delegation
	val := sk.Validator(ctx, Acc6)
	del := sk.Delegation(ctx, Acc5, Acc6)

	// end period
	endingPeriod := k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards
	rewards := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

	// rewards should be zero
	require.True(t, rewards.IsZero())

	// start out block height
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)

	// slash the validator by 50%

	addr := sdk.GetConsAddress(Acc6pubk)
	sk.Slash(ctx, addr, ctx.BlockHeight(), power, sdk.NewDecWithPrec(5, 1))

	// fetch the validator again
	val = sk.Validator(ctx, Acc6)

	// increase block height
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)

	// allocate some rewards
	initial := sdk.TokensFromConsensusPower(10)
	tokens := chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial.ToDec()}}
	k.AllocateTokensToValidator(ctx, val, tokens)

	// slash the validator by 50% again
	sk.Slash(ctx, addr, ctx.BlockHeight(), power/2, sdk.NewDecWithPrec(5, 1))

	// fetch the validator again
	val = sk.Validator(ctx, Acc6)

	// increase block height
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)

	// allocate some more rewards
	k.AllocateTokensToValidator(ctx, val, tokens)

	// end period
	endingPeriod = k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards
	rewards = k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

	// rewards should be half the tokens
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial.ToDec()}}, rewards)

	// commission should be the other half
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial.ToDec()}},
		k.GetValidatorAccumulatedCommission(ctx, Acc6).Commission)
}

func TestCalculateRewardsMultiDelegator(t *testing.T) {
	ctx, ak, k, sk, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)

	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))

	description := GetDescription()

	Acc7Name, _ := Acc7.ToName()
	Acc7Auth, _ := ak.GetAuth(ctx, Acc7Name)
	Acc8Name, _ := Acc8.ToName()
	Acc8pubk := AccPubk[Acc8Name.String()]
	Acc9Name, _ := Acc9.ToName()
	Acc9Auth, _ := ak.GetAuth(ctx, Acc9Name)

	kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
	msg := sktypes.NewKuMsgCreateValidator(Acc7Auth, Acc8, Acc8pubk, description, commission.MaxRate, Acc7)
	kuCtx = kuCtx.WithTransfMsg(msg)
	res, err := sh(kuCtx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	intNum, _ := sdk.NewIntFromString("100000000000000000000")
	initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, intNum))
	err = ask.Transfer(ctx, Acc7, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
	require.NoError(t, err)

	msg1 := sktypes.NewKuMsgDelegate(Acc7Auth, Acc7, Acc8, chainType.NewCoin(constants.DefaultBondDenom, intNum))
	kuCtx = kuCtx.WithTransfMsg(msg1)

	res, err = sh(kuCtx, msg1)
	require.NoError(t, err)
	require.NotNil(t, res)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator and delegation
	val := sk.Validator(ctx, Acc8)
	del1 := sk.Delegation(ctx, Acc7, Acc8)

	// allocate some rewards
	tokens := chainType.NewDecCoins(chainType.NewDecCoin(constants.DefaultBondDenom, intNum))
	k.AllocateTokensToValidator(ctx, val, tokens)

	// second delegation
	err = ask.Transfer(ctx, Acc9, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
	require.NoError(t, err)
	msg2 := sktypes.NewKuMsgDelegate(Acc9Auth, Acc9, Acc8, chainType.NewCoin(constants.DefaultBondDenom, intNum))
	res, err = sh(kuCtx, msg2)
	require.NoError(t, err)
	require.NotNil(t, res)

	del2 := sk.Delegation(ctx, Acc9, Acc8)

	// fetch updated validator
	val = sk.Validator(ctx, Acc8)

	// end block
	staking.EndBlocker(ctx, sk)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// allocate some more rewards
	k.AllocateTokensToValidator(ctx, val, tokens)

	// end period
	endingPeriod := k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards for del1
	rewards := k.CalculateDelegationRewards(ctx, val, del1, endingPeriod)

	// rewards for del1 should be 3/4 initial
	rIntNum, _ := sdk.NewIntFromString("75000000000000000000")
	rInitCoins := chainType.NewDecCoins(chainType.NewDecCoin(constants.DefaultBondDenom, rIntNum))
	require.Equal(t, rInitCoins, rewards)

	//calculate delegation rewards for del2
	rewards = k.CalculateDelegationRewards(ctx, val, del2, endingPeriod)

	// rewards for del2 should be 1/4 initial
	rIntNum2, _ := sdk.NewIntFromString("25000000000000000000")
	rInitCoins2 := chainType.NewDecCoins(chainType.NewDecCoin(constants.DefaultBondDenom, rIntNum2))
	require.Equal(t, rInitCoins2, rewards)

	// commission should be equal to initial (50% twice)
	rIntNum3, _ := sdk.NewIntFromString("100000000000000000000")
	rInitCoins3 := chainType.NewDecCoins(chainType.NewDecCoin(constants.DefaultBondDenom, rIntNum3))
	require.Equal(t, rInitCoins3, k.GetValidatorAccumulatedCommission(ctx, Acc8).Commission)
}

func TestWithdrawDelegationRewardsBasic(t *testing.T) {
	balancePower := int64(1000000001000000)
	power := int64(1000000)
	balanceTokens := sdk.TokensFromConsensusPower(balancePower)

	ctx, ak, k, sk, supplyKeeper, ask := CreateTestInputDefault(t, false, balancePower)
	sh := staking.NewHandler(sk)

	// set module Account coins
	distrAcc := supplyKeeper.GetModuleAccount(ctx, types.ModuleName)

	intNum, _ := sdk.NewIntFromString("1000000000000000000")
	initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, intNum))
	_, err := ask.IssueCoinPower(ctx, distrAcc.GetID(), initCoins)
	require.Nil(t, err)

	valTokens := sdk.TokensFromConsensusPower(power)
	// create validator with 50% commission
	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))

	description := GetDescription()

	Acc9Name, _ := Acc9.ToName()
	Acc9Auth, _ := ak.GetAuth(ctx, Acc9Name)
	Acc10Name, _ := Acc10.ToName()
	Acc10pubk := AccPubk[Acc10Name.String()]

	kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
	msg := sktypes.NewKuMsgCreateValidator(Acc9Auth, Acc10, Acc10pubk, description, commission.MaxRate, Acc9)
	kuCtx = kuCtx.WithTransfMsg(msg)
	res, err := sh(kuCtx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	err = ask.Transfer(ctx, Acc7, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
	require.NoError(t, err)

	msg1 := sktypes.NewKuMsgDelegate(Acc9Auth, Acc9, Acc10, chainType.NewCoin(constants.DefaultBondDenom, intNum))
	kuCtx = kuCtx.WithTransfMsg(msg1)

	res, err = sh(kuCtx, msg1)
	require.NoError(t, err)
	require.NotNil(t, res)

	// assert correct initial balance
	expTokens := balanceTokens.Sub(valTokens)
	RCoins, _ := ask.GetCoins(ctx, Acc10)

	require.Equal(t, chainType.Coins{chainType.NewCoin(constants.DefaultBondDenom, expTokens)}, RCoins)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator and delegation
	val := sk.Validator(ctx, Acc10)

	// allocate some rewards
	initial := sdk.TokensFromConsensusPower(power)
	tokens := chainType.DecCoins{chainType.NewDecCoin(constants.DefaultBondDenom, initial)}

	k.AllocateTokensToValidator(ctx, val, tokens)

	// historical count should be 2 (initial + latest for delegation)
	rc := k.GetValidatorHistoricalReferenceCount(ctx)
	require.Equal(t, uint64(2), rc)

	// withdraw rewards
	_, err = k.WithdrawDelegationRewards(ctx, Acc9, Acc10)
	require.Nil(t, err)

	// historical count should still be 2 (added one record, cleared one)
	require.Equal(t, uint64(2), k.GetValidatorHistoricalReferenceCount(ctx))

	// assert correct balance
	exp := balanceTokens.Sub(valTokens).Add(initial.QuoRaw(2))
	RCoins, _ = ask.GetCoins(ctx, Acc9)
	RCoinsPower := ask.GetCoinPowers(ctx, Acc9)
	for _, c := range RCoinsPower {
		RCoins = RCoins.Add(c)
	}
	require.Equal(t, assettypes.Coins{assettypes.NewCoin(constants.DefaultBondDenom, exp)}, RCoins)

	// withdraw commission
	_, err = k.WithdrawValidatorCommission(ctx, Acc10)
	require.Nil(t, err)

	// assert correct balance
	exp = balanceTokens.Sub(valTokens).Add(initial.QuoRaw(2))
	RCoins, _ = ask.GetCoins(ctx, Acc10)
	RCoinsPower = ask.GetCoinPowers(ctx, Acc10)
	for _, c := range RCoinsPower {
		RCoins = RCoins.Add(c)
	}
	require.Equal(t, assettypes.Coins{assettypes.NewCoin(constants.DefaultBondDenom, exp)}, RCoins)
}

func TestCalculateRewardsAfterManySlashesInSameBlock(t *testing.T) {
	ctx, ak, k, sk, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)

	// create validator with 50% commission
	power := int64(100)
	//valTokens := sdk.TokensFromConsensusPower(power)
	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))

	description := GetDescription()

	Acc11Name, _ := Acc11.ToName()
	Acc11Auth, _ := ak.GetAuth(ctx, Acc11Name)
	Acc11pubk := AccPubk[Acc11Name.String()]

	kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
	msg := sktypes.NewKuMsgCreateValidator(Acc11Auth, Acc11, Acc11pubk, description, commission.MaxRate, Acc11)
	kuCtx = kuCtx.WithTransfMsg(msg)
	res, err := sh(kuCtx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	intNum, _ := sdk.NewIntFromString("100000000000000000000")
	initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, intNum))
	err = ask.Transfer(ctx, Acc11, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
	require.NoError(t, err)

	msg1 := sktypes.NewKuMsgDelegate(Acc11Auth, Acc11, Acc11, chainType.NewCoin(constants.DefaultBondDenom, intNum))
	kuCtx = kuCtx.WithTransfMsg(msg1)

	res, err = sh(kuCtx, msg1)
	require.NoError(t, err)
	require.NotNil(t, res)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	//fetch validator and delegation
	val := sk.Validator(ctx, Acc11)
	del := sk.Delegation(ctx, Acc11, Acc11)

	// end period
	endingPeriod := k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards
	rewards := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

	// rewards should be zero
	require.True(t, rewards.IsZero())

	// start out block height
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)

	// allocate some rewards
	initial := sdk.TokensFromConsensusPower(10).ToDec()
	tokens := chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial}}
	k.AllocateTokensToValidator(ctx, val, tokens)

	//slash the validator by 50%
	addr := sdk.GetConsAddress(Acc11pubk)
	sk.Slash(ctx, addr, ctx.BlockHeight(), power, sdk.NewDecWithPrec(5, 1))

	// slash the validator by 50% again
	sk.Slash(ctx, addr, ctx.BlockHeight(), power/2, sdk.NewDecWithPrec(5, 1))

	// fetch the validator again
	val = sk.Validator(ctx, Acc11)

	// increase block height
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)

	// allocate some more rewards
	k.AllocateTokensToValidator(ctx, val, tokens)

	// end period
	endingPeriod = k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards
	rewards = k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

	// rewards should be half the tokens
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial}}, rewards)

	// commission should be the other half
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial}}, k.GetValidatorAccumulatedCommission(ctx, Acc11).Commission)
}

func TestCalculateRewardsMultiDelegatorMultiSlash(t *testing.T) {
	ctx, ak, k, sk, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)

	// create validator with 50% commission
	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	power := int64(100)
	//valTokens := sdk.TokensFromConsensusPower(power)

	description := GetDescription()

	Acc12Name, _ := Acc12.ToName()
	Acc12Auth, _ := ak.GetAuth(ctx, Acc12Name)
	Acc12pubk := AccPubk[Acc12Name.String()]

	Acc11Name, _ := Acc11.ToName()
	Acc11Auth, _ := ak.GetAuth(ctx, Acc11Name)

	kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
	msg := sktypes.NewKuMsgCreateValidator(Acc12Auth, Acc12, Acc12pubk, description, commission.MaxRate, Acc12)
	kuCtx = kuCtx.WithTransfMsg(msg)
	res, err := sh(kuCtx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	intNum, _ := sdk.NewIntFromString("100000000000000000000")
	initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, intNum))
	err = ask.Transfer(ctx, Acc12, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
	require.NoError(t, err)
	err = ask.Transfer(ctx, Acc11, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
	require.NoError(t, err)

	msg1 := sktypes.NewKuMsgDelegate(Acc12Auth, Acc12, Acc12, chainType.NewCoin(constants.DefaultBondDenom, intNum))
	kuCtx = kuCtx.WithTransfMsg(msg1)

	res, err = sh(kuCtx, msg1)
	require.NoError(t, err)
	require.NotNil(t, res)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator and delegation
	val := sk.Validator(ctx, Acc12)
	del1 := sk.Delegation(ctx, Acc12, Acc12)

	// allocate some rewards
	initial := sdk.TokensFromConsensusPower(30).ToDec()
	tokens := chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial}}
	k.AllocateTokensToValidator(ctx, val, tokens)

	// slash the validator
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)

	addr := sdk.GetConsAddress(Acc12pubk)
	sk.Slash(ctx, addr, ctx.BlockHeight(), power, sdk.NewDecWithPrec(5, 1))
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)

	// second delegation
	//delTokens := sdk.TokensFromConsensusPower(100)

	msg2 := sktypes.NewKuMsgDelegate(Acc11Auth, Acc11, Acc12, chainType.NewCoin(constants.DefaultBondDenom, intNum))
	kuCtx = kuCtx.WithTransfMsg(msg2)
	res, err = sh(kuCtx, msg2)
	require.NoError(t, err)
	require.NotNil(t, res)

	del2 := sk.Delegation(ctx, Acc11, Acc12)

	// end block
	staking.EndBlocker(ctx, sk)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// allocate some more rewards
	k.AllocateTokensToValidator(ctx, val, tokens)

	// slash the validator again
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)
	sk.Slash(ctx, addr, ctx.BlockHeight(), power, sdk.NewDecWithPrec(5, 1))
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)

	// fetch updated validator
	val = sk.Validator(ctx, Acc12)

	// end period
	endingPeriod := k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards for del1
	rewards := k.CalculateDelegationRewards(ctx, val, del1, endingPeriod)

	// rewards for del1 should be 2/3 initial (half initial first period, 1/6 initial second period)
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial.QuoInt64(2).Add(initial.QuoInt64(6))}}, rewards)

	// calculate delegation rewards for del2
	rewards = k.CalculateDelegationRewards(ctx, val, del2, endingPeriod)

	// rewards for del2 should be initial / 3
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial.QuoInt64(3)}}, rewards)

	// commission should be equal to initial (twice 50% commission, unaffected by slashing)
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: initial}}, k.GetValidatorAccumulatedCommission(ctx, Acc12).Commission)
}

func TestCalculateRewardsMultiDelegatorMultWithdraw(t *testing.T) {
	ctx, ak, k, sk, supplyKeeper, ask := CreateTestInputDefault(t, false, 1000)

	sh := staking.NewHandler(sk)
	initial := int64(20)

	// set module Account coins
	distrAcc := supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	intNum, _ := sdk.NewIntFromString("1000000000000000000")
	initCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, intNum))
	_, err := ask.IssueCoinPower(ctx, distrAcc.GetID(), initCoins)
	require.Nil(t, err)

	TDec := chainType.NewDecCoin(constants.DefaultBondDenom, chainType.NewInt(initial))
	tokens := chainType.DecCoins{TDec}

	// create validator with 50% commission
	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))

	description := GetDescription()

	Acc13Name, _ := Acc13.ToName()
	Acc13Auth, _ := ak.GetAuth(ctx, Acc13Name)
	Acc13pubk := AccPubk[Acc13Name.String()]

	kuCtx := chainType.NewKuMsgCtx(ctx, nil, nil)
	msg := sktypes.NewKuMsgCreateValidator(Acc13Auth, Acc13, Acc13pubk, description, commission.MaxRate, Acc13)
	kuCtx = kuCtx.WithTransfMsg(msg)
	res, err := sh(kuCtx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	err = ask.Transfer(ctx, Acc13, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
	require.NoError(t, err)
	err = ask.Transfer(ctx, Acc12, supplyKeeper.GetModuleAccount(ctx, staking.ModuleName).GetID(), initCoins)
	require.NoError(t, err)

	msg1 := sktypes.NewKuMsgDelegate(Acc13Auth, Acc13, Acc13, chainType.NewCoin(constants.DefaultBondDenom, intNum))
	kuCtx = kuCtx.WithTransfMsg(msg1)

	res, err = sh(kuCtx, msg1)
	require.NoError(t, err)
	require.NotNil(t, res)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator and delegation
	val := sk.Validator(ctx, Acc13)
	del1 := sk.Delegation(ctx, Acc13, Acc13)

	// allocate some rewards
	k.AllocateTokensToValidator(ctx, val, tokens)

	// historical count should be 2 (validator init, delegation init)
	require.Equal(t, uint64(2), k.GetValidatorHistoricalReferenceCount(ctx))

	// second delegation
	Acc12Name, _ := Acc12.ToName()
	Acc12Auth, _ := ak.GetAuth(ctx, Acc12Name)
	msg2 := sktypes.NewKuMsgDelegate(Acc12Auth, Acc12, Acc13, chainType.NewCoin(constants.DefaultBondDenom, intNum))
	kuCtx = kuCtx.WithTransfMsg(msg2)

	res, err = sh(kuCtx, msg2)
	require.NoError(t, err)
	require.NotNil(t, res)

	// historical count should be 3 (second delegation init)
	require.Equal(t, uint64(3), k.GetValidatorHistoricalReferenceCount(ctx))

	// fetch updated validator
	val = sk.Validator(ctx, Acc13)
	del2 := sk.Delegation(ctx, Acc12, Acc13)

	// end block
	staking.EndBlocker(ctx, sk)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// allocate some more rewards
	k.AllocateTokensToValidator(ctx, val, tokens)

	// first delegator withdraws
	k.WithdrawDelegationRewards(ctx, Acc13, Acc13)

	// second delegator withdraws
	k.WithdrawDelegationRewards(ctx, Acc12, Acc13)

	// historical count should be 3 (validator init + two delegations)
	require.Equal(t, uint64(3), k.GetValidatorHistoricalReferenceCount(ctx))

	// validator withdraws commission
	k.WithdrawValidatorCommission(ctx, Acc13)

	// end period
	endingPeriod := k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards for del1
	rewards := k.CalculateDelegationRewards(ctx, val, del1, endingPeriod)

	// rewards for del1 should be zero
	require.True(t, rewards.IsZero())

	// calculate delegation rewards for del2
	rewards = k.CalculateDelegationRewards(ctx, val, del2, endingPeriod)

	// rewards for del2 should be zero
	require.True(t, rewards.IsZero())

	// commission should be zero
	require.True(t, k.GetValidatorAccumulatedCommission(ctx, Acc13).Commission.IsZero())

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// allocate some more rewards
	k.AllocateTokensToValidator(ctx, val, tokens)

	// first delegator withdraws again
	k.WithdrawDelegationRewards(ctx, Acc13, Acc13)

	// end period
	endingPeriod = k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards for del1
	rewards = k.CalculateDelegationRewards(ctx, val, del1, endingPeriod)

	// rewards for del1 should be zero
	require.True(t, rewards.IsZero())

	// calculate delegation rewards for del2
	rewards = k.CalculateDelegationRewards(ctx, val, del2, endingPeriod)

	// rewards for del2 should be 1/4 initial
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDec(initial / 4)}}, rewards)

	// commission should be half initial
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDec(initial / 2)}}, k.GetValidatorAccumulatedCommission(ctx, Acc13).Commission)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// allocate some more rewards
	k.AllocateTokensToValidator(ctx, val, tokens)

	// withdraw commission
	k.WithdrawValidatorCommission(ctx, Acc13)

	// end period
	endingPeriod = k.IncrementValidatorPeriod(ctx, val)

	// calculate delegation rewards for del1
	rewards = k.CalculateDelegationRewards(ctx, val, del1, endingPeriod)

	// rewards for del1 should be 1/4 initial
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDec(initial / 4)}}, rewards)

	// calculate delegation rewards for del2
	rewards = k.CalculateDelegationRewards(ctx, val, del2, endingPeriod)

	// rewards for del2 should be 1/2 initial
	require.Equal(t, chainType.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDec(initial / 2)}}, rewards)

	// commission should be zero
	require.True(t, k.GetValidatorAccumulatedCommission(ctx, Acc13).Commission.IsZero())
}
