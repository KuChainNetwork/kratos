package coin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Int = sdk.Int
	Dec = sdk.Dec
)

var (
	NewInt           = sdk.NewInt
	ZeroInt          = sdk.ZeroInt
	NewIntFromString = sdk.NewIntFromString

	MinDec        = sdk.MinDec
	ZeroDec       = sdk.ZeroDec
	NewDecFromStr = sdk.NewDecFromStr
)
