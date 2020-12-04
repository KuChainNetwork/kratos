package keeper

import (
	"strings"

	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

// CreateSymbol create symbol
func (a DexKeeper) CreateSymbol(ctx sdk.Context,
	creator types.Name, symbol *types.Symbol) error {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		return errors.Wrapf(types.ErrDexNotExists,
			"create symbol dex %s not exists",
			creator.String())
	}

	// check base and quote are exists
	values := strings.Split(symbol.Base.Code, "/")
	if len(values) != 2 {
		return errors.Wrapf(types.ErrSymbolFormat,
			"create symbol dex %s coin symbol %s format error",
			creator.String(), symbol.Base.Code)
	}

	baseCode := values[1]
	values = strings.Split(symbol.Quote.Code, "/")
	if len(values) != 2 {
		return errors.Wrapf(types.ErrSymbolFormat,
			"create symbol dex %s coin symbol %s format error",
			creator.String(), symbol.Quote.Code)
	}

	quoteCode := values[1]
	var baseCodeFound, quoteCodeFound bool
	a.assetKeeper.IterateAllCoins(ctx, func(_ chainType.AccountID, balance Coins) bool {
		for _, coin := range balance {
			values = strings.Split(coin.Denom, "/")
			c := values[0]
			s := values[1]
			if c == symbol.Base.Creator && s == baseCode {
				baseCodeFound = true
			}

			if c == symbol.Quote.Creator && s == quoteCode {
				quoteCodeFound = true
			}

			if baseCodeFound && quoteCodeFound {
				return true
			}
		}
		return false
	})

	if !baseCodeFound || !quoteCodeFound {
		return errors.Wrapf(types.ErrSymbolNotSupply,
			"create symbol dex %s coin symbol %s/%s not supply",
			creator.String(), creator.String(), symbol.Base.Code)
	}

	if dex, ok = dex.WithSymbol(symbol); !ok {
		return types.ErrSymbolExists
	}

	a.setDex(ctx, dex)
	return nil
}

// UpdateSymbol update symbol
func (a DexKeeper) UpdateSymbol(ctx sdk.Context, creator types.Name, update *types.Symbol) error {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		return errors.Wrapf(types.ErrDexNotExists,
			"update symbol dex %s not exists",
			creator.String())
	}

	symbol, ok := dex.Symbol(update.Base.Creator, update.Base.Code,
		update.Quote.Creator, update.Quote.Code)

	if !ok {
		return errors.Wrapf(types.ErrSymbolNotExists,
			"update symbol not exists, dex %s", creator.String())
	}

	updated := false
	for _, pair := range []struct {
		Dst *string
		Src string
	}{
		{&symbol.Base.Name, update.Base.Name},
		{&symbol.Base.FullName, update.Base.FullName},
		{&symbol.Base.IconURL, update.Base.IconURL},
		{&symbol.Base.TxURL, update.Base.TxURL},
		{&symbol.Quote.Name, update.Quote.Name},
		{&symbol.Quote.FullName, update.Quote.FullName},
		{&symbol.Quote.IconURL, update.Quote.IconURL},
		{&symbol.Quote.TxURL, update.Quote.TxURL},
	} {
		if 0 < len(pair.Src) && *pair.Dst != pair.Src {
			*pair.Dst = pair.Src
			updated = true
		}
	}

	if !updated {
		return nil
	}

	if !dex.UpdateSymbol(update.Base.Creator,
		update.Base.Code,
		update.Quote.Creator,
		update.Quote.Code,
		&symbol) {
		return errors.Wrapf(types.ErrSymbolNotExists,
			"update symbol (%s/%s) not exists",
			update.Base.Code, update.Quote.Code)
	}

	a.setDex(ctx, dex)
	return nil
}

// PauseSymbol pause symbol
func (a DexKeeper) PauseSymbol(ctx sdk.Context,
	creator types.Name, baseCreator, baseCode, quoteCreator, quoteCode string) error {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		return errors.Wrapf(types.ErrDexNotExists,
			"pause symbol, dex %s not exists",
			creator.String())
	}
	var symbol types.Symbol
	if symbol, ok = dex.Symbol(baseCreator,
		baseCode,
		quoteCreator,
		quoteCode); !ok {
		return errors.Wrapf(types.ErrSymbolNotExists,
			"pause symbol, symbol (%s/%s) not exists",
			baseCode, quoteCode)
	}
	if !dex.UpdateSymbol(baseCreator,
		baseCode,
		quoteCreator,
		quoteCode,
		(&symbol).WithPaused(true)) {
		return errors.Wrapf(types.ErrSymbolNotExists,
			"pause symbol, symbol (%s/%s) not exists",
			baseCode, quoteCode)
	}

	a.setDex(ctx, dex)
	return nil
}

// RestoreSymbol restore symbol
func (a DexKeeper) RestoreSymbol(ctx sdk.Context,
	creator types.Name, baseCreator, baseCode, quoteCreator, quoteCode string) error {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		return errors.Wrapf(types.ErrDexNotExists,
			"restore symbol, dex %s not exists",
			creator.String())
	}

	var symbol types.Symbol
	if symbol, ok = dex.Symbol(baseCreator,
		baseCode,
		quoteCreator,
		quoteCode); !ok {
		return errors.Wrapf(types.ErrSymbolNotExists,
			"restore symbol, symbol (%s/%s) not exists",
			baseCode, quoteCode)
	}

	if !dex.UpdateSymbol(baseCreator,
		baseCode,
		quoteCreator,
		quoteCode,
		(&symbol).WithPaused(false)) {
		return errors.Wrapf(types.ErrSymbolNotExists,
			"restore symbol, symbol (%s/%s) not exists",
			baseCode, quoteCode)
	}

	a.setDex(ctx, dex)
	return nil
}

// ShutdownSymbol shutdown symbol
func (a DexKeeper) ShutdownSymbol(ctx sdk.Context,
	creator types.Name, baseCreator, baseCode, quoteCreator, quoteCode string) error {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		return errors.Wrapf(types.ErrDexNotExists,
			"delete symbol, dex %s not exists",
			creator.String())
	}

	if !dex.DeleteSymbol(baseCreator, baseCode, quoteCreator, quoteCode) {
		return errors.Wrapf(types.ErrSymbolNotExists,
			"delete symbol, symbol (%s/%s) not exists",
			baseCode,
			quoteCode)
	}

	a.setDex(ctx, dex)
	return nil
}
