package gov_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"

	"encoding/hex"
	"github.com/KuChainNetwork/kuchain/chain/config"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/msg"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	govTypes "github.com/KuChainNetwork/kuchain/x/gov/types"
	"github.com/KuChainNetwork/kuchain/x/params/client/utils"
	paramproposal "github.com/KuChainNetwork/kuchain/x/params/types/proposal"
	stakingTypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"time"
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
	asset1 := types.Coins{
		types.NewInt64Coin(otherCoinDenom, 6700000000000),
		initAsset}

	asset2 := types.Coins{
		types.NewInt64Coin(otherCoinDenom, 6700000000000),
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

func signBlock(app *simapp.SimApp) {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
	blocktime := time.Now().Add(app.GovKeeper().GetVotingParams(ctxCheck).VotingPeriod * 2)
	header := abci.Header{Height: app.LastBlockHeight() + 1, Time: blocktime}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()
}

func submitProposal(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice types.AccountID, content govTypes.Content, amount, transferAmount types.Coins, passed bool) error {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addAlice)
	So(err, ShouldBeNil)
	msg := customizeKuMsgProposal(addAlice, content, amount, transferAmount, accAlice)
	fee := types.Coins{types.NewInt64Coin(constants.DefaultBondDenom, 1000000)}
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
		header, accAlice, fee,
		[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		passed, passed, wallet.PrivKey(addAlice))
	ctxCheck.Logger().Info("submitProposal error log", "err", err)

	return err
}

func vote(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice types.AccountID, proposalID uint64, voteOption govTypes.VoteOption, passed bool) error {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addAlice)
	So(err, ShouldBeNil)
	msg := govTypes.NewKuMsgVote(addAlice, accAlice, proposalID, voteOption)
	fee := types.Coins{types.NewInt64Coin(constants.DefaultBondDenom, 1000000)}
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
		header, accAlice, fee,
		[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		passed, passed, wallet.PrivKey(addAlice))
	ctxCheck.Logger().Info("vote error log", "err", err)

	return err
}

func customizeKuMsgDeposit(auth sdk.AccAddress, depositor types.AccountID, proposalID uint64, amount, transferAmount types.Coins) govTypes.KuMsgDeposit {
	return govTypes.KuMsgDeposit{
		*msg.MustNewKuMsg(
			govTypes.RouterKeyName,
			msg.WithAuth(auth),
			msg.WithTransfer(depositor, govTypes.ModuleAccountID, transferAmount),
			msg.WithData(govTypes.Cdc(), &govTypes.MsgDeposit{proposalID, depositor, amount}),
		),
	}
}

func customizeKuMsgProposal(auth sdk.AccAddress, content govTypes.Content, initialDeposit, transferAmount types.Coins, proposer types.AccountID) govTypes.KuMsgSubmitProposal {
	return govTypes.KuMsgSubmitProposal{
		*msg.MustNewKuMsg(
			govTypes.RouterKeyName,
			msg.WithAuth(auth),
			msg.WithTransfer(proposer, govTypes.ModuleAccountID, transferAmount),
			msg.WithData(govTypes.Cdc(), &govTypes.MsgSubmitProposalBase{
				InitialDeposit: initialDeposit,
				Proposer:       proposer,
			}),
		), content,
	}
}

func disposit(t *testing.T, wallet *simapp.Wallet, app *simapp.SimApp, addAlice sdk.AccAddress, accAlice types.AccountID, proposalID uint64, amount, transferAmount types.Coins, passed bool) error {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	origAuthSeq, origAuthNum, err := app.AccountKeeper().GetAuthSequence(ctxCheck, addAlice)
	So(err, ShouldBeNil)
	msg := customizeKuMsgDeposit(addAlice, accAlice, proposalID, amount, transferAmount)

	fee := types.Coins{types.NewInt64Coin(constants.DefaultBondDenom, 1000000)}
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, app.Codec(), app.BaseApp,
		header, accAlice, fee,
		[]sdk.Msg{msg}, []uint64{origAuthNum}, []uint64{origAuthSeq},
		passed, passed, wallet.PrivKey(addAlice))
	ctxCheck.Logger().Info("disposit error log", "err", err)
	return err
}

func logProposal(app *simapp.SimApp, proposalID uint64) {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
	proposal, ok := app.GovKeeper().GetProposal(ctxCheck, proposalID)
	if ok {
		ctxCheck.Logger().Info("logProposal", "proposal", proposal)
	}
}

func logVote(app *simapp.SimApp, proposalID uint64) {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
	vote := app.GovKeeper().GetVotes(ctxCheck, proposalID)
	ctxCheck.Logger().Info("logVote", "vote", vote)
}

