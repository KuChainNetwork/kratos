//nolint
package types

import (
	"github.com/KuChain-io/kuchain/chain/msg"
	chainType "github.com/KuChain-io/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Verify interface at compile time
var _, _, _ sdk.Msg = &MsgSetWithdrawAccountId{}, &MsgWithdrawDelegatorReward{}, &MsgWithdrawValidatorCommission{}

type MsgSetWithdrawAccountIdData struct {
	DelegatorAccountid chainType.AccountID `protobuf:"bytes,1,opt,name=delegator_accountid,json=delegatorAccountid,proto3" json:"delegator_accountid" yaml:"delegator_accountid"`
	WithdrawAccountid  chainType.AccountID `protobuf:"bytes,2,opt,name=withdraw_accountid,json=withdrawAccountid,proto3" json:"withdraw_accountid" yaml:"withdraw_accountid"`
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
	DelegatorAccountId chainType.AccountID `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address" yaml:"delegator_address"`
	ValidatorAccountId chainType.AccountID `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address" yaml:"validator_address"`
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
	ValidatorAccountId chainType.AccountID `protobuf:"bytes,1,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address" yaml:"validator_address"`
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

var TypeMsgFundCommunityPool = MustName("fundcommpool")

type MsgFundCommunityPoolData struct {
	Amount    Coins     `json:"amount" yaml:"amount"`
	Depositor AccountID `json:"depositor" yaml:"depositor"`
}

func (MsgFundCommunityPoolData) Type() Name { return MustName("fundcommpool") }

func (m MsgFundCommunityPoolData) Marshal() ([]byte, error) {
	return ModuleCdc.MarshalJSON(m)
}

func (m *MsgFundCommunityPoolData) Unmarshal(b []byte) error {
	return ModuleCdc.UnmarshalJSON(b, m)
}

type MsgFundCommunityPool struct {
	KuMsg
}

// NewMsgFundCommunityPool returns a new MsgFundCommunityPool with a sender and
// a funding amount.
func NewMsgFundCommunityPool(auth AccAddress, amount Coins, depositor AccountID) MsgFundCommunityPool {
	return MsgFundCommunityPool{
		*msg.MustNewKuMsg(
			MustName(RouterKey),
			msg.WithAuth(auth),
			msg.WithTransfer(depositor, ModuleAccountID, amount),
			msg.WithData(Cdc(), &MsgFundCommunityPoolData{
				Amount:    amount,
				Depositor: depositor,
			}),
		),
	}
}

func (m MsgFundCommunityPool) GetData() (MsgFundCommunityPoolData, error) {
	res := MsgFundCommunityPoolData{}
	if err := m.UnmarshalData(Cdc(), &res); err != nil {
		return MsgFundCommunityPoolData{}, sdkerrors.Wrapf(chainType.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

/*
FIXME: support Validate msg
// ValidateBasic performs basic MsgFundCommunityPool message validation.
func (msg MsgFundCommunityPool) ValidateBasic() error {
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	if msg.Depositor.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Depositor.String())
	}

	return nil
}
*/
