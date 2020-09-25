package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc auth module wide codec
var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*GenesisAsset)(nil), nil)
	cdc.RegisterInterface((*GenesisCoin)(nil), nil)
	cdc.RegisterConcrete(&GenesisState{}, "asset/genesisState", nil)
	cdc.RegisterConcrete(&BaseGensisAssetCoin{}, "asset/genesisCoin", nil)
	cdc.RegisterConcrete(&BaseGenesisAsset{}, "asset/genesisAsset", nil)

	cdc.RegisterConcrete(&KuMsg{}, "kuchain/msg", nil)

	cdc.RegisterConcrete(&MsgTransfer{}, "asset/transfer", nil)
	cdc.RegisterConcrete(&MsgCreateCoinData{}, "asset/createData", nil)
	cdc.RegisterConcrete(&MsgCreateCoin{}, "asset/create", nil)
	cdc.RegisterConcrete(&MsgIssueCoinData{}, "asset/issueData", nil)
	cdc.RegisterConcrete(&MsgIssueCoin{}, "asset/issue", nil)
	cdc.RegisterConcrete(&MsgBurnCoinData{}, "asset/burnData", nil)
	cdc.RegisterConcrete(&MsgBurnCoin{}, "asset/burn", nil)
	cdc.RegisterConcrete(&MsgLockCoinData{}, "asset/lockData", nil)
	cdc.RegisterConcrete(&MsgLockCoin{}, "asset/lock", nil)
	cdc.RegisterConcrete(&MsgUnlockCoinData{}, "asset/unlockData", nil)
	cdc.RegisterConcrete(&MsgUnlockCoin{}, "asset/unlock", nil)
	cdc.RegisterConcrete(&MsgExerciseCoinData{}, "asset/exerciseData", nil)
	cdc.RegisterConcrete(&MsgExerciseCoin{}, "asset/exercise", nil)
	cdc.RegisterConcrete(&MsgApproveData{}, "asset/approveData", nil)
	cdc.RegisterConcrete(&MsgApprove{}, "asset/approve", nil)
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
