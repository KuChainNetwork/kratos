package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrDexHadCreated        = sdkerrors.Register(ModuleName, 1, "dex has created")
	ErrDexDescTooLong       = sdkerrors.Register(ModuleName, 2, "dex description too long")
	ErrDexNotExists         = sdkerrors.Register(ModuleName, 3, "dex not exists")
	ErrDexCanNotBeDestroyed = sdkerrors.Register(ModuleName, 4, "dex can not be destroyed")
)
