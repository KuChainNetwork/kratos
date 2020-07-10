package types

import (
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc auth module wide codec
var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.AccountAuthKeeper)(nil), nil)
	cdc.RegisterInterface((*exported.AccountStatKeeper)(nil), nil)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterInterface((*exported.GenesisAccount)(nil), nil)
	cdc.RegisterInterface((*exported.AuthAccountKeeper)(nil), nil)

	cdc.RegisterConcrete(&MsgCreateAccountData{}, "account/createData", nil)
	cdc.RegisterConcrete(&MsgCreateAccount{}, "account/createMsg", nil)

	cdc.RegisterConcrete(&MsgUpdateAccountAuthData{}, "account/upAuthData", nil)
	cdc.RegisterConcrete(&MsgUpdateAccountAuth{}, "account/upAuth", nil)

	cdc.RegisterConcrete(&KuAccount{}, "kuchain/Account", nil)
	cdc.RegisterConcrete(&ModuleAccount{}, "kuchain/ModuleAccount", nil)

	cdc.RegisterConcrete(&AuthAccounts{}, "account/authAccounts", nil)
}

// RegisterAccountTypeCodec registers an external account type defined in
// another module for the internal ModuleCdc.
func RegisterAccountTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

// Cdc get codec for types
func Cdc() *codec.Codec {
	return ModuleCdc
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
