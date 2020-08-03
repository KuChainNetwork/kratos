package ante_test

import (
	"fmt"
	"math/rand"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/constants/keys"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
)

var (
	// values for test
	wallet      = simapp.NewWallet()
	addr1       = wallet.NewAccAddressByName(name1)
	addr2       = wallet.NewAccAddressByName(name2)
	addr3       = wallet.NewAccAddressByName(name3)
	addr4       = wallet.NewAccAddressByName(name4)
	addr5       = wallet.NewAccAddressByName(name5)
	name1       = types.MustName("test01@chain")
	name2       = types.MustName("aaaeeebbbccc")
	name3       = types.MustName("aaaeeebbbcc2")
	name4       = types.MustName("test")
	name5       = types.MustName("foo")
	account1    = types.NewAccountIDFromName(name1)
	account2    = types.NewAccountIDFromName(name2)
	account3    = types.NewAccountIDFromName(name3)
	account4    = types.NewAccountIDFromName(name4)
	account5    = types.NewAccountIDFromName(name5)
	addAccount1 = types.NewAccountIDFromAccAdd(addr1)
	addAccount2 = types.NewAccountIDFromAccAdd(addr2)
	addAccount3 = types.NewAccountIDFromAccAdd(addr3)
	addAccount4 = types.NewAccountIDFromAccAdd(addr4)
)

func createAppForTest() (*simapp.SimApp, sdk.Context) {
	asset1 := types.NewCoins(
		types.NewInt64Coin("foo/coin", 1000000000000000),
		types.NewInt64Coin(constants.DefaultBondDenom, 10000000000))
	asset2 := types.NewCoins(
		types.NewInt64Coin("foo/coin", 1000000000000000),
		types.NewInt64Coin(constants.DefaultBondDenom, 10000000000))
	asset3 := types.NewCoins(
		types.NewInt64Coin("foo/coin", 1000000000000000),
		types.NewInt64Coin(constants.DefaultBondDenom, 10000000000))

	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(),
		simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1),
		simapp.NewSimGenesisAccount(account2, addr2).WithAsset(asset2),
		simapp.NewSimGenesisAccount(addAccount2, addr2).WithAsset(types.NewInt64CoreCoins(10)),
		simapp.NewSimGenesisAccount(account3, addr3).WithAsset(asset3),
		simapp.NewSimGenesisAccount(account4, addr4).WithAsset(asset2),
		simapp.NewSimGenesisAccount(addAccount4, addr4).WithAsset(asset2),
		simapp.NewSimGenesisAccount(account5, addr5).WithAsset(asset2),
	)
	app := simapp.SetupWithGenesisAccounts(genAccs)

	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
	return app, ctxCheck
}

func testStdTx(app *simapp.SimApp, ids ...types.AccountID) types.StdTx {
	msgs := make([]sdk.Msg, 0, len(ids))
	privs := make([]crypto.PrivKey, 0, len(ids))

	for _, id := range ids {
		var (
			from = id
			addr = wallet.GetAuth(id)
			to   = account1
			amt  = types.NewCoins(types.NewCoin("foo/coin", types.NewInt(10)))
		)

		// a createMsg
		msg := assetTypes.NewMsgTransfer(addr, from, to, amt)
		msgs = append(msgs, &msg)

		privs = append(privs, wallet.PrivKey(addr))
	}

	return simapp.NewTxForTest(ids[0], msgs, privs...).GetTx(app)
}

func generatePubKeysAndSignatures(n int, msg []byte, keyTypeed25519 bool) (pubkeys []crypto.PubKey, signatures [][]byte) {
	pubkeys = make([]crypto.PubKey, n)
	signatures = make([][]byte, n)
	for i := 0; i < n; i++ {
		var privkey crypto.PrivKey
		if rand.Int63()%2 == 0 {
			privkey = ed25519.GenPrivKey()
		} else {
			privkey = secp256k1.GenPrivKey()
		}
		pubkeys[i] = privkey.PubKey()
		signatures[i], _ = privkey.Sign(msg)
	}
	return
}

func expectedGasCostByKeys(pubkeys []crypto.PubKey) uint64 {
	cost := uint64(0)
	for _, pubkey := range pubkeys {
		pubkeyType := strings.ToLower(fmt.Sprintf("%T", pubkey))
		switch {
		case strings.Contains(pubkeyType, "ed25519"):
			cost += keys.DefaultSigVerifyCostED25519
		case strings.Contains(pubkeyType, "secp256k1"):
			cost += keys.DefaultSigVerifyCostSecp256k1
		default:
			panic("unexpected key type")
		}
	}
	return cost
}
