package types

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

var (
	RouterKeyName = chainTypes.MustName(RouterKey)
)

type KuMsgCreateValidator struct {
	chainTypes.KuMsg
}

//func (KuMsgCreateValidator) Type() chainTypes.Name { return types.MustName("createvalidator") }

// NewMsgCreate new create coin msg
func NewKuMsgCreateValidator(auth sdk.AccAddress, valAddr chainTypes.AccountID, pubKey crypto.PubKey,
	description Description, commission sdk.Dec, delAcc chainTypes.AccountID) KuMsgCreateValidator {

	var pkStr string
	if pubKey != nil {
		pkStr = sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubKey)
	}
	return KuMsgCreateValidator{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgCreateValidator{
				Description:      description,
				ValidatorAccount: valAddr,
				Pubkey:           pkStr,
				DelegatorAccount: delAcc,
				CommissionRates:  commission,
			}),
		),
	}
}

func (msg KuMsgCreateValidator) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	msgData := MsgCreateValidator{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return err
	}
	return msgData.ValidateBasic()
}

type KuMsgDelegate struct {
	chainTypes.KuMsg
}

// NewKuMsgDelegate create kuMsgDelegate
func NewKuMsgDelegate(auth sdk.AccAddress, delAddr chainTypes.AccountID, valAddr chainTypes.AccountID, amount chainTypes.Coin) KuMsgDelegate {
	return KuMsgDelegate{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithTransfer(delAddr, ModuleAccountID, chainTypes.Coins{amount}),
			msg.WithData(Cdc(), &MsgDelegate{
				DelegatorAccount: delAddr,
				ValidatorAccount: valAddr,
				Amount:           amount,
			}),
		),
	}
}

func (msg KuMsgDelegate) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	msgData := MsgDelegate{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return err
	}

	if err := msg.KuMsg.ValidateTransferRequire(ModuleAccountID, chainTypes.NewCoins(msgData.Amount)); err != nil {
		return chainTypes.ErrKuMsgInconsistentAmount
	}

	return msgData.ValidateBasic()
}

type KuMsgEditValidator struct {
	chainTypes.KuMsg
}

func NewKuMsgEditValidator(auth sdk.AccAddress, valAddr chainTypes.AccountID, description Description, newRate *sdk.Dec) KuMsgEditValidator {

	return KuMsgEditValidator{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgEditValidator{
				Description:      description,
				CommissionRate:   newRate,
				ValidatorAccount: valAddr,
			}),
		),
	}
}

func (msg KuMsgEditValidator) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	msgData := MsgEditValidator{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return err
	}
	return msgData.ValidateBasic()
}

type KuMsgRedelegate struct {
	chainTypes.KuMsg
}

func NewKuMsgRedelegate(auth sdk.AccAddress, delAddr chainTypes.AccountID, valSrcAddr, valDstAddr chainTypes.AccountID, amount chainTypes.Coin) KuMsgRedelegate {

	return KuMsgRedelegate{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgBeginRedelegate{
				DelegatorAccount:    delAddr,
				ValidatorSrcAccount: valSrcAddr,
				ValidatorDstAccount: valDstAddr,
				Amount:              amount,
			}),
		),
	}
}

func (msg KuMsgRedelegate) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	msgData := MsgBeginRedelegate{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return err
	}
	return msgData.ValidateBasic()
}

type KuMsgUnbond struct {
	chainTypes.KuMsg
}

func NewKuMsgUnbond(auth sdk.AccAddress, delAddr chainTypes.AccountID, valAddr chainTypes.AccountID, amount chainTypes.Coin) KuMsgUnbond {

	return KuMsgUnbond{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgUndelegate{
				DelegatorAccount: delAddr,
				ValidatorAccount: valAddr,
				Amount:           amount,
			}),
		),
	}
}

func (msg KuMsgUnbond) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	msgData := MsgUndelegate{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return err
	}
	return msgData.ValidateBasic()
}
