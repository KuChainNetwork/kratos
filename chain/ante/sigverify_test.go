package ante_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/KuChainNetwork/kuchain/chain/ante"
	"github.com/KuChainNetwork/kuchain/chain/constants/keys"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/account"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestSetPubKey(t *testing.T) {
	// setup
	app, ctx := createAppForTest()
	ak := app.AccountKeeper()
	accountQuerier := account.NewQuerier(*ak)
	path := []string{accountTypes.QueryAuthByAddress}

	Convey("test set pubkey", t, func() {
		accounts := []types.AccountID{
			account1, account2, account3,
		}
		addrs := []types.AccAddress{
			addr1, addr2, addr3,
		}

		tx := testStdTx(app, accounts...)

		antehandler := sdk.ChainAnteDecorators(ante.NewSetPubKeyDecorator(*ak))

		ctx, err := antehandler(ctx, tx, true)
		So(err, ShouldBeNil)

		// Require that all accounts have pubkey set after Decorator runs
		for _, auth := range addrs {
			pubKey := wallet.PrivKey(auth).PubKey()

			req := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", accountTypes.QuerierRoute, accountTypes.QueryAuthByAddress),
				Data: []byte{},
			}
			req.Data = app.Codec().MustMarshalJSON(accountTypes.NewQueryAddAuthParams(auth))

			res, err := accountQuerier(ctx, path, req)
			So(err, ShouldBeNil)

			var authData accountTypes.Auth
			err = app.Codec().UnmarshalJSON(res, &authData)
			So(err, ShouldBeNil)

			So(authData.GetAddress(), simapp.ShouldEq, auth)
			So(authData.GetPubKey(), simapp.ShouldEq, pubKey)

		}
	})
}

func TestConsumeSignatureVerificationGas(t *testing.T) {
	msg := []byte{1, 2, 3, 4}

	pkSet1, sigSet1 := generatePubKeysAndSignatures(5, msg, false)
	multisigKey1 := multisig.NewPubKeyMultisigThreshold(2, pkSet1)
	multisignature1 := multisig.NewMultisig(len(pkSet1))
	expectedCost1 := expectedGasCostByKeys(pkSet1)
	for i := 0; i < len(pkSet1); i++ {
		multisignature1.AddSignatureFromPubKey(sigSet1[i], pkSet1[i], pkSet1)
	}

	Convey("test consume sig gas verification", t, func() {
		type args struct {
			meter  sdk.GasMeter
			sig    []byte
			pubkey crypto.PubKey
		}
		tests := []struct {
			name        string
			args        args
			gasConsumed uint64
			shouldErr   bool
		}{
			{"PubKeyEd25519", args{sdk.NewInfiniteGasMeter(), nil, ed25519.GenPrivKey().PubKey()}, keys.DefaultSigVerifyCostED25519, true},
			{"PubKeySecp256k1", args{sdk.NewInfiniteGasMeter(), nil, secp256k1.GenPrivKey().PubKey()}, keys.DefaultSigVerifyCostSecp256k1, false},
			{"Multisig", args{sdk.NewInfiniteGasMeter(), multisignature1.Marshal(), multisigKey1}, expectedCost1, false},
			{"unknown key", args{sdk.NewInfiniteGasMeter(), nil, nil}, 0, true},
		}
		for _, tt := range tests {
			err := ante.DefaultSigVerificationGasConsumer(tt.args.meter, tt.args.sig, tt.args.pubkey)

			if tt.shouldErr {
				Convey(fmt.Sprintf("%s should err", tt.name), func() {
					So(err, ShouldNotBeNil)
				})
			} else {
				Convey(fmt.Sprintf("%s %d != %d", tt.name, tt.gasConsumed, tt.args.meter.GasConsumed()), func() {
					So(err, ShouldBeNil)
					So(tt.gasConsumed, ShouldEqual, tt.args.meter.GasConsumed())
				})
			}
		}
	})
}

