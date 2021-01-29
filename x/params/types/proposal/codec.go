package proposal

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc module codec
var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()

	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers all necessary param module types with a given codec.
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(ParameterChangeProposal{}, "kuchain/ParameterChangeProposal", nil)
}
