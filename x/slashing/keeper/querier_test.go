package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	slashKeeper "github.com/KuChainNetwork/kuchain/x/slashing/keeper"

	"github.com/KuChainNetwork/kuchain/test/simapp"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewQuerier(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("queryparameters", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.SlashKeeper()
		querier := slashKeeper.NewQuerier(*keeper)
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

		query := abci.RequestQuery{
			Path: "",
			Data: []byte{},
		}

		_, err := querier(ctx, []string{"parameters"}, query)
		require.NoError(t, err)
	})
}
