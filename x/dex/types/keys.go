package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	// ModuleName is the name of the module
	ModuleName = "dex"

	// StoreKey is the store key string for slashing
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute is the querier route for slashing
	QuerierRoute = ModuleName
)

var (
	ModuleAccountID = NewAccountIDFromName(MustName(ModuleName))
)

var (
	loggerName = fmt.Sprintf("x/%s", ModuleName)

	// RouterKeyName is the name type of router
	RouterKeyName = MustName(RouterKey)

	// ModuleDexKeyPrefix prefix for asset store
	ModuleDexKeyPrefix = []byte{0x0B}

	DexNumberStoreKey = MustName("dex.number").Bytes()

	DexStoreKeyPrefix          = MustName("dex").Bytes()
	DexSigInStoreKeyPrefix     = MustName("dex.sigin").Bytes()
	DexSigSumStoreKeyPrefix    = MustName("dex.sigsum").Bytes()
	DexSigOutReqStoreKeyPrefix = MustName("dex.sigoutreq").Bytes()
)

func Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", loggerName)
}

func genStoreKey(pre []byte, keys ...[]byte) []byte {
	res := make([]byte, 0, 256)
	res = append(res, ModuleDexKeyPrefix...)
	res = append(res, pre...)
	for _, k := range keys {
		res = append(res, k...)
	}
	return res
}

func GenStoreKey(pre []byte, keys ...[]byte) []byte {
	return genStoreKey(pre, keys...)
}

// DexStoreKey get the key of coin state store keeper for asset
func DexStoreKey(creator Name) []byte {
	return genStoreKey(DexStoreKeyPrefix, creator.Bytes())
}

// GetNumberStoreKey get the key of next dex number
func GetDexNumberStoreKey() []byte {
	return genStoreKey(DexNumberStoreKey)
}

// DexSigOutReqStoreKey get the key of coin state store keeper for asset
func DexSigOutReqStoreKey(user AccountID) []byte {
	return genStoreKey(DexSigOutReqStoreKeyPrefix, user.Bytes())
}
