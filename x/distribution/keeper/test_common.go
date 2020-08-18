package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/fee"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account"
	"github.com/KuChainNetwork/kuchain/x/asset"
	sktypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	//"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	Accounttypes "github.com/KuChainNetwork/kuchain/x/account/types"
	assettypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	"github.com/KuChainNetwork/kuchain/x/params"
	"github.com/KuChainNetwork/kuchain/x/staking"
	"github.com/KuChainNetwork/kuchain/x/supply"

	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	"github.com/tendermint/tendermint/crypto"
)

//nolint:deadcode,unused

func CreateAccount(keeper account.Keeper, ctx sdk.Context, id AccountID) {
	Name, _ := id.ToName()
	newAccount := keeper.NewAccountByName(ctx, Name)

	privKey := secp256k1.GenPrivKey()
	address := types.AccAddress(privKey.PubKey().Address())

	newAccount.SetAuth(address)

	// set Account
	keeper.SetAccount(ctx, newAccount)
	AccPubk[Name.String()] = privKey.PubKey()
}

var (
	TestMaster    = constants.ChainMainNameStr
	MasterName, _ = chainTypes.NewName(TestMaster)
	Master        = chainTypes.NewAccountIDFromName(MasterName)

	name1, _ = chainTypes.NewName("acc1")
	Acc1     = chainTypes.NewAccountIDFromName(name1)

	name2, _ = chainTypes.NewName("acc2")
	Acc2     = chainTypes.NewAccountIDFromName(name2)

	name3, _ = chainTypes.NewName("acc3")
	Acc3     = chainTypes.NewAccountIDFromName(name3)

	name4, _ = chainTypes.NewName("acc4")
	Acc4     = chainTypes.NewAccountIDFromName(name4)

	name5, _ = chainTypes.NewName("acc5")
	Acc5     = chainTypes.NewAccountIDFromName(name5)

	name6, _ = chainTypes.NewName("acc6")
	Acc6     = chainTypes.NewAccountIDFromName(name6)

	name7, _ = chainTypes.NewName("acc7")
	Acc7     = chainTypes.NewAccountIDFromName(name7)

	name8, _ = chainTypes.NewName("acc8")
	Acc8     = chainTypes.NewAccountIDFromName(name8)

	name9, _ = chainTypes.NewName("acc9")
	Acc9     = chainTypes.NewAccountIDFromName(name9)

	name10, _ = chainTypes.NewName("acc10")
	Acc10     = chainTypes.NewAccountIDFromName(name10)

	name11, _ = chainTypes.NewName("acc11")
	Acc11     = chainTypes.NewAccountIDFromName(name11)

	name12, _ = chainTypes.NewName("acc12")
	Acc12     = chainTypes.NewAccountIDFromName(name12)

	name13, _ = chainTypes.NewName("acc13")
	Acc13     = chainTypes.NewAccountIDFromName(name13)

	name14, _ = chainTypes.NewName("acc14")
	Acc14     = chainTypes.NewAccountIDFromName(name14)

	TestAddrs = []AccountID{
		Acc1, Acc2, Acc3, Acc4, Acc5, Acc6, Acc7, Acc8, Acc9, Acc10, Acc11, Acc12, Acc13, Acc14,
	}

	AccPubk = make(map[string]crypto.PubKey)
)

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()

	Accounttypes.RegisterCodec(cdc)
	assettypes.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc) // distr
	return cdc
}

// test input with default values
func CreateTestInputDefault(t *testing.T, isCheckTx bool, initPower int64) (
	sdk.Context, account.Keeper, Keeper, staking.Keeper, supply.Keeper, asset.Keeper) {

	communityTax := sdk.NewDecWithPrec(2, 2)

	ctx, ak, ask, dk, sk, _, supplyKeeper := CreateTestInputAdvanced(t, isCheckTx, initPower, communityTax)
	return ctx, ak, dk, sk, supplyKeeper, ask
}

