package types

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
)

const (
	// module name
	ModuleName = "account"

	// StoreKey is string representation of the store key for auth
	StoreKey = "kuacc"

	// QuerierRoute is the querier route for acc
	QuerierRoute = ModuleName
)

var (
	// AddressStoreKeyPrefix prefix for account-by-address store
	AddressStoreKeyPrefix = []byte{0x0A}

	// AuthSeqStoreKeyPerfix seq num store prefix
	AuthSeqStoreKeyPerfix = []byte{0x0B}

	// AuthAccountsStoreKeyPerfix - Accounts store prefix
	AuthAccountsStoreKeyPerfix = []byte{0x0C}

	// GlobalAccountNumberKey param key for global account number
	GlobalAccountNumberKey = types.MustName("g.account.number").Value
)

// AccountIDStoreKey turn an address to key used to get it from the account store
func AccountIDStoreKey(addr types.AccountID) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}

// NameStoreKey turn an name to key used to get it from the account store
func NameStoreKey(name types.Name) []byte {
	return append(AddressStoreKeyPrefix, name.Bytes()...)
}

// AuthSeqStoreKey seq key for store
func AuthSeqStoreKey(addr types.AccAddress) []byte {
	return append(AuthSeqStoreKeyPerfix, addr.Bytes()...)
}

// Auth - Accounts key for store
func AuthAccountsStoreKey(auth types.AccAddress) []byte {
	return append(AuthAccountsStoreKeyPerfix, auth.Bytes()...)
}
