package keeper_test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/staking/keeper"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/staking/exported"
)

// TODO integrate with test_common.go helper (CreateTestInput)
// setup helper function - creates two validators
func setupHelper(t *testing.T, power int64) (sdk.Context, *keeper.Keeper, types.Params) {

	// setup
	wallet := simapp.NewWallet()
	_, _, _, _, _, _, app := NewTestApp(wallet)
	keeper := app.StakeKeeper()
	keeper = keeper.EmptyHooks()
	ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	params := keeper.GetParams(ctx)
	numVals := int64(3)
	amt := exported.TokensFromConsensusPower(power)
	bondAmount := exported.TokensFromConsensusPower(power * numVals)

	bondedPool := keeper.GetBondedPool(ctx)
	notBondedPool := keeper.GetNotBondedPool(ctx)
	app.AssetKeeper().IssueCoinPower(ctx, bondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), bondAmount)))
	app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), bondAmount)))
	app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), chainTypes.NewCoins(chainTypes.NewCoin(keeper.BondDenom(ctx), bondAmount)))

	// add numVals validators
	for i := int64(0); i < numVals; i++ {
		validator := types.NewValidator(Accd[i], PKs[i], types.Description{})
		validator, _ = validator.AddTokensFromDel(amt)
		validator = TestingUpdateValidator(app, ctx, validator, true)
		keeper.SetValidatorByConsAddr(ctx, validator)
	}

	return ctx, keeper, params
}

//_________________________________________________________________________________

