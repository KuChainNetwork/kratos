package types

import (
	"fmt"
	"strings"

	"github.com/KuChainNetwork/kuchain/chain/constants/keys"
	"github.com/KuChainNetwork/kuchain/chain/types/coin"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type (
	Coin     = coin.Coin
	Coins    = coin.Coins
	Int      = coin.Int
	Dec      = coin.Dec
	DecCoin  = coin.DecCoin
	DecCoins = coin.DecCoins
)

var (
	NewCoin              = coin.NewCoin
	NewCoins             = coin.NewCoins
	NewDecCoinsFromCoins = coin.NewDecCoinsFromCoins
	NewDecCoinFromDec    = coin.NewDecCoinFromDec
	NewInt64Coin         = coin.NewInt64Coin
	NewInt64Coins        = coin.NewInt64Coins
	NewInt               = coin.NewInt
	ParseCoin            = coin.ParseCoin
	ParseCoins           = coin.ParseCoins
	NewDecCoin           = coin.NewDecCoin
	NewDecCoins          = coin.NewDecCoins
	ParseDecCoins        = coin.ParseDecCoins
	NewDec               = sdk.NewDec
	ValidateDenom        = coin.ValidateDenom
	ErrCoinDenomInvalid  = coin.ErrCoinDenomInvalid
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

// NewInt64CoreCoin create a default bond denom coin amount by int64 amt
func NewInt64CoreCoin(amt int64) Coin {
	return NewCoin(keys.DefaultBondDenom, NewInt(amt))
}

// NewInt64CoreCoins create a default bond denom coins type amount by int64 amt
func NewInt64CoreCoins(amt int64) Coins {
	return Coins{NewInt64CoreCoin(amt)}
}

// NewIntCoreCoin create a default bond denom coin amount by Int amt
func NewIntCoreCoin(val Int) Coin {
	return NewCoin(keys.DefaultBondDenom, val)
}

// NewIntCoreCoins create a default bond denom coins type amount by Int amt
func NewIntCoreCoins(val Int) Coins {
	return Coins{NewIntCoreCoin(val)}
}
