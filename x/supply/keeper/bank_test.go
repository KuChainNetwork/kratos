package keeper_test

import (
	"fmt"
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account"
	"github.com/KuChainNetwork/kuchain/x/asset"
	keep "github.com/KuChainNetwork/kuchain/x/supply/keeper"
	"github.com/KuChainNetwork/kuchain/x/supply/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const initialPower = int64(100)

// create module accounts for testing
var (
	holderAcc     = types.NewEmptyModuleAccount(holder)
	burnerAcc     = types.NewEmptyModuleAccount(types.Burner, types.Burner)
	minterAcc     = types.NewEmptyModuleAccount(types.Minter, types.Minter)
	multiPermAcc  = types.NewEmptyModuleAccount(multiPerm, types.Burner, types.Minter, types.Staking)
	randomPermAcc = types.NewEmptyModuleAccount(randomPerm, "random")

	initTokens = sdk.TokensFromConsensusPower(initialPower)
	initCoins  = chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, initTokens))

	AccPubk = make(map[string]crypto.PubKey)
)

func CreateAccount(keeper account.Keeper, ctx sdk.Context, id chainType.AccountID) {
	Name, _ := id.ToName()
	newAccount := keeper.NewAccountByName(ctx, Name)

	privKey := secp256k1.GenPrivKey()
	address := chainType.AccAddress(privKey.PubKey().Address())

	newAccount.SetAuth(address)

	// set Account
	keeper.SetAccount(ctx, newAccount)
	AccPubk[Name.String()] = privKey.PubKey()
}

func getCoinsByName(ctx sdk.Context, sk keep.Keeper, moduleName string, ask asset.Keeper) chainType.Coins {
	//moduleAddress := sk.GetModuleAddress(moduleName)
	//macc := ak.GetAccount(ctx, moduleAddress)

	//ask.GetCoinsTotalSupply()
	//
	//if macc == nil {
	//	return sdk.Coins(nil)
	//}

	mAcc := sk.GetModuleAccount(ctx, moduleName)
	return ask.GetCoinPowers(ctx, mAcc.GetID())
}

