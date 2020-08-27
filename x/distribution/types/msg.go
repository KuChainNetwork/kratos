//nolint
package types

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type findAccount func(acc chainType.AccountID) bool

var FindAcc findAccount

// Verify interface at compile time
var _, _, _ sdk.Msg = &MsgSetWithdrawAccountId{}, &MsgWithdrawDelegatorReward{}, &MsgWithdrawValidatorCommission{}

type MsgSetWithdrawAccountIdData struct {
	DelegatorAccountid chainType.AccountID `json:"delegator_accountid" yaml:"delegator_accountid"`
	WithdrawAccountid  chainType.AccountID `json:"withdraw_accountid" yaml:"withdraw_accountid"`
}

func (m MsgSetWithdrawAccountIdData) Sender() AccountID {
	return m.DelegatorAccountid
}

func (MsgSetWithdrawAccountIdData) Type() Name { return MustName("withdrawcccid") }

func (m MsgSetWithdrawAccountIdData) Marshal() ([]byte, error) {
	return ModuleCdc.MarshalJSON(m)
}

func (m *MsgSetWithdrawAccountIdData) Unmarshal(b []byte) error {
	return ModuleCdc.UnmarshalJSON(b, m)
}

type MsgSetWithdrawAccountId struct {
	KuMsg
}

func (m MsgSetWithdrawAccountId) GetData() (MsgSetWithdrawAccountIdData, error) {
	res := MsgSetWithdrawAccountIdData{}
	if err := m.UnmarshalData(Cdc(), &res); err != nil {
		return MsgSetWithdrawAccountIdData{}, sdkerrors.Wrapf(chainType.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

func (m MsgSetWithdrawAccountId) ValidateBasic() error {
	data, err := m.GetData()
	if err == nil {
		_, ok := data.WithdrawAccountid.ToName()
		if ok {
			if data.DelegatorAccountid.Equal(&data.WithdrawAccountid) {
				return chainType.ErrKuMsgDataSameAccount
			}

			if FindAcc != nil {
				if !FindAcc(data.WithdrawAccountid) {
					return chainType.ErrKuMsgDataNotFindAccount
				}
			}
		}
	}

	return m.KuMsg.ValidateTransfer()
}

func NewMsgSetWithdrawAccountId(auth AccAddress, delAddr, withdrawAddr AccountID) MsgSetWithdrawAccountId {
	return MsgSetWithdrawAccountId{
		*msg.MustNewKuMsg(
			MustName(RouterKey),
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgSetWithdrawAccountIdData{
				DelegatorAccountid: delAddr,
				WithdrawAccountid:  withdrawAddr,
			}),
		),
	}
}

type MsgWithdrawDelegatorRewardData struct {
	DelegatorAccountId chainType.AccountID `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAccountId chainType.AccountID `json:"validator_address" yaml:"validator_address"`
}

func (m MsgWithdrawDelegatorRewardData) Sender() AccountID {
	return m.DelegatorAccountId
}

func (MsgWithdrawDelegatorRewardData) Type() Name { return MustName("withdrawdelreward") }

func (m MsgWithdrawDelegatorRewardData) Marshal() ([]byte, error) {
	return ModuleCdc.MarshalJSON(m)
}

func (m *MsgWithdrawDelegatorRewardData) Unmarshal(b []byte) error {
	return ModuleCdc.UnmarshalJSON(b, m)
}

type MsgWithdrawDelegatorReward struct {
	KuMsg
}

func (m MsgWithdrawDelegatorReward) ValidateBasic() error {
	data, err := m.GetData()
	if err == nil {
		_, ok := data.DelegatorAccountId.ToName()
		if ok {
			if FindAcc != nil {
				if !FindAcc(data.DelegatorAccountId) {
					return chainType.ErrKuMsgDataNotFindAccount
				}
			}
		}
		_, ok = data.ValidatorAccountId.ToName()
		if ok {
			if FindAcc != nil {
				if !FindAcc(data.ValidatorAccountId) {
					return chainType.ErrKuMsgDataNotFindAccount
				}
			}
		}
	}

	return m.KuMsg.ValidateTransfer()
}

func (m MsgWithdrawDelegatorReward) GetData() (MsgWithdrawDelegatorRewardData, error) {
	res := MsgWithdrawDelegatorRewardData{}
	if err := m.UnmarshalData(Cdc(), &res); err != nil {
		return MsgWithdrawDelegatorRewardData{}, sdkerrors.Wrapf(chainType.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

func NewMsgWithdrawDelegatorReward(auth AccAddress, delAddr, valAddr AccountID) MsgWithdrawDelegatorReward {
	return MsgWithdrawDelegatorReward{
		*msg.MustNewKuMsg(
			MustName(RouterKey),
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgWithdrawDelegatorRewardData{
				DelegatorAccountId: delAddr,
				ValidatorAccountId: valAddr,
			}),
		),
	}
}

type MsgWithdrawValidatorCommissionData struct {
	ValidatorAccountId chainType.AccountID `json:"validator_address" yaml:"validator_address"`
}

func (m MsgWithdrawValidatorCommissionData) Sender() AccountID {
	return m.ValidatorAccountId
}

func (MsgWithdrawValidatorCommissionData) Type() Name { return MustName("withdrawvalcom") }

func (m MsgWithdrawValidatorCommissionData) Marshal() ([]byte, error) {
	return ModuleCdc.MarshalJSON(m)
}

func (m *MsgWithdrawValidatorCommissionData) Unmarshal(b []byte) error {
	return ModuleCdc.UnmarshalJSON(b, m)
}

type MsgWithdrawValidatorCommission struct {
	KuMsg
}

func (m MsgWithdrawValidatorCommission) GetData() (MsgWithdrawValidatorCommissionData, error) {
	res := MsgWithdrawValidatorCommissionData{}
	if err := m.UnmarshalData(Cdc(), &res); err != nil {
		return MsgWithdrawValidatorCommissionData{}, sdkerrors.Wrapf(chainType.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

func (m MsgWithdrawValidatorCommission) ValidateBasic() error {
	data, err := m.GetData()
	if err == nil {
		_, ok := data.ValidatorAccountId.ToName()
		if ok {
			if FindAcc != nil {
				if !FindAcc(data.ValidatorAccountId) {
					return chainType.ErrKuMsgDataNotFindAccount
				}
			}
		}
	}

	return m.KuMsg.ValidateTransfer()
}

func NewMsgWithdrawValidatorCommission(auth AccAddress, valAddr AccountID) MsgWithdrawValidatorCommission {
	return MsgWithdrawValidatorCommission{
		*msg.MustNewKuMsg(
			MustName(RouterKey),
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgWithdrawValidatorCommissionData{
				ValidatorAccountId: valAddr,
			}),
		),
	}
}
