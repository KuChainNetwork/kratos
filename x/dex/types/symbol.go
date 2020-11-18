package types

import "time"

// CurrencyBase
type CurrencyBase struct {
	Creator  string `json:"creator" yaml:"creator"`
	Code     string `json:"code" yaml:"code"`
	Name     string `json:"name" yaml:"name"`
	FullName string `json:"full_name" yaml:"full_name"`
	IconURL  string `json:"icon_url" yaml:"icon_url"`
	TxURL    string `json:"tx_url" yaml:"tx_url"`
}

// Validate validate
func (object *CurrencyBase) Validate() bool {
	return 0 < len(object.Creator) &&
		0 < len(object.Code) &&
		0 < len(object.Name) &&
		0 < len(object.FullName) &&
		0 < len(object.IconURL) &&
		0 < len(object.TxURL)
}

// Equals check whether self equals other
func (object *CurrencyBase) Equal(other *CurrencyBase) bool {
	return object.Creator == other.Creator &&
		object.Code == other.Code &&
		object.Name == other.Name &&
		object.FullName == other.FullName &&
		object.IconURL == other.IconURL &&
		object.TxURL == other.TxURL
}

// Empty whether all members are invalid
func (object *CurrencyBase) Empty(checkCode bool) bool {
	return (checkCode && 0 >= len(object.Creator) || 0 >= len(object.Code)) &&
		0 >= len(object.Name) &&
		0 >= len(object.FullName) &&
		0 >= len(object.IconURL) &&
		0 >= len(object.TxURL)
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

// Symbol
type Symbol struct {
	Base       BaseCurrency  `json:"base" yaml:"base"`
	Quote      QuoteCurrency `json:"quote" yaml:"quote"`
	Height     int64         `json:"height" yaml:"height"`
	CreateTime time.Time     `json:"create_time" yaml:"create_time"`
	IsPaused   bool          `json:"is_paused" yaml:"is_paused"`
}

// WithBase set base currency param
func (object *Symbol) WithBase(base *BaseCurrency) *Symbol {
	object.Base = *base
	return object
}

// WithQuote set quote currency param
func (object *Symbol) WithQuote(quote *QuoteCurrency) *Symbol {
	object.Quote = *quote
	return object
}

// WithHeight set block height
func (object *Symbol) WithHeight(height int64) *Symbol {
	object.Height = height
	return object
}

// WithCreateTime set create time
func (object *Symbol) WithCreateTime(createTime time.Time) *Symbol {
	object.CreateTime = createTime
	return object
}

// WithPaused set pause flag
func (object *Symbol) WithPaused(paused bool) *Symbol {
	object.IsPaused = paused
	return object
}

// Validate validate
func (object *Symbol) Validate() bool {
	return object.Quote.Validate() &&
		object.Base.Validate() &&
		0 < object.Height &&
		!object.CreateTime.IsZero()
}

// Paused check currency pause flag
func (object *Symbol) Paused() bool {
	return object.IsPaused
}

// Equal check whether self equals other
func (object *Symbol) Equal(other *Symbol) bool {
	return object.Base.Equal(&other.Base) &&
		object.Quote.Equal(&other.Quote) &&
		object.Height == other.Height &&
		object.CreateTime.Equal(other.CreateTime) &&
		object.IsPaused == other.IsPaused
}

// NewEmptySymbol new empty currency
func NewEmptySymbol() *Symbol {
	return &Symbol{}
}

// NewSymbol new currency with base currency and quote currency
func NewSymbol(base *BaseCurrency, quote *QuoteCurrency) *Symbol {
	return &Symbol{Base: *base, Quote: *quote}
}
