package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers all the necessary types and interfaces for the
// governance module.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Content)(nil), nil)
	cdc.RegisterConcrete(&MsgSubmitProposalBase{}, "kuchain/MsgSubmitProposalBase", nil)
	cdc.RegisterConcrete(MsgSubmitProposal{}, "kuchain/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(&MsgDeposit{}, "kuchain/MsgDeposit", nil)
	cdc.RegisterConcrete(&MsgVote{}, "kuchain/MsgVote", nil)
	cdc.RegisterConcrete(TextProposal{}, "kuchain/TextProposal", nil)

	cdc.RegisterConcrete(KuMsgSubmitProposal{}, "kuchain/kuMsgSubmitProposal", nil)
	cdc.RegisterConcrete(KuMsgDeposit{}, "kuchain/kuMsgDeposit", nil)
	cdc.RegisterConcrete(KuMsgVote{}, "kuchain/kuMsgVote", nil)
}

// RegisterProposalTypeCodec registers an external proposal content type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitProposal
// to be correctly Amino encoded and decoded.
//
// NOTE: This should only be used for applications that are still using a concrete
// Amino codec for serialization.
func RegisterProposalTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

var (
	// ModuleCdc references the global x/gov module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/gov and
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
}
