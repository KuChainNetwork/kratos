package types

import (
	chaintype "github.com/KuChain-io/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// verify interface at compile time
//var _ sdk.Msg = &MsgUnjail{}

// NewMsgUnjail creates a new MsgUnjail instance
func NewMsgUnjail(validatorAddr chaintype.AccountID) MsgUnjail {
	return MsgUnjail{
		ValidatorAddr: validatorAddr,
	}
}

//nolint
func (msg MsgUnjail) Route() string        { return RouterKey }
func (msg MsgUnjail) Type() chaintype.Name { return chaintype.MustName("unjail") }
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
