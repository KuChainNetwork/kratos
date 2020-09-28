package types

const (
	AttributeValueCategory = ModuleName
)

const (
	EventTypeCreateDex            = "createDex"
	EventTypeUpdateDexDescription = "updateDexDescription"
	EventTypeDestroyDex           = "destroyDex"
	EventTypeDexSigIn             = "dexSigIn"
	EventTypeDexSigOut            = "dexSigOut"
)

// TODO: use one in all modules

const (
	AttributeKeyFrom          = "from"
	AttributeKeyTo            = "to"
	AttributeKeyAmount        = "amount"
	AttributeKeyCreator       = "creator"
	AttributeKeySymbol        = "symbol"
	AttributeKeyStakings      = "stakings"
	AttributeKeyMaxSupply     = "max"
	AttributeKeyAccount       = "id"
	AttributeKeyUnlockHeight  = "unlockHeight"
	AttributeKeyCanIssue      = "canIssue"
	AttributeKeyCanLock       = "canLock"
	AttributeKeyIssueToHeight = "issueToHeight"
	AttributeKeyInit          = "init"
	AttributeKeyDescription   = "desc"
	AttributeKeyUser          = "user"
	AttributeKeyDex           = "dex"
	AttributeKeyIsTimeout     = "isTimeout"
)
