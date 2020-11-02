package keeper

import (
	"fmt"
	"time"

	"github.com/KuChainNetwork/kuchain/chain/store"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	params "github.com/KuChainNetwork/kuchain/x/params/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/log"
)

const key = "NotDistributionTimePoint"

// Keeper of the distribution store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	paramSpace       params.Subspace
	BankKeeper       types.BankKeeperAccountID
	stakingKeeper    types.StakingKeeperAccountID
	supplyKeeper     types.SupplyKeeperAccountID
	AccKeeper        account.Keeper
	blacklistedAddrs map[string]bool

	feeCollectorName string // name of the FeeCollector ModuleAccount

	startNotDistriTimePoint time.Time
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey,
	paramSpace params.Subspace,
	bk types.BankKeeperAccountID,
	sk types.StakingKeeperAccountID,
	supplyKeeper types.SupplyKeeperAccountID,
	accKeeper account.Keeper,
	feeCollectorName string, blacklistedAddrs map[string]bool,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:                key,
		cdc:                     cdc,
		paramSpace:              paramSpace,
		BankKeeper:              bk,
		stakingKeeper:           sk,
		supplyKeeper:            supplyKeeper,
		AccKeeper:               accKeeper,
		feeCollectorName:        feeCollectorName,
		blacklistedAddrs:        blacklistedAddrs,
		startNotDistriTimePoint: time.Time{},
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetWithdrawAddr sets a new address that will receive the rewards upon withdrawal
func (k Keeper) SetWithdrawAddr(ctx sdk.Context, delegatorId chainTypes.AccountID, withdrawId chainTypes.AccountID) error {
	if k.blacklistedAddrs[withdrawId.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is blacklisted from receiving external funds", withdrawId)
	}

	if !k.GetWithdrawAddrEnabled(ctx) {
		return types.ErrSetWithdrawAddrDisabled
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetWithdrawAddress,
			sdk.NewAttribute(types.AttributeKeyWithdrawAddress, withdrawId.String()),
		),
	)

	k.SetDelegatorWithdrawAddr(ctx, delegatorId, withdrawId)
	return nil
}

// withdraw rewards from a delegation
func (k Keeper) WithdrawDelegationRewards(ctx sdk.Context, delAddr chainTypes.AccountID, valAddr chainTypes.AccountID) (Coins, error) {

	val := k.stakingKeeper.Validator(ctx, valAddr)
	ctx.Logger().Debug("WithdrawDelegationRewards", "val:", val)
	if val == nil {
		return nil, types.ErrNoValidatorDistInfo
	}

	del := k.stakingKeeper.Delegation(ctx, delAddr, valAddr)
	ctx.Logger().Debug("WithdrawDelegationRewards", "del:", del)
	if del == nil {
		return nil, types.ErrEmptyDelegationDistInfo
	}

	// withdraw rewards
	rewards, err := k.withdrawDelegationRewards(ctx, val, del)
	if err != nil {
		ctx.Logger().Debug("WithdrawDelegationRewards", "err:", err)
		return nil, err
	}
	ctx.Logger().Debug("WithdrawDelegationRewards", "rewards:", rewards)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, rewards.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
		),
	)

	// reinitialize the delegation
	k.initializeDelegation(ctx, valAddr, delAddr)
	return rewards, nil
}

// withdraw validator commission
func (k Keeper) WithdrawValidatorCommission(ctx sdk.Context, valAddr chainTypes.AccountID) (Coins, error) {
	// fetch validator accumulated commission
	accumCommission := k.GetValidatorAccumulatedCommission(ctx, valAddr)
	if accumCommission.Commission.IsZero() {
		return nil, types.ErrNoValidatorCommission
	}

	commission, remainder := accumCommission.Commission.TruncateDecimal()

	k.SetValidatorAccumulatedCommission(ctx, valAddr, types.ValidatorAccumulatedCommission{Commission: remainder}) // leave remainder to withdraw later

	// update outstanding
	outstanding := k.GetValidatorOutstandingRewards(ctx, valAddr).Rewards
	k.SetValidatorOutstandingRewards(ctx, valAddr, types.ValidatorOutstandingRewards{Rewards: outstanding.Sub(chainTypes.NewDecCoinsFromCoins(commission...))})

	if !commission.IsZero() {
		//accAddr := sdk.AccAddress(valAddr)
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, valAddr)
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, commission)
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
		),
	)

	return commission, nil
}

// GetTotalRewards returns the total amount of fee distribution rewards held in the store
func (k Keeper) GetTotalRewards(ctx sdk.Context) (totalRewards chainTypes.DecCoins) {
	k.IterateValidatorOutstandingRewards(ctx,
		func(_ AccountID, rewards types.ValidatorOutstandingRewards) (stop bool) {
			totalRewards = totalRewards.Add(rewards.Rewards...)
			return false
		},
	)

	return totalRewards
}

// FundCommunityPool allows an account to directly fund the community fund pool.
// The amount is first added to the distribution module account and then directly
// added to the pool. An error is returned if the amount cannot be sent to the
// module account.
func (k Keeper) FundCommunityPool(ctx sdk.Context, amount Coins, sender chainTypes.AccountID) error {
	ctx.Logger().Debug("FundCommunityPool", "amount", amount, "sender", sender)

	// module name to coin power
	if err := k.BankKeeper.CoinsToPower(ctx, sender, types.ModuleAccountID, amount); err != nil {

		return sdkerrors.Wrap(err, "FundCommunityPool to power")
	}

	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(chainTypes.NewDecCoinsFromCoins(amount...)...)
	k.SetFeePool(ctx, feePool)

	return nil
}

func (k Keeper) SetStartNotDistributionTimePoint(ctx sdk.Context, t time.Time) {
	k.startNotDistriTimePoint = t

	store := store.NewStore(ctx, k.storeKey)
	bz := k.cdc.MustMarshalJSON(k.startNotDistriTimePoint)

	store.Set([]byte(key), bz)
	ctx.Logger().Debug("SetStartNotDistributionTimePoint",
		"time", k.startNotDistriTimePoint.Nanosecond())
}

func (k *Keeper) GetStartNotDistributionTimePoint(ctx sdk.Context) {
	store := store.NewStore(ctx, k.storeKey)
	bz := store.Get([]byte(key))

	k.cdc.UnmarshalJSON(bz, &k.startNotDistriTimePoint)
	ctx.Logger().Debug("GetStartNotDistributionTimePoint",
		"time", k.startNotDistriTimePoint.Nanosecond())
}

func (k Keeper) CanDistribution(ctx sdk.Context) (bool, time.Time) {
	if k.startNotDistriTimePoint.Nanosecond() <= 0 {
		return true, k.startNotDistriTimePoint
	}

	nt := ctx.BlockHeader().Time
	tEnd := k.startNotDistriTimePoint.Add(24 * 3600 * 1e9)
	if nt.Before(tEnd) && nt.After(k.startNotDistriTimePoint) {
		return false, k.startNotDistriTimePoint
	} else {
		k.SetStartNotDistributionTimePoint(ctx, time.Time{})
		ctx.Logger().Info("time CanDistribution",
			"time", k.startNotDistriTimePoint.Nanosecond())
	}

	return true, k.startNotDistriTimePoint
}

func (k Keeper) GetStoreKey() sdk.StoreKey {
	return k.storeKey
}

func (k Keeper) GetSupplyKeeper() types.SupplyKeeperAccountID {
	return k.supplyKeeper
}

func (k Keeper) GetFeeCollectorName() string {
	return k.feeCollectorName
}
