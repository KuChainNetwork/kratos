package types

const (
	AttributeValueCategory = ModuleName
)

const (
	EventTypeCreateDex            = "createDex"
	EventTypeUpdateDexDescription = "updateDexDescription"
	EventTypeDestroyDex           = "destroyDex"

	EventTypeCreateSymbol   = "createSymbol"
	EventTypeUpdateSymbol   = "updateSymbol"
	EventTypePauseSymbol    = "pauseSymbol"
	EventTypeRestoreSymbol  = "restoreSymbol"
	EventTypeShutdownSymbol = "shutdownSymbol"

	EventTypeDexSigIn  = "dexSigIn"
	EventTypeDexSigOut = "dexSigOut"
	EventTypeDexDeal   = "dexDeal"
)

// TODO: use one in all modules

const (
	AttributeKeyFrom                = "from"
	AttributeKeyTo                  = "to"
	AttributeKeyAmount              = "amount"
	AttributeKeyCreator             = "creator"
	AttributeKeySymbol              = "symbol"
	AttributeKeyStakings            = "stakings"
	AttributeKeyMaxSupply           = "max"
	AttributeKeyAccount             = "id"
	AttributeKeyUnlockHeight        = "unlockHeight"
	AttributeKeyCanIssue            = "canIssue"
	AttributeKeyCanLock             = "canLock"
	AttributeKeyIssueToHeight       = "issueToHeight"
	AttributeKeyInit                = "init"
	AttributeKeyDescription         = "desc"
	AttributeKeySymbolCreateHeight  = "symbolCreateHeight"
	AttributeKeySymbolBaseCode      = "baseCode"
	AttributeKeySymbolBaseName      = "baseName"
	AttributeKeySymbolBaseFullName  = "baseFullName"
	AttributeKeySymbolBaseIconUrl   = "baseIconUrl"
	AttributeKeySymbolBaseTxUrl     = "baseTxUrl"
	AttributeKeySymbolQuoteCode     = "quoteCode"
	AttributeKeySymbolQuoteName     = "quoteName"
	AttributeKeySymbolQuoteFullName = "quoteFullName"
	AttributeKeySymbolQuoteIconUrl  = "quoteIconUrl"
	AttributeKeySymbolQuoteTxUrl    = "quoteTxUrl"
	AttributeKeySymbolDomainAddress = "domainAddress"
	AttributeKeyUser                = "user"
)

const (
	AttributeKeyDex       = "dex"
	AttributeKeyIsTimeout = "isTimeout"
)

const (
	AttributeKeyDealRole1  = "role1"
	AttributeKeyDealFee1   = "fee1"
	AttributeKeyDealToken1 = "token1"
	AttributeKeyDealRole2  = "role2"
	AttributeKeyDealFee2   = "fee2"
	AttributeKeyDealToken2 = "token2"
)
