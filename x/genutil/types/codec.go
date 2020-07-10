package types

import (
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	"github.com/KuChainNetwork/kuchain/chain/types"
	accounttypes "github.com/KuChainNetwork/kuchain/x/account/types"
	assettypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	evidencetypes "github.com/KuChainNetwork/kuchain/x/evidence/types"
	stakingtypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ModuleCdc defines a generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

// required for genesis transactions
func init() {
	ModuleCdc = codec.New()
	stakingtypes.RegisterCodec(ModuleCdc)
	evidencetypes.RegisterCodec(ModuleCdc)
	accounttypes.RegisterCodec(ModuleCdc)
	assettypes.RegisterCodec(ModuleCdc)
	sdk.RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	txutil.RegisterCodec(ModuleCdc)
	types.RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
