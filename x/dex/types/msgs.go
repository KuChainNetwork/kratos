package types

import (
	"time"

	"github.com/KuChainNetwork/kuchain/chain/msg"
	"github.com/KuChainNetwork/kuchain/chain/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ types.KuTransfMsg = &MsgCreateDex{}
)

// MsgCreateDex msg for create dex
type MsgCreateDex struct {
	types.KuMsg
}

const MaxDexDescriptorLen = 512

// MsgCreateDexData msg data for create dex
type MsgCreateDexData struct {
	Creator  Name   `json:"creator" yaml:"creator"` // Creator coin creator account name
	Stakings Coins  `json:"staking" yaml:"staking"` // Staking for dex
	Desc     []byte `json:"desc" yaml:"desc"`       // Description
}

func (MsgCreateDexData) Type() types.Name { return types.MustName("create@dex") }

func (m MsgCreateDexData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgCreateDex new create dex msg
func NewMsgCreateDex(auth types.AccAddress, creator types.Name, stakings types.Coins, desc []byte) MsgCreateDex {
	return MsgCreateDex{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithTransfer(types.NewAccountIDFromName(creator), ModuleAccountID, stakings),
			msg.WithData(ModuleCdc, &MsgCreateDexData{
				Creator:  creator,
				Stakings: stakings,
				Desc:     desc,
			}),
		),
	}
}

func (msg MsgCreateDex) GetData() (MsgCreateDexData, error) {
	res := MsgCreateDexData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgCreateDexData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

func (msg MsgCreateDex) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.Creator.Empty() {
		return sdkerrors.Wrap(types.ErrKuMSgNameEmpty, "creator name not empty")
	}

	if data.Stakings.IsAnyNegative() {
		return sdkerrors.Wrap(types.ErrKuMsgCoinsHasNegative, "stakings should be positive")
	}

	if len(data.Desc) >= MaxDexDescriptorLen {
		return ErrDexDescTooLong
	}

	if err := msg.ValidateTransferTo(types.NewAccountIDFromName(data.Creator), ModuleAccountID, data.Stakings); err != nil {
		return sdkerrors.Wrap(err, "create dex stakings transfer error")
	}

	return nil
}

// MsgUpdateDexDescription msg for update dex description
type MsgUpdateDexDescription struct {
	types.KuMsg
}

func (msg MsgUpdateDexDescription) GetData() (MsgUpdateDexDescriptionData, error) {
	res := MsgUpdateDexDescriptionData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgUpdateDexDescriptionData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgUpdateDexDescriptionData msg data for update dex description
type MsgUpdateDexDescriptionData struct {
	Creator Name   `json:"creator" yaml:"creator"` // Creator coin creator account name
	Desc    []byte `json:"desc" yaml:"desc"`       // Description
}

func (MsgUpdateDexDescriptionData) Type() types.Name { return types.MustName("updatedesc@dex") }

func (m MsgUpdateDexDescriptionData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgUpdateDexDescription new update dex description msg
func NewMsgUpdateDexDescription(auth types.AccAddress,
	creator types.Name,
	desc []byte) MsgUpdateDexDescription {
	return MsgUpdateDexDescription{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgUpdateDexDescriptionData{
				Creator: creator,
				Desc:    desc,
			})),
	}
}

func (msg MsgUpdateDexDescription) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.Creator.Empty() {
		return sdkerrors.Wrap(types.ErrKuMSgNameEmpty, "creator name not empty")
	}

	if len(data.Desc) >= MaxDexDescriptorLen {
		return ErrDexDescTooLong
	}

	return nil
}

// MsgDestroyDex msg for delete dex
type MsgDestroyDex struct {
	types.KuMsg
}

func (msg MsgDestroyDex) GetData() (MsgDestroyDexData, error) {
	res := MsgDestroyDexData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgDestroyDexData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgDestroyDexData msg data for delete dex
type MsgDestroyDexData struct {
	Creator Name `json:"creator" yaml:"creator"`
}

func (MsgDestroyDexData) Type() types.Name { return types.MustName("destroy@dex") }

func (m MsgDestroyDexData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgDestroyDex new destroy dex msg
func NewMsgDestroyDex(auth types.AccAddress, creator types.Name) MsgDestroyDex {
	return MsgDestroyDex{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgDestroyDexData{
				Creator: creator,
			})),
	}
}

func (msg MsgDestroyDex) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}

	data, err := msg.GetData()
	if err != nil {
		return err
	}

	if data.Creator.Empty() {
		return sdkerrors.Wrap(types.ErrKuMSgNameEmpty, "creator name not empty")
	}

	return nil
}

