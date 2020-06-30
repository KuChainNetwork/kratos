package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Context for plugin ctx
type Context struct {
	chainID string
	logger  log.Logger
}

func (c Context) ChainID() string    { return c.chainID }
func (c Context) Logger() log.Logger { return c.logger }

// NewContext create a new context
func NewContext(logger log.Logger) Context {
	return Context{
		logger: logger,
	}
}

func (c Context) WithChainID(chainID string) Context {
	c.chainID = chainID
	return c
}

func NewCtx(ctx sdk.Context) Context {
	return NewContext(ctx.Logger()).WithChainID(ctx.ChainID())
}
