package types

import (
	supplyexported "github.com/KuChainNetwork/kuchain/x/supply/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	InitModuleAccount(ctx sdk.Context, moduleName string) error
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) supplyexported.ModuleAccountI
}
