package types

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	"github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgDexSigIn struct {
	types.KuMsg
}

type MsgDexSigInData struct {
	User   AccountID `json:"user" yaml:"user"`
	Dex    AccountID `json:"dex" yaml:"dex"`
	Amount Coins     `json:"amount" yaml:"amount"`
}

// Type imp for data KuMsgData
func (MsgDexSigInData) Type() types.Name { return types.MustName("sigin") }

func (msg MsgDexSigInData) Sender() AccountID {
	return msg.User
}

// NewMsgDexSigIn new dex sig in msg
func NewMsgDexSigIn(auth types.AccAddress, user, dex AccountID, amount Coins) MsgDexSigIn {
	return MsgDexSigIn{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgDexSigInData{
				User:   user,
				Dex:    dex,
				Amount: amount,
			}),
		),
	}
}

func (msg MsgDexSigIn) GetData() (MsgDexSigInData, error) {
	res := MsgDexSigInData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgDexSigInData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

func (msg MsgDexSigIn) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.User.Empty() {
		return sdkerrors.Wrap(types.ErrKuMsgAccountIDNil, "user accountID empty")
	}

	if data.Dex.Empty() {
		return sdkerrors.Wrap(types.ErrKuMsgAccountIDNil, "dex accountID empty")
	}

	if data.User.Eq(data.Dex) {
		return sdkerrors.Wrap(types.ErrKuMsgSpenderShouldNotEqual, "dex should not be equal to user")
	}

	if data.Amount.IsAnyNegative() {
		return types.ErrKuMsgCoinsHasNegative
	}

	return nil
}

type MsgDexSigOut struct {
	types.KuMsg
}

type MsgDexSigOutData struct {
	User      AccountID `json:"user" yaml:"user"`
	Dex       AccountID `json:"dex" yaml:"dex"`
	Amount    Coins     `json:"amount" yaml:"amount"`
	IsTimeout bool      `json:"is_timeout" yaml:"is_timeout"`
}

// Type imp for data KuMsgData
func (MsgDexSigOutData) Type() types.Name { return types.MustName("sigout") }

func (msg MsgDexSigOutData) Sender() AccountID {
	return msg.User
}

// NewMsgDexSigOut new dex sig in msg
func NewMsgDexSigOut(auth types.AccAddress, isTimeout bool, user, dex AccountID, amount Coins) MsgDexSigOut {
	return MsgDexSigOut{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithTransfer(user, dex, Coins{}),
			msg.WithData(Cdc(), &MsgDexSigOutData{
				User:      user,
				Dex:       dex,
				Amount:    amount,
				IsTimeout: isTimeout,
			}),
		),
	}
}

func (msg MsgDexSigOut) GetData() (MsgDexSigOutData, error) {
	res := MsgDexSigOutData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgDexSigOutData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

func (msg MsgDexSigOut) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.User.Empty() {
		return sdkerrors.Wrap(types.ErrKuMsgAccountIDNil, "user accountID empty")
	}

	if data.Dex.Empty() {
		return sdkerrors.Wrap(types.ErrKuMsgAccountIDNil, "dex accountID empty")
	}

	if data.User.Eq(data.Dex) {
		return sdkerrors.Wrap(types.ErrKuMsgSpenderShouldNotEqual, "dex should not be equal to user")
	}

	if data.Amount.IsAnyNegative() {
		return types.ErrKuMsgCoinsHasNegative
	}

	return nil
}

type MsgDexDeal struct {
	types.KuMsg
}

type MsgDexDealData struct {
	Dex     AccountID `json:"dex" yaml:"dex"`
	ExtData []byte    `json:"ext" yaml:"ext"`       // some external datas
	ID      []byte    `json:"id" yaml:"id"`         // L2 block id for prove
	Hash    []byte    `json:"hash" yaml:"hash"`     // L2 hash for prove
	Proves  []byte    `json:"proves" yaml:"proves"` // L2 prove datas
}

