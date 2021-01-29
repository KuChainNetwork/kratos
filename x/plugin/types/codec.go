package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

// ModuleCdc auth module wide codec
var ModuleCdc = codec.NewLegacyAmino()

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&GenesisState{}, "plugin/genesisState", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	cryptocodec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
