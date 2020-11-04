package types

import (
	"bytes"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/pkg/errors"
)

const (
	// module name
	ModuleName = "asset"

	// StoreKey is string representation of the store key for auth
	StoreKey = "asset"

	// QuerierRoute is the querier route for acc
	QuerierRoute = StoreKey
)

const (
	MaxDescriptionLength = 256
)

var (
	// AssetModuleKeyPrefix prefix for asset store
	// TODO: all Store key prefix should be more to a go package
	AssetModuleKeyPrefix = []byte{0x02}

	CoinStoreKeyPrefix           = chainTypes.MustName("coin").Bytes()
	CoinPowerStoreKeyPrefix      = chainTypes.MustName("coin.power").Bytes()
	CoinLockedStoreKeyPrefix     = chainTypes.MustName("coin.lock").Bytes()
	CoinLockedStatStoreKeyPrefix = chainTypes.MustName("coin.locks").Bytes()
	CoinStatStoreKeyPrefix       = chainTypes.MustName("coin.stat").Bytes()
	CoinDescStoreKeyPrefix       = chainTypes.MustName("coin.desc").Bytes()
	CoinApproveStoreKeyPrefix    = chainTypes.MustName("coin.approve").Bytes()
	CoinApproveSumStoreKeyPrefix = chainTypes.MustName("coin.approvesum").Bytes()

	coinStoreKeyPreLen = len(AssetModuleKeyPrefix)
)

func genCoinStoreKey(pre []byte, keys ...[]byte) []byte {
	res := make([]byte, 0, 256)
	res = append(res, AssetModuleKeyPrefix...)
	res = append(res, pre...)
	for _, k := range keys {
		res = append(res, k...)
	}
	return res
}

// GetKeyPrefix get key prefix for coins
func GetKeyPrefix(pre []byte) []byte {
	res := make([]byte, 0, 256)
	res = append(res, AssetModuleKeyPrefix...)
	res = append(res, pre...)
	return res
}

func coinStoreKey2AccountID(pre []byte, key []byte) AccountID {
	startIdx := coinStoreKeyPreLen + len(pre)
	keyBytes := key[startIdx:]
	if len(keyBytes) > chainTypes.AccIDStoreKeyLen {
		panic(errors.Errorf("coinStoreKey2AccountID key too len %d", len(keyBytes)))
	}

	if !bytes.Equal(key[coinStoreKeyPreLen:coinStoreKeyPreLen+len(pre)], pre) {
		panic(errors.Errorf("coinStoreKey2AccountID pre not equal"))
	}

	return chainTypes.NewAccountIDFromByte(key[startIdx:])
}

// CoinStoreKey get the key of coin store keeper for asset
func CoinStoreKey(account chainTypes.AccountID) []byte {
	return genCoinStoreKey(CoinStoreKeyPrefix, account.StoreKey())
}

// AccountIDFromCoinStoreKey get accountID from key
func AccountIDFromCoinStoreKey(key []byte) chainTypes.AccountID {
	return coinStoreKey2AccountID(CoinStoreKeyPrefix, key)
}

// CoinPowerStoreKey get the key of coin store keeper for asset
func CoinPowerStoreKey(account chainTypes.AccountID) []byte {
	return genCoinStoreKey(CoinPowerStoreKeyPrefix, account.StoreKey())
}

// AccountIDFromCoinPowerStoreKey get accountID from key
func AccountIDFromCoinPowerStoreKey(key []byte) chainTypes.AccountID {
	return coinStoreKey2AccountID(CoinPowerStoreKeyPrefix, key)
}

// CoinLockedStoreKey get the key of coin store keeper for asset
func CoinLockedStoreKey(account chainTypes.AccountID) []byte {
	return genCoinStoreKey(CoinLockedStoreKeyPrefix, account.StoreKey())
}

// AccountIDFromCoinLockedStoreKey get accountID from key
func AccountIDFromCoinLockedStoreKey(key []byte) chainTypes.AccountID {
	return coinStoreKey2AccountID(CoinLockedStoreKeyPrefix, key)
}

// CoinLockedStatStoreKey get the key of coin store keeper for asset
func CoinLockedStatStoreKey(account chainTypes.AccountID) []byte {
	return genCoinStoreKey(CoinLockedStatStoreKeyPrefix, account.StoreKey())
}

// AccountIDFromCoinLockedStatStoreKey get accountID from key
func AccountIDFromCoinLockedStatStoreKey(key []byte) chainTypes.AccountID {
	return coinStoreKey2AccountID(CoinLockedStatStoreKeyPrefix, key)
}

// CoinStatStoreKey get the key of coin state store keeper for asset
func CoinStatStoreKey(creator, symbol chainTypes.Name) []byte {
	if creator.Empty() {
		return genCoinStoreKey(CoinStatStoreKeyPrefix, symbol.Bytes())
	}
	return genCoinStoreKey(CoinStatStoreKeyPrefix, creator.Bytes(), symbol.Bytes())
}

// CoinDescStoreKey get the key of coin desc store keeper for asset
func CoinDescStoreKey(creator, symbol chainTypes.Name) []byte {
	if creator.Empty() {
		return genCoinStoreKey(CoinDescStoreKeyPrefix, symbol.Bytes())
	}
	return genCoinStoreKey(CoinDescStoreKeyPrefix, creator.Bytes(), symbol.Bytes())
}

// ApproveStoreKey get the key of coin approve store keeper for asset
func ApproveStoreKey(account, spender AccountID) []byte {
	return genCoinStoreKey(CoinApproveStoreKeyPrefix, account.Bytes(), spender.Bytes())
}

func ApproveSumStoreKey(account AccountID) []byte {
	return genCoinStoreKey(CoinApproveSumStoreKeyPrefix, account.Bytes())
}
