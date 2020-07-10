package types

import (
	//"bytes"

	"github.com/tendermint/tendermint/crypto"

	chaintype "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"github.com/KuChainNetwork/kuchain/chain/types"
)

// NewMsgCreateValidator creates a new MsgCreateValidator instance.
// Delegator address and validator address are the same.
func NewMsgCreateValidator(
	valAddr chaintype.AccountID, pubKey crypto.PubKey,
	description Description, commission sdk.Dec, delAcc chaintype.AccountID,
) MsgCreateValidator {

	var pkStr string
	if pubKey != nil {
		pkStr = sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubKey)
	}

	return MsgCreateValidator{
		Description:      description,
		ValidatorAccount: valAddr,
		Pubkey:           pkStr,
		DelegatorAccount: delAcc,
		CommissionRates:  commission,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgCreateValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (MsgCreateValidator) Type() chaintype.Name { return chaintype.MustName("create@validator") }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgCreateValidator) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	addrs := []sdk.AccAddress{}
	delegatorAccAddress, ok := msg.DelegatorAccount.ToAccAddress()
	if ok { //name   ctx
		addrs = append(addrs, delegatorAccAddress)
	}

	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateValidator) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures either empty or valid
	if msg.DelegatorAccount.Empty() {
		return ErrEmptyDelegatorAddr
	}
	if msg.ValidatorAccount.Empty() {
		return ErrEmptyValidatorAddr
	}

	if msg.Pubkey == "" {
		return ErrEmptyValidatorPubKey
	}
	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.CommissionRates.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "commission_rate is negative")
	}

	if msg.CommissionRates.GT(sdk.OneDec()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "commission_rate is greater then 1")
	}

	return nil
}

// NewMsgEditValidator creates a new MsgEditValidator instance
func NewMsgEditValidator(valAddr chaintype.AccountID, description Description, newRate *sdk.Dec) MsgEditValidator {
	return MsgEditValidator{
		Description:      description,
		CommissionRate:   newRate,
		ValidatorAccount: valAddr,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgEditValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (MsgEditValidator) Type() chaintype.Name { return chaintype.MustName("edit@validator") }

// GetSigners implements the sdk.Msg interface.
func (msg MsgEditValidator) GetSigners() []sdk.AccAddress {
	validatorAccAddress, _ := msg.ValidatorAccount.ToAccAddress()
	return []sdk.AccAddress{validatorAccAddress}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgEditValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgEditValidator) ValidateBasic() error {
	if msg.ValidatorAccount.Empty() {
		return ErrEmptyValidatorAddr
	}
	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}
	if msg.CommissionRate != nil {
		if msg.CommissionRate.GT(sdk.OneDec()) || msg.CommissionRate.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "commission rate must be between 0 and 1 (inclusive)")
		}
	}

	return nil
}

// NewMsgDelegate creates a new MsgDelegate instance.
func NewMsgDelegate(delAddr chaintype.AccountID, valAddr chaintype.AccountID, amount sdk.Coin) MsgDelegate {
	return MsgDelegate{
		DelegatorAccount: delAddr,
		ValidatorAccount: valAddr,
		Amount:           amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgDelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (MsgDelegate) Type() chaintype.Name { return chaintype.MustName("delegate") }

// GetSigners implements the sdk.Msg interface.
func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	delegatorAccAddress, _ := msg.DelegatorAccount.ToAccAddress()
	return []sdk.AccAddress{delegatorAccAddress}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDelegate) ValidateBasic() error {
	if msg.DelegatorAccount.Empty() {
		return ErrEmptyDelegatorAddr
	}
	if msg.ValidatorAccount.Empty() {
		return ErrEmptyValidatorAddr
	}
	if !msg.Amount.Amount.IsPositive() {
		return ErrBadDelegationAmount
	}
	return nil
}

// NewMsgBeginRedelegate creates a new MsgBeginRedelegate instance.
func NewMsgBeginRedelegate(
	delAddr chaintype.AccountID, valSrcAddr, valDstAddr chaintype.AccountID, amount sdk.Coin,
) MsgBeginRedelegate {
	return MsgBeginRedelegate{
		DelegatorAccount:    delAddr,
		ValidatorSrcAccount: valSrcAddr,
		ValidatorDstAccount: valDstAddr,
		Amount:              amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (MsgBeginRedelegate) Type() chaintype.Name { return chaintype.MustName("beginredelegate") }

// GetSigners implements the sdk.Msg interface
func (msg MsgBeginRedelegate) GetSigners() []sdk.AccAddress {
	delegatorAccAddress, _ := msg.DelegatorAccount.ToAccAddress()
	return []sdk.AccAddress{delegatorAccAddress}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) ValidateBasic() error {
	if msg.DelegatorAccount.Empty() {
		return ErrEmptyDelegatorAddr
	}
	if msg.ValidatorSrcAccount.Empty() {
		return ErrEmptyValidatorAddr
	}
	if msg.ValidatorDstAccount.Empty() {
		return ErrEmptyValidatorAddr
	}
	if !msg.Amount.Amount.IsPositive() {
		return ErrBadSharesAmount
	}
	return nil
}

// NewMsgUndelegate creates a new MsgUndelegate instance.
func NewMsgUndelegate(delAddr chaintype.AccountID, valAddr chaintype.AccountID, amount sdk.Coin) MsgUndelegate {
	return MsgUndelegate{
		DelegatorAccount: delAddr,
		ValidatorAccount: valAddr,
		Amount:           amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUndelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (MsgUndelegate) Type() chaintype.Name { return chaintype.MustName("beginunbonding") }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUndelegate) GetSigners() []sdk.AccAddress {
	delegatorAccAddress, _ := msg.DelegatorAccount.ToAccAddress()
	return []sdk.AccAddress{delegatorAccAddress}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUndelegate) ValidateBasic() error {
	if msg.DelegatorAccount.Empty() {
		return ErrEmptyDelegatorAddr
	}
	if msg.ValidatorAccount.Empty() {
		return ErrEmptyValidatorAddr
	}
	if !msg.Amount.Amount.IsPositive() {
		return ErrBadSharesAmount
	}
	return nil
}