func TestSendCoins(t *testing.T) {
	app, ctx := createTestApp(false)
	ak := app.AccountKeeper()

	Convey("test send coins to supply", t, func() {

		bName, _ := chainType.NewName("baseacc")
		baseAcc := ak.NewAccountByName(ctx, bName)

		app.SupplyKeeper().SetModuleAccount(ctx, holderAcc)
		app.SupplyKeeper().SetModuleAccount(ctx, burnerAcc)
		ak.SetAccount(ctx, baseAcc)

		{
			SymbolName, _ := chainType.NewName(constants.DefaultBondSymbol)
			TestMaster := constants.ChainMainNameStr
			MasterName, _ := chainType.NewName(TestMaster)
			Master := chainType.NewAccountIDFromName(MasterName)

			intNum1, _ := sdk.NewIntFromString("20000000000000000000000")
			intNum3, _ := sdk.NewIntFromString("60000000000000000000000")
			intMaxNum, _ := sdk.NewIntFromString("100000000000000000000000")

			err := app.AssetKeeper().Create(ctx, MasterName, SymbolName, chainType.NewIntCoreCoin(intMaxNum),
				true, true, true, 0, chainType.NewInt64CoreCoin(0), []byte("create"))
			So(err, ShouldBeNil)

			_, err = app.AssetKeeper().IssueCoinPower(ctx, Master, chainType.NewCoins(chainType.NewIntCoreCoin(intNum3)))
			So(err, ShouldBeNil)

			Coins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, intNum1))
			supplyAcc := app.SupplyKeeper().GetModuleAccount(ctx, types.ModuleName)
			err = app.AssetKeeper().SendCoinPower(ctx, Master, supplyAcc.GetID(), Coins)
			So(err, ShouldBeNil)

			HCoins := chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, initTokens.Mul(sdk.NewInt(3))))
			err = app.AssetKeeper().SendCoinPower(ctx, Master, holderAcc.GetID(), HCoins)
			So(err, ShouldBeNil)
		}

		require.Panics(t, func() {
			app.SupplyKeeper().SendCoinsFromModuleToModule(ctx, "", holderAcc.String(), initCoins)
		})

		require.Panics(t, func() {
			app.SupplyKeeper().SendCoinsFromModuleToModule(ctx, types.Burner, "", initCoins)
		})

		require.Panics(t, func() {
			app.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, "", baseAcc.GetID(), initCoins)
		})

		bAcc := baseAcc.GetID()
		hAccStr := holderAcc.GetName().String()
		err := app.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, hAccStr, bAcc, initCoins.Add(initCoins...))
		So(err, ShouldBeNil)

		err = app.SupplyKeeper().SendCoinsFromModuleToModule(ctx, holderAcc.GetName().String(), types.Burner, initCoins)
		So(err, ShouldBeNil)

		holderAccCoins := getCoinsByName(ctx, *app.SupplyKeeper(), holderAcc.GetName().String(), *app.AssetKeeper())
		fmt.Println(holderAccCoins.String())
		require.Equal(t, chainType.Coins{}, holderAccCoins)

		BurnerAccCoins := getCoinsByName(ctx, *app.SupplyKeeper(), types.Burner, *app.AssetKeeper())
		require.Equal(t, initCoins, BurnerAccCoins)

		err = app.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.Burner, bAcc, initCoins)
		So(err, ShouldBeNil)
		BurnerAccCoins = getCoinsByName(ctx, *app.SupplyKeeper(), types.Burner, *app.AssetKeeper())
		require.Equal(t, chainType.Coins{}, BurnerAccCoins)

		BAccCoins := app.AssetKeeper().GetCoinPowers(ctx, baseAcc.GetID())
		require.Equal(t, chainType.NewCoins(chainType.NewCoin(constants.DefaultBondDenom, initTokens.Mul(sdk.NewInt(3)))), BAccCoins)

		err = app.SupplyKeeper().SendCoinsFromAccountToModule(ctx, baseAcc.GetID(), types.Burner, initCoins)
		So(err, ShouldBeNil)

		BAccCoins = app.AssetKeeper().GetCoinPowers(ctx, baseAcc.GetID())
		require.Equal(t, chainType.Coins{chainType.NewCoin(constants.DefaultBondDenom, initTokens.Mul(sdk.NewInt(2)))}, BAccCoins)

		BurnerAccCoins = getCoinsByName(ctx, *app.SupplyKeeper(), types.Burner, *app.AssetKeeper())
		require.Equal(t, initCoins, BurnerAccCoins)
	})
}

func TestMintCoins(t *testing.T) {
	app, ctx := createTestApp(false)
	keeper := app.SupplyKeeper()

	Convey("test mint coins to supply", t, func() {

		keeper.SetModuleAccount(ctx, burnerAcc)
		keeper.SetModuleAccount(ctx, minterAcc)
		keeper.SetModuleAccount(ctx, multiPermAcc)
		keeper.SetModuleAccount(ctx, randomPermAcc)

		initialSupply := keeper.GetSupply(ctx)

		{
			SymbolName, _ := chainType.NewName(constants.DefaultBondSymbol)
			TestMaster := constants.ChainMainNameStr
			MasterName, _ := chainType.NewName(TestMaster)

			intMaxNum, _ := sdk.NewIntFromString("100000000000000000000000")

			err := app.AssetKeeper().Create(ctx, MasterName, SymbolName, chainType.NewIntCoreCoin(intMaxNum),
				true, true, true, 0, chainType.NewInt64CoreCoin(0), []byte("create"))

			So(err, ShouldBeNil)
		}

		require.Panics(t, func() { keeper.MintCoins(ctx, "", &initCoins) }, "no module account")
		require.Panics(t, func() { keeper.MintCoins(ctx, types.Burner, &initCoins) }, "invalid permission")

		err := keeper.MintCoins(ctx, types.Minter, &chainType.Coins{chainType.Coin{Denom: "denom", Amount: chainType.NewInt(-10)}})
		require.Error(t, err, "insufficient coins")

		require.Panics(t, func() { keeper.MintCoins(ctx, randomPerm, &initCoins) })

		err = keeper.MintCoins(ctx, types.Minter, &initCoins)
		So(err, ShouldBeNil)

		mintCoins := getCoinsByName(ctx, *keeper, types.Minter, *app.AssetKeeper())
		require.Equal(t, initCoins, mintCoins)
		require.Equal(t, initialSupply.GetTotal().Add(initCoins...), keeper.GetSupply(ctx).GetTotal())

		// test same functionality on module account with multiple permissions
		initialSupply = keeper.GetSupply(ctx)

		err = keeper.MintCoins(ctx, multiPermAcc.GetName().String(), &initCoins)
		So(err, ShouldBeNil)
		require.Equal(t, initCoins, getCoinsByName(ctx, *keeper, multiPermAcc.GetName().String(), *app.AssetKeeper()))
		require.Equal(t, initialSupply.GetTotal().Add(initCoins...), keeper.GetSupply(ctx).GetTotal())

		require.Panics(t, func() { keeper.MintCoins(ctx, types.Burner, &initCoins) })
	})
}

