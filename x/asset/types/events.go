package types

const (
	AttributeValueCategory = ModuleName
)

const (
	EventTypeCreate   = "create"
	EventTypeIssue    = "issue"
	EventTypeTransfer = "transfer"
	EventTypeLock     = "lock"
	EventTypeUnlock   = "unlock"
	EventTypeExercise = "exercise"
	EventTypeApprove  = "approve"
)

const (
	AttributeKeyFrom          = "from"
	AttributeKeyTo            = "to"
	AttributeKeySpender       = "spender"
	AttributeKeyAmount        = "amount"
	AttributeKeyCreator       = "creator"
	AttributeKeySymbol        = "symbol"
	AttributeKeyMaxSupply     = "max"
	AttributeKeyAccount       = "id"
	AttributeKeyUnlockHeight  = "unlockHeight"
	AttributeKeyCanIssue      = "canIssue"
	AttributeKeyCanLock       = "canLock"
	AttributeKeyIssueToHeight = "issueToHeight"
	AttributeKeyInit          = "init"
	AttributeKeyDescription   = "desc"
)
