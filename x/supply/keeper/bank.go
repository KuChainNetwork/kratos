package keeper

import (
	"fmt"

	stakingTypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/KuChainNetwork/kuchain/x/supply/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress.
// It will panic if the module account does not exist.
func (k Keeper) SendCoinsFromModuleToAccount(
	ctx sdk.Context, senderModule string, recipient types.AccountID, amt Coins) error {
	senderAcc := k.GetModuleAccount(ctx, senderModule)
	if senderAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress,
			"module account %s does not exist", senderModule))
	}

	if err := k.bk.SendCoinPower(ctx, senderAcc.GetID(), recipient, amt); err != nil {
		return sdkerrors.Wrapf(err,
			"SendCoinsFromModuleToAccount %s to %s by %s",
			senderAcc.String(), recipient.String(), amt.String())
	}

	return nil
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another.
// It will panic if either module account does not exist.
func (k Keeper) SendCoinsFromModuleToModule(
	ctx sdk.Context, senderModule, recipientModule string, amt Coins,
) error {
	senderAcc := k.GetModuleAccount(ctx, senderModule)
	if senderAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	if err := k.bk.SendCoinPower(ctx, senderAcc.GetID(), recipientAcc.GetID(), amt); err != nil {
		return sdkerrors.Wrapf(err,
			"SendCoinsFromModuleToModule %s to %s by %s", senderAcc.String(), recipientAcc.String(), amt.String())
	}

	return nil
}

// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount.
// It will panic if the module account does not exist.
func (k Keeper) SendCoinsFromAccountToModule(
	ctx sdk.Context, sender types.AccountID, recipientModule string, amt Coins,
) error {
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	ctx.Logger().Debug("send coins power to module", "from", sender, "to", recipientAcc.GetID().String(), "amount", amt)

	if err := k.bk.SendCoinPower(ctx, sender, recipientAcc.GetID(), amt); err != nil {
		return sdkerrors.Wrapf(err,
			"SendCoinsFromAccountToModule %s to %s by %s", sender.String(), recipientModule, amt.String())
	}

	return nil
}

// DelegateCoinsFromAccountToModule delegates coins and transfers them from a
// delegator account to a module account. It will panic if the module account
// does not exist or is unauthorized.
func (k Keeper) DelegateCoinsFromAccountToModule(
	ctx sdk.Context, recipientModule string, amt Coins,
) error {
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress,
			"module account %s does not exist", recipientModule))
	}

	if !recipientAcc.HasPermission(types.Staking) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized,
			"module account %s does not have permissions to receive delegated coins", recipientModule))
	}

	// Delegate will first send coins to ModuleAccountID
	if err := k.bk.CoinsToPower(ctx, stakingTypes.ModuleAccountID, recipientAcc.GetID(), amt); err != nil {
		return sdkerrors.Wrapf(err,
			"DelegateCoinsFromAccountToModule %s by %s", recipientModule, amt.String())
	}

	return nil
}

// UndelegateCoinsFromModuleToAccount undelegates the unbonding coins and transfers
// them from a module account to the delegator account. It will panic if the
// module account does not exist or is unauthorized.
func (k Keeper) UndelegateCoinsFromModuleToAccount(
	ctx sdk.Context, senderModule string, recipientAcc types.AccountID, amt Coins,
) error {
	acc := k.GetModuleAccount(ctx, senderModule)
	if acc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	if !acc.HasPermission(types.Staking) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to undelegate coins", senderModule))
	}

	// Delegate will first send coins to ModuleAccountID
	if err := k.bk.SendCoinPower(ctx, acc.GetID(), recipientAcc, amt); err != nil {
		return sdkerrors.Wrapf(err,
			"UndelegateCoinsFromModuleToAccount %s by %s", recipientAcc, amt.String())
	}

	return nil
}

// MintCoins creates new coins from thin air and adds it to the module account.
// It will panic if the module account does not exist or is unauthorized.
func (k Keeper) MintCoins(ctx sdk.Context, moduleName string, amt *Coins) error {
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress,
			"module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(types.Minter) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized,
			"module account %s does not have permissions to mint tokens", moduleName))
	}

	_, err := k.bk.IssueCoinPower(ctx, acc.GetID(), *amt)
	if err != nil {
		return err
	}

	k.Logger(ctx).Info(fmt.Sprintf("minted %s from %s module account", amt.String(), moduleName))

	return nil
}

// BurnCoins burns coins deletes coins from the balance of the module account.
// It will panic if the module account does not exist or is unauthorized.
func (k Keeper) BurnCoins(ctx sdk.Context, moduleName types.AccountID, amt Coins) error {
	acc := k.GetModuleAccount(ctx, moduleName.String())
	if acc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress,
			"module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(types.Burner) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized,
			"module account %s does not have permissions to burn tokens", moduleName))
	}

	err := k.SendCoinsFromModuleToModule(ctx, moduleName.String(), types.BlackHole, amt)
	if err != nil {
		return sdkerrors.Wrapf(err, "burn coins error by sub coin power")
	}

	k.Logger(ctx).Info(fmt.Sprintf("burned %s from %s module account", amt.String(), moduleName))

	return nil
}

func (k Keeper) ModuleCoinsToPower(ctx sdk.Context, recipientModule string, amt Coins) error {
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	if err := k.bk.CoinsToPower(ctx, recipientAcc.GetID(), recipientAcc.GetID(), amt); err != nil {
		return sdkerrors.Wrapf(err,
			"ModuleCoinsToPower error %s by %s", recipientModule, amt.String())
	}

	return nil
}
