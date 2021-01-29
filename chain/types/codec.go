package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(StdTx{}, "kuchain/Tx", nil)
	cdc.RegisterInterface((*KuMsgData)(nil), nil)
}

// ModuleCdc module wide codec
var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	RegisterCodec(ModuleCdc)
	cryptocodec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

// Cdc get codec for types
func Cdc() *codec.LegacyAmino {
	return ModuleCdc
}
