package staking_test

import (
	//"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChainNetwork/kuchain/chain/config"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	stakingTypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/tendermint/tendermint/crypto"
)

func NewTestApp(wallet *simapp.Wallet) (addAlice, addJack, addValidator sdk.AccAddress, accAlice, accJack, accValidator types.AccountID, app *simapp.SimApp) {
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

func CreateValidator(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice types.AccountID, rate sdk.Dec, pk crypto.PubKey, passed bool) error {
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
	ctxCheck.Logger().Info("CreateValidator error log", "err", err)
	return err
}

func ExitValidator(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice types.AccountID, rate sdk.Dec, passed bool) error {
	//auth sdk.AccAddress, valAddr chainTypes.AccountID, description Description, newRate *sdk.Dec
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addAlice)
	So(err, ShouldBeNil)
	description := stakingTypes.NewDescription("Newmoniker", "Newidentity", "Newwebsite", "NewsecurityContact", "Newdetails")
	msg := stakingTypes.NewKuMsgEditValidator(addAlice, accAlice, description, &rate)
	fee := types.Coins{types.NewInt64Coin(constants.DefaultBondDenom, 1000000)}
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
		header, accAlice, fee,
		[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		passed, passed, wallet.PrivKey(addAlice))
	ctxCheck.Logger().Info("ExitValidator error log", "err", err)
	return err
}

func DelegationValidator(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice, accValidator types.AccountID, amount types.Coin, passed bool) error {
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
	ctxCheck.Logger().Info("ExitValidator error log", "err", err)

	return err
}

func RedelegateValidator(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice, accJack, accValidator types.AccountID, amount types.Coin, passed bool) error {
	//NewKuMsgRedelegate(auth sdk.AccAddress, delAddr chainTypes.AccountID, valSrcAddr, valDstAddr chainTypes.AccountID, amount chainTypes.Coin)
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
	ctxCheck.Logger().Info("ExitValidator error log", "err", err)

	return err
}

func UnbondValidator(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice, accJack types.AccountID, amount types.Coin, passed bool) error {
	//NewKuMsgUnbond(auth sdk.AccAddress, delAddr chainTypes.AccountID, valAddr chainTypes.AccountID, amount chainTypes.Coin)
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
	ctxCheck.Logger().Info("UnbondValidator error log", "err", err)
	return err
}

