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

// MsgCreateCurrency msg for create dex currency
type MsgCreateCurrency struct {
	types.KuMsg
}

func (msg MsgCreateCurrency) ValidateBasic() error {
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

func (msg MsgCreateCurrency) GetData() (MsgCreateCurrencyData, error) {
	res := MsgCreateCurrencyData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgCreateCurrencyData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgCreateCurrencyData msg data for delete dex
type MsgCreateCurrencyData struct {
	Creator       Name          `json:"creator" yaml:"creator"`
	Base          BaseCurrency  `json:"base" yaml:"base"`
	Quote         QuoteCurrency `json:"quote" yaml:"quote"`
	DomainAddress string        `json:"domain_address" yaml:"domain_address"`
	CreateTime    time.Time     `json:"create_time" yaml:"create_time"`
}

func (MsgCreateCurrencyData) Type() types.Name { return types.MustName("create@currency") }

func (m MsgCreateCurrencyData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgCreateCurrency new destroy dex msg
func NewMsgCreateCurrency(auth types.AccAddress,
	creator types.Name,
	base *BaseCurrency,
	quote *QuoteCurrency,
	domainAddress string,
	createTime time.Time,
) MsgCreateCurrency {
	return MsgCreateCurrency{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgCreateCurrencyData{
				Creator:       creator,
				Base:          *base,
				Quote:         *quote,
				DomainAddress: domainAddress,
				CreateTime:    createTime,
			})),
	}
}

// MsgUpdateCurrency msg for update dex currency
type MsgUpdateCurrency struct {
	types.KuMsg
}

func (msg MsgUpdateCurrency) ValidateBasic() error {
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

func (msg MsgUpdateCurrency) GetData() (MsgUpdateCurrencyData, error) {
	res := MsgUpdateCurrencyData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgUpdateCurrencyData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgUpdateCurrencyData msg data for update dex currency
type MsgUpdateCurrencyData struct {
	Creator Name          `json:"creator" yaml:"creator"`
	Base    BaseCurrency  `json:"base" yaml:"base"`
	Quote   QuoteCurrency `json:"quote" yaml:"quote"`
}

func (MsgUpdateCurrencyData) Type() types.Name { return types.MustName("update@currency") }

func (m MsgUpdateCurrencyData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgUpdateCurrency new update dex currency msg
func NewMsgUpdateCurrency(auth types.AccAddress,
	creator types.Name,
	base *BaseCurrency,
	quote *QuoteCurrency,
) MsgUpdateCurrency {
	return MsgUpdateCurrency{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgUpdateCurrencyData{
				Creator: creator,
				Base:    *base,
				Quote:   *quote,
			})),
	}
}

// MsgPauseCurrency msg for pause dex currency
type MsgPauseCurrency struct {
	types.KuMsg
}

func (msg MsgPauseCurrency) ValidateBasic() error {
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

func (msg MsgPauseCurrency) GetData() (MsgPauseCurrencyData, error) {
	res := MsgPauseCurrencyData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgPauseCurrencyData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgPauseCurrencyData msg data for pause dex currency
type MsgPauseCurrencyData struct {
	Creator   Name   `json:"creator" yaml:"creator"`
	BaseCode  string `json:"base_code" yaml:"base_code"`
	QuoteCode string `json:"quote_code" yaml:"quote_code"`
}

func (MsgPauseCurrencyData) Type() types.Name { return types.MustName("pause@currency") }

func (m MsgPauseCurrencyData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgPauseCurrency new pause dex currency msg
func NewMsgPauseCurrency(auth types.AccAddress,
	creator types.Name,
	baseCode,
	quoteCode string,
) MsgPauseCurrency {
	return MsgPauseCurrency{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgPauseCurrencyData{
				Creator:   creator,
				BaseCode:  baseCode,
				QuoteCode: quoteCode,
			})),
	}
}

// MsgRestoreCurrency msg for restore dex currency
type MsgRestoreCurrency struct {
	types.KuMsg
}

func (msg MsgRestoreCurrency) ValidateBasic() error {
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

func (msg MsgRestoreCurrency) GetData() (MsgRestoreCurrencyData, error) {
	res := MsgRestoreCurrencyData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgRestoreCurrencyData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgRestoreCurrencyData msg data for restore dex currency
type MsgRestoreCurrencyData struct {
	Creator   Name   `json:"creator" yaml:"creator"`
	BaseCode  string `json:"base_code" yaml:"base_code"`
	QuoteCode string `json:"quote_code" yaml:"quote_code"`
}

func (MsgRestoreCurrencyData) Type() types.Name { return types.MustName("restore@currency") }

func (m MsgRestoreCurrencyData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgRestoreCurrency new restore dex currency msg
func NewMsgRestoreCurrency(auth types.AccAddress,
	creator types.Name,
	baseCode,
	quoteCode string,
) MsgRestoreCurrency {
	return MsgRestoreCurrency{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgRestoreCurrencyData{
				Creator:   creator,
				BaseCode:  baseCode,
				QuoteCode: quoteCode,
			})),
	}
}

// MsgShutdownCurrency msg for shutdown dex currency
type MsgShutdownCurrency struct {
	types.KuMsg
}

func (msg MsgShutdownCurrency) ValidateBasic() error {
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

func (msg MsgShutdownCurrency) GetData() (MsgShutdownCurrencyData, error) {
	res := MsgShutdownCurrencyData{}
	if err := msg.UnmarshalData(Cdc(), &res); err != nil {
		return MsgShutdownCurrencyData{}, sdkerrors.Wrapf(types.ErrKuMsgDataUnmarshal, "%s", err.Error())
	}
	return res, nil
}

// MsgUpdateCurrencyData msg data for update dex currency
type MsgShutdownCurrencyData struct {
	Creator   Name   `json:"creator" yaml:"creator"`
	BaseCode  string `json:"base_code" yaml:"base_code"`
	QuoteCode string `json:"quote_code" yaml:"quote_code"`
}

func (MsgShutdownCurrencyData) Type() types.Name { return types.MustName("shutdown@currency") }

func (m MsgShutdownCurrencyData) Sender() AccountID {
	return types.NewAccountIDFromName(m.Creator)
}

// NewMsgUpdateCurrency new update dex currency msg
func NewMsgShutdownCurrency(auth types.AccAddress,
	creator types.Name,
	baseCode,
	quoteCode string,
) MsgShutdownCurrency {
	return MsgShutdownCurrency{
		KuMsg: *msg.MustNewKuMsg(RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(ModuleCdc, &MsgShutdownCurrencyData{
				Creator:   creator,
				BaseCode:  baseCode,
				QuoteCode: quoteCode,
			})),
	}
}
