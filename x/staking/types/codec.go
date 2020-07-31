package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers the necessary x/staking interfaces and concrete types
// on the provided Amino codec. These types are used for Amino JSON serialization.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgCreateValidator{}, "kuchain/MsgCreateValidator", nil)
	cdc.RegisterConcrete(&MsgEditValidator{}, "kuchain/MsgEditValidator", nil)
	cdc.RegisterConcrete(&MsgDelegate{}, "kuchain/MsgDelegate", nil)
	cdc.RegisterConcrete(&MsgUndelegate{}, "kuchain/MsgUndelegate", nil)
	cdc.RegisterConcrete(&MsgBeginRedelegate{}, "kuchain/MsgBeginRedelegate", nil)

	cdc.RegisterConcrete(KuMsgCreateValidator{}, "kuchain/KuMsgCreateValidator", nil)
	cdc.RegisterConcrete(KuMsgDelegate{}, "kuchain/KuMsgDelegate", nil)
	cdc.RegisterConcrete(KuMsgEditValidator{}, "kuchain/KuMsgEditValidator", nil)
	cdc.RegisterConcrete(KuMsgRedelegate{}, "kuchain/KuMsgRedelegate", nil)
	cdc.RegisterConcrete(KuMsgUnbond{}, "kuchain/KuMsgUnbond", nil)
}

var (
	// ModuleCdc references the global x/staking module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.New()
)

// Cdc get codec for types
func Cdc() *codec.Codec {
	return ModuleCdc
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
