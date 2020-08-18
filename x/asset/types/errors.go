package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrAssetHasCreated                       = sdkerrors.Register(ModuleName, 1, "asset has created")
	ErrAssetDenom                            = sdkerrors.Register(ModuleName, 2, "denom format error")
	ErrAssetDescriptorTooLarge               = sdkerrors.Register(ModuleName, 3, "description for coin too large")
	ErrAssetIssueGTMaxSupply                 = sdkerrors.Register(ModuleName, 4, "asset issue cannot great then max supply")
	ErrAssetNoCreator                        = sdkerrors.Register(ModuleName, 5, "asset creator error")
	ErrAssetCoinNoExit                       = sdkerrors.Register(ModuleName, 6, "asset no exit")
	ErrAssetCoinNoEnough                     = sdkerrors.Register(ModuleName, 7, "account coins no enough")
	ErrAssetFromAccountEmpty                 = sdkerrors.Register(ModuleName, 8, "from account empty")
	ErrAssetToAccountEmpty                   = sdkerrors.Register(ModuleName, 9, "to account empty")
	ErrAssetLockCoinsNoEnough                = sdkerrors.Register(ModuleName, 10, "current coins need great or equal to lock coins")
	ErrAssetLockUnlockBlockHeightErr         = sdkerrors.Register(ModuleName, 11, "unlock block height error")
	ErrAssetUnLockCoins                      = sdkerrors.Register(ModuleName, 12, "unlock coins error")
	ErrAssetCoinsLocked                      = sdkerrors.Register(ModuleName, 13, "coins has locked")
	ErrAssetCoinCannotBeLock                 = sdkerrors.Register(ModuleName, 14, "coin state not allowe lock")
	ErrAssetCoinCannotBeIssue                = sdkerrors.Register(ModuleName, 15, "coin state not allowe issue")
	ErrAssetCoinCannotBeIssueInHeight        = sdkerrors.Register(ModuleName, 16, "coin state not allowe issue that in this height")
	ErrAssetCoinMustCanIssueWhenIssueByBlock = sdkerrors.Register(ModuleName, 17, "coin state must can issue when set to issue by height")
	ErrAssetCoinMustSupplyNeedGTInitSupply   = sdkerrors.Register(ModuleName, 18, "coin max_supply need > init_supply")
	ErrAssetIssueToHeightMustGTCurrentHeight = sdkerrors.Register(ModuleName, 19, "coin issue to height must > current height")
	ErrAssetSymbolError                      = sdkerrors.Register(ModuleName, 20, "asset symbol error")
	ErrAssetCoinNoZero                       = sdkerrors.Register(ModuleName, 21, "amount should not be zero")
	ErrAssetCoinCannotBeBurn                 = sdkerrors.Register(ModuleName, 22, "coin state not allowed burn")
)
