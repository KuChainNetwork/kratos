package simapp

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&noExistMsgData{}, "testNoExist/Data", nil)
	cdc.RegisterConcrete(&noMsg{}, "testNoExist/Msg", nil)
}

// ModuleCdc module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

// Cdc get codec for types
func Cdc() *codec.Codec {
	return ModuleCdc
}
