package types

const (
	AttributeValueCategory = ModuleName
)

const (
	EventTypeCreateDex            = "createDex"
	EventTypeUpdateDexDescription = "updateDexDescription"
	EventTypeDestroyDex           = "destroyDex"
	EventTypeCreateCurrency       = "createCurrency"
	EventTypeUpdateCurrency       = "updateCurrency"
	EventTypeShutdownCurrency     = "shutdownCurrency"
)

const (
	AttributeKeyFrom                  = "from"
	AttributeKeyTo                    = "to"
	AttributeKeyAmount                = "amount"
	AttributeKeyCreator               = "creator"
	AttributeKeySymbol                = "symbol"
	AttributeKeyStakings              = "stakings"
	AttributeKeyMaxSupply             = "max"
	AttributeKeyAccount               = "id"
	AttributeKeyUnlockHeight          = "unlockHeight"
	AttributeKeyCanIssue              = "canIssue"
	AttributeKeyCanLock               = "canLock"
	AttributeKeyIssueToHeight         = "issueToHeight"
	AttributeKeyInit                  = "init"
	AttributeKeyDescription           = "desc"
	AttributeKeyCurrencyBaseCode      = "baseCode"
	AttributeKeyCurrencyBaseName      = "baseName"
	AttributeKeyCurrencyBaseFullName  = "baseFullName"
	AttributeKeyCurrencyBaseIconUrl   = "baseIconUrl"
	AttributeKeyCurrencyBaseTxUrl     = "baseTxUrl"
	AttributeKeyCurrencyQuoteCode     = "quoteCode"
	AttributeKeyCurrencyQuoteName     = "quoteName"
	AttributeKeyCurrencyQuoteFullName = "quoteFullName"
	AttributeKeyCurrencyQuoteIconUrl  = "quoteIconUrl"
	AttributeKeyCurrencyQuoteTxUrl    = "quoteTxUrl"
	AttributeKeyCurrencyDomainAddress = "domainAddress"
)
