package types

import (
	"encoding/binary"

	"github.com/KuChainNetwork/kuchain/chain/types"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "kudistribution"

	// StoreKey is the store key string for distribution
	StoreKey = ModuleName

	// RouterKey is the message route for distribution
	RouterKey = ModuleName

	// QuerierRoute is the querier route for distribution
	QuerierRoute = ModuleName
)

var (
	ModuleAccountName = MustName(ModuleName)
	ModuleAccountID   = types.NewAccountIDFromName(ModuleAccountName)
)

// Keys for distribution store
// Items are stored with the following key: values
//
// - 0x00<proposalID_Bytes>: FeePol
//
// - 0x01: sdk.ConsAddress
//
// - 0x02<valAddr_Bytes>: ValidatorOutstandingRewards
//
// - 0x03<accAddr_Bytes>: sdk.AccAddress
//
// - 0x04<valAddr_Bytes><accAddr_Bytes>: DelegatorStartingInfo
//
// - 0x05<valAddr_Bytes><period_Bytes>: ValidatorHistoricalRewards
//
// - 0x06<valAddr_Bytes>: ValidatorCurrentRewards
//
// - 0x07<valAddr_Bytes>: ValidatorCurrentRewards
//
// - 0x08<valAddr_Bytes><height>: ValidatorSlashEvent
var (
	FeePoolKey                        = []byte{0x00} // key for global distribution state
	ProposerKey                       = []byte{0x01} // key for the proposer operator address
	ValidatorOutstandingRewardsPrefix = []byte{0x02} // key for outstanding rewards

	DelegatorWithdrawAddrPrefix          = []byte{0x03} // key for delegator withdraw address
	DelegatorStartingInfoPrefix          = []byte{0x04} // key for delegator starting info
	ValidatorHistoricalRewardsPrefix     = []byte{0x05} // key for historical validators rewards / stake
	ValidatorCurrentRewardsPrefix        = []byte{0x06} // key for current validator rewards
	ValidatorAccumulatedCommissionPrefix = []byte{0x07} // key for accumulated validator commission
	ValidatorSlashEventPrefix            = []byte{0x08} // key for validator slash fraction
)

// gets an address from a validator's outstanding rewards key
func GetValidatorOutstandingRewardsAddress(key []byte) AccountID {
	return NewAccountIDFromStoreKey(key)
}

// gets an address from a delegator's withdraw info key
func GetDelegatorWithdrawInfoAddress(key []byte) AccountID {
	return NewAccountIDFromStoreKey(key)
}

// gets an address from a delegator's withdraw info key, by cancer
func GetDelegatorWithdrawInfoAddressUseAccountId(key []byte) AccountID {
	return NewAccountIDFromStoreKey(key)
}

// gets the addresses from a delegator starting info key
func GetDelegatorStartingInfoAddresses(key []byte) (valAddr AccountID, delAddr AccountID) {
	addr := key[:1+AccIDStoreKeyLen]
	valAddr = NewAccountIDFromStoreKey(addr)
	addr = key[1+AccIDStoreKeyLen:]
	delAddr = NewAccountIDFromStoreKey(addr)
	return
}

// gets the address & period from a validator's historical rewards key
func GetValidatorHistoricalRewardsAddressPeriod(key []byte) (valAddr AccountID, period uint64) {
	addr := key[:1+AccIDStoreKeyLen]
	valAddr = NewAccountIDFromStoreKey(addr)
	b := key[1+AccIDStoreKeyLen:]
	if len(b) != 8 {
		panic("unexpected key length")
	}
	period = binary.LittleEndian.Uint64(b)
	return
}

// gets the address from a validator's current rewards key
func GetValidatorCurrentRewardsAddress(key []byte) AccountID {
	return NewAccountIDFromStoreKey(key)
}

// gets the address from a validator's accumulated commission key
func GetValidatorAccumulatedCommissionAddress(key []byte) AccountID {
	return NewAccountIDFromStoreKey(key)
}

// gets the height from a validator's slash event key
func GetValidatorSlashEventAddressHeight(key []byte) (valAddr AccountID, height uint64) {
	addr := key[:1+AccIDStoreKeyLen]
	valAddr = NewAccountIDFromStoreKey(addr)
	startB := 1 + AccIDStoreKeyLen
	b := key[startB : startB+8] // the next 8 bytes represent the height
	height = binary.BigEndian.Uint64(b)
	return
}

// gets the outstanding rewards key for a validator
func GetValidatorOutstandingRewardsKey(val AccountID) []byte {
	return append(ValidatorOutstandingRewardsPrefix, val.StoreKey()...)
}

// gets the key for a delegator's withdraw addr
func GetDelegatorWithdrawAddrKey(del AccountID) []byte {
	return append(DelegatorWithdrawAddrPrefix, del.StoreKey()...)
}

// gets the key for a delegator's starting info
func GetDelegatorStartingInfoKey(v AccountID, d AccountID) []byte {
	return append(append(DelegatorStartingInfoPrefix, v.StoreKey()...), d.StoreKey()...)
}

// gets the prefix key for a validator's historical rewards
func GetValidatorHistoricalRewardsPrefix(v AccountID) []byte {
	return append(ValidatorHistoricalRewardsPrefix, v.StoreKey()...)
}

// gets the key for a validator's historical rewards
func GetValidatorHistoricalRewardsKey(v AccountID, k uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, k)
	return append(append(ValidatorHistoricalRewardsPrefix, v.StoreKey()...), b...)
}

// gets the key for a validator's current rewards
func GetValidatorCurrentRewardsKey(v AccountID) []byte {
	return append(ValidatorCurrentRewardsPrefix, v.StoreKey()...)
}

// gets the key for a validator's current commission
func GetValidatorAccumulatedCommissionKey(v AccountID) []byte {
	return append(ValidatorAccumulatedCommissionPrefix, v.StoreKey()...)
}

// gets the prefix key for a validator's slash fractions
func GetValidatorSlashEventPrefix(v AccountID) []byte {
	return append(ValidatorSlashEventPrefix, v.StoreKey()...)
}

// gets the prefix key for a validator's slash fraction (ValidatorSlashEventPrefix + height)
func GetValidatorSlashEventKeyPrefix(v AccountID, height uint64) []byte {
	heightBz := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBz, height)
	return append(
		ValidatorSlashEventPrefix,
		append(v.StoreKey(), heightBz...)...,
	)
}

// gets the key for a validator's slash fraction
func GetValidatorSlashEventKey(v AccountID, height, period uint64) []byte {
	periodBz := make([]byte, 8)
	binary.BigEndian.PutUint64(periodBz, period)
	prefix := GetValidatorSlashEventKeyPrefix(v, height)
	return append(prefix, periodBz...)
}
