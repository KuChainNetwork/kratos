package plugin

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) {
	logger := Logger(ctx)

	logger.Info("begin block", "req", req)
}

func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	logger := Logger(ctx)

	logger.Info("end block", "req", req)

	return []abci.ValidatorUpdate{}
}
