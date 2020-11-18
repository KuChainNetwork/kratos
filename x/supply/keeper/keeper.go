package keeper

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/x/supply/exported"
	"github.com/KuChainNetwork/kuchain/x/supply/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the supply store
type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey

	accountKeeper types.AccountKeeper
	bk            types.BankKeeper
	permAddrs     map[string]types.PermissionsForAddress
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	maccPerms map[string][]string,
) Keeper {

	permAddrs := make(map[string]types.PermissionsForAddress)
	for name, perms := range maccPerms {
		permAddrs[name] = types.NewPermissionsForAddress(name, perms)
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      key,
		accountKeeper: ak,
		bk:            bk,
		permAddrs:     permAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetSupply retrieves the Supply from store
func (k Keeper) GetSupply(ctx sdk.Context) exported.SupplyI {
	return types.NewSupply(k.bk.GetCoinsTotalSupply(ctx))
}

// ValidatePermissions validates that the module account has been granted
// permissions within its set of allowed permissions.
func (k Keeper) ValidatePermissions(macc exported.ModuleAccountI) error {
	permAddr := k.permAddrs[macc.GetName().String()]
	for _, perm := range macc.GetPermissions() {
		if !permAddr.HasPermission(perm) {
			return fmt.Errorf("invalid module permission %s", perm)
		}
	}

	return nil
}

func (k Keeper) GetStoreKey() sdk.StoreKey {
	return k.storeKey
}
