package simapp

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/types"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

func CommitTransferTx(t *testing.T, app *SimApp, wallet *Wallet, isSuccess bool,
	from, to types.AccountID, amt types.Coins, payer types.AccountID) error {
	ctx := app.NewTestContext()
	auth := app.AccountKeeper().GetAccount(ctx, from).GetAuth()

	msg := assetTypes.NewMsgTransfer(auth, from, to, amt)
	tx := NewTxForTest(
		from,
		[]sdk.Msg{
			&msg,
		}, wallet.PrivKey(auth))

	if !isSuccess {
		tx = tx.WithCannotPass()
	}

	return CheckTxs(t, app, ctx, tx)
}

func NewTestTx(ctx sdk.Context, msgs []sdk.Msg, privs []crypto.PrivKey, accNums []uint64, seqs []uint64, fee types.StdFee) sdk.Tx {
	sigs := make([]types.StdSignature, len(privs))
	for i, priv := range privs {
		signBytes := types.StdSignBytes(ctx.ChainID(), accNums[i], seqs[i], fee, msgs, "")

		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = types.StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := types.NewStdTx(msgs, fee, sigs, "")
	return tx
}

func NewTestTxWithMemo(ctx sdk.Context, msgs []sdk.Msg, privs []crypto.PrivKey, accNums []uint64, seqs []uint64, fee types.StdFee, memo string) sdk.Tx {
	sigs := make([]types.StdSignature, len(privs))
	for i, priv := range privs {
		signBytes := types.StdSignBytes(ctx.ChainID(), accNums[i], seqs[i], fee, msgs, memo)

		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = types.StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := types.NewStdTx(msgs, fee, sigs, memo)
	return tx
}
