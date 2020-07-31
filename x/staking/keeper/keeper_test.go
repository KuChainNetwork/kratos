package keeper_test

import (
	"testing"
	"github.com/stretchr/testify/require"

	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestKeeper(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("TestParams", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.StakeKeeper()
		keeper = keeper.EmptyHooks()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		expParams := types.DefaultParams()
		resParams := keeper.GetParams(ctx)
		require.Equal(t, expParams, resParams)
		expParams.MaxValidators = 777
		keeper.SetParams(ctx, expParams)
		resParams = keeper.GetParams(ctx)
		require.Equal(t, expParams, resParams)

	})
}
