package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrAccountHasCreated             = sdkerrors.Register(ModuleName, 1, "account has created")
	ErrAccountNoFound                = sdkerrors.Register(ModuleName, 2, "account no found")
	ErrAccountCannotCreateSysAccount = sdkerrors.Register(ModuleName, 3, "cannot create system account by create")
	ErrAccountNameInvalid            = sdkerrors.Register(ModuleName, 4, "account name is invalid")
	ErrAccountNameLenInvalid         = sdkerrors.Register(ModuleName, 5, "account name length is invalid")
)
