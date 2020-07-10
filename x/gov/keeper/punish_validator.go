package keeper

import (
	chaintype "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func (keeper Keeper) SetPunishValidator(ctx sdk.Context, validator_to_punish types.PunishValidator) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryBare(&validator_to_punish)
	store.Set(types.GetValidatorKey(validator_to_punish.GetValidatorAccount()), bz)
}

func (keeper Keeper) GetPunishValidator(ctx sdk.Context, validatorAccount chaintype.AccountID) (punishValidator types.PunishValidator, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.GetValidatorKey(validatorAccount))
	if bz == nil {
		return punishValidator, false
	}
	keeper.cdc.MustUnmarshalBinaryBare(bz, &punishValidator)
	return punishValidator, true
}

func (keeper Keeper) deletePunishValidator(ctx sdk.Context, validatorAccount chaintype.AccountID) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.GetValidatorKey(validatorAccount))
}

func (keeper Keeper) IterateAllPunishValidators(ctx sdk.Context, cb func(validator types.PunishValidator) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
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
