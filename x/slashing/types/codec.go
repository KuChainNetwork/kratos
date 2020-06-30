package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgUnjail{}, "kuchain/MsgUnjail", nil)
	cdc.RegisterConcrete(KuMsgUnjail{}, "kuchain/KuMsgUnjail", nil)
}

var (
	amino = codec.New()

	// ModuleCdc references the global x/slashing module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding as Amino
	// is still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/slashing and
	// defined at the application level.
	ModuleCdc = codec.NewHybridCodec(amino)
)

// Cdc get codec for types
func Cdc() *codec.Codec {
	return amino
}

func init() {
	RegisterCodec(amino)
	codec.RegisterCrypto(amino)
	amino.Seal()
}
