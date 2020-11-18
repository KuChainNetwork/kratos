package keeper

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/types/coin"
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/log"
)

// AssetCoinsKeeper keeper interface for asset module
type AssetCoinsKeeper interface {
	AssetViewKeeper
	AssetTransfer

	Create(ctx sdk.Context,
		creator, symbol types.Name, maxSupply types.Coin,
		canIssue, canLock, canBurn bool,
		issue2Height int64, initSupply types.Coin, desc []byte) error
	Issue(ctx sdk.Context, creator, symbol types.Name, amount types.Coin) error
	Burn(ctx sdk.Context, id types.AccountID, amt types.Coin) error
	LockCoins(ctx sdk.Context, account types.AccountID, unlockBlockHeight int64, coins types.Coins) error
	UnLockCoins(ctx sdk.Context, account types.AccountID, coins types.Coins) error
	ExerciseCoinPower(ctx sdk.Context, id types.AccountID, amt types.Coin) error
	Approve(ctx sdk.Context, id, spender types.AccountID, amt types.Coins, isLock bool) error
}

// AssetViewKeeper keeper view interface for asset module
type AssetViewKeeper interface {
	Cdc() *codec.Codec

	GetCoins(ctx sdk.Context, account types.AccountID) (types.Coins, error)
	GetCoin(ctx sdk.Context, account types.AccountID, creator, symbol types.Name) (types.Coin, error)
	GetCoinPowers(ctx sdk.Context, account types.AccountID) types.Coins
	GetCoinPower(ctx sdk.Context, account types.AccountID, creator, symbol types.Name) (types.Coin, error)
	GetCoinDesc(ctx sdk.Context, creator, symbol types.Name) (*types.CoinDescription, error)
	GetCoinStat(ctx sdk.Context, creator, symbol types.Name) (*types.CoinStat, error)
	GetLockCoins(ctx sdk.Context, account types.AccountID) (types.Coins, []LockedCoins, error)
	GetApproveCoins(ctx sdk.Context, account, spender types.AccountID) (*ApproveData, error)
	GetApproveSum(ctx sdk.Context, account types.AccountID) (types.Coins, error)
}

type AccountEnsurer interface {
	EnsureAccount(ctx sdk.Context, account types.AccountID) error
}

// AssetKeeper for asset state
type AssetKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	// AccountKeeper interface
	ak AccountEnsurer
}

var _ AssetCoinsKeeper = AssetKeeper{}

// NewAssetKeeper new asset keeper
func NewAssetKeeper(cdc *codec.Codec, key sdk.StoreKey, ak AccountEnsurer) AssetKeeper {
	return AssetKeeper{
		key: key,
		cdc: cdc,
		ak:  ak,
	}
}

// Cdc get cdc
func (a AssetKeeper) Cdc() *codec.Codec {
	return a.cdc
}

