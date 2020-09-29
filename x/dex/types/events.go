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
	EventTypeDexDeal              = "dexDeal"
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
