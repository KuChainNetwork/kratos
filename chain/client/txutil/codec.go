package txutil

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc auth module wide codec
var ModuleCdc = codec.NewLegacyAmino()

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.LegacyAmino) {
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
