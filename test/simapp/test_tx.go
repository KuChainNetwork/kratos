package simapp

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tendermint/tendermint/crypto"
)

var (
	DefaultTestFee = types.NewInt64Coins(constants.DefaultBondDenom, 100000)
)

type TestTx struct {
	expSimPass bool
	expPass    bool
	payer      types.AccountID
	fee        types.Coins
	msgs       []sdk.Msg
	priv       []crypto.PrivKey
}

func NewTxForTest(payer types.AccountID, msgs []sdk.Msg, priv ...crypto.PrivKey) *TestTx {
	return &TestTx{
		expSimPass: true,
		expPass:    true,
		payer:      payer,
		fee:        DefaultTestFee,
		msgs:       msgs,
		priv:       priv,
	}
}

func (t *TestTx) WithExpSimPass(expSimPass bool) *TestTx {
	t.expSimPass = expSimPass
	return t
}

func (t *TestTx) WithExpPass(expPass bool) *TestTx {
	t.expPass = expPass
	return t
}

func (t *TestTx) WithCannotPass() *TestTx {
	t.expPass = false
	t.expSimPass = false
	return t
}

func (t *TestTx) WithPayer(payer types.AccountID) *TestTx {
	t.payer = payer
	return t
}

func (t *TestTx) WithFee(fee types.Coins) *TestTx {
	t.fee = fee
	return t
}

func (t *TestTx) GetTx(app *SimApp) types.StdTx {
	var (
		nums = make([]uint64, 0, len(t.priv))
		seqs = make([]uint64, 0, len(t.priv))
	)

	for _, p := range t.priv {
		auth := types.AccAddress(p.PubKey().Address())
		seq, num, err := app.AccountKeeper().GetAuthSequence(app.NewTestContext(), auth)
		So(err, ShouldBeNil)

		nums = append(nums, num)
		seqs = append(seqs, seq)
	}

	return helpers.GenTx(
		t.msgs,
		t.fee,
		helpers.DefaultGenTxGas,
		t.payer,
		"", // all use chain id to ""
		nums,
		seqs,
		t.priv...,
	)
}

func CheckTxs(t *testing.T, app *SimApp, ctx sdk.Context, txs ...*TestTx) error {
	for _, tx := range txs {
		var (
			nums = make([]uint64, 0, len(tx.priv))
			seqs = make([]uint64, 0, len(tx.priv))
		)

		for _, p := range tx.priv {
			auth := types.AccAddress(p.PubKey().Address())

			seq, num, err := app.AccountKeeper().GetAuthSequence(ctx, auth)
			So(err, ShouldBeNil)

			nums = append(nums, num)
			seqs = append(seqs, seq)
		}

		header := ctx.BlockHeader()
		_, _, err := SignCheckDeliver(t, app.Codec(), app.BaseApp,
			header, tx.payer, tx.fee,
			tx.msgs, nums, seqs,
			tx.expSimPass, tx.expPass, tx.priv...)
		if err != nil {
			return err
		}
	}

	return nil
}
