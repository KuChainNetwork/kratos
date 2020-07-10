package types

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	chaintype "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

var (
	RouterKeyName = chaintype.MustName(RouterKey)
)

type KuMsgCreateValidator struct {
	chaintype.KuMsg
}

//func (KuMsgCreateValidator) Type() chaintype.Name { return types.MustName("createvalidator") }

// NewMsgCreate new create coin msg
func NewKuMsgCreateValidator(auth sdk.AccAddress, valAddr chaintype.AccountID, pubKey crypto.PubKey,
	description Description, commission sdk.Dec, delAcc chaintype.AccountID) KuMsgCreateValidator {

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

type KuMsgDelegate struct {
	chaintype.KuMsg
}

// NewKuMsgDelegate create kuMsgDelegate
func NewKuMsgDelegate(auth sdk.AccAddress, delAddr chaintype.AccountID, valAddr chaintype.AccountID, amount sdk.Coin) KuMsgDelegate {
	return KuMsgDelegate{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithTransfer(delAddr, ModuleAccountID, sdk.Coins{amount}),
			msg.WithData(Cdc(), &MsgDelegate{
				DelegatorAccount: delAddr,
				ValidatorAccount: valAddr,
				Amount:           amount,
			}),
		),
	}
}

type KuMsgEditValidator struct {
	chaintype.KuMsg
}

func NewKuMsgEditValidator(auth sdk.AccAddress, valAddr chaintype.AccountID, description Description, newRate *sdk.Dec) KuMsgEditValidator {

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

type KuMsgRedelegate struct {
	chaintype.KuMsg
}

func NewKuMsgRedelegate(auth sdk.AccAddress, delAddr chaintype.AccountID, valSrcAddr, valDstAddr chaintype.AccountID, amount sdk.Coin) KuMsgRedelegate {

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

type KuMsgUnbond struct {
	chaintype.KuMsg
}

func NewKuMsgUnbond(auth sdk.AccAddress, delAddr chaintype.AccountID, valAddr chaintype.AccountID, amount sdk.Coin) KuMsgUnbond {

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
