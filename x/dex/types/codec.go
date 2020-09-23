package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc auth module wide codec
var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&GenesisState{}, "dex/genesisState", nil)

	cdc.RegisterConcrete(&Dex{}, "dex/dexData", nil)
	cdc.RegisterConcrete(&MsgCreateDex{}, "dex/MsgCreateDex", nil)
	cdc.RegisterConcrete(&MsgCreateDexData{}, "dex/MsgCreateDexData", nil)
	cdc.RegisterConcrete(&MsgUpdateDexDescription{}, "dex/MsgUpdateDexDescription", nil)
	cdc.RegisterConcrete(&MsgUpdateDexDescriptionData{}, "dex/MsgUpdateDexDescriptionData", nil)
	cdc.RegisterConcrete(&MsgDestroyDex{}, "dex/MsgDestroyDex", nil)
	cdc.RegisterConcrete(&MsgDestroyDexData{}, "dex/MsgDestroyDexData", nil)
	cdc.RegisterConcrete(&MsgCreateCurrency{}, "dex/MsgCreateCurrency", nil)
	cdc.RegisterConcrete(&MsgCreateCurrencyData{}, "dex/MsgCreateCurrencyData", nil)
	cdc.RegisterConcrete(&MsgUpdateCurrency{}, "dex/MsgUpdateCurrency", nil)
	cdc.RegisterConcrete(&MsgUpdateCurrencyData{}, "dex/MsgUpdateCurrencyData", nil)
	cdc.RegisterConcrete(&MsgShutdownCurrency{}, "dex/MsgShutdownCurrency", nil)
	cdc.RegisterConcrete(&MsgShutdownCurrencyData{}, "dex/MsgShutdownCurrencyData", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

// Cdc get codec for dex
func Cdc() *codec.Codec {
	return ModuleCdc
}