// MsgCreateSymbol msg for create dex symbol
type MsgCreateSymbol struct {
	types.KuMsg
}

func (msg MsgCreateSymbol) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	data, err := msg.GetData()
	if err != nil {
		return err
	}
	if data.Creator.Empty() {
		return sdkerrors.Wrap(types.ErrKuMSgNameEmpty, "creator name not empty")
	}
	if !data.Base.Validate() {
		return sdkerrors.Wrap(ErrSymbolBaseInvalid, "base part invalid")
	}
	if !data.Quote.Validate() {
		return sdkerrors.Wrap(ErrSymbolQuoteInvalid, "quote part invalid")
	}
	return nil
}

func (msg MsgCreateSymbol) GetData() (MsgCreateSymbolData, error) {
	res := MsgCreateSymbolData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgCreateSymbolData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgCreateSymbolData msg data for delete dex
type MsgCreateSymbolData struct {
	Creator    Name          `json:"creator" yaml:"creator"`
	Base       BaseCurrency  `json:"base" yaml:"base"`
	Quote      QuoteCurrency `json:"quote" yaml:"quote"`
	CreateTime time.Time     `json:"create_time" yaml:"create_time"`
}

func (MsgCreateSymbolData) Type() types.Name { return types.MustName("create@symbol") }

func (m MsgCreateSymbolData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgCreateSymbol new destroy dex msg
func NewMsgCreateSymbol(auth types.AccAddress,
	creator types.Name,
	base *BaseCurrency,
	quote *QuoteCurrency,
	createTime time.Time,
) MsgCreateSymbol {
	return MsgCreateSymbol{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgCreateSymbolData{
				Creator:    creator,
				Base:       *base,
				Quote:      *quote,
				CreateTime: createTime,
			})),
	}
}

// MsgUpdateSymbol msg for update dex symbol
type MsgUpdateSymbol struct {
	types.KuMsg
}

func (msg MsgUpdateSymbol) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	data, err := msg.GetData()
	if err != nil {
		return err
	}
	if data.Creator.Empty() {
		return sdkerrors.Wrap(types.ErrKuMSgNameEmpty, "creator name not empty")
	}
	if 0 >= len(data.Base.Code) {
		return sdkerrors.Wrap(ErrSymbolBaseCodeEmpty, "base code not empty")
	}
	if 0 >= len(data.Quote.Code) {
		return sdkerrors.Wrap(ErrSymbolQuoteCodeEmpty, "quote code not empty")
	}
	if data.Base.Empty(false) && data.Quote.Empty(false) {
		return sdkerrors.Wrap(ErrSymbolUpdateFieldsInvalid, "update fields not empty at least one")
	}
	return nil
}

func (msg MsgUpdateSymbol) GetData() (MsgUpdateSymbolData, error) {
	res := MsgUpdateSymbolData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgUpdateSymbolData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgUpdateSymbolData msg data for update dex symbol
type MsgUpdateSymbolData struct {
	Creator Name          `json:"creator" yaml:"creator"`
	Base    BaseCurrency  `json:"base" yaml:"base"`
	Quote   QuoteCurrency `json:"quote" yaml:"quote"`
}

func (MsgUpdateSymbolData) Type() types.Name { return types.MustName("update@symbol") }

func (m MsgUpdateSymbolData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgUpdateSymbol new update dex symbol msg
func NewMsgUpdateSymbol(auth types.AccAddress,
	creator types.Name,
	base *BaseCurrency,
	quote *QuoteCurrency,
) MsgUpdateSymbol {
	return MsgUpdateSymbol{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgUpdateSymbolData{
				Creator: creator,
				Base:    *base,
				Quote:   *quote,
			})),
	}
}

// MsgPauseSymbol msg for pause dex symbol
type MsgPauseSymbol struct {
	types.KuMsg
}

func (msg MsgPauseSymbol) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	data, err := msg.GetData()
	if err != nil {
		return err
	}
	if data.Creator.Empty() {
		return sdkerrors.Wrap(types.ErrKuMSgNameEmpty, "creator name not empty")
	}
	if 0 >= len(data.BaseCode) {
		return sdkerrors.Wrap(ErrSymbolBaseCodeEmpty, "base code not empty")
	}
	if 0 >= len(data.QuoteCode) {
		return sdkerrors.Wrap(ErrSymbolQuoteCodeEmpty, "quote code not empty")
	}
	return nil
}