func TestStakingHandler(t *testing.T) {
	Convey("TestValidatorHandler", t, func() {
		config.SealChainConfig()
		wallet := simapp.NewWallet()

		addAlice, addJack, _, accAlice, accJack, _, app := NewTestApp(wallet)
		wrongrate, _ := sdk.NewDecFromStr("1.1")
		pk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepqn4usdx22zdntysj7n795xj77wrc62sytheeevr7zlna4yhwppdrs8mpds3")
		//wrong rate
		err := CreateValidator(t, wallet, app, addAlice, accAlice, wrongrate, pk, false)
		So(err, ShouldNotBeNil)
		//right
		rightRate, _ := sdk.NewDecFromStr("0.65")
		err = CreateValidator(t, wallet, app, addAlice, accAlice, rightRate, pk, true)
		So(err, ShouldBeNil)
		//pubkey and account already exist
		err = CreateValidator(t, wallet, app, addAlice, accAlice, rightRate, pk, false)
		So(err, ShouldNotBeNil)
		//account  already exist
		Newpk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepq0cm4j2wtny3x435zuc53zffk9fndj7f37xjkxv4lqdx4w4z3mayqmf9aef")
		err = CreateValidator(t, wallet, app, addAlice, accAlice, rightRate, Newpk, false)
		So(err, ShouldNotBeNil)
		//pubkey already exist
		err = CreateValidator(t, wallet, app, addJack, accJack, rightRate, pk, false)
		So(err, ShouldNotBeNil)
		//right
		err = CreateValidator(t, wallet, app, addJack, accJack, rightRate, Newpk, true)
		So(err, ShouldBeNil)
		//wrong not 24 hour
		err = ExitValidator(t, wallet, app, addJack, accJack, rightRate, false)
		So(err, ShouldNotBeNil)
		//wrong not 24 hour
		err = ExitValidator(t, wallet, app, addJack, accJack, wrongrate, false)
		So(err, ShouldNotBeNil)
	})
	Convey("TestDelegateHandler", t, func() {
		//config.SealChainConfig() 	//just once
		wallet := simapp.NewWallet()
		addAlice, addJack, addValidator, accAlice, accJack, accValidator, app := NewTestApp(wallet)
		rightRate, _ := sdk.NewDecFromStr("0.65")
		pk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepqn4usdx22zdntysj7n795xj77wrc62sytheeevr7zlna4yhwppdrs8mpds3")
		Newpk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepq0cm4j2wtny3x435zuc53zffk9fndj7f37xjkxv4lqdx4w4z3mayqmf9aef")
		err := CreateValidator(t, wallet, app, addJack, accJack, rightRate, Newpk, true)
		So(err, ShouldBeNil)
		err = CreateValidator(t, wallet, app, addAlice, accAlice, rightRate, pk, true)
		So(err, ShouldBeNil)
		//right delegate
		delegateAmount := types.NewInt64Coin(constants.DefaultBondDenom, 50000000)
		smallAmount := types.NewInt64Coin(constants.DefaultBondDenom, 50000)
		bigAmount := types.NewInt64Coin(constants.DefaultBondDenom, 2100000000000000000)
		//alice D jack 50000000
		err = DelegationValidator(t, wallet, app, addAlice, accAlice, accJack, delegateAmount, true)
		So(err, ShouldBeNil)

		// coin not enought
		err = DelegationValidator(t, wallet, app, addValidator, accValidator, accJack, delegateAmount, false)
		So(err, ShouldNotBeNil)
		// not exist validator
		err = DelegationValidator(t, wallet, app, addAlice, accAlice, accValidator, delegateAmount, false)
		So(err, ShouldNotBeNil)
		//right   alice R jack T alice 50000
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, smallAmount, true)
		So(err, ShouldBeNil)
		//invalid shares amount
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, bigAmount, false)
		So(err, ShouldNotBeNil)
		//redelegation destination validator not found
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accValidator, smallAmount, false)
		So(err, ShouldNotBeNil)
		//src validator does not exist
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accValidator, accJack, smallAmount, false)
		So(err, ShouldNotBeNil)
		// //because not bonded   alice R alice jack 50000
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		//alice U jack 50000
		err = UnbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		//alice U jack 50000000
		err = UnbondValidator(t, wallet, app, addAlice, accAlice, accJack, bigAmount, false)
		So(err, ShouldNotBeNil)
		// unbond 7 times
		err = UnbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = UnbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = UnbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = UnbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = UnbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = UnbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, true)
		So(err, ShouldBeNil)
		err = UnbondValidator(t, wallet, app, addAlice, accAlice, accJack, smallAmount, false)
		So(err, ShouldNotBeNil)
		// jack U jack
		err = UnbondValidator(t, wallet, app, addJack, accJack, accJack, smallAmount, false)
		So(err, ShouldNotBeNil)
		// jack U Validator
		err = UnbondValidator(t, wallet, app, addJack, accJack, accValidator, smallAmount, false)
		So(err, ShouldNotBeNil)
	})
	Convey("TestReDelegateHandler", t, func() {
		wallet := simapp.NewWallet()
		addAlice, addJack, _, accAlice, accJack, _, app := NewTestApp(wallet)
		rightRate, _ := sdk.NewDecFromStr("0.65")
		pk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepqn4usdx22zdntysj7n795xj77wrc62sytheeevr7zlna4yhwppdrs8mpds3")
		Newpk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepq0cm4j2wtny3x435zuc53zffk9fndj7f37xjkxv4lqdx4w4z3mayqmf9aef")
		err := CreateValidator(t, wallet, app, addJack, accJack, rightRate, Newpk, true)
		So(err, ShouldBeNil)
		err = CreateValidator(t, wallet, app, addAlice, accAlice, rightRate, pk, true)
		So(err, ShouldBeNil)
		bigAmount := types.NewInt64Coin(constants.DefaultBondDenom, 2100000000000000000)
		//alice D jack
		err = DelegationValidator(t, wallet, app, addAlice, accAlice, accJack, bigAmount, true)
		So(err, ShouldBeNil)
		//jack D alice
		err = DelegationValidator(t, wallet, app, addJack, accJack, accAlice, bigAmount, true)
		So(err, ShouldBeNil)
		delegateAmount := types.NewInt64Coin(constants.DefaultBondDenom, 50000000)
		// alice R jack ->alice
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		// alice R alice ->jack
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accAlice, accJack, delegateAmount, false)
		So(err, ShouldNotBeNil)
		//  alice R jack ->alice 7 times
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, true)
		So(err, ShouldBeNil)
		err = RedelegateValidator(t, wallet, app, addAlice, accAlice, accJack, accAlice, delegateAmount, false)
		So(err, ShouldNotBeNil)
		//
	})
}
