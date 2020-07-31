package mint_test

import (
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/mint"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, types.Header{})

	app.InitChain(
		types.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)

	acc := app.SupplyKeeper().GetModuleAccount(ctx, mint.ModuleName).GetID()
	require.NotNil(t, acc)
}
