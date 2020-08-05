package staking_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"

	abci "github.com/tendermint/tendermint/abci/types"

	"encoding/hex"

	"github.com/KuChainNetwork/kuchain/chain/config"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/msg"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	stakingTypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/tendermint/tendermint/crypto"
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
	otherCoinDenom := types.CoinDenom(types.MustName("foo"), types.MustName("coin"))
	initAsset := types.NewCoin(constants.DefaultBondDenom, resInt)

	asset1 := types.NewCoins(
		types.NewInt64Coin(otherCoinDenom, 67),
		initAsset)

	asset2 := types.NewCoins(
		types.NewInt64Coin(otherCoinDenom, 67),
		types.NewInt64Coin(constants.DefaultBondDenom, 10000000))

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

func exitValidator(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice types.AccountID, rate sdk.Dec, passed bool) error {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addAlice)
	So(err, ShouldBeNil)
	description := stakingTypes.NewDescription("Newmoniker", "Newidentity", "Newwebsite", "NewsecurityContact", "Newdetails")
	msg := stakingTypes.NewKuMsgEditValidator(addAlice, accAlice, description, &rate)
	fee := types.NewInt64Coins(constants.DefaultBondDenom, 1000000)
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
		header, accAlice, fee,
		[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		passed, passed, wallet.PrivKey(addAlice))
	ctxCheck.Logger().Info("exitValidator error log", "err", err)
	return err
}

// NewKuMsgDelegate create kuMsgDelegate
func customizeKuMsgDelegate(auth sdk.AccAddress, delAddr types.AccountID, valAddr types.AccountID, amount, transferAmount types.Coin) stakingTypes.KuMsgDelegate {
	return stakingTypes.KuMsgDelegate{
		*msg.MustNewKuMsg(
			stakingTypes.RouterKeyName,
			msg.WithAuth(auth),
			msg.WithTransfer(delAddr, stakingTypes.ModuleAccountID, types.Coins{transferAmount}),
			msg.WithData(stakingTypes.Cdc(), &stakingTypes.MsgDelegate{
				DelegatorAccount: delAddr,
				ValidatorAccount: valAddr,
				Amount:           amount,
			}),
		),
	}
}

func delegationValidator(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice, accValidator types.AccountID, amount, transferAmount types.Coin, passed bool) error {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addAlice)
	So(err, ShouldBeNil)
	msg := customizeKuMsgDelegate(addAlice, accAlice, accValidator, amount, transferAmount)
	fee := types.NewInt64Coins(constants.DefaultBondDenom, 1000000)
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
		header, accAlice, fee,
		[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		passed, passed, wallet.PrivKey(addAlice))
	ctxCheck.Logger().Info("exitValidator error log", "err", err)

	return err
}

func redelegateValidator(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice, accJack, accValidator types.AccountID, amount types.Coin, passed bool) error {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addAlice)
	So(err, ShouldBeNil)
	msg := stakingTypes.NewKuMsgRedelegate(addAlice, accAlice, accJack, accValidator, amount)
	fee := types.Coins{types.NewInt64Coin(constants.DefaultBondDenom, 1000000)}
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
		header, accAlice, fee,
		[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		passed, passed, wallet.PrivKey(addAlice))
	ctxCheck.Logger().Info("exitValidator error log", "err", err)

	return err
}

func unbondValidator(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice, accJack types.AccountID, amount types.Coin, passed bool) error {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addAlice)
	So(err, ShouldBeNil)
	msg := stakingTypes.NewKuMsgUnbond(addAlice, accAlice, accJack, amount)
	fee := types.Coins{types.NewInt64Coin(constants.DefaultBondDenom, 1000000)}
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
		header, accAlice, fee,
		[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		passed, passed, wallet.PrivKey(addAlice))
	ctxCheck.Logger().Info("unbondValidator error log", "err", err)
	return err
}

func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes)
	return pkEd
}

