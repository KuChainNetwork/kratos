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

// Equals check whether self equals other
func (object *CurrencyBase) Equal(other *CurrencyBase) bool {
	return object.Code == other.Code &&
		object.Name == other.Name &&
		object.FullName == other.FullName &&
		object.IconUrl == other.IconUrl &&
		object.TxUrl == other.TxUrl
}

// Empty whether all members are invalid
func (object *CurrencyBase) Empty(checkCode bool) bool {
	return (checkCode && 0 >= len(object.Code)) &&
		0 >= len(object.Name) &&
		0 >= len(object.FullName) &&
		0 >= len(object.IconUrl) &&
		0 >= len(object.TxUrl)
}

// BaseCurrency
type BaseCurrency struct {
	CurrencyBase
}

// Validate validate
func (object *BaseCurrency) Validate() bool {
	return object.CurrencyBase.Validate()
}

// Equals check whether self equals other
func (object *BaseCurrency) Equal(other *BaseCurrency) bool {
	return object.CurrencyBase.Equal(&other.CurrencyBase)
}

// QuoteCurrency
type QuoteCurrency struct {
	CurrencyBase
}

// Validate validate
func (object *QuoteCurrency) Validate() bool {
	return object.CurrencyBase.Validate()
}

// Equals check whether self equals other
func (object *QuoteCurrency) Equal(other *QuoteCurrency) bool {
	return object.CurrencyBase.Equal(&other.CurrencyBase)
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

// Equal check whether self equals other
func (object *Currency) Equal(other *Currency) bool {
	return object.Base.Equal(&other.Base) &&
		object.Quote.Equal(&other.Quote) &&
		object.DomainAddress == other.DomainAddress &&
		object.CreateTime.Equal(other.CreateTime) &&
		object.IsPaused == other.IsPaused
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
