package constants

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// some fix height

var (
	// FixAssetHeight fix asset bugs height
	FixAssetHeight       string = ""
	FixAssetHeightVal, _        = strconv.ParseInt(FixAssetHeight, 10, 64)
)

// LogVersion log version info
func LogVersion(logger log.Logger) {
	logger.Info("FixAsset", "height", GetFixAssetHeight())
}

func GetFixAssetHeight() int64 {
	return FixAssetHeightVal
}

// IsFixAssetHeight is fix asset
func IsFixAssetHeight(ctx sdk.Context) bool {
	return ctx.BlockHeight() > FixAssetHeightVal
}
