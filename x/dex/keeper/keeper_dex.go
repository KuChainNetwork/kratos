package keeper

import (
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

// CreateDex create a dex by creator
func (a DexKeeper) CreateDex(ctx sdk.Context, creator types.Name, staking types.Coins, description string) error {
	if _, ok := a.getDex(ctx, creator); ok {
		return errors.Wrapf(types.ErrDexHadCreated, "dex %s already exists", creator.String())
	}

	dex := types.NewDex(creator, staking, description).WithNumber(a.nextNumber(ctx))
	a.setDex(ctx, dex)
	return nil
}

// DestroyDex delete a dex by creator
func (a DexKeeper) DestroyDex(ctx sdk.Context, creator types.Name) error {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		return errors.Wrapf(types.ErrDexNotExists, "dex %s not exists", creator.String())
	}

	// check the dex can be destroyed
	if !dex.CanDestroy(func() chainTypes.Coins {
		return a.GetSigInSumForDex(ctx, accountTypes.NewAccountIDFromName(creator))
	}) {
		return errors.Wrapf(types.ErrDexCanNotBeDestroyed,
			"dex %s can not be destroy", creator.String())
	}

	// transfer asset to coin power
	if err := a.assetKeeper.CoinsToPower(ctx,
		types.ModuleAccountID,
		chainTypes.NewAccountIDFromName(creator),
		dex.Staking); err != nil {
		return errors.Wrapf(err, "transfer coin power error")
	}

	a.deleteDex(ctx, dex)
	return nil
}

// UpdateDexDescription update a dex description
func (a DexKeeper) UpdateDexDescription(ctx sdk.Context,
	creator types.Name, description string) error {
	var (
		dex   *types.Dex
		found bool
	)

	if dex, found = a.getDex(ctx, creator); !found {
		return errors.Wrapf(types.ErrDexNotExists, "dex %s not exists", creator.String())
	}

	if dex.Description == description {
		return errors.Wrapf(types.ErrDexDescriptionSame,
			"msg update dex %s description", creator.String())
	}

	dex.Description = description
	a.setDex(ctx, dex)

	return nil
}

// getDex get dex data, if no found, return false
func (a DexKeeper) getDex(ctx sdk.Context, creator types.Name) (*types.Dex, bool) {
	store := ctx.KVStore(a.key)
	bz := store.Get(types.DexStoreKey(creator))
	if bz == nil {
		return nil, false
	}

	res := &types.Dex{}
	if err := a.cdc.UnmarshalBinaryBare(bz, res); err != nil {
		panic(errors.Wrap(err, "get stat unmarshal"))
	}

	return res, true
}

// setDex set dex data
func (a DexKeeper) setDex(ctx sdk.Context, dex *types.Dex) {
	store := ctx.KVStore(a.key)
	bz, err := a.cdc.MarshalBinaryBare(*dex)
	if err != nil {
		panic(errors.Wrap(err, "marshal dex error"))
	}

	store.Set(types.DexStoreKey(dex.Creator), bz)
}

// deleteDex delete dex data
func (a DexKeeper) deleteDex(ctx sdk.Context, dex *types.Dex) {
	store := ctx.KVStore(a.key)
	store.Delete(types.DexStoreKey(dex.Creator))
}

// nextNumber next dex number
func (a DexKeeper) nextNumber(ctx sdk.Context) (n uint64) {
	var err error
	store := ctx.KVStore(a.key)
	bz := store.Get(types.GetDexNumberStoreKey())
	if nil != bz {
		if err = a.cdc.UnmarshalBinaryLengthPrefixed(bz, &n); nil != err {
			panic(err)
		}
	}
	bz = a.cdc.MustMarshalBinaryLengthPrefixed(n + 1)
	store.Set(types.GetDexNumberStoreKey(), bz)
	return
}