// Type imp for data KuMsgData
func (MsgDexDealData) Type() types.Name { return types.MustName("deal") }

func (msg MsgDexDealData) Sender() AccountID {
	return msg.Dex
}

// NewMsgDexSigOut new dex sig in msg
func NewMsgDexDeal(auth types.AccAddress, dex AccountID, from, to AccountID, fromAsset, toAsset, feeFromFrom, feeFromTo Coin, ext []byte) MsgDexDeal {
	return MsgDexDeal{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			// from -> dex -> to
			msg.WithTransfer(from, dex, types.NewCoins(fromAsset.Add(feeFromFrom))),
			msg.WithTransfer(dex, to, types.NewCoins(fromAsset)),
			// to -> dex -> from
			msg.WithTransfer(to, dex, types.NewCoins(toAsset.Add(feeFromTo))),
			msg.WithTransfer(dex, from, types.NewCoins(toAsset)),
			msg.WithData(Cdc(), &MsgDexDealData{
				Dex:     dex,
				ExtData: ext,
			}),
		),
	}
}

func (msg MsgDexDeal) GetDealByDex() (AccountID, Coins, AccountID, Coins) {
	trs := msg.Transfers
	return trs[0].From, trs[0].Amount, trs[2].From, trs[2].Amount
}

func (msg MsgDexDeal) GetDealFeeByDex() (Coins, Coins) {
	trs := msg.Transfers
	return trs[0].Amount.Sub(trs[1].Amount), trs[2].Amount.Sub(trs[3].Amount)
}

func (msg MsgDexDeal) GetData() (MsgDexDealData, error) {
	res := MsgDexDealData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgDexDealData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

func (msg MsgDexDeal) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.Dex.Empty() {
		return sdkerrors.Wrap(types.ErrKuMsgAccountIDNil, "dex accountID empty")
	}

	// check transfer
	trs := msg.Transfers

	if len(trs) != 4 {
		return sdkerrors.Wrap(types.ErrKuMsgTransferError, "transfer in deal should be 4 trans")
	}

	// a -> dex -> b
	if trs[0].From.Eq(data.Dex) {
		return sdkerrors.Wrap(types.ErrKuMsgTransferError, "dex should not be from in deal0")
	}

	if !trs[0].To.Eq(data.Dex) {
		return sdkerrors.Wrap(types.ErrKuMsgTransferError, "dex should be to in deal0")
	}

	if !trs[1].From.Eq(data.Dex) {
		return sdkerrors.Wrap(types.ErrKuMsgTransferError, "dex should be from in deal1")
	}

	if trs[1].To.Eq(data.Dex) {
		return sdkerrors.Wrap(types.ErrKuMsgTransferError, "dex should not be to in deal1")
	}

	if !trs[0].Amount.IsAllGTE(trs[1].Amount) {
		return sdkerrors.Wrap(types.ErrKuMsgTransferError, "fee should larger than 0")
	}

	// b -> dex -> a
	if trs[2].From.Eq(data.Dex) {
		return sdkerrors.Wrap(types.ErrKuMsgTransferError, "dex should not be from in deal2")
	}

	if !trs[2].To.Eq(data.Dex) {
		return sdkerrors.Wrap(types.ErrKuMsgTransferError, "dex should be to in deal2")
	}

	if !trs[3].From.Eq(data.Dex) {
		return sdkerrors.Wrap(types.ErrKuMsgTransferError, "dex should be from in deal3")
	}

	if trs[3].To.Eq(data.Dex) {
		return sdkerrors.Wrap(types.ErrKuMsgTransferError, "dex should not be to in deal3")
	}

	if !trs[2].Amount.IsAllGTE(trs[3].Amount) {
		return sdkerrors.Wrap(types.ErrKuMsgTransferError, "fee should larger than 2")
	}

	return nil
}

// GetSigners
func (msg MsgDexDeal) GetSigners() []sdk.AccAddress {
	return msg.Auth
}
