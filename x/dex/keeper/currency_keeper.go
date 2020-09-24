package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

// CreateCurrency create currency
func (a DexKeeper) CreateCurrency(ctx sdk.Context,
	creator types.Name, currency *types.Currency) (err error) {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		err = errors.Wrapf(types.ErrDexNotExists,
			"create currency dex %s not exists",
			creator.String())
		return
	}
	if dex, ok = dex.WithCurrency(currency); !ok {
		err = types.ErrCurrencyExists
		return
	}
	a.setDex(ctx, dex)
	return
}

// UpdateCurrencyInfo update concurrency basic info
func (a DexKeeper) UpdateCurrencyInfo(ctx sdk.Context,
	creator types.Name, update *types.Currency) (err error) {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		err = errors.Wrapf(types.ErrDexNotExists,
			"update currency dex %s not exists",
			creator.String())
		return
	}
	var currency types.Currency
	if currency, ok = dex.Currency(update.Base.Code, update.Quote.Code); !ok {
		err = errors.Wrapf(types.ErrCurrencyNotExists,
			"update currency not exists, dex %s", creator.String())
		return
	}
	updated := false
	for _, pair := range []struct {
		Dst *string
		Src string
	}{
		{&currency.Base.Name, update.Base.Name},
		{&currency.Base.FullName, update.Base.FullName},
		{&currency.Base.IconUrl, update.Base.IconUrl},
		{&currency.Base.TxUrl, update.Base.TxUrl},
		{&currency.Quote.Name, update.Quote.Name},
		{&currency.Quote.FullName, update.Quote.FullName},
		{&currency.Quote.IconUrl, update.Quote.IconUrl},
		{&currency.Quote.TxUrl, update.Quote.TxUrl},
	} {
		if *pair.Dst != pair.Src {
			*pair.Dst = pair.Src
			updated = true
		}
	}
	if !updated {
		return
	}
	if !dex.UpdateCurrency(update.Base.Code, update.Quote.Code, &currency) {
		err = errors.Wrapf(types.ErrCurrencyNotExists,
			"update currency (%s/%s) not exists",
			update.Base.Code, update.Quote.Code)
		return
	}
	a.setDex(ctx, dex)
	return
}

// PauseCurrency pause currency
func (a DexKeeper) PauseCurrency(ctx sdk.Context,
	creator types.Name, baseCode, quoteCode string) (err error) {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		err = errors.Wrapf(types.ErrDexNotExists,
			"pause currency, dex %s not exists",
			creator.String())
		return
	}
	var currency types.Currency
	if currency, ok = dex.Currency(baseCode, quoteCode); !ok {
		err = errors.Wrapf(types.ErrCurrencyNotExists,
			"pause currency, currency (%s/%s) not exists",
			baseCode,
			quoteCode)
		return
	}
	if !dex.UpdateCurrency(baseCode, quoteCode, (&currency).WithPaused(true)) {
		err = errors.Wrapf(types.ErrCurrencyNotExists,
			"pause currency, currency (%s/%s) not exists",
			baseCode,
			quoteCode)
		return
	}
	a.setDex(ctx, dex)
	return
}

// RestoreCurrency pause currency
func (a DexKeeper) RestoreCurrency(ctx sdk.Context,
	creator types.Name, baseCode, quoteCode string) (err error) {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		err = errors.Wrapf(types.ErrDexNotExists,
			"restore currency, dex %s not exists",
			creator.String())
		return
	}
	var currency types.Currency
	if currency, ok = dex.Currency(baseCode, quoteCode); !ok {
		err = errors.Wrapf(types.ErrCurrencyNotExists,
			"restore currency, currency (%s/%s) not exists",
			baseCode,
			quoteCode)
		return
	}
	if !dex.UpdateCurrency(baseCode, quoteCode, (&currency).WithPaused(false)) {
		err = errors.Wrapf(types.ErrCurrencyNotExists,
			"restore currency, currency (%s/%s) not exists",
			baseCode,
			quoteCode)
		return
	}
	a.setDex(ctx, dex)
	return
}

// ShutdownCurrency shutdown currency
func (a DexKeeper) ShutdownCurrency(ctx sdk.Context,
	creator types.Name, baseCode, quoteCode string) (err error) {
	dex, ok := a.getDex(ctx, creator)
	if !ok {
		err = errors.Wrapf(types.ErrDexNotExists,
			"delete currency, dex %s not exists",
			creator.String())
		return
	}
	if !dex.DeleteCurrency(baseCode, quoteCode) {
		err = errors.Wrapf(types.ErrCurrencyNotExists,
			"delete currency, currency (%s/%s) not exists",
			baseCode,
			quoteCode)
		return
	}
	a.setDex(ctx, dex)
	return
}