func (msg MsgPauseSymbol) GetData() (MsgPauseSymbolData, error) {
	res := MsgPauseSymbolData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgPauseSymbolData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgPauseSymbolData msg data for pause dex symbol
type MsgPauseSymbolData struct {
	Creator   Name   `json:"creator" yaml:"creator"`
	BaseCode  string `json:"base_code" yaml:"base_code"`
	QuoteCode string `json:"quote_code" yaml:"quote_code"`
}

func (MsgPauseSymbolData) Type() types.Name { return types.MustName("pause@symbol") }

func (m MsgPauseSymbolData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgPauseSymbol new pause dex symbol msg
func NewMsgPauseSymbol(auth types.AccAddress,
	creator types.Name,
	baseCode,
	quoteCode string,
) MsgPauseSymbol {
	return MsgPauseSymbol{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgPauseSymbolData{
				Creator:   creator,
				BaseCode:  baseCode,
				QuoteCode: quoteCode,
			})),
	}
}

// MsgRestoreSymbol msg for restore dex symbol
type MsgRestoreSymbol struct {
	types.KuMsg
}

func (msg MsgRestoreSymbol) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	data, err := msg.GetData()
	if err != nil {
		return err
	}
	if data.Creator.Empty() {
		return sdkerrors.Wrap(types.ErrKuMSgNameEmpty, "creator name not empty")
	}
	if 0 >= len(data.BaseCode) {
		return sdkerrors.Wrap(ErrSymbolBaseCodeEmpty, "base code not empty")
	}
	if 0 >= len(data.QuoteCode) {
		return sdkerrors.Wrap(ErrSymbolQuoteCodeEmpty, "quote code not empty")
	}
	return nil
}

func (msg MsgRestoreSymbol) GetData() (MsgRestoreSymbolData, error) {
	res := MsgRestoreSymbolData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgRestoreSymbolData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgRestoreSymbolData msg data for restore dex symbol
type MsgRestoreSymbolData struct {
	Creator   Name   `json:"creator" yaml:"creator"`
	BaseCode  string `json:"base_code" yaml:"base_code"`
	QuoteCode string `json:"quote_code" yaml:"quote_code"`
}

func (MsgRestoreSymbolData) Type() types.Name { return types.MustName("restore@symbol") }

func (m MsgRestoreSymbolData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgRestoreSymbol new restore dex symbol msg
func NewMsgRestoreSymbol(auth types.AccAddress,
	creator types.Name,
	baseCode,
	quoteCode string,
) MsgRestoreSymbol {
	return MsgRestoreSymbol{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgRestoreSymbolData{
				Creator:   creator,
				BaseCode:  baseCode,
				QuoteCode: quoteCode,
			})),
	}
}

// MsgShutdownSymbol msg for shutdown dex symbol
type MsgShutdownSymbol struct {
	types.KuMsg
}

func (msg MsgShutdownSymbol) ValidateBasic() error {
	if err := msg.KuMsg.ValidateTransfer(); err != nil {
		return err
	}
	data, err := msg.GetData()
	if err != nil {
		return err
	}
	if data.Creator.Empty() {
		return sdkerrors.Wrap(types.ErrKuMSgNameEmpty, "creator name not empty")
	}
	if 0 >= len(data.BaseCode) {
		return sdkerrors.Wrap(ErrSymbolBaseCodeEmpty, "base code not empty")
	}
	if 0 >= len(data.QuoteCode) {
		return sdkerrors.Wrap(ErrSymbolQuoteCodeEmpty, "quote code not empty")
	}
	return nil
}

func (msg MsgShutdownSymbol) GetData() (MsgShutdownSymbolData, error) {
	res := MsgShutdownSymbolData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgShutdownSymbolData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgUpdateSymbolData msg data for update dex symbol
type MsgShutdownSymbolData struct {
	Creator   Name   `json:"creator" yaml:"creator"`
	BaseCode  string `json:"base_code" yaml:"base_code"`
	QuoteCode string `json:"quote_code" yaml:"quote_code"`
}

func (MsgShutdownSymbolData) Type() types.Name { return types.MustName("shutdown@symbol") }

func (m MsgShutdownSymbolData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgUpdateSymbol new update dex symbol msg
func NewMsgShutdownSymbol(auth types.AccAddress,
	creator types.Name,
	baseCode,
	quoteCode string,
) MsgShutdownSymbol {
	return MsgShutdownSymbol{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgShutdownSymbolData{
				Creator:   creator,
				BaseCode:  baseCode,
				QuoteCode: quoteCode,
			})),
	}
}