func logTally(app *simapp.SimApp, proposalID uint64) {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
	proposal, ok := app.GovKeeper().GetProposal(ctxCheck, proposalID)
	if ok {
		_, _, tallyResult, _, _, _ := app.GovKeeper().Tally(ctxCheck, proposal)
		ctxCheck.Logger().Info("logTally", "tallyResult", tallyResult)
	}
}

func getGovParams(app *simapp.SimApp, proposalID uint64) {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
	vote := app.GovKeeper().GetVotes(ctxCheck, proposalID)
	ctxCheck.Logger().Info("logVote", "vote", vote)
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

func TestGovHandler(t *testing.T) {
	Convey("TestGovHandler", t, func() {
		config.SealChainConfig()
		wallet := simapp.NewWallet()

		addAlice, _, addValidator, accAlice, _, accValidator, app := newTestApp(wallet)
		rate, _ := sdk.NewDecFromStr("0.6")
		pk := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF100")
		err := createValidator(t, wallet, app, addAlice, accAlice, rate, pk, true)
		So(err, ShouldBeNil)
		initdepost := types.NewInt64Coin(constants.DefaultBondDenom, 1000000000000000000)
		otherCoinDenom := types.CoinDenom(types.MustName("foo"), types.MustName("coin"))
		otherdeposit := types.NewInt64Coin(otherCoinDenom, 1000000)

		depositInt, succ := sdk.NewIntFromString("500000000000000000000")
		if !succ {
			depositInt = sdk.NewInt(500000000000000000)
		}
		depositAsset := types.NewCoin(constants.DefaultBondDenom, depositInt)

		err = delegationValidator(t, wallet, app, addAlice, accAlice, accAlice, initdepost, true)
		So(err, ShouldBeNil)
		//proposal normal text proposal
		textContent := govTypes.ContentFromProposalType("test title", "test decription", govTypes.ProposalTypeText)
		err = submitProposal(t, wallet, app, addAlice, accAlice, textContent, types.Coins{initdepost}, types.Coins{initdepost}, true)
		So(err, ShouldBeNil)
		//deposit other coin
		err = disposit(t, wallet, app, addAlice, accAlice, 1, types.Coins{otherdeposit}, types.Coins{otherdeposit}, true)
		So(err, ShouldBeNil)
		//vote before depost enough coin
		err = vote(t, wallet, app, addAlice, accAlice, 1, govTypes.OptionYes, false)
		So(err, ShouldNotBeNil)
		//deposit but account does not have enough coin
		err = disposit(t, wallet, app, addValidator, accValidator, 1, types.Coins{depositAsset}, types.Coins{depositAsset}, false)
		So(err, ShouldNotBeNil)
		//deposit
		err = disposit(t, wallet, app, addAlice, accAlice, 1, types.Coins{depositAsset}, types.Coins{depositAsset}, true)
		So(err, ShouldBeNil)
		//wrong coins
		depositCoins := types.Coins{depositAsset}
		initCoins := types.Coins{initdepost}
		wrongCoin, _ := initCoins.SafeSub(depositCoins)
		err = disposit(t, wallet, app, addAlice, accAlice, 1, wrongCoin, wrongCoin, false)
		So(err, ShouldNotBeNil)
		//deposit after proposal have enough coin
		err = disposit(t, wallet, app, addAlice, accAlice, 1, types.Coins{depositAsset}, types.Coins{depositAsset}, true)
		So(err, ShouldBeNil)
		//vote
		err = vote(t, wallet, app, addAlice, accAlice, 1, govTypes.OptionYes, true)
		So(err, ShouldBeNil)
		//vote after proposal close
		err = vote(t, wallet, app, addAlice, accAlice, 1, govTypes.OptionYes, false)
		So(err, ShouldNotBeNil)
		//deposit after proposal close
		err = disposit(t, wallet, app, addAlice, accAlice, 1, types.Coins{depositAsset}, types.Coins{depositAsset}, false)
		So(err, ShouldNotBeNil)
		// transfer amount not equal deposit amount
		textContent = govTypes.ContentFromProposalType("test title2", "test decription2", govTypes.ProposalTypeText)
		err = submitProposal(t, wallet, app, addAlice, accAlice, textContent, types.Coins{initdepost}, types.Coins{depositAsset}, false)
		So(err, ShouldNotBeNil)
	})
	Convey("TestTextProposal", t, func() {
		wallet := simapp.NewWallet()

		addAlice, addJack, addValidator, accAlice, accJack, accValidator, app := newTestApp(wallet)
		rate, _ := sdk.NewDecFromStr("0.6")
		pk := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF100")
		err := createValidator(t, wallet, app, addAlice, accAlice, rate, pk, true)
		So(err, ShouldBeNil)
		Newpk := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF200")
		err = createValidator(t, wallet, app, addJack, accJack, rate, Newpk, true)
		So(err, ShouldBeNil)
		initdepost := types.NewInt64Coin(constants.DefaultBondDenom, 1000000000000000000)
		jackdepost := types.NewInt64Coin(constants.DefaultBondDenom, 1500000000000000000)

		err = delegationValidator(t, wallet, app, addAlice, accAlice, accAlice, initdepost, true)
		So(err, ShouldBeNil)
		err = delegationValidator(t, wallet, app, addJack, accJack, accJack, jackdepost, true)
		So(err, ShouldBeNil)
		//proposaler has no enough coin
		textContent := govTypes.ContentFromProposalType("test title", "test decription", govTypes.ProposalTypeText)
		err = submitProposal(t, wallet, app, addValidator, accValidator, textContent, types.Coins{initdepost}, types.Coins{initdepost}, false)
		So(err, ShouldNotBeNil)

		longtitle := "test title test title test title test title test title test title test title test title test title test title test title test title test title test title "
		longdecription := "test decription test decription test decription test decription test decription test decription test decription test decription "
		for i := 0; i < 20; i++ {
			longdecription = longdecription + longdecription
		}
		// title too long
		textContent = govTypes.ContentFromProposalType(longtitle, "longdecription", govTypes.ProposalTypeText)
		err = submitProposal(t, wallet, app, addAlice, accAlice, textContent, types.Coins{initdepost}, types.Coins{initdepost}, false)
		So(err, ShouldNotBeNil)
		// decription too long
		textContent = govTypes.ContentFromProposalType("longtitle", longdecription, govTypes.ProposalTypeText)
		err = submitProposal(t, wallet, app, addAlice, accAlice, textContent, types.Coins{initdepost}, types.Coins{initdepost}, false)
		So(err, ShouldNotBeNil)

		//right proposal
		textContent = govTypes.ContentFromProposalType("test title", "test decription", govTypes.ProposalTypeText)
		err = submitProposal(t, wallet, app, addAlice, accAlice, textContent, types.Coins{initdepost}, types.Coins{initdepost}, true)
		So(err, ShouldBeNil)
		//deposit
		depositInt, succ := sdk.NewIntFromString("500000000000000000000")
		if !succ {
			depositInt = sdk.NewInt(500000000000000000)
		}
		depositAsset := types.NewCoin(constants.DefaultBondDenom, depositInt)
		// transfer amount not equal deposit amount
		err = disposit(t, wallet, app, addAlice, accAlice, 1, types.Coins{depositAsset}, types.Coins{initdepost}, false)
		So(err, ShouldNotBeNil)
		err = disposit(t, wallet, app, addAlice, accAlice, 1, types.Coins{depositAsset}, types.Coins{depositAsset}, true)
		So(err, ShouldBeNil)
		logProposal(app, 1)
		//no exist proposal
		err = disposit(t, wallet, app, addAlice, accAlice, 2, types.Coins{depositAsset}, types.Coins{depositAsset}, false)
		So(err, ShouldNotBeNil)
		//vote
		err = vote(t, wallet, app, addAlice, accAlice, 1, govTypes.OptionYes, true)
		So(err, ShouldBeNil)
		//already vote
		err = vote(t, wallet, app, addAlice, accAlice, 1, govTypes.OptionNo, true)
		So(err, ShouldBeNil)
		// //jack vote
		err = vote(t, wallet, app, addJack, accJack, 1, govTypes.OptionYes, true)
		So(err, ShouldBeNil)
		logVote(app, 1)
		logTally(app, 1)
		signBlock(app)
		logProposal(app, 1)
		//second proposal
		err = submitProposal(t, wallet, app, addAlice, accAlice, textContent, types.Coins{initdepost}, types.Coins{initdepost}, true)
		So(err, ShouldBeNil)
		err = disposit(t, wallet, app, addAlice, accAlice, 2, types.Coins{depositAsset}, types.Coins{depositAsset}, true)
		So(err, ShouldBeNil)
		logProposal(app, 2)
		err = vote(t, wallet, app, addJack, accJack, 2, govTypes.OptionNo, true)
		So(err, ShouldBeNil)
		err = vote(t, wallet, app, addAlice, accAlice, 2, govTypes.OptionYes, true)
		So(err, ShouldBeNil)
		logVote(app, 2)
		logTally(app, 2)
		signBlock(app)
		logProposal(app, 2)
		//third proposal
		err = submitProposal(t, wallet, app, addAlice, accAlice, textContent, types.Coins{initdepost}, types.Coins{initdepost}, true)
		So(err, ShouldBeNil)
		err = disposit(t, wallet, app, addAlice, accAlice, 3, types.Coins{depositAsset}, types.Coins{depositAsset}, true)
		So(err, ShouldBeNil)
		logProposal(app, 3)
		err = vote(t, wallet, app, addJack, accJack, 3, govTypes.OptionYes, true)
		So(err, ShouldBeNil)
		err = vote(t, wallet, app, addAlice, accAlice, 3, govTypes.OptionNoWithVeto, true)
		So(err, ShouldBeNil)
		logVote(app, 3)
		logTally(app, 3)
		signBlock(app)
		logProposal(app, 3)
		//fouth proposal
		err = submitProposal(t, wallet, app, addAlice, accAlice, textContent, types.Coins{initdepost}, types.Coins{initdepost}, true)
		So(err, ShouldBeNil)
		signBlock(app)
		logProposal(app, 4)
		//fivth proposal
		err = submitProposal(t, wallet, app, addAlice, accAlice, textContent, types.Coins{initdepost}, types.Coins{initdepost}, true)
		So(err, ShouldBeNil)
		logProposal(app, 5)
		err = disposit(t, wallet, app, addAlice, accAlice, 5, types.Coins{depositAsset}, types.Coins{depositAsset}, true)
		So(err, ShouldBeNil)
		err = vote(t, wallet, app, addJack, accJack, 5, govTypes.OptionAbstain, true)
		So(err, ShouldBeNil)
		signBlock(app)
		logProposal(app, 5)
		//sixth proposal
		err = submitProposal(t, wallet, app, addAlice, accAlice, textContent, types.Coins{initdepost}, types.Coins{initdepost}, true)
		So(err, ShouldBeNil)
		err = disposit(t, wallet, app, addAlice, accAlice, 6, types.Coins{depositAsset}, types.Coins{depositAsset}, true)
		So(err, ShouldBeNil)
		signBlock(app)
		logProposal(app, 6)
	})
	Convey("TestParamsChangeProposal", t, func() {
		wallet := simapp.NewWallet()

		addAlice, _, _, accAlice, _, _, app := newTestApp(wallet)
		rate, _ := sdk.NewDecFromStr("0.6")
		pk := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF100")
		err := createValidator(t, wallet, app, addAlice, accAlice, rate, pk, true)
		So(err, ShouldBeNil)
		initdepost := types.NewInt64Coin(constants.DefaultBondDenom, 1000000000000000000)

		depositInt, succ := sdk.NewIntFromString("500000000000000000000")
		if !succ {
			depositInt = sdk.NewInt(500000000000000000)
		}
		depositAsset := types.NewCoin(constants.DefaultBondDenom, depositInt)

		err = delegationValidator(t, wallet, app, addAlice, accAlice, accAlice, initdepost, true)
		So(err, ShouldBeNil)

		jsonChangevotingparams := `
		{
		"title": "Staking Param Change",
		"description": "Update voting period",
		"changes": [
			{
			"subspace": "kugov",
			"key": "votingparams",
			"value": {
				"voting_period": "1209800000000000"
			}
			}
		],
		"deposit": "1000kuchain/kcs"
		}
		`
		proposal := utils.ParamChangeProposalJSON{}

		app.Codec().UnmarshalJSON([]byte(jsonChangevotingparams), &proposal)
		content := paramproposal.NewParameterChangeProposal(proposal.Title, proposal.Description, proposal.Changes.ToParamChanges())

		err = submitProposal(t, wallet, app, addAlice, accAlice, content, types.Coins{initdepost}, types.Coins{initdepost}, true)
		So(err, ShouldBeNil)

		//deposit other coin
		err = disposit(t, wallet, app, addAlice, accAlice, 1, types.Coins{depositAsset}, types.Coins{depositAsset}, true)
		So(err, ShouldBeNil)
		err = vote(t, wallet, app, addAlice, accAlice, 1, govTypes.OptionYes, true)
		So(err, ShouldBeNil)

		ctxCheck := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		votingParam := app.GovKeeper().GetVotingParams(ctxCheck)
		ctxCheck.Logger().Info("logVote", "votingParam", votingParam)
		So(int64(1209800000000000) == int64(votingParam.VotingPeriod), ShouldBeTrue)
	})
}
