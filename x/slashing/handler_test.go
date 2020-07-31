package slashing_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChainNetwork/kuchain/chain/config"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	slashingTypes "github.com/KuChainNetwork/kuchain/x/slashing/types"
	stakingTypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/tendermint/tendermint/crypto"
	"time"
	"encoding/hex"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func newTestApp(wallet *simapp.Wallet) (addAlice, addJack, addValidator sdk.AccAddress, accAlice, accJack, accValidator types.AccountID, app *simapp.SimApp) {
	addAlice = wallet.NewAccAddress()
	addJack = wallet.NewAccAddress()
	addValidator = wallet.NewAccAddress()

	accAlice = types.MustAccountID("alice@ok")
	accJack = types.MustAccountID("jack@ok")
	accValidator = types.MustAccountID("validator@ok")

	resInt, succ := sdk.NewIntFromString("100000000000000000000000")
	if !succ {
		resInt = sdk.NewInt(10000000000000000)
	}
	initAsset := types.NewCoin(constants.DefaultBondDenom, resInt)
	asset1 := types.Coins{
		types.NewInt64Coin("foo/coin", 67),
		initAsset}

	asset2 := types.Coins{
		types.NewInt64Coin("foo/coin", 67),
		types.NewInt64Coin(constants.DefaultBondDenom, 10000000)}

	genAlice := simapp.NewSimGenesisAccount(accAlice, addAlice).WithAsset(asset1)
	genJack := simapp.NewSimGenesisAccount(accJack, addJack).WithAsset(asset1)
	genValidator := simapp.NewSimGenesisAccount(accValidator, addValidator).WithAsset(asset2)

	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAlice, genJack, genValidator)
	app = simapp.SetupWithGenesisAccounts(genAccs)

	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	accountAlice := app.AccountKeeper().GetAccount(ctxCheck, accAlice)
	accountJack := app.AccountKeeper().GetAccount(ctxCheck, accJack)
	accountValidator := app.AccountKeeper().GetAccount(ctxCheck, accValidator)

	So(accountAlice, ShouldNotBeNil)
	So(genAlice.GetID().Eq(accountAlice.GetID()), ShouldBeTrue)
	So(genAlice.GetAuth().Equals(accountAlice.GetAuth()), ShouldBeTrue)

	So(accountJack, ShouldNotBeNil)
	So(genJack.GetID().Eq(accountJack.GetID()), ShouldBeTrue)
	So(genJack.GetAuth().Equals(accountJack.GetAuth()), ShouldBeTrue)

	So(accountValidator, ShouldNotBeNil)
	So(genValidator.GetID().Eq(accountValidator.GetID()), ShouldBeTrue)
	So(genValidator.GetAuth().Equals(accountValidator.GetAuth()), ShouldBeTrue)

	return addAlice, addJack, addValidator, accAlice, accJack, accValidator, app
}

func createValidator(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice types.AccountID, rate sdk.Dec, pk crypto.PubKey, passed bool) error {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addAlice)
	So(err, ShouldBeNil)

	ctxCheck.Logger().Info("auth nums", "seq", origAuthSeq, "num", origAuthNum)

	description := stakingTypes.NewDescription("moniker", "identity", "website", "securityContact", "details")

	msg := stakingTypes.NewKuMsgCreateValidator(addAlice, accAlice, pk, description, rate, accAlice)

	fee := types.Coins{types.NewInt64Coin(constants.DefaultBondDenom, 1000000)}
	header := abci.Header{Height: app.LastBlockHeight() + 1}

	_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
		header, accAlice, fee,
		[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		passed, passed, wallet.PrivKey(addAlice))
	ctxCheck.Logger().Info("createValidator error log", "err", err)
	return err
}

func delegationValidator(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice, accValidator types.AccountID, amount types.Coin, passed bool) error {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addAlice)
	So(err, ShouldBeNil)
	msg := stakingTypes.NewKuMsgDelegate(addAlice, accAlice, accValidator, amount)
	fee := types.Coins{types.NewInt64Coin(constants.DefaultBondDenom, 1000000)}
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
		header, accAlice, fee,
		[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		passed, passed, wallet.PrivKey(addAlice))
	ctxCheck.Logger().Info("delegationValidator error log", "err", err)

	return err
}

func unjail(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice types.AccountID, passed bool) error {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addAlice)
	So(err, ShouldBeNil)
	msg := slashingTypes.NewKuMsgUnjail(addAlice, accAlice)
	fee := types.Coins{types.NewInt64Coin(constants.DefaultBondDenom, 1000000)}
	blockTime := time.Now()
	if passed {
		blockTime = time.Now().Add(app.SlashKeeper().DowntimeJailDuration(ctxCheck) * 2)
	}

	header := abci.Header{Height: app.LastBlockHeight() + 1, Time: blockTime}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()
	header = abci.Header{Height: app.LastBlockHeight() + 1, Time: blockTime}
	_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
		header, accAlice, fee,
		[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		passed, passed, wallet.PrivKey(addAlice))
	ctxCheck.Logger().Info("unjail error log", "err", err)

	return err
}

func signBlock(app *simapp.SimApp, pk crypto.PubKey, missed bool, accAlice types.AccountID) {
	consAddr := pk.Address()
	validator := abci.Validator{Address: consAddr, Power: 1}
	voteInfo := abci.VoteInfo{Validator: validator, SignedLastBlock: missed}
	var votes []abci.VoteInfo
	votes = append(votes, voteInfo)
	lastcommitinfo := abci.LastCommitInfo{Round: 1, Votes: votes}
	header := abci.Header{Height: app.LastBlockHeight() + 1, Time: time.Now()}
	app.BeginBlock(abci.RequestBeginBlock{Header: header, LastCommitInfo: lastcommitinfo})
	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

}

func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	//res, err = crypto.PubKeyFromBytes(pkBytes)
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes)
	return pkEd
}

func TestSlashHandler(t *testing.T) {
	Convey("TestValidatorHandler", t, func() {
		config.SealChainConfig()
		wallet := simapp.NewWallet()

		addAlice, _, _, accAlice, _, _, app := newTestApp(wallet)
		rate, _ := sdk.NewDecFromStr("0.6")
		pk := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF100")
		err := createValidator(t, wallet, app, addAlice, accAlice, rate, pk, true)
		So(err, ShouldBeNil)
		bigAmount := types.NewInt64Coin(constants.DefaultBondDenom, 2100000000000000000)
		err = delegationValidator(t, wallet, app, addAlice, accAlice, accAlice, bigAmount, true)
		So(err, ShouldBeNil)
		err = unjail(t, wallet, app, addAlice, accAlice, false)
		So(err, ShouldNotBeNil)
		for i := 0; i < 30; i++ {
			signBlock(app, pk, true, accAlice)
		}
		for i := 0; i < 53; i++ {
			signBlock(app, pk, false, accAlice)
		}
		for i := 0; i < 18; i++ {
			signBlock(app, pk, true, accAlice)
		}
		err = unjail(t, wallet, app, addAlice, accAlice, false)
		So(err, ShouldNotBeNil)
		err = unjail(t, wallet, app, addAlice, accAlice, true)
		So(err, ShouldBeNil)
		err = unjail(t, wallet, app, addAlice, accAlice, false)
		So(err, ShouldNotBeNil)
	})

}
