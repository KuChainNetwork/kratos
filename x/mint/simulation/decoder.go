package simulation

import (
	"bytes"
	"fmt"

	"github.com/KuChainNetwork/kuchain/x/mint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding mint type
func DecodeStore(cdc *codec.LegacyAmino, kvA, kvB kv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key, types.MinterKey):
		var minterA, minterB types.Minter
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &minterA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &minterB)
		return fmt.Sprintf("%v\n%v", minterA, minterB)
	default:
		panic(fmt.Sprintf("invalid mint key %X", kvA.Key))
	}
}
