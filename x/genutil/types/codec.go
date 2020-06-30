package types

import (
	"github.com/KuChain-io/kuchain/chain/client/txutil"
	"github.com/KuChain-io/kuchain/chain/types"
	accounttypes "github.com/KuChain-io/kuchain/x/account/types"
	assettypes "github.com/KuChain-io/kuchain/x/asset/types"
	evidencetypes "github.com/KuChain-io/kuchain/x/evidence/types"
	stakingtypes "github.com/KuChain-io/kuchain/x/staking/types"
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
