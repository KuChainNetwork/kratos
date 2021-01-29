package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUnjail{}, "kuchain/MsgUnjail", nil)
	cdc.RegisterConcrete(KuMsgUnjail{}, "kuchain/KuMsgUnjail", nil)
}

var (
	// ModuleCdc references the global x/slashing module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding as Amino
	// is still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/slashing and
	// defined at the application level.
	ModuleCdc = codec.NewLegacyAmino()
)

// Cdc get codec for types
func Cdc() *codec.LegacyAmino {
	return ModuleCdc
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
