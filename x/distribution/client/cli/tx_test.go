package cli

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// decCoins := chainTypes.NewDecCoins(chainTypes.NewDecCoin(constants.DefaultBondDenom, sdk.NewInt(sdk.OneDec().Int64())))

func createFakeTxBuilder() txutil.TxBuilder {
	cdc := codec.New()
	return txutil.NewTxBuilder(
		txutil.GetTxEncoder(cdc),
		123,
		9876,
		0,
		1.2,
		false,
		"test_chain",
		"hello",
		chainTypes.NewCoins(chainTypes.NewCoin(constants.DefaultBondDenom, sdk.NewInt(sdk.OneDec().Int64()))),
		chainTypes.NewDecCoins(chainTypes.DecCoin{Denom: constants.DefaultBondDenom, Amount: sdk.NewDecWithPrec(10000, sdk.Precision)}),
	)
}

func Test_splitAndCall_NoMessages(t *testing.T) {
	ctx := context.CLIContext{}
	txBldr := createFakeTxBuilder()
	ctxl := txutil.NewKuCLICtx(ctx)

	err := splitAndApply(nil, ctxl, txBldr, nil, 10)
	assert.NoError(t, err, "")
}

func Test_splitAndCall_Splitting(t *testing.T) {
	ctx := context.CLIContext{}
	txBldr := createFakeTxBuilder()
	ctxl := txutil.NewKuCLICtx(ctx)

	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// Add five messages
	msgs := []sdk.Msg{
		sdk.NewTestMsg(addr),
		sdk.NewTestMsg(addr),
		sdk.NewTestMsg(addr),
		sdk.NewTestMsg(addr),
		sdk.NewTestMsg(addr),
	}

	// Keep track of number of calls
	const chunkSize = 2

	callCount := 0
	err := splitAndApply(
		func(ctx txutil.KuCLIContext, txBldr txutil.TxBuilder, msgs []sdk.Msg) error {
			callCount++

			assert.NotNil(t, ctx)
			assert.NotNil(t, txBldr)
			assert.NotNil(t, msgs)

			if callCount < 3 {
				assert.Equal(t, len(msgs), 2)
			} else {
				assert.Equal(t, len(msgs), 1)
			}

			return nil
		},
		ctxl, txBldr, msgs, chunkSize)

	assert.NoError(t, err, "")
	assert.Equal(t, 3, callCount)
}
