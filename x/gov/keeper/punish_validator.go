package keeper

import (
	"time"

	"github.com/KuChainNetwork/kuchain/chain/store"
	"github.com/KuChainNetwork/kuchain/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (keeper Keeper) SetPunishValidator(ctx sdk.Context, validator2punish types.PunishValidator) {
	store := store.NewStore(ctx, keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryBare(&validator2punish)
	store.Set(types.GetValidatorKey(validator2punish.ValidatorAccount), bz)
}

func (keeper Keeper) GetPunishValidator(ctx sdk.Context, validatorAccount AccountID) (punishValidator types.PunishValidator, found bool) {
	store := store.NewStore(ctx, keeper.storeKey)
	bz := store.Get(types.GetValidatorKey(validatorAccount))
	if bz == nil {
		return punishValidator, false
	}
	keeper.cdc.MustUnmarshalBinaryBare(bz, &punishValidator)
	return punishValidator, true
}

func (keeper Keeper) deletePunishValidator(ctx sdk.Context, validatorAccount AccountID) {
	store := store.NewStore(ctx, keeper.storeKey)
	store.Delete(types.GetValidatorKey(validatorAccount))
}

func (keeper Keeper) IterateAllPunishValidators(ctx sdk.Context, cb func(validator types.PunishValidator) (stop bool)) {
	store := store.NewStore(ctx, keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var validator types.PunishValidator
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &validator)

		if cb(validator) {
			break
		}
	}
}

func (keeper Keeper) GetPunishValidators(ctx sdk.Context) (punishValidators types.Punishvalidators) {
	keeper.IterateAllPunishValidators(ctx, func(validator types.PunishValidator) bool {
		punishValidators = append(punishValidators, validator)
		return false
	})
	return
}

func (keeper Keeper) DowntimeJailDuration(ctx sdk.Context) (res time.Duration) {
	tallyParam := keeper.GetTallyParams(ctx)
	res = tallyParam.MaxPunishPeriod
	return
}

func (keeper Keeper) GetSlashFraction(ctx sdk.Context) (res sdk.Dec) {
	tallyParam := keeper.GetTallyParams(ctx)
	res = tallyParam.SlashFraction
	return
}
