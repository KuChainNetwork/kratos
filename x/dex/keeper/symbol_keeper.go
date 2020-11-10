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
	creator types.Name, symbol *types.Symbol) (err error) {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		err = errors.Wrapf(types.ErrDexNotExists,
			"create symbol dex %s not exists",
			creator.String())
		return
	}
	// check base and quote are exists
	values := strings.Split(symbol.Base.Code, "/")
	if 2 != len(values) {
		err = errors.Wrapf(types.ErrSymbolFormat,
			"create symbol dex %s coin symbol %s format error",
			creator.String(),
			symbol.Base.Code)
		return
	}
	baseCode := values[1]
	values = strings.Split(symbol.Quote.Code, "/")
	if 2 != len(values) {
		err = errors.Wrapf(types.ErrSymbolFormat,
			"create symbol dex %s coin symbol %s format error",
			creator.String(),
			symbol.Quote.Code)
		return
	}
	quoteCode := values[1]
	var baseCodeFound, quoteCodeFound bool
	a.assetKeeper.IterateAllCoins(ctx, func(_ chainType.AccountID, balance Coins) (stop bool) {
		for _, coin := range balance {
			s := strings.Split(coin.Denom, "/")[1]
			if s == baseCode {
				baseCodeFound = true
			}
			if s == quoteCode {
				quoteCodeFound = true
			}
			if baseCodeFound && quoteCodeFound {
				stop = true
				return
			}
		}
		return
	})
	if !baseCodeFound || !quoteCodeFound {
		err = errors.Wrapf(types.ErrSymbolNotSupply,
			"create symbol dex %s coin symbol %s/%s not supply",
			creator.String(),
			creator.String(),
			symbol.Base.Code)
		return
	}
	if dex, ok = dex.WithSymbol(symbol); !ok {
		err = types.ErrSymbolExists
		return
	}
	a.setDex(ctx, dex)
	return
}

// UpdateSymbol update symbol
func (a DexKeeper) UpdateSymbol(ctx sdk.Context,
	creator types.Name, update *types.Symbol) (err error) {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		err = errors.Wrapf(types.ErrDexNotExists,
			"update symbol dex %s not exists",
			creator.String())
		return
	}
	var symbol types.Symbol
	if symbol, ok = dex.Symbol(update.Base.Creator,
		update.Base.Code,
		update.Quote.Creator,
		update.Quote.Code); !ok {
		err = errors.Wrapf(types.ErrSymbolNotExists,
			"update symbol not exists, dex %s", creator.String())
		return
	}
	updated := false
	for _, pair := range []struct {
		Dst *string
		Src string
	}{
		{&symbol.Base.Name, update.Base.Name},
		{&symbol.Base.FullName, update.Base.FullName},
		{&symbol.Base.IconUrl, update.Base.IconUrl},
		{&symbol.Base.TxUrl, update.Base.TxUrl},
		{&symbol.Quote.Name, update.Quote.Name},
		{&symbol.Quote.FullName, update.Quote.FullName},
		{&symbol.Quote.IconUrl, update.Quote.IconUrl},
		{&symbol.Quote.TxUrl, update.Quote.TxUrl},
	} {
		if 0 < len(pair.Src) && *pair.Dst != pair.Src {
			*pair.Dst = pair.Src
			updated = true
		}
	}
	if !updated {
		return
	}
	if !dex.UpdateSymbol(update.Base.Creator,
		update.Base.Code,
		update.Quote.Creator,
		update.Quote.Code,
		&symbol) {
		err = errors.Wrapf(types.ErrSymbolNotExists,
			"update symbol (%s/%s) not exists",
			update.Base.Code, update.Quote.Code)
		return
	}
	a.setDex(ctx, dex)
	return
}

// PauseSymbol pause symbol
func (a DexKeeper) PauseSymbol(ctx sdk.Context,
	creator types.Name, baseCreator, baseCode, quoteCreator, quoteCode string) (err error) {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		err = errors.Wrapf(types.ErrDexNotExists,
			"pause symbol, dex %s not exists",
			creator.String())
		return
	}
	var symbol types.Symbol
	if symbol, ok = dex.Symbol(baseCreator,
		baseCode,
		quoteCreator,
		quoteCode); !ok {
		err = errors.Wrapf(types.ErrSymbolNotExists,
			"pause symbol, symbol (%s/%s) not exists",
			baseCode,
			quoteCode)
		return
	}
	if !dex.UpdateSymbol(baseCreator,
		baseCode,
		quoteCreator,
		quoteCode,
		(&symbol).WithPaused(true)) {
		err = errors.Wrapf(types.ErrSymbolNotExists,
			"pause symbol, symbol (%s/%s) not exists",
			baseCode,
			quoteCode)
		return
	}
	a.setDex(ctx, dex)
	return
}

// RestoreSymbol restore symbol
func (a DexKeeper) RestoreSymbol(ctx sdk.Context,
	creator types.Name, baseCreator, baseCode, quoteCreator, quoteCode string) (err error) {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		err = errors.Wrapf(types.ErrDexNotExists,
			"restore symbol, dex %s not exists",
			creator.String())
		return
	}
	var symbol types.Symbol
	if symbol, ok = dex.Symbol(baseCreator,
		baseCode,
		quoteCreator,
		quoteCode); !ok {
		err = errors.Wrapf(types.ErrSymbolNotExists,
			"restore symbol, symbol (%s/%s) not exists",
			baseCode,
			quoteCode)
		return
	}
	if !dex.UpdateSymbol(baseCreator,
		baseCode,
		quoteCreator,
		quoteCode,
		(&symbol).WithPaused(false)) {
		err = errors.Wrapf(types.ErrSymbolNotExists,
			"restore symbol, symbol (%s/%s) not exists",
			baseCode,
			quoteCode)
		return
	}
	a.setDex(ctx, dex)
	return
}

// ShutdownSymbol shutdown symbol
func (a DexKeeper) ShutdownSymbol(ctx sdk.Context,
	creator types.Name, baseCreator, baseCode, quoteCreator, quoteCode string) (err error) {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		err = errors.Wrapf(types.ErrDexNotExists,
			"delete symbol, dex %s not exists",
			creator.String())
		return
	}
	if !dex.DeleteSymbol(baseCreator, baseCode, quoteCreator, quoteCode) {
		err = errors.Wrapf(types.ErrSymbolNotExists,
			"delete symbol, symbol (%s/%s) not exists",
			baseCode,
			quoteCode)
		return
	}
	a.setDex(ctx, dex)
	return
}
