package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/supply/exported"
	"github.com/KuChainNetwork/kuchain/x/supply/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetModuleAddress returns an address based on the module name
func (k Keeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	permAddr, ok := k.permAddrs[moduleName]
	if !ok {
		return nil
	}
	return permAddr.GetAddress()
}

// GetModuleAddressAndPermissions returns an address and permissions based on the module name
func (k Keeper) GetModuleAddressAndPermissions(moduleName string) (addr sdk.AccAddress, permissions []string) {
	permAddr, ok := k.permAddrs[moduleName]
	if !ok {
		return addr, permissions
	}
	return permAddr.GetAddress(), permAddr.GetPermissions()
}

// GetModuleAccountAndPermissions gets the module account from the auth account store and its
// registered permissions
func (k Keeper) GetModuleAccountAndPermissions(ctx sdk.Context, moduleName string) (exported.ModuleAccountI, []string) {
	addr, perms := k.GetModuleAddressAndPermissions(moduleName)
	if addr == nil {
		return nil, []string{}
	}

	name := types.MustName(moduleName)

	acc := k.accountKeeper.GetAccount(ctx, types.NewAccountIDFromName(name))
	if acc != nil {
		macc, ok := acc.(exported.ModuleAccountI)
		if !ok {
			panic("account is not a module account")
		}
		return macc, perms
	}

	// create a new module account
	macc := types.NewEmptyModuleAccount(moduleName, perms...)
	maccI := (k.accountKeeper.NewAccount(ctx, macc)).(exported.ModuleAccountI) // set the account number
	k.SetModuleAccount(ctx, maccI)

	return maccI, perms
}

// GetModuleAccount gets the module account from the auth account store, if the account does not
// exist in the AccountKeeper, then it is created.
func (k Keeper) GetModuleAccount(ctx sdk.Context, moduleName string) exported.ModuleAccountI {
	acc, _ := k.GetModuleAccountAndPermissions(ctx, moduleName)
	return acc
}

// SetModuleAccount sets the module account to the auth account store
func (k Keeper) SetModuleAccount(ctx sdk.Context, macc exported.ModuleAccountI) { //nolint:interfacer
	k.accountKeeper.SetAccount(ctx, macc)
}

// InitModuleAccount init module account
func (k Keeper) InitModuleAccount(ctx sdk.Context, moduleName string) error {
	_, perms := k.GetModuleAddressAndPermissions(moduleName)

	name := types.MustName(moduleName)
	acc := k.accountKeeper.GetAccount(ctx, types.NewAccountIDFromName(name))
	if acc != nil {
		_, ok := acc.(exported.ModuleAccountI)
		if !ok {
			panic("account is not a module account")
		}
		return nil
	}

	// create a new module account
	macc := types.NewEmptyModuleAccount(moduleName, perms...)
	maccI := (k.accountKeeper.NewAccount(ctx, macc)).(exported.ModuleAccountI) // set the account number
	k.SetModuleAccount(ctx, maccI)

	return nil
}
