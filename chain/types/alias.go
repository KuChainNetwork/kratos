package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	AccAddress = sdk.AccAddress
	Tx         = sdk.Tx
)

var (
	AccAddressFromBech32 = sdk.AccAddressFromBech32
)

type (
	StoreKey          = sdk.StoreKey
	KVStoreKey        = sdk.KVStoreKey
	TransientStoreKey = sdk.TransientStoreKey
)

// MustAccAddressFromBech32 AccAddressFromBech32 if error then panic
func MustAccAddressFromBech32(str string) AccAddress {
	res, err := AccAddressFromBech32(str)
	if err != nil {
		panic(err)
	}
	return res
}
