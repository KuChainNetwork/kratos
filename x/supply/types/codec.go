package types

import (
	"github.com/KuChainNetwork/kuchain/x/supply/exported"
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&Supply{}, "kuchain/Supply", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.LegacyAmino

func init() {
	cdc := codec.NewLegacyAmino()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}