// Logger returns a module-specific logger.
func (a AssetKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (a AssetKeeper) Create(ctx sdk.Context,
	creator, symbol types.Name, maxSupply types.Coin,
	canIssue, canLock, canBurn bool,
	issue2Height int64, initSupply types.Coin, desc []byte) error {
	stat, _ := a.getStat(ctx, creator, symbol)
	if stat != nil {
		return types.ErrAssetHasCreated
	}

	// check denom
	denom := types.CoinDenom(creator, symbol)
	if err := coin.ValidateDenom(denom); err != nil {
		return sdkerrors.Wrapf(types.ErrAssetDenom, "denom %s", denom)
	}

	// init state
	newStat := types.NewCoinStat(ctx, creator, symbol, maxSupply)
	if err := newStat.SetOpt(canIssue, canLock, canBurn, issue2Height, initSupply); err != nil {
		return sdkerrors.Wrapf(err, "set stat opt")
	}
	if err := a.setStat(ctx, &newStat); err != nil {
		return sdkerrors.Wrap(err, "create token set stat")
	}

	if len(desc) > types.CoinDescriptionLen {
		return types.ErrAssetDescriptorTooLarge
	}

	// init desc
	newDesc := types.NewCoinDescription(creator, symbol, desc)
	if err := a.setDescription(ctx, &newDesc); err != nil {
		return sdkerrors.Wrap(err, "create token set desc")
	}

	return nil
}

func (a AssetKeeper) Issue(ctx sdk.Context, creator, symbol types.Name, amount types.Coin) error {
	if err := a.issueCoinStat(ctx, amount); err != nil {
		return err
	}

	creatorAccount := types.NewAccountIDFromName(creator)
	coins, err := a.getCoins(ctx, creatorAccount)
	if err != nil {
		return sdkerrors.Wrap(err, "get coins")
	}

	if err := a.setCoins(ctx, creatorAccount, coins.Add(amount)); err != nil {
		return sdkerrors.Wrap(err, "issue set coins")
	}

	return nil
}

func (a AssetKeeper) Burn(ctx sdk.Context, id types.AccountID, amount types.Coin) error {
	if err := a.burnCoinStat(ctx, amount); err != nil {
		return err
	}

	coins, err := a.getCoins(ctx, id)
	if err != nil {
		return sdkerrors.Wrap(err, "burn get coins")
	}

	newCoins, isNegative := coins.SafeSub(NewCoins(amount))
	if isNegative {
		return sdkerrors.Wrap(types.ErrAssetCoinNoEnough, "burn coins error")
	}

	if err := a.checkIsCanUseCoins(ctx, id, NewCoins(amount), coins, false); err != nil {
		return sdkerrors.Wrap(err, "burn")
	}

	if err := a.setCoins(ctx, id, newCoins); err != nil {
		return sdkerrors.Wrap(err, "burn set coins")
	}

	return nil
}
func (a AssetKeeper) Transfer(ctx sdk.Context, from, to types.AccountID, amount types.Coins) error {
	return a.TransferDetail(ctx, from, to, amount, false)
}

func (a AssetKeeper) TransferDetail(ctx sdk.Context, from, to types.AccountID, amount types.Coins, isApplyApprove bool) error {
	logger := a.Logger(ctx)

	logger.Debug("transfer coins", "from", from, "to", to, "amount", amount)

	if from.Empty() {
		return types.ErrAssetFromAccountEmpty
	}

	if to.Empty() {
		return types.ErrAssetToAccountEmpty
	}

	if amount.IsZero() {
		return nil
	}

	if err := a.ak.EnsureAccount(ctx, to); err != nil {
		return sdkerrors.Wrapf(err, "ensure account %s error", to)
	}

	fromCoins, err := a.getCoins(ctx, from)
	if err != nil {
		return sdkerrors.Wrap(err, "get from coins")
	}

	toCoins, err := a.getCoins(ctx, to)
	if err != nil {
		return sdkerrors.Wrap(err, "get to coins")
	}

	coinSubed, hasNeg := fromCoins.SafeSub(amount)
	if hasNeg {
		return sdkerrors.Wrap(types.ErrAssetCoinNoEnough, "transfer")
	}

	if err := a.checkIsCanUseCoins(ctx, from, amount, fromCoins, isApplyApprove); err != nil {
		return sdkerrors.Wrap(err, "transfer")
	}

	if err := a.setCoins(ctx, to, toCoins.Add(amount...)); err != nil {
		return sdkerrors.Wrap(err, "set to coins")
	}

	if err := a.setCoins(ctx, from, coinSubed); err != nil {
		return sdkerrors.Wrap(err, "set from coins")
	}

	return nil
}

func (a AssetKeeper) GenesisCoins(ctx sdk.Context, account types.AccountID, coins types.Coins) error {
	for _, coin := range coins {
		if err := a.issueCoinStat(ctx, coin); err != nil {
			return err
		}
	}
	return a.setCoins(ctx, account, coins)
}

func (a AssetKeeper) GetStoreKey() sdk.StoreKey {
	return a.key
}
