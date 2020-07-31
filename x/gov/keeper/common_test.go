// nolint
package keeper_test // noalias

// DONTCOVER

import (
	"bytes"
	"encoding/hex"
	govKeeper "github.com/KuChainNetwork/kuchain/x/gov/keeper"
	govTypes "github.com/KuChainNetwork/kuchain/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"strconv"

	"github.com/KuChainNetwork/kuchain/x/staking"
	stakeTypes "github.com/KuChainNetwork/kuchain/x/staking/types"

	"github.com/KuChainNetwork/kuchain/chain/config"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/staking/exported"
	"github.com/KuChainNetwork/kuchain/x/supply"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	Addrs  = createTestAddrs(500)
	PKs    = createTestPubKeys(500)
	Accd   = createTestAccount(500)
	Accdel = createTestAccount(500)

	addrDels = []sdk.AccAddress{
		Addrs[0],
		Addrs[1],
	}
	addrVals = []sdk.ValAddress{
		sdk.ValAddress(Addrs[2]),
		sdk.ValAddress(Addrs[3]),
		sdk.ValAddress(Addrs[4]),
		sdk.ValAddress(Addrs[5]),
		sdk.ValAddress(Addrs[6]),
	}
)

// dummy addresses used for testing
var (
	accAlice     = types.MustAccountID("alice@ok")
	accJack      = types.MustAccountID("jack@ok")
	accValidator = types.MustAccountID("validator@ok")

	delPk1   = PKs[1]
	delPk2   = PKs[2]
	delPk3   = PKs[3]
	delAddr1 = accAlice
	delAddr2 = accJack
	delAddr3 = accValidator

	valOpPk1    = PKs[4]
	valOpPk2    = PKs[5]
	valOpPk3    = PKs[6]
	valOpAddr1  = accAlice
	valOpAddr2  = accJack
	valOpAddr3  = accValidator
	valAccAddr1 = accAlice
	valAccAddr2 = accJack
	valAccAddr3 = accValidator

	TestAddrs = []types.AccountID{
		delAddr1, delAddr2, delAddr3,
		valAccAddr1, valAccAddr2, valAccAddr3,
	}

	powers = []int64{1, 2, 3}

	emptyDelAddr sdk.AccAddress
	emptyValAddr sdk.ValAddress
	emptyPubkey  crypto.PubKey
)

// TODO: remove dependency with staking
var (
	TestProposal        = govTypes.NewTextProposal("Test", "description")
	TestDescription     = staking.NewDescription("T", "E", "S", "T", "Z")
	TestCommissionRates = staking.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
)

// Hogpodge of all sorts of input required for testing.
// `initPower` is converted to an amount of tokens.
// If `initPower` is 0, no addrs get created.
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
	otherCoinDenom := types.CoinDenom(types.MustName("foo"),types.MustName("coin"))
	initAsset := types.NewCoin(constants.DefaultBondDenom, resInt)
	asset1 := types.Coins{
		types.NewInt64Coin(otherCoinDenom, 67),
		initAsset}

	asset2 := types.Coins{
		types.NewInt64Coin(otherCoinDenom, 67),
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

	ModuleName, _ := govTypes.ModuleAccountID.ToName()
	StakeModuleName, _ := stakeTypes.ModuleAccountID.ToName()

	test0Name, _ := TestAddrs[0].ToName()
	test1Name, _ := TestAddrs[1].ToName()
	app.AssetKeeper().Issue(ctxCheck, StakeModuleName, StakeModuleName, types.NewCoin(app.StakeKeeper().BondDenom(ctxCheck), exported.TokensFromConsensusPower(1000000)))
	app.AssetKeeper().Issue(ctxCheck, ModuleName, ModuleName, types.NewCoin(app.StakeKeeper().BondDenom(ctxCheck), exported.TokensFromConsensusPower(1000000)))
	app.AssetKeeper().Issue(ctxCheck, test0Name, test0Name, types.NewCoin(app.StakeKeeper().BondDenom(ctxCheck), exported.TokensFromConsensusPower(1000000)))
	app.AssetKeeper().Issue(ctxCheck, test1Name, test1Name, types.NewCoin(app.StakeKeeper().BondDenom(ctxCheck), exported.TokensFromConsensusPower(1000000)))
	app.AssetKeeper().IssueCoinPower(ctxCheck, govTypes.ModuleAccountID, types.NewCoins(types.NewCoin(app.StakeKeeper().BondDenom(ctxCheck), exported.TokensFromConsensusPower(1000000))))

	return addAlice, addJack, addValidator, accAlice, accJack, accValidator, app
}

func NewPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	//res, err = crypto.PubKeyFromBytes(pkBytes)
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes)
	return pkEd
}

// for incode address generation
func GetTestAddr(addr string, bech string) sdk.AccAddress {

	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		panic(err)
	}
	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}

	return res
}

// nolint: unparam
func createTestAddrs(numAddrs int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addresses = append(addresses, GetTestAddr(buffer.String(), bech))
		buffer.Reset()
	}
	return addresses
}

func createTestAccount(numAddrs int) []types.AccountID {
	var accountID []types.AccountID
	var buffer bytes.Buffer

	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("testaccount") //base address string
		buffer.WriteString(numString)     //adding on final two digits to make addresses unique

		tmpAccount := types.MustAccountID(buffer.String())
		accountID = append(accountID, tmpAccount)
		buffer.Reset()
	}
	return accountID
}

// nolint: unparam
func createTestPubKeys(numPubKeys int) []crypto.PubKey {
	var publicKeys []crypto.PubKey
	var buffer bytes.Buffer

	//start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") //base pubkey string
		buffer.WriteString(numString)                                                       //adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKey(buffer.String()))
		buffer.Reset()
	}
	return publicKeys
}

func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// ProposalEqual checks if two proposals are equal (note: slow, for tests only)
func ProposalEqual(keeper *govKeeper.Keeper, proposalA govTypes.Proposal, proposalB govTypes.Proposal) bool {
	return bytes.Equal(keeper.MustMarshalProposal(proposalA),
		keeper.MustMarshalProposal(proposalB))
}

func createValidators(app *simapp.SimApp, ctx sdk.Context, sk *staking.Keeper, powers []int64) {
	val1 := staking.NewValidator(valOpAddr1, valOpPk1, staking.Description{})
	val2 := staking.NewValidator(valOpAddr2, valOpPk2, staking.Description{})
	val3 := staking.NewValidator(valOpAddr3, valOpPk3, staking.Description{})

	sk.SetValidator(ctx, val1)
	sk.SetValidator(ctx, val2)
	sk.SetValidator(ctx, val3)
	sk.SetValidatorByConsAddr(ctx, val1)
	sk.SetValidatorByConsAddr(ctx, val2)
	sk.SetValidatorByConsAddr(ctx, val3)
	sk.SetNewValidatorByPowerIndex(ctx, val1)
	sk.SetNewValidatorByPowerIndex(ctx, val2)
	sk.SetNewValidatorByPowerIndex(ctx, val3)

	_, _ = sk.Delegate(ctx, valAccAddr1, exported.TokensFromConsensusPower(powers[0]), exported.Unbonded, val1, true)
	_, _ = sk.Delegate(ctx, valAccAddr2, exported.TokensFromConsensusPower(powers[1]), exported.Unbonded, val2, true)
	_, _ = sk.Delegate(ctx, valAccAddr3, exported.TokensFromConsensusPower(powers[2]), exported.Unbonded, val3, true)

	notBondedPool := sk.GetNotBondedPool(ctx)
	app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), types.NewCoins(types.NewCoin(sk.BondDenom(ctx), exported.TokensFromConsensusPower(powers[0]))))
	app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), types.NewCoins(types.NewCoin(sk.BondDenom(ctx), exported.TokensFromConsensusPower(powers[1]))))
	app.AssetKeeper().IssueCoinPower(ctx, notBondedPool.GetID(), types.NewCoins(types.NewCoin(sk.BondDenom(ctx), exported.TokensFromConsensusPower(powers[2]))))

	_ = staking.EndBlocker(ctx, *sk)
}

func TestInit(t *testing.T) {
	config.SealChainConfig()
}
