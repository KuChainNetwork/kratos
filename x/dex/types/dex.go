package types

import chainTypes "github.com/KuChainNetwork/kuchain/chain/types"

// Dex model
type Dex struct {
	Creator     Name     `json:"creator" yaml:"creator"`         // Creator
	Staking     Coins    `json:"staking" yaml:"staking"`         // Dex Staking
	Description string   `json:"description" yaml:"description"` // Dex Description
	Number      uint64   `json:"number" yaml:"number"`           // Dex number
	Sequence    uint64   `json:"sequence" yaml:"sequence"`       // Dex sequence
	Symbols     []Symbol `json:"symbols" yaml:"symbols"`         // Dex symbols
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
	return len(sum) == 0 || sum.IsZero()
}

// WithSymbol dex add symbol
func (d *Dex) WithSymbol(symbol *Symbol) (dex *Dex, ok bool) {
	if nil == symbol || !symbol.Validate() {
		return
	}
	dex = d
	for i := 0; i < len(d.Symbols); i++ {
		s := &d.Symbols[i]
		if s.Base.Creator == symbol.Base.Creator &&
			s.Base.Code == symbol.Base.Code &&
			s.Quote.Creator == symbol.Quote.Creator &&
			s.Quote.Code == symbol.Quote.Code {
			return
		}

		if s.Base.Creator == symbol.Quote.Creator &&
			s.Base.Code == symbol.Quote.Code &&
			s.Quote.Creator == symbol.Base.Creator &&
			s.Quote.Code == symbol.Base.Code {
			return
		}
	}
	d.Symbols = append(d.Symbols, *symbol)
	ok = true
	return
}

// Symbol get dex symbol
func (d *Dex) Symbol(baseCreator, baseCode, quoteCreator, quoteCode string) (symbol Symbol, ok bool) {
	if 0 >= len(baseCode) || 0 >= len(quoteCode) {
		return
	}
	for i := 0; i < len(d.Symbols); i++ {
		s := &d.Symbols[i]
		if baseCreator == s.Base.Creator &&
			baseCode == s.Base.Code &&
			quoteCreator == s.Quote.Creator &&
			quoteCode == s.Quote.Code {
			ok = true
			symbol = d.Symbols[i]
			return
		}
	}
	return
}

// UpdateSymbol update symbol
func (d *Dex) UpdateSymbol(baseCreator, baseCode, quoteCreator, quoteCode string,
	symbol *Symbol) (ok bool) {
	if 0 >= len(baseCode) || 0 >= len(quoteCode) || nil == symbol || !symbol.Validate() {
		return
	}
	for i := 0; i < len(d.Symbols); i++ {
		s := &d.Symbols[i]
		if baseCreator == s.Base.Creator &&
			baseCode == s.Base.Code &&
			quoteCreator == s.Quote.Creator &&
			quoteCode == s.Quote.Code {
			ok = true
			*s = *symbol
			return
		}
	}
	return
}

// DeleteSymbol delete symbol
func (d *Dex) DeleteSymbol(baseCreator, baseCode, quoteCreator, quoteCode string) bool {
	if 0 >= len(baseCreator) || 0 >= len(baseCode) || 0 >= len(quoteCreator) || 0 >= len(quoteCode) {
		return false
	}

	for i := 0; i < len(d.Symbols); i++ {
		s := &d.Symbols[i]
		if baseCreator == s.Base.Creator &&
			baseCode == s.Base.Code &&
			quoteCreator == s.Quote.Creator &&
			quoteCode == s.Quote.Code {
			d.Symbols = append(d.Symbols[:i], d.Symbols[i+1:]...)
			return true
		}
	}

	return false
}
