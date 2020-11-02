package keeper_test

import (
	"fmt"
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset"
	assettypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	keep "github.com/KuChainNetwork/kuchain/x/supply/keeper"
	"github.com/KuChainNetwork/kuchain/x/supply/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func initCoin(t *testing.T, ctx sdk.Context, assetKeeper asset.Keeper, coin chainType.Coin, id chainType.AccountID) {
	intMaxNum, _ := sdk.NewIntFromString("100000000000000000000000")

	denom := coin.Denom
	creater, symbol, err := chainType.CoinAccountsFromDenom(denom)
	So(err, ShouldBeNil)

	err = assetKeeper.Create(ctx, creater, symbol, chainType.NewCoin(denom, intMaxNum),
		true, true, true, 0, chainType.NewInt64CoreCoin(0), []byte("create"))
	So(err, ShouldBeNil)

	err = assetKeeper.Issue(ctx, creater, symbol, assettypes.NewCoin(coin.Denom, coin.Amount))
	So(err, ShouldBeNil)

	Coins := chainType.NewCoins(coin)
	err = assetKeeper.Transfer(ctx, chainType.NewAccountIDFromName(creater), id, Coins)
	So(err, ShouldBeNil)
}

func fInitCoins(t *testing.T, ctx sdk.Context, ask asset.Keeper, coins chainType.Coins, id chainType.AccountID) {
	for _, c := range coins {
		initCoin(t, ctx, ask, c, id)
	}
}

func TestNewQuerier(t *testing.T) {
	app, ctx := createTestApp(false)
	keeper := app.SupplyKeeper()
	cdc := app.Codec()

	Convey("test new query", t, func() {

		supplyCoins := chainType.NewCoins(
			chainType.NewCoin(constants.DefaultBondDenom, sdk.NewInt(100)),
			chainType.NewCoin(constants.ChainMainNameStr+"/"+"photon", sdk.NewInt(50)),
			chainType.NewCoin(constants.ChainMainNameStr+"/"+"atom", sdk.NewInt(2000)),
			chainType.NewCoin(constants.ChainMainNameStr+"/"+"btc", sdk.NewInt(21000000)),
		)

		supplyAcc := keeper.GetModuleAccount(ctx, types.ModuleName).GetID()
		fInitCoins(t, ctx, *app.AssetKeeper(), supplyCoins, supplyAcc)

		query := abci.RequestQuery{
			Path: "",
			Data: []byte{},
		}
		//
		querier := keep.NewQuerier(*keeper)

		bz, err := querier(ctx, []string{"other"}, query)
		require.Error(t, err)
		require.Nil(t, bz)

		queryTotalSupplyParams := types.NewQueryTotalSupplyParams(1, 20)
		bz, errRes := cdc.MarshalJSON(queryTotalSupplyParams)
		require.Nil(t, errRes)

		query.Path = fmt.Sprintf("/custom/supply/%s", types.QueryTotalSupply)
		query.Data = bz

		_, err = querier(ctx, []string{types.QueryTotalSupply}, query)
		require.Nil(t, err)

		querySupplyParams := types.NewQuerySupplyOfParams(constants.DefaultBondDenom)
		bz, errRes = cdc.MarshalJSON(querySupplyParams)
		require.Nil(t, errRes)

		query.Path = fmt.Sprintf("/custom/supply/%s", types.QuerySupplyOf)
		query.Data = bz

		_, err = querier(ctx, []string{types.QuerySupplyOf}, query)
		require.Nil(t, err)
	})
}

func TestQuerySupply(t *testing.T) {
	app, ctx := createTestApp(false)
	keeper := *app.SupplyKeeper()
	cdc := app.Codec()

	Convey("test query supply", t, func() {

		supplyCoins := chainType.NewCoins(
			chainType.NewCoin(constants.DefaultBondDenom, sdk.NewInt(100)),
			chainType.NewCoin(constants.ChainMainNameStr+"/"+"photon", sdk.NewInt(50)),
			chainType.NewCoin(constants.ChainMainNameStr+"/"+"atom", sdk.NewInt(2000)),
			chainType.NewCoin(constants.ChainMainNameStr+"/"+"btc", sdk.NewInt(21000000)),
		)

		supplyAcc := keeper.GetModuleAccount(ctx, types.ModuleName).GetID()
		fInitCoins(t, ctx, *app.AssetKeeper(), supplyCoins, supplyAcc)

		query := abci.RequestQuery{
			Path: "",
			Data: []byte{},
		}

		querier := keep.NewQuerier(keeper)

		//keeper.SetSupply(ctx, types.NewSupply(supplyCoins))

		queryTotalSupplyParams := types.NewQueryTotalSupplyParams(1, 10)
		bz, errRes := cdc.MarshalJSON(queryTotalSupplyParams)
		require.Nil(t, errRes)

		query.Path = fmt.Sprintf("/custom/supply/%s", types.QueryTotalSupply)
		query.Data = bz

		res, err := querier(ctx, []string{types.QueryTotalSupply}, query)
		require.Nil(t, err)

		var totalCoins chainType.Coins
		errRes = cdc.UnmarshalJSON(res, &totalCoins)
		require.Nil(t, errRes)
		require.Equal(t, supplyCoins, totalCoins)

		querySupplyParams := types.NewQuerySupplyOfParams(constants.DefaultBondDenom)
		bz, errRes = cdc.MarshalJSON(querySupplyParams)
		require.Nil(t, errRes)

		query.Path = fmt.Sprintf("/custom/supply/%s", types.QuerySupplyOf)
		query.Data = bz

		res, err = querier(ctx, []string{types.QuerySupplyOf}, query)
		require.Nil(t, err)

		var supply sdk.Int
		errRes = supply.UnmarshalJSON(res)
		require.Nil(t, errRes)
		require.True(sdk.IntEq(t, sdk.NewInt(100), supply))

	})
}