func TestBurnCoins(t *testing.T) {
	app, ctx := createTestApp(false)
	keeper := *app.SupplyKeeper()

	Convey("test burn coins to supply", t, func() {
		{
			SymbolName, _ := chainType.NewName(constants.DefaultBondSymbol)
			TestMaster := constants.ChainMainNameStr
			MasterName, _ := chainType.NewName(TestMaster)

			intNum2, _ := sdk.NewIntFromString("80000000000000000000000")

			err := app.AssetKeeper().Create(ctx, MasterName, SymbolName, chainType.NewIntCoreCoin(intNum2),
				true, true, true, 0, chainType.NewInt64CoreCoin(0), []byte("create"))

			So(err, ShouldBeNil)
		}

		_, err := app.AssetKeeper().IssueCoinPower(ctx, burnerAcc.GetID(), initCoins)
		So(err, ShouldBeNil)

		keeper.SetModuleAccount(ctx, burnerAcc)
		initialSupply := keeper.GetSupply(ctx)

		nName, _ := chainType.NewName("")
		nAcc := app.AccountKeeper().NewAccountByName(ctx, nName)
		require.Panics(t, func() { keeper.BurnCoins(ctx, nAcc.GetID(), initCoins) }, "no module account")

		MintName, _ := chainType.NewName(types.Minter)
		mintAcc := app.AccountKeeper().NewAccountByName(ctx, MintName)
		require.Panics(t, func() { keeper.BurnCoins(ctx, mintAcc.GetID(), initCoins) }, "invalid permission")

		randomPermName, _ := chainType.NewName(randomPerm)
		randomPermNameAcc := app.AccountKeeper().NewAccountByName(ctx, randomPermName)
		require.Panics(t, func() { keeper.BurnCoins(ctx, randomPermNameAcc.GetID(), initialSupply.GetTotal()) }, "random permission")

		bName, _ := chainType.NewName(types.Burner)
		bAcc := app.AccountKeeper().NewAccountByName(ctx, bName)

		fmt.Println(initialSupply.GetTotal())
		err = keeper.BurnCoins(ctx, bAcc.GetID(), initialSupply.GetTotal())
		So(err, ShouldBeNil)

		bAccCoins := getCoinsByName(ctx, keeper, types.Burner, *app.AssetKeeper())

		_, err = app.AssetKeeper().IssueCoinPower(ctx, burnerAcc.GetID(), initCoins)
		So(err, ShouldBeNil)
		err = keeper.BurnCoins(ctx, bAcc.GetID(), initCoins)
		So(err, ShouldBeNil)

		bAccCoins = getCoinsByName(ctx, keeper, types.Burner, *app.AssetKeeper())
		require.Equal(t, chainType.Coins{}, bAccCoins)
		require.Equal(t, chainType.Coins{chainType.NewCoin(constants.DefaultBondDenom, sdk.NewInt(200000000))}, keeper.GetSupply(ctx).GetTotal())

		// test same functionality on module account with multiple permissions
		initialSupply = keeper.GetSupply(ctx)
		keeper.SetModuleAccount(ctx, multiPermAcc)

		_, err = app.AssetKeeper().IssueCoinPower(ctx, multiPermAcc.GetID(), initCoins)
		So(err, ShouldBeNil)

		err = keeper.BurnCoins(ctx, multiPermAcc.GetID(), initCoins)
		So(err, ShouldBeNil)

		multiPermAccCoins := getCoinsByName(ctx, keeper, multiPermAcc.GetName().String(), *app.AssetKeeper())
		require.Equal(t, chainType.Coins{}, multiPermAccCoins)
		require.Equal(t, chainType.Coins{chainType.NewCoin(constants.DefaultBondDenom, sdk.NewInt(300000000))}, keeper.GetSupply(ctx).GetTotal())
	})
}
