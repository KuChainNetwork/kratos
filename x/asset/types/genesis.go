package types

import (
	"encoding/json"
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/types/coin"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	GenesisAssets     []GenesisAsset `json:"genesisAssets"`
	GenesisCoins      []GenesisCoin  `json:"genesisCoins"`
	GenesisCoinPowers []GenesisAsset `json:"genesisCoinPowers"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		GenesisAssets:     make([]GenesisAsset, 0),
		GenesisCoins:      make([]GenesisCoin, 0),
		GenesisCoinPowers: make([]GenesisAsset, 0),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	res := NewGenesisState()

	return res
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func (g GenesisState) ValidateGenesis(bz json.RawMessage) error {
	gs := DefaultGenesisState()
	if err := ModuleCdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	// TODO: check genesis

	return nil
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
