package types

import (
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// verify interface at compile time
//var _ sdk.Msg = &MsgUnjail{}
var _ chainType.KuMsgData = (*MsgUnjail)(nil)

// MsgUnjail - struct for unjailing jailed validator
type MsgUnjail struct {
	ValidatorAddr AccountID `json:"address" yaml:"address"`
}

// NewMsgUnjail creates a new MsgUnjail instance
func NewMsgUnjail(validatorAddr AccountID) MsgUnjail {
	return MsgUnjail{
		ValidatorAddr: validatorAddr,
	}
}

//nolint
func (msg MsgUnjail) Route() string     { return RouterKey }
func (msg MsgUnjail) Type() Name        { return MustName("unjail") }
func (msg MsgUnjail) Sender() AccountID { return msg.ValidatorAddr }
func (msg MsgUnjail) GetSigners() []sdk.AccAddress {
	valAccAddress, ok := msg.ValidatorAddr.ToAccAddress()
	if ok {
		return []sdk.AccAddress{valAccAddress}
	}
	return []sdk.AccAddress{}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgUnjail) ValidateBasic() error {
	if msg.ValidatorAddr.Empty() {
		return ErrBadValidatorAddr
	}

	return nil
}