func TestSlash(t *testing.T) {
	Convey("TestRevocation", t, func() {
		// setup
		ctx, keeper, _ := setupHelper(t, 10)
		addr := Accd[0]
		consAddr := sdk.ConsAddress(PKs[0].Address())

		// initial state
		val, found := keeper.GetValidator(ctx, addr)
		require.True(t, found)
		require.False(t, val.IsJailed())

		// test jail
		keeper.Jail(ctx, consAddr)
		val, found = keeper.GetValidator(ctx, addr)
		require.True(t, found)
		require.True(t, val.IsJailed())

		// test unjail
		keeper.Unjail(ctx, consAddr)
		val, found = keeper.GetValidator(ctx, addr)
		require.True(t, found)
		require.False(t, val.IsJailed())
	})
	Convey("TestSlashUnbondingDelegation", t, func() {
		ctx, keeper, _ := setupHelper(t, 10)
		fraction := sdk.NewDecWithPrec(5, 1)

		// set an unbonding delegation with expiration timestamp (beyond which the
		// unbonding delegation shouldn't be slashed)
		ubd := types.NewUnbondingDelegation(Accd[50], Accd[0], 0,
			time.Unix(5, 0), sdk.NewInt(10))

		keeper.SetUnbondingDelegation(ctx, ubd)

		// unbonding started prior to the infraction height, stakw didn't contribute
		slashAmount := keeper.SlashUnbondingDelegation(ctx, ubd, 1, fraction)
		require.Equal(t, int64(0), slashAmount.Int64())

		// after the expiration time, no longer eligible for slashing
		ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(10, 0)})
		keeper.SetUnbondingDelegation(ctx, ubd)
		slashAmount = keeper.SlashUnbondingDelegation(ctx, ubd, 0, fraction)
		require.Equal(t, int64(0), slashAmount.Int64())

		// test valid slash, before expiration timestamp and to which stake contributed
		oldUnbondedPool := keeper.TotalNotBondedTokens(ctx)

		ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(0, 0)})
		keeper.SetUnbondingDelegation(ctx, ubd)
		slashAmount = keeper.SlashUnbondingDelegation(ctx, ubd, 0, fraction) // TODO coin not enough
		require.Equal(t, int64(5), slashAmount.Int64())
		ubd, found := keeper.GetUnbondingDelegation(ctx, Accd[50], Accd[0])
		require.True(t, found)
		require.Len(t, ubd.Entries, 1)

		// initial balance unchanged
		require.Equal(t, sdk.NewInt(10), ubd.Entries[0].InitialBalance)

		// balance decreased
		require.Equal(t, sdk.NewInt(5), ubd.Entries[0].Balance)
		newUnbondedPool := keeper.TotalNotBondedTokens(ctx)
		diffTokens := oldUnbondedPool.Sub(newUnbondedPool).Int64()
		require.Equal(t, int64(5), diffTokens)
	})
	Convey("TestSlashRedelegation", t, func() {
		ctx, keeper, _ := setupHelper(t, 10)
		fraction := sdk.NewDecWithPrec(5, 1)
		oldBondedAmount := keeper.TotalBondedTokens(ctx)

		// set a redelegation with an expiration timestamp beyond which the
		// redelegation shouldn't be slashed
		rd := types.NewRedelegation(Accd[50], Accd[0], Accd[1], 0,
			time.Unix(5, 0), sdk.NewInt(10), sdk.NewDec(10))
		keeper.SetRedelegation(ctx, rd)

		// set the associated delegation
		del := types.NewDelegation(Accd[50], Accd[1], sdk.NewDec(10))
		keeper.SetDelegation(ctx, del)

		// started redelegating prior to the current height, stake didn't contribute to infraction
		validator, found := keeper.GetValidator(ctx, Accd[1])
		require.True(t, found)
		slashAmount := keeper.SlashRedelegation(ctx, validator, rd, 1, fraction)
		require.Equal(t, int64(0), slashAmount.Int64())

		// after the expiration time, no longer eligible for slashing
		ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(10, 0)})
		keeper.SetRedelegation(ctx, rd)
		validator, found = keeper.GetValidator(ctx, Accd[1])
		require.True(t, found)
		slashAmount = keeper.SlashRedelegation(ctx, validator, rd, 0, fraction)
		require.Equal(t, int64(0), slashAmount.Int64())

		// test valid slash, before expiration timestamp and to which stake contributed
		ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(0, 0)})
		keeper.SetRedelegation(ctx, rd)
		validator, found = keeper.GetValidator(ctx, Accd[1])
		require.True(t, found)
		slashAmount = keeper.SlashRedelegation(ctx, validator, rd, 0, fraction)
		require.Equal(t, int64(5), slashAmount.Int64())
		rd, found = keeper.GetRedelegation(ctx, Accd[50], Accd[0], Accd[1])
		require.True(t, found)
		require.Len(t, rd.Entries, 1)

		// end block
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		require.Equal(t, 1, len(updates))

		// initialbalance unchanged
		require.Equal(t, sdk.NewInt(10), rd.Entries[0].InitialBalance)

		// shares decreased
		del, found = keeper.GetDelegation(ctx, Accd[50], Accd[1])
		require.True(t, found)
		require.Equal(t, int64(5), del.Shares.RoundInt64())
		newBondedAmount := keeper.TotalBondedTokens(ctx)
		require.Equal(t, oldBondedAmount.Sub(slashAmount), newBondedAmount)
	})
	Convey("TestSlashAtFutureHeight", t, func() {
		ctx, keeper, _ := setupHelper(t, 10)
		consAddr := sdk.ConsAddress(PKs[0].Address())
		fraction := sdk.NewDecWithPrec(5, 1)
		require.Panics(t, func() { keeper.Slash(ctx, consAddr, 10, 10, fraction) })
	})
	Convey("TestSlashAtNegativeHeight", t, func() {
		ctx, keeper, _ := setupHelper(t, 10)
		consAddr := sdk.ConsAddress(PKs[0].Address())
		fraction := sdk.NewDecWithPrec(5, 1)

		//oldBondedAmount := keeper.TotalBondedTokens(ctx)
		oldBondedAmount := keeper.TotalBondedTokens(ctx)
		validator, found := keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)
		keeper.Slash(ctx, consAddr, -2, 10, fraction)

		// read updated state
		validator, found = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)

		// end block
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		require.Equal(t, 1, len(updates), "cons addr: %v, updates: %v", []byte(consAddr), updates)

		validator, _ = keeper.GetValidator(ctx, validator.OperatorAccount)
		// power decreased
		require.Equal(t, int64(5), validator.GetConsensusPower())
		// pool bonded shares decreased
		newBondedAmount := keeper.TotalBondedTokens(ctx)
		diffTokens := oldBondedAmount.Sub(newBondedAmount)
		require.Equal(t, exported.TokensFromConsensusPower(5), diffTokens)
	})
	Convey("TestSlashValidatorAtCurrentHeight", t, func() {
		ctx, keeper, _ := setupHelper(t, 10)
		consAddr := sdk.ConsAddress(PKs[0].Address())
		fraction := sdk.NewDecWithPrec(5, 1)

		oldBondedAmount := keeper.TotalBondedTokens(ctx)
		validator, found := keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)
		keeper.Slash(ctx, consAddr, ctx.BlockHeight(), 10, fraction)

		// read updated state
		validator, found = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)
		newBondedAmount := keeper.TotalBondedTokens(ctx)

		// end block
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		require.Equal(t, 1, len(updates), "cons addr: %v, updates: %v", []byte(consAddr), updates)

		validator, _ = keeper.GetValidator(ctx, validator.OperatorAccount)
		// power decreased
		require.Equal(t, int64(5), validator.GetConsensusPower())
		// pool bonded shares decreased
		diffTokens := oldBondedAmount.Sub(newBondedAmount)
		require.Equal(t, exported.TokensFromConsensusPower(5), diffTokens)
	})
	Convey("TestSlashWithUnbondingDelegation", t, func() {
		ctx, keeper, _ := setupHelper(t, 10)
		consAddr := sdk.ConsAddress(PKs[0].Address())
		fraction := sdk.NewDecWithPrec(5, 1)

		// set an unbonding delegation with expiration timestamp beyond which the
		// unbonding delegation shouldn't be slashed
		ubdTokens := exported.TokensFromConsensusPower(4)
		ubd := types.NewUnbondingDelegation(Accd[50], Accd[0], 11,
			time.Unix(0, 0), ubdTokens)
		keeper.SetUnbondingDelegation(ctx, ubd)

		// slash validator for the first time
		ctx = ctx.WithBlockHeight(12)
		oldBondedAmount := keeper.TotalBondedTokens(ctx)
		validator, found := keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)
		keeper.Slash(ctx, consAddr, 10, 10, fraction)

		// end block
		updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		require.Equal(t, 1, len(updates))

		// read updating unbonding delegation
		ubd, found = keeper.GetUnbondingDelegation(ctx, Accd[50], Accd[0])
		require.True(t, found)
		require.Len(t, ubd.Entries, 1)
		// balance decreased
		require.Equal(t, exported.TokensFromConsensusPower(2), ubd.Entries[0].Balance)
		// read updated pool
		newBondedAmount := keeper.TotalBondedTokens(ctx)
		// bonded tokens burned
		diffTokens := oldBondedAmount.Sub(newBondedAmount)
		require.Equal(t, exported.TokensFromConsensusPower(3), diffTokens)
		// read updated validator
		validator, found = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)
		// power decreased by 3 - 6 stake originally bonded at the time of infraction
		// was still bonded at the time of discovery and was slashed by half, 4 stake
		// bonded at the time of discovery hadn't been bonded at the time of infraction
		// and wasn't slashed
		require.Equal(t, int64(7), validator.GetConsensusPower())

		// slash validator again
		ctx = ctx.WithBlockHeight(13)
		keeper.Slash(ctx, consAddr, 9, 10, fraction)
		ubd, found = keeper.GetUnbondingDelegation(ctx, Accd[50], Accd[0])
		require.True(t, found)
		require.Len(t, ubd.Entries, 1)
		// balance decreased again
		require.Equal(t, sdk.NewInt(0), ubd.Entries[0].Balance)
		// read updated pool
		newBondedAmount = keeper.TotalBondedTokens(ctx)
		// bonded tokens burned again
		diffTokens = oldBondedAmount.Sub(newBondedAmount)
		require.Equal(t, exported.TokensFromConsensusPower(6), diffTokens)
		// read updated validator
		validator, found = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)
		// power decreased by 3 again
		require.Equal(t, int64(4), validator.GetConsensusPower())

		// slash validator again
		// all originally bonded stake has been slashed, so this will have no effect
		// on the unbonding delegation, but it will slash stake bonded since the infraction
		// this may not be the desirable behaviour, ref https://github.com/cosmos/cosmos-sdk/issues/1440
		ctx = ctx.WithBlockHeight(13)
		keeper.Slash(ctx, consAddr, 9, 10, fraction)
		ubd, found = keeper.GetUnbondingDelegation(ctx, Accd[50], Accd[0])
		require.True(t, found)
		require.Len(t, ubd.Entries, 1)
		// balance unchanged
		require.Equal(t, sdk.NewInt(0), ubd.Entries[0].Balance)
		// read updated pool
		newBondedAmount = keeper.TotalBondedTokens(ctx)
		// bonded tokens burned again
		diffTokens = oldBondedAmount.Sub(newBondedAmount)
		require.Equal(t, exported.TokensFromConsensusPower(9), diffTokens)
		// read updated validator
		validator, found = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)
		// power decreased by 3 again
		require.Equal(t, int64(1), validator.GetConsensusPower())

		// slash validator again
		// all originally bonded stake has been slashed, so this will have no effect
		// on the unbonding delegation, but it will slash stake bonded since the infraction
		// this may not be the desirable behaviour, ref https://github.com/cosmos/cosmos-sdk/issues/1440
		ctx = ctx.WithBlockHeight(13)
		keeper.Slash(ctx, consAddr, 9, 10, fraction)
		ubd, found = keeper.GetUnbondingDelegation(ctx, Accd[50], Accd[0])
		require.True(t, found)
		require.Len(t, ubd.Entries, 1)
		// balance unchanged
		require.Equal(t, sdk.NewInt(0), ubd.Entries[0].Balance)
		// read updated pool
		newBondedAmount = keeper.TotalBondedTokens(ctx)
		// just 1 bonded token burned again since that's all the validator now has
		diffTokens = oldBondedAmount.Sub(newBondedAmount)
		require.Equal(t, exported.TokensFromConsensusPower(10), diffTokens)
		// apply TM updates
		keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		// read updated validator
		// power decreased by 1 again, validator is out of stake
		// validator should be in unbonding period
		validator, _ = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.Equal(t, validator.GetStatus(), exported.Unbonding)
	})
	Convey("TestSlashWithRedelegation", t, func() {
		ctx, keeper, _ := setupHelper(t, 10)
		consAddr := sdk.ConsAddress(PKs[0].Address())
		fraction := sdk.NewDecWithPrec(5, 1)
		// set a redelegation
		rdTokens := exported.TokensFromConsensusPower(6)
		rd := types.NewRedelegation(Accd[50], Accd[0], Accd[1], 11,
			time.Unix(0, 0), rdTokens, rdTokens.ToDec())
		keeper.SetRedelegation(ctx, rd)

		// set the associated delegation
		del := types.NewDelegation(Accd[50], Accd[1], rdTokens.ToDec())
		keeper.SetDelegation(ctx, del)

		oldBondedAmount := keeper.TotalBondedTokens(ctx)
		oldNotBondedAmount := keeper.TotalNotBondedTokens(ctx)

		// slash validator
		ctx = ctx.WithBlockHeight(12)
		validator, found := keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)

		require.NotPanics(t, func() { keeper.Slash(ctx, consAddr, 10, 10, fraction) })
		burnAmount := exported.TokensFromConsensusPower(10).ToDec().Mul(fraction).TruncateInt()

		// burn bonded tokens from only from delegations
		require.True(sdk.IntEq(t, oldBondedAmount.Sub(burnAmount), keeper.TotalBondedTokens(ctx)))
		require.True(sdk.IntEq(t, oldNotBondedAmount, keeper.TotalNotBondedTokens(ctx)))
		oldBondedAmount = keeper.TotalBondedTokens(ctx)

		// read updating redelegation
		rd, found = keeper.GetRedelegation(ctx, Accd[50], Accd[0], Accd[1])
		require.True(t, found)
		require.Len(t, rd.Entries, 1)
		// read updated validator
		validator, found = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)
		// power decreased by 2 - 4 stake originally bonded at the time of infraction
		// was still bonded at the time of discovery and was slashed by half, 4 stake
		// bonded at the time of discovery hadn't been bonded at the time of infraction
		// and wasn't slashed
		require.Equal(t, int64(8), validator.GetConsensusPower())

		// slash the validator again
		validator, found = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)

		require.NotPanics(t, func() { keeper.Slash(ctx, consAddr, 10, 10, sdk.OneDec()) })
		burnAmount = exported.TokensFromConsensusPower(7)

		// seven bonded tokens burned
		require.True(sdk.IntEq(t, oldBondedAmount.Sub(burnAmount), keeper.TotalBondedTokens(ctx)))
		require.True(sdk.IntEq(t, oldNotBondedAmount, keeper.TotalNotBondedTokens(ctx)))
		oldBondedAmount = keeper.TotalBondedTokens(ctx)

		// read updating redelegation
		rd, found = keeper.GetRedelegation(ctx, Accd[50], Accd[0], Accd[1])
		require.True(t, found)
		require.Len(t, rd.Entries, 1)
		// read updated validator
		validator, found = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)
		// power decreased by 4
		require.Equal(t, int64(4), validator.GetConsensusPower())

		// slash the validator again, by 100%
		ctx = ctx.WithBlockHeight(12)
		validator, found = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, found)

		require.NotPanics(t, func() { keeper.Slash(ctx, consAddr, 10, 10, sdk.OneDec()) })

		burnAmount = exported.TokensFromConsensusPower(10).ToDec().Mul(sdk.OneDec()).TruncateInt()
		burnAmount = burnAmount.Sub(sdk.OneDec().MulInt(rdTokens).TruncateInt())

		require.True(sdk.IntEq(t, oldBondedAmount.Sub(burnAmount), keeper.TotalBondedTokens(ctx)))
		require.True(sdk.IntEq(t, oldNotBondedAmount, keeper.TotalNotBondedTokens(ctx)))
		oldBondedAmount = keeper.TotalBondedTokens(ctx)

		// read updating redelegation
		rd, found = keeper.GetRedelegation(ctx, Accd[50], Accd[0], Accd[1])
		require.True(t, found)
		require.Len(t, rd.Entries, 1)
		// apply TM updates
		keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		// read updated validator
		// validator decreased to zero power, should be in unbonding period
		validator, _ = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.Equal(t, validator.GetStatus(), exported.Unbonding)

		// slash the validator again, by 100%
		// no stake remains to be slashed
		ctx = ctx.WithBlockHeight(12)
		// validator still in unbonding period
		validator, _ = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.Equal(t, validator.GetStatus(), exported.Unbonding)

		require.NotPanics(t, func() { keeper.Slash(ctx, consAddr, 10, 10, sdk.OneDec()) })

		// read updated pool
		require.True(sdk.IntEq(t, oldBondedAmount, keeper.TotalBondedTokens(ctx)))
		require.True(sdk.IntEq(t, oldNotBondedAmount, keeper.TotalNotBondedTokens(ctx)))

		// read updating redelegation
		rd, found = keeper.GetRedelegation(ctx, Accd[50], Accd[0], Accd[1])
		require.True(t, found)
		require.Len(t, rd.Entries, 1)
		// read updated validator
		// power still zero, still in unbonding period
		validator, _ = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.Equal(t, validator.GetStatus(), exported.Unbonding)
	})
	Convey("TestSlashBoth", t, func() {
		ctx, keeper, _ := setupHelper(t, 10)
		fraction := sdk.NewDecWithPrec(5, 1)

		// set a redelegation with expiration timestamp beyond which the
		// redelegation shouldn't be slashed
		rdATokens := exported.TokensFromConsensusPower(6)
		rdA := types.NewRedelegation(Accd[50], Accd[0], Accd[1], 11,
			time.Unix(0, 0), rdATokens,
			rdATokens.ToDec())
		keeper.SetRedelegation(ctx, rdA)

		// set the associated delegation
		delA := types.NewDelegation(Accd[50], Accd[1], rdATokens.ToDec())
		keeper.SetDelegation(ctx, delA)

		// set an unbonding delegation with expiration timestamp (beyond which the
		// unbonding delegation shouldn't be slashed)
		ubdATokens := exported.TokensFromConsensusPower(4)
		ubdA := types.NewUnbondingDelegation(Accd[50], Accd[0], 11,
			time.Unix(0, 0), ubdATokens)
		keeper.SetUnbondingDelegation(ctx, ubdA)

		oldBondedAmount := keeper.TotalBondedTokens(ctx)
		oldNotBondedAmount := keeper.TotalNotBondedTokens(ctx)

		// slash validator
		ctx = ctx.WithBlockHeight(12)
		validator, found := keeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(PKs[0]))
		require.True(t, found)
		consAddr0 := sdk.ConsAddress(PKs[0].Address())
		keeper.Slash(ctx, consAddr0, 10, 10, fraction)

		burnedNotBondedAmount := fraction.MulInt(ubdATokens).TruncateInt()
		burnedBondAmount := exported.TokensFromConsensusPower(10).ToDec().Mul(fraction).TruncateInt()
		burnedBondAmount = burnedBondAmount.Sub(burnedNotBondedAmount)

		// read updated pool
		require.True(sdk.IntEq(t, oldBondedAmount.Sub(burnedBondAmount), keeper.TotalBondedTokens(ctx)))
		require.True(sdk.IntEq(t, oldNotBondedAmount.Sub(burnedNotBondedAmount), keeper.TotalNotBondedTokens(ctx)))

		// read updating redelegation
		rdA, found = keeper.GetRedelegation(ctx, Accd[50], Accd[0], Accd[1])
		require.True(t, found)
		require.Len(t, rdA.Entries, 1)
		// read updated validator
		validator, found = keeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(PKs[0]))
		require.True(t, found)
		// power not decreased, all stake was bonded since
		require.Equal(t, int64(10), validator.GetConsensusPower())
	})
}
