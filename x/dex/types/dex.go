package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Dex model
type Dex struct {
	Creator     Name   `json:"creator" yaml:"creator"`         // Creator
	Staking     Coins  `json:"staking" yaml:"staking"`         // Dex Staking
	Description string `json:"description" yaml:"description"` // Dex Description
	Number      uint64 `json:"number" yaml:"number"`           // Dex number
	Sequence    uint64 `json:"sequence" yaml:"sequence"`       // Dex sequence
	//DestroyFlag  bool       `json:"destroy_flag" yaml:"destroy_flag"`   // Dex destroy flag
	CurrencyList []Currency `json:"currency_list" yaml:"currency_list"` // Dex currency list
}

// NewDex creator a new dex
func NewDex(creator Name, staking Coins, description string) *Dex {
	return &Dex{
		Creator:     creator,
		Staking:     staking,
		Description: description,
	}
}

// WithNumber set dex number
func (d *Dex) WithNumber(n uint64) *Dex {
	d.Number = n
	return d
}

// WithDestroyFlag set dex destroy flag
//func (d *Dex) WithDestroyFlag() *Dex {
//	d.DestroyFlag = true
//	return d
//}

// CanDestroy check whether dex can destroy
func (d *Dex) CanDestroy(ctx *sdk.Context) (ok bool) {
	// TODO
	ok = true
	return
}

// WithCurrency dex add currency
func (d *Dex) WithCurrency(currency *Currency) (dex *Dex, ok bool) {
	if nil == currency || !currency.Validate() {
		return
	}
	dex = d
	for i := 0; i < len(d.CurrencyList); i++ {
		if d.CurrencyList[i].Base.Code == currency.Base.Code &&
			d.CurrencyList[i].Quote.Code == currency.Quote.Code {
			return
		}
	}
	d.CurrencyList = append(d.CurrencyList, *currency)
	ok = true
	return
}

// Currency get dex currency
func (d *Dex) Currency(baseCode, quoteCode string) (currency Currency, ok bool) {
	if 0 >= len(baseCode) || 0 >= len(quoteCode) {
		return
	}
	for i := 0; i < len(d.CurrencyList); i++ {
		if baseCode == d.CurrencyList[i].Base.Code &&
			quoteCode == d.CurrencyList[i].Quote.Code {
			ok = true
			currency = d.CurrencyList[i]
			return
		}
	}
	return
}

// UpdateCurrency update currency
func (d *Dex) UpdateCurrency(baseCode, quoteCode string,
	currency *Currency) (ok bool) {
	if 0 >= len(baseCode) || 0 >= len(quoteCode) || nil == currency || !currency.Validate() {
		return
	}
	for i := 0; i < len(d.CurrencyList); i++ {
		if baseCode == d.CurrencyList[i].Base.Code &&
			quoteCode == d.CurrencyList[i].Quote.Code {
			ok = true
			d.CurrencyList[i] = *currency
			return
		}
	}
	return
}

// DeleteCurrency delete currency
func (d *Dex) DeleteCurrency(baseCode, quoteCode string) (ok bool) {
	if 0 >= len(baseCode) || 0 >= len(quoteCode) {
		return
	}
	for i := 0; i < len(d.CurrencyList); i++ {
		if baseCode == d.CurrencyList[i].Base.Code &&
			quoteCode == d.CurrencyList[i].Quote.Code {
			ok = true
			d.CurrencyList = append(d.CurrencyList[:i], d.CurrencyList[i+1:]...)
			return
		}
	}
	return
}
