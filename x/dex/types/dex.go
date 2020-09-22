package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Dex model
type Dex struct {
	Creator     Name   `json:"creator" yaml:"creator"`           // Creator
	Staking     Coins  `json:"staking" yaml:"staking"`           // Dex Staking
	Description string `json:"description" yaml:"description"`   // Dex Description
	Number      uint64 `json:"number" yaml:"number"`             // Dex number
	Sequence    uint64 `json:"sequence" yaml:"sequence"`         // Dex sequence
	DestroyFlag bool   `json:"destroy_flag" yaml:"destroy_flag"` // Dex destroy flag
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
func (d *Dex) WithDestroyFlag() *Dex {
	d.DestroyFlag = true
	return d
}

// CanDestroy check whether dex can destroy
func (d *Dex) CanDestroy(ctx *sdk.Context) (ok bool) {
	// TODO
	ok = true
	return
}
