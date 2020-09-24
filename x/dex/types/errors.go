package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrDexHadCreated        = sdkerrors.Register(ModuleName, 1, "dex has created")
	ErrDexDescTooLong       = sdkerrors.Register(ModuleName, 2, "dex description too long")
	ErrDexNotExists         = sdkerrors.Register(ModuleName, 3, "dex not exists")
	ErrDexCanNotBeDestroyed = sdkerrors.Register(ModuleName, 4, "dex can not be destroyed")
	ErrDexStakingsNotMatch  = sdkerrors.Register(ModuleName, 5, "dex stakings not match")
	ErrDexWasDestroy        = sdkerrors.Register(ModuleName, 6, "dex was destroy")
	ErrCurrencyExists       = sdkerrors.Register(ModuleName, 7, "dex currency exists")
	ErrCurrencyNotExists    = sdkerrors.Register(ModuleName, 8, "dex currency not exists")
	ErrCurrencyIncorrect    = sdkerrors.Register(ModuleName, 9, "dex currency data incorrect")
)
