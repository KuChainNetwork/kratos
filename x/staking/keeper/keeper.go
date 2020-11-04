package keeper

import (
	"container/list"
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/store"
	"github.com/KuChainNetwork/kuchain/x/staking/external"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

const aminoCacheSize = 500

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Implements DelegationSet interface
var _ types.DelegationSet = Keeper{}

// keeper of the staking store
type Keeper struct {
	storeKey           sdk.StoreKey
	cdc                *codec.Codec
	bankKeeper         types.BankKeeper
	supplyKeeper       types.SupplyKeeper
	hooks              types.StakingHooks
	accountKeeper      types.AccountStatKeeper
	paramstore         external.ParamsSubspace
	validatorCache     map[string]cachedValidator
	validatorCacheList *list.List
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, bk types.BankKeeper, sk types.SupplyKeeper, ps external.ParamsSubspace, ak types.AccountStatKeeper,
) Keeper {

	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(ParamKeyTable())
	}

	// ensure bonded and not bonded module accounts are set
	if addr := sk.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := sk.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	return Keeper{
		storeKey:           key,
		cdc:                cdc,
		bankKeeper:         bk,
		supplyKeeper:       sk,
		paramstore:         ps,
		hooks:              nil,
		accountKeeper:      ak,
		validatorCache:     make(map[string]cachedValidator, aminoCacheSize),
		validatorCacheList: list.New(),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Set the validator hooks
func (k *Keeper) SetHooks(sh types.StakingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}
	k.hooks = sh
	return k
}

func (k *Keeper) EmptyHooks() *Keeper {
	k.hooks = nil
	return k
}

// Load the last total validator power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) sdk.Int {
	store := store.NewStore(ctx, k.storeKey)
	bz := store.Get(types.LastTotalPowerKey)
	if bz == nil {
		return sdk.ZeroInt()
	}

	ip := sdk.Int{}
	k.cdc.MustUnmarshalBinaryBare(bz, &ip)
	return ip
}

// Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := store.NewStore(ctx, k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(power)
	store.Set(types.LastTotalPowerKey, bz)
}

func (k Keeper) ValidatorAccount(ctx sdk.Context, id AccountID) bool {
	return k.accountKeeper.GetAccount(ctx, id) != nil
}

func (k Keeper) GetStoreKey() sdk.StoreKey {
	return k.storeKey
}
