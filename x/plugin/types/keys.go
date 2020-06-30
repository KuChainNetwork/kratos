package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	// ModuleName is the name of the module
	ModuleName = "plugin"

	// StoreKey is the store key string for slashing
	// StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute is the querier route for slashing
	QuerierRoute = ModuleName
)

var (
	loggerName = fmt.Sprintf("x/%s", ModuleName)
)

func Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", loggerName)
}