func TestStakingHandler(t *testing.T) {
	Convey("TestValidatorHandler", t, func() {
		config.SealChainConfig()
		wallet := simapp.NewWallet()

		addAlice, addJack, _, accAlice, accJack, _, app := newTestApp(wallet)
		wrongrate, _ := sdk.NewDecFromStr("1.1")
		pk := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF100")
		Newpk := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF200")
		//wrong rate
		err := createValidator(t, wallet, app, addAlice, accAlice, wrongrate, pk, false)
		So(err, ShouldNotBeNil)
		//right
		rightRate, _ := sdk.NewDecFromStr("0.65")
		err = createValidator(t, wallet, app, addAlice, accAlice, rightRate, pk, true)
		So(err, ShouldBeNil)
		//pubkey and account already exist
		err = createValidator(t, wallet, app, addAlice, accAlice, rightRate, pk, false)
		So(err, ShouldNotBeNil)
		//account  already exist
		err = createValidator(t, wallet, app, addAlice, accAlice, rightRate, Newpk, false)
		So(err, ShouldNotBeNil)
		//pubkey already exist
		err = createValidator(t, wallet, app, addJack, accJack, rightRate, pk, false)
		So(err, ShouldNotBeNil)
		//right
		err = createValidator(t, wallet, app, addJack, accJack, rightRate, Newpk, true)
		So(err, ShouldBeNil)
		//wrong not 24 hour
		err = exitValidator(t, wallet, app, addJack, accJack, rightRate, false)
		So(err, ShouldNotBeNil)
		//wrong not 24 hour
		err = exitValidator(t, wallet, app, addJack, accJack, wrongrate, false)
		So(err, ShouldNotBeNil)
	})
	Convey("TestDelegateHandler", t, func() {
		//config.SealChainConfig() 	//just once
		wallet := simapp.NewWallet()
		addAlice, addJack, addValidator, accAlice, accJack, accValidator, app := newTestApp(wallet)
		rightRate, _ := sdk.NewDecFromStr("0.65")
		pk := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF100")
		Newpk := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF200")
		err := createValidator(t, wallet, app, addJack, accJack, rightRate, Newpk, true)
		So(err, ShouldBeNil)
		err = createValidator(t, wallet, app, addAlice, accAlice, rightRate, pk, true)
		So(err, ShouldBeNil)
		//right delegate
		delegateAmount := types.NewInt64Coin(constants.DefaultBondDenom, 50000000)
		smallAmount := types.NewInt64Coin(constants.DefaultBondDenom, 50000)
		bigAmount := types.NewInt64Coin(constants.DefaultBondDenom, 2100000000000000000)
		//transferAmount not equal  delegateAmount
		err = delegationValidator(t, wallet, app, addAlice, accAlice, accJack, delegateAmount, smallAmount, false)
		So(err, ShouldNotBeNil)
		//alice D jack 50000000
		err = delegationValidator(t, wallet, app, addAlice, accAlice, accJack, delegateAmount, delegateAmount, true)
		So(err, ShouldBeNil)
		//wrong coin delegate
		wrongAmount := types.Coin{Denom: constants.DefaultBondDenom, Amount: sdk.NewInt(-5000000)}
		err = delegationValidator(t, wallet, app, addValidator, accValidator, accJack, wrongAmount, wrongAmount, false)
		So(err, ShouldNotBeNil)
		// coin not enought
		err = delegationValidator(t, wallet, app, addValidator, accValidator, accJack, delegateAmount, delegateAmount, false)
		So(err, ShouldNotBeNil)
		// not exist validator
		err = delegationValidator(t, wallet, app, addAlice, accAlice, accValidator, delegateAmount, delegateAmount, false)
		So(err, ShouldNotBeNil)
		//wrong coin redelegate
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, wrongAmount, false)
		So(err, ShouldNotBeNil)
		//right   alice R jack T alice 50000
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, smallAmount, true)
		So(err, ShouldBeNil)
		//invalid shares amount
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, bigAmount, false)
		So(err, ShouldNotBeNil)
		//redelegation destination validator not found
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accValidator, smallAmount, false)
		So(err, ShouldNotBeNil)
		//src validator does not exist
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accValidator, accJack, smallAmount, false)
		So(err, ShouldNotBeNil)
		// //because not bonded   alice R alice jack 50000
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		//alice U jack 50000
		err = unbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		//wrong coin unbond
		err = unbondValidator(t, wallet, app, addAlice, accAlice, accJack, wrongAmount, false)
		So(err, ShouldNotBeNil)
		//alice U jack 50000000
		err = unbondValidator(t, wallet, app, addAlice, accAlice, accJack, bigAmount, false)
		So(err, ShouldNotBeNil)
		// unbond 7 times
		err = unbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = unbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = unbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = unbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = unbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = unbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = unbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, false)
		So(err, ShouldNotBeNil)
		// jack U jack
		err = unbondValidator(t, wallet, app, addJack, accJack, accJack, smallAmount, false)
		So(err, ShouldNotBeNil)
		// jack U Validator
		err = unbondValidator(t, wallet, app, addJack, accJack, accValidator, smallAmount, false)
		So(err, ShouldNotBeNil)
	})
	Convey("TestReDelegateHandler", t, func() {
		wallet := simapp.NewWallet()
		addAlice, addJack, _, accAlice, accJack, _, app := newTestApp(wallet)
		rightRate, _ := sdk.NewDecFromStr("0.65")
		pk := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF100")
		Newpk := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF200")
		err := createValidator(t, wallet, app, addJack, accJack, rightRate, Newpk, true)
		So(err, ShouldBeNil)
		err = createValidator(t, wallet, app, addAlice, accAlice, rightRate, pk, true)
		So(err, ShouldBeNil)
		bigAmount := types.NewInt64Coin(constants.DefaultBondDenom, 2100000000000000000)
		//alice D jack
		err = delegationValidator(t, wallet, app, addAlice, accAlice, accJack, bigAmount, bigAmount, true)
		So(err, ShouldBeNil)
		//jack D alice
		err = delegationValidator(t, wallet, app, addJack, accJack, accAlice, bigAmount, bigAmount, true)
		So(err, ShouldBeNil)
		delegateAmount := types.NewInt64Coin(constants.DefaultBondDenom, 50000000)
		// alice R jack ->alice
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		// alice R alice ->jack
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accAlice, accJack, delegateAmount, false)
		So(err, ShouldNotBeNil)
		//  alice R jack ->alice 7 times
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = redelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, false)
		So(err, ShouldNotBeNil)
		//
	})
}
