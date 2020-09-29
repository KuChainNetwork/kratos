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

	cdc.RegisterConcrete(&MsgCreateSymbol{}, "dex/MsgCreateSymbol", nil)
	cdc.RegisterConcrete(&MsgCreateSymbolData{}, "dex/MsgCreateSymbolData", nil)
	cdc.RegisterConcrete(&MsgUpdateSymbol{}, "dex/MsgUpdateSymbol", nil)
	cdc.RegisterConcrete(&MsgUpdateSymbolData{}, "dex/MsgUpdateSymbolData", nil)
	cdc.RegisterConcrete(&MsgPauseSymbol{}, "dex/MsgPauseSymbol", nil)
	cdc.RegisterConcrete(&MsgPauseSymbolData{}, "dex/MsgPauseSymbolData", nil)
	cdc.RegisterConcrete(&MsgRestoreSymbol{}, "dex/MsgRestoreSymbol", nil)
	cdc.RegisterConcrete(&MsgRestoreSymbolData{}, "dex/MsgRestoreSymbolData", nil)
	cdc.RegisterConcrete(&MsgShutdownSymbol{}, "dex/MsgShutdownSymbol", nil)
	cdc.RegisterConcrete(&MsgShutdownSymbolData{}, "dex/MsgShutdownSymbolData", nil)

	cdc.RegisterConcrete(&MsgDexSigIn{}, "dex/sigin", nil)
	cdc.RegisterConcrete(&MsgDexSigInData{}, "dex/siginData", nil)
	cdc.RegisterConcrete(&MsgDexSigOut{}, "dex/sigout", nil)
	cdc.RegisterConcrete(&MsgDexSigOutData{}, "dex/sigoutData", nil)
	cdc.RegisterConcrete(&MsgDexDeal{}, "dex/deal", nil)
	cdc.RegisterConcrete(&MsgDexDealData{}, "dex/dealData", nil)
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
