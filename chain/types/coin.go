package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type (
	Coin     = sdk.Coin
	Coins    = sdk.Coins
	Int      = sdk.Int
	DecCoin  = sdk.DecCoin
	DecCoins = sdk.DecCoins
)

var (
	NewCoin       = sdk.NewCoin
	NewInt        = sdk.NewInt
	ParseCoin     = sdk.ParseCoin
	ParseCoins    = sdk.ParseCoins
	NewDecCoin    = sdk.NewDecCoin
	NewDecCoins   = sdk.NewDecCoins
	ParseDecCoins = sdk.ParseDecCoins
)

const (
	coinDenomSeq = "/"
	coinDenomFmt = string("%s") + coinDenomSeq + "%s"
)

// CoinDenom get denom for coin, in kuchain, all denom for a coin is by creator and symbol
func CoinDenom(creator, symbol Name) string {
	if creator.Empty() {
		return symbol.String()
	}
	return fmt.Sprintf(coinDenomFmt, creator.String(), symbol.String())
}

// CoinAccountsFromDenom get creator and symbol from denom
func CoinAccountsFromDenom(denom string) (Name, Name, error) {
	strs := strings.Split(denom, coinDenomSeq)
	if len(strs) == 1 {
		// no `/`
		symbol, err := NewName(strs[0])
		if err != nil {
			return Name{}, Name{}, sdkerrors.Wrap(err, "parse root symbol name error")
		}
		return Name{}, symbol, nil
	}

	if len(strs) != 2 {
		return Name{}, Name{}, fmt.Errorf("denom format error by no '/'")
	}

	creator, err := NewName(strs[0])
	if err != nil {
		return Name{}, Name{}, sdkerrors.Wrap(err, "parse creator name error")
	}

	symbol, err := NewName(strs[1])
	if err != nil {
		return Name{}, Name{}, sdkerrors.Wrap(err, "parse symbol name error")
	}

	return creator, symbol, nil
}
