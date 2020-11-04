package store

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	KVStore = sdk.KVStore
)

type (
	TraceContext = store.TraceContext
	CacheWrap    = store.CacheWrap
	StoreType    = store.StoreType
	Iterator     = store.Iterator
)