// hogpodge of all sorts of input required for testing
func CreateTestInputAdvanced(t *testing.T, isCheckTx bool, initPower int64,
	communityTax sdk.Dec) (sdk.Context, account.Keeper, asset.Keeper,
	Keeper, staking.Keeper, params.Keeper, supply.Keeper) {

	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	cdc := MakeTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	AccountKeeper := account.NewAccountKeeper(cdc, sdk.NewKVStoreKey(account.StoreKey))

	mAccPerms := map[string][]string{
		fee.CollectorName:         nil,
		supply.BlackHole:          nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		types.ModuleName:          nil,
		staking.ModuleName:        nil,
	}

	assetKeeper := asset.NewAssetKeeper(cdc, sdk.NewKVStoreKey(asset.StoreKey), AccountKeeper)
	supplyKeeper := supply.NewKeeper(cdc, sdk.NewKVStoreKey(supply.StoreKey), AccountKeeper, assetKeeper, mAccPerms)

	distrAcc := supply.NewEmptyModuleAccount(types.ModuleName)
	feeCollectorAcc := supply.NewEmptyModuleAccount(fee.CollectorName)
	skModuleAcc := supply.NewEmptyModuleAccount(staking.ModuleName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[supply.NewModuleAddress(feeCollectorAcc.GetName().String()).String()] = true
	blacklistedAddrs[supply.NewModuleAddress(notBondedPool.GetName().String()).String()] = true
	blacklistedAddrs[supply.NewModuleAddress(bondPool.GetName().String()).String()] = true
	blacklistedAddrs[supply.NewModuleAddress(distrAcc.GetName().String()).String()] = true
	blacklistedAddrs[supply.NewModuleAddress(skModuleAcc.GetName().String()).String()] = true

	sk := staking.NewKeeper(
		cdc, sdk.NewKVStoreKey(staking.StoreKey), assetKeeper, supplyKeeper, pk.Subspace(staking.DefaultParamspace), AccountKeeper)

	keeper := NewKeeper(cdc, sdk.NewKVStoreKey(types.StoreKey), pk.Subspace(types.DefaultParamspace),
		assetKeeper, sk, supplyKeeper, AccountKeeper, fee.CollectorName, blacklistedAddrs)

	ms.MountStoreWithDB(keeper.storeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(sk.GetStoreKey(), sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(supplyKeeper.GetStoreKey(), sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(assetKeeper.GetStoreKey(), sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(AccountKeeper.GetStoreKey(), sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	for _, addr := range TestAddrs {
		CreateAccount(AccountKeeper, ctx, addr)
	}
	CreateAccount(AccountKeeper, ctx, Master)

	sk.SetParams(ctx, staking.DefaultParams())

	intNum1, _ := sdk.NewIntFromString("1000000000000000000000")
	intNum2, _ := sdk.NewIntFromString("80000000000000000000000")
	intNum3, _ := sdk.NewIntFromString("60000000000000000000000")
	intNumFee, _ := sdk.NewIntFromString("20000000000000000000000")
	intMaxNum, _ := sdk.NewIntFromString("100000000000000000000000")

	SymbolName, _ := chainTypes.NewName(constants.DefaultBondSymbol)

	assetKeeper.Create(ctx, MasterName, SymbolName, assettypes.NewCoin(constants.DefaultBondDenom, intNum2),
		true, true, true, 0, assettypes.NewCoin(constants.DefaultBondDenom, intMaxNum), []byte("create"))

	assetKeeper.Issue(ctx, MasterName, SymbolName,
		assettypes.NewCoin(constants.DefaultBondDenom, intNum3))

	{
		for _, addr := range TestAddrs {
			Coins := chainTypes.NewCoins(chainTypes.NewCoin(constants.DefaultBondDenom, intNum1))
			err := assetKeeper.Transfer(ctx, Master, addr, Coins)
			//fmt.Println("id", id, "account", addr.String())
			require.Nil(t, err)
		}
		Coins := chainTypes.NewCoins(chainTypes.NewCoin(constants.DefaultBondDenom, intNumFee))
		err := assetKeeper.Transfer(ctx, Master, supplyKeeper.GetModuleAccount(ctx, keeper.feeCollectorName).GetID(), Coins)
		require.Nil(t, err)
	}

	// set module Accounts
	keeper.supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	keeper.supplyKeeper.SetModuleAccount(ctx, bondPool)
	keeper.supplyKeeper.SetModuleAccount(ctx, distrAcc)
	keeper.supplyKeeper.SetModuleAccount(ctx, skModuleAcc)

	// set the distribution hooks on staking
	sk.SetHooks(keeper.Hooks())

	// set genesis items required for distribution
	keeper.SetFeePool(ctx, types.InitialFeePool())

	params := types.DefaultParams()
	params.CommunityTax = communityTax
	params.BaseProposerReward = sdk.NewDecWithPrec(1, 2)
	params.BonusProposerReward = sdk.NewDecWithPrec(4, 2)
	keeper.SetParams(ctx, params)

	return ctx, AccountKeeper, assetKeeper, keeper, sk, pk, supplyKeeper
}

func GetDescription() sktypes.Description {
	FlagMoniker := "moniker"
	FlagIdentity := "identity"
	FlagWebsite := "website"
	FlagSecurityContact := "security-contact"
	FlagDetails := "details"

	description := sktypes.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagSecurityContact),
		viper.GetString(FlagDetails),
	)

	return description
}
