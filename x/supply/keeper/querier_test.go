package keeper_test

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset"
	assettypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	"github.com/KuChainNetwork/kuchain/x/supply/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"strings"
	"testing"

	//"github.com/stretchr/testify/require"
	//abci "github.com/tendermint/tendermint/abci/types"

	keep "github.com/KuChainNetwork/kuchain/x/supply/keeper"
	//"github.com/KuChainNetwork/kuchain/x/supply/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func initCoin(t *testing.T, ctx sdk.Context, assetKeeper asset.Keeper, coin chainType.Coin, id chainType.AccountID) {

	intNum, _ := sdk.NewIntFromString("80000000000000000000000")
	intMaxNum, _ := sdk.NewIntFromString("100000000000000000000000")

	TestMaster := constants.ChainMainNameStr
	MasterName, _ := chainType.NewName(TestMaster)
	Master := chainType.NewAccountIDFromName(MasterName)

	Symbol := strings.Split(coin.Denom, "/")
	SymbolName, _ := chainType.NewName(Symbol[1])

	assetKeeper.Create(ctx, MasterName, SymbolName, assettypes.NewCoin(coin.Denom, intNum),
		true, true, 0, assettypes.NewCoin(coin.Denom, intMaxNum), []byte("create"))

	assetKeeper.Issue(ctx, MasterName, SymbolName, assettypes.NewCoin(coin.Denom, coin.Amount))

	Coins := chainType.NewCoins(chainType.NewCoin(coin.Denom, coin.Amount))
	err := assetKeeper.Transfer(ctx, Master, id, Coins)
	require.Nil(t, err)

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
}

func TestQuerySupply(t *testing.T) {
	app, ctx := createTestApp(false)
	keeper := *app.SupplyKeeper()
	cdc := app.Codec()


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
}
