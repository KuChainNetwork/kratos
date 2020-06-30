package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers the necessary x/distribution interfaces and concrete types
// on the provided Amino codec. These types are used for Amino JSON serialization.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgWithdrawDelegatorReward{}, "kuchain/MsgWithdrawDelegationReward", nil)
	cdc.RegisterConcrete(&MsgWithdrawDelegatorRewardData{}, "kuchain/MsgWithdrawDelegationRewardData", nil)

	cdc.RegisterConcrete(MsgWithdrawValidatorCommission{}, "kuchain/MsgWithdrawValidatorCommission", nil)
	cdc.RegisterConcrete(&MsgWithdrawValidatorCommissionData{}, "kuchain/MsgWithdrawValidatorCommissionData", nil)

	cdc.RegisterConcrete(&MsgSetWithdrawAccountIdData{}, "kuchain/MsgSetWithdrawAccountIdData", nil)
	cdc.RegisterConcrete(MsgSetWithdrawAccountId{}, "kuchain/MsgSetWithdrawAccountId", nil)

	cdc.RegisterConcrete(MsgFundCommunityPool{}, "cosmos-sdk/MsgFundCommunityPool", nil)
	cdc.RegisterConcrete(&MsgFundCommunityPoolData{}, "kuchain/MsgFundCommunityPoolData", nil)

	cdc.RegisterConcrete(CommunityPoolSpendProposal{}, "cosmos-sdk/CommunityPoolSpendProposal", nil)
}

var (
	amino = codec.New()

	// ModuleCdc references the global x/distribution module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding as Amino
	// is still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/distribution and
	// defined at the application level.
	ModuleCdc = codec.NewHybridCodec(amino)
)

// Cdc get codec for types
func Cdc() *codec.Codec {
	return amino
}

func init() {
	RegisterCodec(amino)
	codec.RegisterCrypto(amino)
	amino.Seal()
}
