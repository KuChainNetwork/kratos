package simulation

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding supply type
func DecodeStore(cdc *codec.LegacyAmino, kvA, kvB kv.Pair) string {
	switch {
	default:
		panic(fmt.Sprintf("invalid supply key %X", kvA.Key))
	}
}
