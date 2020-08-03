package simapp

import (
	"os"
	"testing"
	"time"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp/helpers"
	"github.com/KuChainNetwork/kuchain/x/account"
	"github.com/KuChainNetwork/kuchain/x/asset"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(isCheckTx bool) *SimApp {
	db := dbm.NewMemDB()
	app := NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, 0)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := NewDefaultGenesisState()
		stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:    []abci.ValidatorUpdate{},
				AppStateBytes: stateBytes,
			},
		)
	}

	return app
}

// SetupWithGenesisAccounts initializes a new SimApp with the passed in
// genesis accounts.
func SetupWithGenesisAccounts(genAccs *GenesisAccounts) *SimApp {
	db := dbm.NewMemDB()
	app := NewSimApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, 0)

	// initialize the chain with the passed in genesis accounts
	genesisState := NewDefaultGenesisState()

	accountGenesis := account.NewGenesisState(genAccs.accounts)
	accountGenesisStateBz := app.Codec().MustMarshalJSON(accountGenesis)
	genesisState[account.ModuleName] = accountGenesisStateBz

	assetGenesis := asset.GenesisState{
		GenesisAssets: genAccs.assets,
		GenesisCoins:  genAccs.coins,
	}
	assetGenesisStateBz := app.Codec().MustMarshalJSON(assetGenesis)
	genesisState[asset.ModuleName] = assetGenesisStateBz

	stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
	if err != nil {
		panic(err)
	}

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
			ChainId:       "",
		},
	)

	app.Commit()

	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: app.LastBlockHeight() + 1}})

	return app
}

// AddTestAddress constructs and returns accNum amount of accounts with an
// initial balance of accAmt
func AddTestAddress(app *SimApp, ctx sdk.Context, accNum int, accAmt sdk.Int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	initCoins := types.NewCoins(types.NewCoin(constants.DefaultBondDenom, accAmt))
	totalSupply := types.NewCoin(constants.DefaultBondDenom, accAmt.MulRaw(int64(len(testAddrs))))

	if err := app.assetKeeper.Issue(ctx, constants.SystemAccount, constants.DefaultBondSymbolName, totalSupply); err != nil {
		panic(err)
	}

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range testAddrs {
		addrID := types.NewAccountIDFromAccAdd(addr)
		err := app.assetKeeper.Transfer(ctx, constants.SystemAccountID, addrID, initCoins)
		if err != nil {
			panic(err)
		}
	}
	return testAddrs
}

// CheckBalance checks the balance of an account.
func CheckBalance(t *testing.T, app *SimApp, id types.AccountID, exp types.Coins) {
	res, err := app.assetKeeper.GetCoins(app.NewTestContext(), id)

	if err != nil {
		panic(err)
	}

	if !exp.IsEqual(res) {
		app.NewTestContext().Logger().Error("check balance err", "id", id, "exp", exp, "res", res)
	}

	So(exp.IsEqual(res), ShouldBeTrue)
}

// SignCheckDeliver checks a generated signed transaction and simulates a
// block commitment with the given transaction. A test assertion is made using
// the parameter 'expPass' against the result. A corresponding result is
// returned.
func SignCheckDeliver(
	t *testing.T, cdc *codec.Codec, app *bam.BaseApp, header abci.Header,
	payer types.AccountID, fee types.Coins,
	msgs []sdk.Msg,
	accNums, seq []uint64, expSimPass, expPass bool, priv ...crypto.PrivKey,
) (sdk.GasInfo, *sdk.Result, error) {

	tx := helpers.GenTx(
		msgs,
		fee,
		helpers.DefaultGenTxGas,
		payer,
		"",
		accNums,
		seq,
		priv...,
	)

	txBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)
	So(err, ShouldBeNil)

	// Must simulate now as CheckTx doesn't run Msgs anymore
	_, res, err := app.Simulate(txBytes, tx)

	if expSimPass {
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	} else {
		So(err, ShouldNotBeNil)
		So(res, ShouldBeNil)
	}

	// Simulate a sending a transaction and committing a block
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	gInfo, res, err := app.Deliver(tx)

	if expPass {
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	} else {
		So(err, ShouldNotBeNil)
		So(res, ShouldBeNil)
	}

	app.EndBlock(abci.RequestEndBlock{
		Height: header.GetHeight(),
	})
	app.Commit()

	return gInfo, res, err
}

// GenSequenceOfTxs generates a set of signed transactions of messages, such
// that they differ only by having the sequence numbers incremented between
// every transaction.
func GenSequenceOfTxs(payer types.AccountID, msgs []sdk.Msg, accNums []uint64, initSeqNums []uint64, numToGenerate int, priv ...crypto.PrivKey) []types.StdTx {
	txs := make([]types.StdTx, numToGenerate)
	for i := 0; i < numToGenerate; i++ {
		txs[i] = helpers.GenTx(
			msgs,
			types.NewInt64Coins(constants.DefaultBondDenom, 200000),
			helpers.DefaultGenTxGas,
			payer,
			"",
			accNums,
			initSeqNums,
			priv...,
		)
		incrementAllSequenceNumbers(initSeqNums)
	}

	return txs
}

func incrementAllSequenceNumbers(initSeqNums []uint64) {
	for i := 0; i < len(initSeqNums); i++ {
		initSeqNums[i]++
	}
}

func (app *SimApp) NewTestContext() sdk.Context {
	return app.BaseApp.NewContext(true,
		abci.Header{
			Height: app.LastBlockHeight() + 1,
			Time:   time.Now(),
		})
}

func AfterBlockCommitted(app *SimApp, blockNum int) {
	for i := 0; i < blockNum; i++ {
		header := abci.Header{
			Height: app.LastBlockHeight() + 1,
			Time:   time.Now(),
		}

		app.BeginBlock(abci.RequestBeginBlock{
			Header: header,
		})
		app.EndBlock(abci.RequestEndBlock{
			Height: header.Height,
		})
		app.Commit()
	}
}
