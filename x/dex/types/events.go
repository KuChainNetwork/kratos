package types

const (
	AttributeValueCategory = ModuleName
)

const (
	EventTypeCreateDex            = "createDex"
	EventTypeUpdateDexDescription = "updateDexDescription"
	EventTypeDestroyDex           = "destroyDex"
	EventTypeCreateSymbol         = "createSymbol"
	EventTypeUpdateSymbol         = "updateSymbol"
	EventTypePauseSymbol          = "pauseSymbol"
	EventTypeRestoreSymbol        = "restoreSymbol"
	EventTypeShutdownSymbol       = "shutdownSymbol"
)

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
)
