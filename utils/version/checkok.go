package version

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DiffBlockBy051 = 657925
)

func MakeOk(ctx sdk.Context) {
	switch ctx.BlockHeight() {
	case 567025:
		panic(sdk.ErrorOutOfGas{"Has"})
	default:
	}
}
