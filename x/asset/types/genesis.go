package types

import (
	"encoding/json"
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/types/coin"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type GenesisLocks interface {
	GetID() AccountID
	LockedByHeight() []LockedCoins
}

type BaseGenesisLocks struct {
	ID      AccountID     `json:"id" yaml:"id"`
	Lockeds []LockedCoins `json:"lockeds" yaml:"lockeds"`
}

func NewBaseGenesisLocks(id AccountID, lockeds []LockedCoins) BaseGenesisLocks {
	return BaseGenesisLocks{
		ID:      id,
		Lockeds: lockeds,
	}
}

func (g BaseGenesisLocks) GetID() AccountID { return g.ID }

func (g BaseGenesisLocks) LockedByHeight() []LockedCoins {
	return g.Lockeds
}

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	GenesisAssets     []GenesisAsset `json:"genesisAssets"`
	GenesisCoins      []GenesisCoin  `json:"genesisCoins"`
	GenesisCoinPowers []GenesisAsset `json:"genesisCoinPowers"`
	GenesisLocks      []GenesisLocks `json:"genesisLocks"`
	GenesisLockAssets []GenesisAsset `json:"genesisLockAssets"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		GenesisAssets:     make([]GenesisAsset, 0),
		GenesisCoins:      make([]GenesisCoin, 0),
		GenesisCoinPowers: make([]GenesisAsset, 0),
		GenesisLocks:      make([]GenesisLocks, 0),
		GenesisLockAssets: make([]GenesisAsset, 0),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	res := NewGenesisState()

	return res
}

func ValidateGenesis(g GenesisState) error {
	// check account numbers
	// TODO: check genesis

	return nil
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func (g GenesisState) ValidateGenesis(bz json.RawMessage) error {
	gs := DefaultGenesisState()
	if err := ModuleCdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	return ValidateGenesis(gs)
}

// GenesisAsset gensis asset for accountID
type GenesisAsset interface {
	GetID() AccountID
	GetCoins() Coins
}

// GenesisCoin gensis coin type
type GenesisCoin interface {
	GetCreator() Name
	GetSymbol() Name
	GetSupply() Coin
	GetMaxSupply() Coin
	GetDescription() string

	Validate() error
}

type BaseGenesisAsset struct {
	ID    AccountID `json:"id"`
	Coins Coins     `json:"coins"`
}

func NewGenesisAssetByCoins(id AccountID, coins Coins) BaseGenesisAsset {
	return BaseGenesisAsset{
		ID:    id,
		Coins: coins,
	}
}

func NewGenesisAsset(id AccountID, coins ...Coin) BaseGenesisAsset {
	return BaseGenesisAsset{
		ID:    id,
		Coins: coins,
	}
}

func (g BaseGenesisAsset) GetID() AccountID { return g.ID }

func (g BaseGenesisAsset) GetCoins() Coins { return g.Coins }

// GensisAssetCoin
type BaseGensisAssetCoin struct {
	Stat        CoinStat `json:"stat"`
	Description string   `json:"description"`
}

func NewGenesisCoin(stat *CoinStat, description []byte) BaseGensisAssetCoin {
	return BaseGensisAssetCoin{
		Stat:        *stat,
		Description: string(description),
	}
}

// Validate imp GenesisCoin
func (g BaseGensisAssetCoin) Validate() error {
	if len(g.Description) >= MaxDescriptionLength {
		return fmt.Errorf("genesis coin description too length")
	}

	denom := CoinDenom(g.Stat.Creator, g.Stat.Symbol)

	if denom != g.Stat.MaxSupply.Denom {
		return fmt.Errorf("genesis max supply coin denom error")
	}

	if err := coin.ValidateDenom(denom); err != nil {
		return sdkerrors.Wrapf(ErrAssetDenom, "denom %s", denom)
	}

	return nil
}

// GetCreator imp GenesisCoin
func (g BaseGensisAssetCoin) GetCreator() Name { return g.Stat.Creator }

// GetSymbol imp GenesisCoin
func (g BaseGensisAssetCoin) GetSymbol() Name { return g.Stat.Symbol }

// GetMaxSupply imp GenesisCoin
func (g BaseGensisAssetCoin) GetMaxSupply() Coin { return g.Stat.MaxSupply }

// GetDescription imp GenesisCoin
func (g BaseGensisAssetCoin) GetDescription() string { return g.Description }

// GetSupply imp GenesisCoin
func (g BaseGensisAssetCoin) GetSupply() Coin { return g.Stat.Supply }
