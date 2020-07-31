package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var (
	// ModuleCdc references the global x/mint module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding as
	// Amino is still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/mint and
	// defined at the application level.
	ModuleCdc = codec.New()
)

func Cdc() *codec.Codec {
	return ModuleCdc
}

func init() {
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
