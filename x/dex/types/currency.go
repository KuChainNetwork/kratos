package types

import "time"

// CurrencyBase
type CurrencyBase struct {
	Code     string `json:"code" yaml:"code"`
	Name     string `json:"name" yaml:"name"`
	FullName string `json:"full_name" yaml:"full_name"`
	IconUrl  string `json:"icon_url" yaml:"icon_url"`
	TxUrl    string `json:"tx_url" yaml:"tx_url"`
}

// Validate validate
func (object *CurrencyBase) Validate() bool {
	return 0 < len(object.Code) &&
		0 < len(object.Name) &&
		0 < len(object.FullName) &&
		0 < len(object.IconUrl) &&
		0 < len(object.TxUrl)
}

// BaseCurrency
type BaseCurrency struct {
	CurrencyBase
}

// Validate validate
func (object *BaseCurrency) Validate() bool {
	return object.CurrencyBase.Validate()
}

// QuoteCurrency
type QuoteCurrency struct {
	CurrencyBase
}

// Validate validate
func (object *QuoteCurrency) Validate() bool {
	return object.CurrencyBase.Validate()
}

// Currency
type Currency struct {
	Base          BaseCurrency  `json:"base" yaml:"base"`
	Quote         QuoteCurrency `json:"quote" yaml:"quote"`
	DomainAddress string        `json:"domain_address" yaml:"domain_address"`
	CreateTime    time.Time     `json:"create_time" yaml:"create_time"`
	IsPaused      bool          `json:"is_paused" yaml:"is_paused"`
}

// WithBase set base currency param
func (object *Currency) WithBase(base *BaseCurrency) *Currency {
	object.Base = *base
	return object
}

// WithQuote set quote currency param
func (object *Currency) WithQuote(quote *QuoteCurrency) *Currency {
	object.Quote = *quote
	return object
}

// WithDomainAddress set domain address
func (object *Currency) WithDomainAddress(address string) *Currency {
	object.DomainAddress = address
	return object
}

// WithCreateTime set create time
func (object *Currency) WithCreateTime(createTime time.Time) *Currency {
	object.CreateTime = createTime
	return object
}

// WithPaused set pause flag
func (object *Currency) WithPaused(paused bool) *Currency {
	object.IsPaused = paused
	return object
}

// Validate validate
func (object *Currency) Validate() bool {
	return object.Quote.Validate() &&
		object.Base.Validate() &&
		0 < len(object.DomainAddress) &&
		!object.CreateTime.IsZero()
}

// Paused check currency pause flag
func (object *Currency) Paused() bool {
	return object.IsPaused
}

// NewEmptyCurrency new empty currency
func NewEmptyCurrency() *Currency {
	return &Currency{}
}

// NewCurrency new currency with base currency and quote currency
func NewCurrency(base *BaseCurrency, quote *QuoteCurrency,
	domainAddress string) *Currency {
	return &Currency{Base: *base, Quote: *quote, DomainAddress: domainAddress}
}
