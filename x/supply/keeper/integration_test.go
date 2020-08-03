package keeper_test

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChainNetwork/kuchain/test/simapp"
	keep "github.com/KuChainNetwork/kuchain/x/supply/keeper"
	"github.com/KuChainNetwork/kuchain/x/supply/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	multiPerm  = "multipleaccount"
	randomPerm = "randompermission"
	holder     = "holder"
	test       = "test"
)

// nolint:deadcode,unused
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)

	// add module accounts to supply keeper
	maccPerms := simapp.GetMaccPerms()
	maccPerms[holder] = nil
	maccPerms[types.Burner] = []string{types.Burner}
	maccPerms[types.Minter] = []string{types.Minter}
	maccPerms[multiPerm] = []string{types.Burner, types.Minter, types.Staking}
	maccPerms[randomPerm] = []string{"random"}
	maccPerms[types.ModuleName] = nil
	maccPerms[test] = []string{types.Minter, types.Staking}

	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	SupplyKeeper := keep.NewKeeper(app.Codec(), app.GetKey(types.StoreKey), app.AccountKeeper(), app.AssetKeeper(), maccPerms)

	app.SetSupplyKeeper(SupplyKeeper)

	return app, ctx
}
