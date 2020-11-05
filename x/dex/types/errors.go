package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrDexHadCreated        = sdkerrors.Register(ModuleName, 2, "dex has created")
	ErrDexDescTooLong       = sdkerrors.Register(ModuleName, 3, "dex description too long")
	ErrDexNotExists         = sdkerrors.Register(ModuleName, 4, "dex not exists")
	ErrDexCanNotBeDestroyed = sdkerrors.Register(ModuleName, 5, "dex can not be destroyed")
	ErrDexStakingsNotMatch  = sdkerrors.Register(ModuleName, 6, "dex stakings not match")
	ErrDexWasDestroy        = sdkerrors.Register(ModuleName, 7, "dex was destroy")
	ErrDexDescriptionSame   = sdkerrors.Register(ModuleName, 8, "dex description is same")
)

var (
	ErrSymbolExists               = sdkerrors.Register(ModuleName, 9, "dex symbol exists")
	ErrSymbolNotExists            = sdkerrors.Register(ModuleName, 10, "dex symbol not exists")
	ErrSymbolIncorrect            = sdkerrors.Register(ModuleName, 11, "dex symbol data incorrect")
	ErrSymbolBaseCodeEmpty        = sdkerrors.Register(ModuleName, 12, "dex symbol base code is empty")
	ErrSymbolQuoteCodeEmpty       = sdkerrors.Register(ModuleName, 13, "dex symbol quote code is empty")
	ErrSymbolBaseInvalid          = sdkerrors.Register(ModuleName, 14, "dex symbol base part invalid")
	ErrSymbolQuoteInvalid         = sdkerrors.Register(ModuleName, 15, "dex symbol quote part invalid")
	ErrSymbolDomainAddressInvalid = sdkerrors.Register(ModuleName, 16, "dex symbol domain address invalid")
	ErrSymbolUpdateFieldsInvalid  = sdkerrors.Register(ModuleName, 17, "dex symbol update fields invalid")
	ErrSymbolFormat               = sdkerrors.Register(ModuleName, 18, "dex symbol base_code/quote_code format error")
	ErrSymbolNotSupply            = sdkerrors.Register(ModuleName, 19, "dex symbol base_code/quote_code not supply")
	ErrSymbolDexDescriptionSame   = sdkerrors.Register(ModuleName, 20, "dex symbol to change is same as existing")
)

var (
	ErrDexSigInChangeToNegative = sdkerrors.Register(ModuleName, 21, "dex sigIn amt should cannot tobe changed to negative")
	ErrDexSigOutByUserNoUnlock  = sdkerrors.Register(ModuleName, 22, "dex sig out by user is locked in current")
	ErrDexSigInAmountNotEnough  = sdkerrors.Register(ModuleName, 23, "dex sig in amount not enough")
)
