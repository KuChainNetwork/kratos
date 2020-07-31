package keeper_test

import (
	"testing"
	"time"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/x/staking"
	"github.com/KuChainNetwork/kuchain/x/staking/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKeeper(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("TestHandleNewValidator", t, func() {

		addAlice, _, _, accAlice, _, _, app := NewTestApp(wallet)
		keeper := app.SlashKeeper()
		stakeKeeper := app.StakeKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

		pk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepqn4usdx22zdntysj7n795xj77wrc62sytheeevr7zlna4yhwppdrs8mpds3")
		rightRate, _ := sdk.NewDecFromStr("0.65")
		err := CreateValidator(t, wallet, app, addAlice, accAlice, rightRate, pk, true)
		So(err, ShouldBeNil)

		bigAmount := chainTypes.NewInt64Coin(constants.DefaultBondDenom, 2100000000000000000)
		//alice D jack 50000000
		err = DelegationValidator(t, wallet, app, addAlice, accAlice, accAlice, bigAmount, true)
		So(err, ShouldBeNil)
		ctx = app.BaseApp.NewContext(true, abci.Header{Height: keeper.SignedBlocksWindow(ctx) + 1})

		// Now a validator, for two blocks
		keeper.HandleValidatorSignature(ctx, pk.Address(), 100, true)
		ctx = app.BaseApp.NewContext(true, abci.Header{Height: keeper.SignedBlocksWindow(ctx) + 2})

		keeper.HandleValidatorSignature(ctx, pk.Address(), 100, false)

		info, found := keeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(pk.Address()))
		require.True(t, found)
		require.Equal(t, int64(2), info.IndexOffset)
		require.Equal(t, int64(1), info.MissedBlocksCounter)
		require.Equal(t, time.Unix(0, 0).UTC(), info.JailedUntil)

		// validator should be bonded still, should not have been jailed or slashed
		validator, _ := stakeKeeper.GetValidatorByConsAddr(ctx, sdk.ConsAddress(pk.Address()))
		require.Equal(t, exported.Bonded, validator.GetStatus())
		require.Equal(t, bigAmount.Amount, stakeKeeper.TotalBondedTokens(ctx))
	})
	Convey("TestHandleAlreadyJailed", t, func() {
		// initial setup
		addAlice, _, _, accAlice, _, _, app := NewTestApp(wallet)
		keeper := app.SlashKeeper()
		stakeKeeper := app.StakeKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		power := int64(100)
		amt := exported.TokensFromConsensusPower(power)
		pk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepqn4usdx22zdntysj7n795xj77wrc62sytheeevr7zlna4yhwppdrs8mpds3")
		rightRate, _ := sdk.NewDecFromStr("0.65")
		err := CreateValidator(t, wallet, app, addAlice, accAlice, rightRate, pk, true)
		So(err, ShouldBeNil)
		initAsset := chainTypes.NewCoin(constants.DefaultBondDenom, amt)
		err = DelegationValidator(t, wallet, app, addAlice, accAlice, accAlice, initAsset, true)
		So(err, ShouldBeNil)
		ctx = app.BaseApp.NewContext(true, abci.Header{Height: keeper.SignedBlocksWindow(ctx) + 1})

		// 1000 first blocks OK
		height := int64(0)
		for ; height < keeper.SignedBlocksWindow(ctx); height++ {
			ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + height})

			keeper.HandleValidatorSignature(ctx, pk.Address(), power, true)
		}

		// 501 blocks missed
		for ; height < keeper.SignedBlocksWindow(ctx)+(keeper.SignedBlocksWindow(ctx)-keeper.MinSignedPerWindow(ctx))+1; height++ {
			ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + height})

			keeper.HandleValidatorSignature(ctx, pk.Address(), power, false)
		}

		// end block
		staking.EndBlocker(ctx, *stakeKeeper)

		// validator should have been jailed and slashed
		validator, _ := stakeKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk))
		require.Equal(t, exported.Unbonding, validator.GetStatus())

		// validator should have been slashed
		resultingTokens := amt.Sub(exported.TokensFromConsensusPower(1).Quo(sdk.NewInt(100)))
		require.Equal(t, resultingTokens, validator.GetTokens())

		// another block missed
		ctx = app.BaseApp.NewContext(true, abci.Header{Height: height})

		keeper.HandleValidatorSignature(ctx, pk.Address(), power, false)

		// validator should not have been slashed twice
		validator, _ = stakeKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk))
		require.Equal(t, resultingTokens, validator.GetTokens())
	})
	Convey("TestValidatorDippingInAndOut", t, func() {
		// initial setup
		// TestParams set the SignedBlocksWindow to 1000 and MaxMissedBlocksPerWindow to 500
		addAlice, addJack, _, accAlice, accJack, _, app := NewTestApp(wallet)
		keeper := app.SlashKeeper()
		stakeKeeper := app.StakeKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		params := stakeKeeper.GetParams(ctx)
		params.MaxValidators = 1
		stakeKeeper.SetParams(ctx, params)
		power := int64(100)
		amt := exported.TokensFromConsensusPower(power)
		pk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepqn4usdx22zdntysj7n795xj77wrc62sytheeevr7zlna4yhwppdrs8mpds3")
		consAddr := sdk.ConsAddress(pk.Address())
		rightRate, _ := sdk.NewDecFromStr("0.65")
		err := CreateValidator(t, wallet, app, addAlice, accAlice, rightRate, pk, true)
		So(err, ShouldBeNil)
		initAsset := chainTypes.NewCoin(constants.DefaultBondDenom, amt)
		err = DelegationValidator(t, wallet, app, addAlice, accAlice, accAlice, initAsset, true)
		So(err, ShouldBeNil)

		// 100 first blocks OK
		height := int64(10)
		for ; height < int64(100); height++ {
			ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + height})
			keeper.HandleValidatorSignature(ctx, pk.Address(), power, true)
		}
		// kick first validator out of validator set
		pk2, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepq0cm4j2wtny3x435zuc53zffk9fndj7f37xjkxv4lqdx4w4z3mayqmf9aef")
		err = CreateValidator(t, wallet, app, addJack, accJack, rightRate, pk2, true)
		So(err, ShouldBeNil)
		newAmt := sdk.TokensFromConsensusPower(101)
		jackAsset := chainTypes.NewCoin(constants.DefaultBondDenom, newAmt)
		err = DelegationValidator(t, wallet, app, addJack, accJack, accJack, jackAsset, true)
		So(err, ShouldBeNil)
		ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + height})

		validator, _ := stakeKeeper.GetValidator(ctx, accAlice)
		require.Equal(t, exported.Bonded, validator.Status)

		// 600 more blocks happened
		height = int64(700)
		ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + height})

		// validator added back in
		delTokens := exported.TokensFromConsensusPower(50)
		err = DelegationValidator(t, wallet, app, addJack, accJack, accJack, chainTypes.NewCoin(constants.DefaultBondDenom, delTokens), true)
		require.NoError(t, err)

		validator, _ = stakeKeeper.GetValidator(ctx, accAlice)
		require.Equal(t, exported.Bonded, validator.Status)
		newPower := int64(150)

		// validator misses a block
		keeper.HandleValidatorSignature(ctx, pk.Address(), newPower, false)
		height++

		// shouldn't be jailed/kicked yet
		validator, _ = stakeKeeper.GetValidator(ctx, accAlice)
		require.Equal(t, exported.Bonded, validator.Status)

		// validator misses 500 more blocks, 501 total
		latest := height
		for ; height < latest+51; height++ {
			ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + height})
			keeper.HandleValidatorSignature(ctx, pk.Address(), newPower, false)
		}

		// should now be jailed & kicked
		staking.EndBlocker(ctx, *stakeKeeper)
		validator, _ = stakeKeeper.GetValidator(ctx, accAlice)
		require.Equal(t, exported.Unbonding, validator.Status)

		// check all the signing information
		signInfo, found := keeper.GetValidatorSigningInfo(ctx, consAddr)
		require.True(t, found)
		require.Equal(t, int64(0), signInfo.MissedBlocksCounter)
		require.Equal(t, int64(0), signInfo.IndexOffset)
		// array should be cleared
		for offset := int64(0); offset < keeper.SignedBlocksWindow(ctx); offset++ {
			missed := keeper.GetValidatorMissedBlockBitArray(ctx, consAddr, offset)
			require.False(t, missed)
		}

		// some blocks pass
		height = int64(5000)
		ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + height})

		staking.EndBlocker(ctx, *stakeKeeper)

		// validator rejoins and starts signing again
		stakeKeeper.Unjail(ctx, consAddr)
		keeper.HandleValidatorSignature(ctx, pk.Address(), newPower, true)

		height++

		// validator should not be kicked since we reset counter/array when it was jailed
		staking.EndBlocker(ctx, *stakeKeeper)
		validator, _ = stakeKeeper.GetValidator(ctx, accAlice)
		require.Equal(t, exported.Bonded, validator.Status)

		// validator misses 501 blocks
		latest = height
		for ; height < latest+501; height++ {
			ctx = app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + height})
			keeper.HandleValidatorSignature(ctx, pk.Address(), newPower, false)
		}

		// validator should now be jailed & kicked
		staking.EndBlocker(ctx, *stakeKeeper)
		validator, _ = stakeKeeper.GetValidator(ctx, accAlice)
		require.Equal(t, exported.Unbonding, validator.Status)
	})

}
