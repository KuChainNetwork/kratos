package staking

import (
	"encoding/json"
	"fmt"

	"github.com/KuChainNetwork/kuchain/x/asset"
	genutiltypes "github.com/KuChainNetwork/kuchain/x/genutil/types"
	stakingtypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ genutiltypes.StakingFuncManager = FuncManager{}
)

type FuncManager struct{}

func NewFuncManager() FuncManager {
	return FuncManager{}
}

func (FuncManager) MsgCreateValidatorVal(msg sdk.Msg) bool {
	if _, ok := msg.(stakingtypes.KuMsgCreateValidator); !ok {
		return false
	}

	return true
}

func (FuncManager) GetMsgCreateValidatorMoniker(m sdk.Msg) (string, error) {
	msg := m.(stakingtypes.KuMsgCreateValidator)

	msgData := stakingtypes.MsgCreateValidator{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return "", sdkerrors.Wrapf(err, "msg CreateValidatorMoniker data unmarshal error")
	}
	return msgData.Description.Moniker, nil
}

func (FuncManager) MsgDelegateWithBalance(m sdk.Msg,
	balancesMap map[string]asset.GenesisAsset) error {
	kumsg := m.(stakingtypes.KuMsgDelegate)
	msgData := stakingtypes.MsgDelegate{}
	if err := kumsg.UnmarshalData(Cdc(), &msgData); err != nil {
		return sdkerrors.Wrapf(err, "msg DelegateWithBalance data unmarshal error")
	}

	delAcc := msgData.DelegatorAccount
	valAcc := msgData.ValidatorAccount

	_, delOk := balancesMap[delAcc.String()]
	if !delOk {
		return fmt.Errorf("account %s balance not in genesis state: %+v", delAcc.String(), balancesMap)
	}

	_, valOk := balancesMap[valAcc.String()]
	if !valOk {
		return fmt.Errorf("account %s balance not in genesis state: %+v", valAcc.String(), balancesMap)
	}

	delBal, delOk := balancesMap[delAcc.String()]
	if !delOk {
		return fmt.Errorf("account %s balance not in genesis state: %+v", delAcc.String(), balancesMap)
	}

	if delBal.GetCoins().AmountOf(msgData.Amount.Denom).LT(msgData.Amount.Amount) {
		return fmt.Errorf("insufficient fund for delegation %v: %v < %v",
			delBal.GetID().String(), delBal.GetCoins().AmountOf(msgData.Amount.Denom), msgData.Amount.Amount,
		)
	}

	return nil
}

// GetBondDenom
func (FuncManager) GetBondDenom(appGenesisState map[string]json.RawMessage) string {
	var stakingData stakingtypes.GenesisState
	stakingtypes.ModuleCdc.MustUnmarshalJSON(appGenesisState[stakingtypes.ModuleName], &stakingData)
	return stakingData.Params.BondDenom
}