func TestSigVerification(t *testing.T) {
	// setup
	app, ctx := createAppForTest()
	// make block height non-zero to ensure account numbers part of signBytes
	ctx = app.NewTestContext()

	Convey("test sig verification", t, func() {
		priv1, priv2, priv3 := wallet.PrivKey(addr1), wallet.PrivKey(addr2), wallet.PrivKey(addr3)
		ttx := testStdTx(app, account1, addAccount2, account3)
		msgs, fee := ttx.Msgs, ttx.Fee

		ak := app.AccountKeeper()
		antehandler := sdk.ChainAnteDecorators(
			ante.NewSetPubKeyDecorator(*ak),
			ante.NewSigVerificationDecorator(*ak))

		type testCase struct {
			name      string
			privs     []crypto.PrivKey
			accNums   []uint64
			seqs      []uint64
			recheck   bool
			shouldErr bool
		}

		testCases := []testCase{
			{"no signers", []crypto.PrivKey{}, []uint64{}, []uint64{}, false, true},
			{"not enough signers", []crypto.PrivKey{priv1, priv2}, []uint64{1, 2}, []uint64{1, 1}, false, true},
			{"wrong order signers", []crypto.PrivKey{priv3, priv2, priv1}, []uint64{3, 2, 1}, []uint64{1, 1, 1}, false, true},
			{"wrong accnums", []crypto.PrivKey{priv1, priv2, priv3}, []uint64{97, 98, 99}, []uint64{1, 1, 1}, false, true},
			{"wrong sequences", []crypto.PrivKey{priv1, priv2, priv3}, []uint64{0, 1, 2}, []uint64{83, 84, 85}, false, true},
			{"valid tx", []crypto.PrivKey{priv1, priv2, priv3}, []uint64{1, 2, 3}, []uint64{1, 1, 1}, false, false},
			{"no err on recheck", []crypto.PrivKey{}, []uint64{}, []uint64{}, true, false},
		}

		for i, tc := range testCases {
			func(ii int, ttc testCase) {
				Convey(fmt.Sprintf("testCase %d: %s", ii, ttc.name), func() {
					ctx = ctx.WithIsReCheckTx(ttc.recheck)
					tx := simapp.NewTestTx(ctx, msgs, ttc.privs, ttc.accNums, ttc.seqs, fee)

					_, err := antehandler(ctx, tx, false)
					if ttc.shouldErr {
						So(err, ShouldNotBeNil)
					} else {
						So(err, ShouldBeNil)
					}
				})
			}(i, tc)
		}
	})
}

func TestIncrementSequenceDecorator(t *testing.T) {
	app, ctx := createAppForTest()

	ak := app.AccountKeeper()
	isd := ante.NewIncrementSequenceDecorator(*ak)
	antehandler := sdk.ChainAnteDecorators(isd)

	type testCase struct {
		ctx         sdk.Context
		simulate    bool
		expectedSeq uint64
	}

	testCases := []testCase{
		{ctx.WithIsReCheckTx(true), false, 1},
		{ctx.WithIsCheckTx(true).WithIsReCheckTx(false), false, 2},
		{ctx.WithIsReCheckTx(true), false, 2},
		{ctx.WithIsReCheckTx(true), false, 2},
		{ctx.WithIsReCheckTx(true), true, 3},
	}

	Convey("test increment seq", t, func() {
		Convey("test account auth increment seq", func() {
			tx := testStdTx(app, account1)

			for ii, ttc := range testCases {
				func(i int, tc testCase) {
					Convey(fmt.Sprintf("testcase %d : %v - %d", i, tc.simulate, tc.expectedSeq), func() {
						_, err := antehandler(tc.ctx, tx, tc.simulate)
						So(err, ShouldBeNil)

						seq, _, err := app.AccountKeeper().GetAuthSequence(tc.ctx, addr1)
						So(err, ShouldBeNil)
						So(seq, ShouldEqual, tc.expectedSeq)
					})
				}(ii, ttc)
			}
		})
	})
}
