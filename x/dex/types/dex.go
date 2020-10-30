package types

import chainTypes "github.com/KuChainNetwork/kuchain/chain/types"

// Dex model
type Dex struct {
	Creator     Name     `json:"creator" yaml:"creator"`         // Creator
	Staking     Coins    `json:"staking" yaml:"staking"`         // Dex Staking
	Description string   `json:"description" yaml:"description"` // Dex Description
	Number      uint64   `json:"number" yaml:"number"`           // Dex number
	Sequence    uint64   `json:"sequence" yaml:"sequence"`       // Dex sequence
	Symbols     []Symbol `json:"symbols" yaml"symbols"`          // Dex symbols
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

// CanDestroy check whether dex can destroy
func (d *Dex) CanDestroy(sumCallback func() chainTypes.Coins) (ok bool) {
	sum := sumCallback()
	return 0 == len(sum) || sum.IsZero()
}

// WithSymbol dex add symbol
func (d *Dex) WithSymbol(symbol *Symbol) (dex *Dex, ok bool) {
	if nil == symbol || !symbol.Validate() {
		return
	}
	dex = d
	for i := 0; i < len(d.Symbols); i++ {
		if d.Symbols[i].Base.Code == symbol.Base.Code &&
			d.Symbols[i].Quote.Code == symbol.Quote.Code {
			return
		}
	}
	d.Symbols = append(d.Symbols, *symbol)
	ok = true
	return
}

// Symbol get dex symbol
func (d *Dex) Symbol(baseCode, quoteCode string) (symbol Symbol, ok bool) {
	if 0 >= len(baseCode) || 0 >= len(quoteCode) {
		return
	}
	for i := 0; i < len(d.Symbols); i++ {
		if baseCode == d.Symbols[i].Base.Code &&
			quoteCode == d.Symbols[i].Quote.Code {
			ok = true
			symbol = d.Symbols[i]
			return
		}
	}
	return
}

// UpdateSymbol update symbol
func (d *Dex) UpdateSymbol(baseCode, quoteCode string,
	symbol *Symbol) (ok bool) {
	if 0 >= len(baseCode) || 0 >= len(quoteCode) || nil == symbol || !symbol.Validate() {
		return
	}
	for i := 0; i < len(d.Symbols); i++ {
		if baseCode == d.Symbols[i].Base.Code &&
			quoteCode == d.Symbols[i].Quote.Code {
			ok = true
			d.Symbols[i] = *symbol
			return
		}
	}
	return
}

// DeleteSymbol delete symbol
func (d *Dex) DeleteSymbol(baseCode, quoteCode string) (ok bool) {
	if 0 >= len(baseCode) || 0 >= len(quoteCode) {
		return
	}
	for i := 0; i < len(d.Symbols); i++ {
		if baseCode == d.Symbols[i].Base.Code &&
			quoteCode == d.Symbols[i].Quote.Code {
			ok = true
			d.Symbols = append(d.Symbols[:i], d.Symbols[i+1:]...)
			return
		}
	}
	return
}
