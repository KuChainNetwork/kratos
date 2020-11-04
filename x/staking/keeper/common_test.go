package keeper_test // noalias

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"strconv"
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/store"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"

	"github.com/KuChainNetwork/kuchain/chain/config"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	stakingTypes "github.com/KuChainNetwork/kuchain/x/staking/types"
)

// dummy addresses used for testing
// nolint:unused, deadcode
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
	otherCoinDenom := types.CoinDenom(types.MustName("foo"), types.MustName("coin"))
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

//_____________________________________________________________________________________

// does a certain by-power index record exist
// func ValidatorByPowerIndexExists(ctx sdk.Context, keeper Keeper, power []byte) bool {
// 	store := store.NewStore(ctx, keeper.storeKey)
// 	return store.Has(power)
// }

// update validator for testing
func TestingUpdateValidator(app *simapp.SimApp, ctx sdk.Context, validator stakingTypes.Validator, apply bool) stakingTypes.Validator {
	keeper := app.StakeKeeper()
	keeper.SetValidator(ctx, validator)

	keeper.SetValidatorByPowerIndex(ctx, validator)
	if apply {
		keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		validator, found := keeper.GetValidator(ctx, validator.OperatorAccount)
		if !found {
			panic("validator expected but not found")
		}
		return validator
	}
	cachectx, _ := ctx.CacheContext()
	keeper.ApplyAndReturnValidatorSetUpdates(cachectx)
	validator, found := keeper.GetValidator(cachectx, validator.OperatorAccount)
	if !found {
		panic("validator expected but not found")
	}
	return validator
}

// nolint:deadcode, unused
func validatorByPowerIndexExists(app *simapp.SimApp, ctx sdk.Context, power []byte) bool {
	storkey := sdk.NewKVStoreKey(stakingTypes.StoreKey)
	store := store.NewStore(ctx, storkey)
	return store.Has(power)
}

// RandomValidator returns a random validator given access to the keeper and ctx
func RandomValidator(r *rand.Rand, app *simapp.SimApp, ctx sdk.Context) (val stakingTypes.Validator, ok bool) {
	keeper := app.StakeKeeper()
	vals := keeper.GetAllValidators(ctx)
	if len(vals) == 0 {
		return stakingTypes.Validator{}, false
	}

	i := r.Intn(len(vals))
	return vals[i], true
}

func TestInit(t *testing.T) {
	config.SealChainConfig()
}
