package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type sdkContext = types.Context

// Context for kuchain msg handler
type Context struct {
	sdkContext

	msg    KuTransfMsg
	auths  []AccAddress
	auther AccountAuther

	authStats *contextAuthStat
}

type contextAuthStat struct {
	requireAuths        []AccAddress
	requireAccountAuths []Name
	allowedAccountAuths []Name
}

func NewKuMsgCtx(ctx types.Context, auther AccountAuther, msg sdk.Msg) Context {
	return Context{
		sdkContext: ctx,
		auther:     auther,
		authStats: &contextAuthStat{
			requireAuths:        make([]AccAddress, 0, 1),
			requireAccountAuths: make([]Name, 0, 1),
			allowedAccountAuths: make([]Name, 0, 1),
		},
	}.WithTransfMsg(msg)
}

// Context get sdk context wrapped
func (c Context) Context() types.Context {
	return c.sdkContext
}

// WithTransfer ctx with transfer
func (c Context) WithTransfMsg(msg sdk.Msg) Context {
	m, ok := msg.(KuTransfMsg)
	if ok {
		c.msg = m
	}
	return c
}

// WithAuths ctx with trx auths
func (c Context) WithAuths(auths []AccAddress) Context {
	c.auths = auths
	return c
}

func (c Context) isHasAuth(auth AccAddress) bool {
	for _, a := range c.auths {
		if a.Equals(auth) {
			return true
		}
	}

	return false
}

func (c Context) isAllowAccountAuth(auth Name) bool {
	for _, a := range c.authStats.allowedAccountAuths {
		if a.Eq(auth) {
			return true
		}
	}

	return false
}

func (c Context) CheckAuths() error {
	for _, auth := range c.authStats.requireAuths {
		if !c.isHasAuth(auth) {
			return sdkerrors.Wrapf(ErrMissingAuth, "missing auth %s", auth)
		}
	}

	for _, n := range c.authStats.requireAccountAuths {
		if c.isAllowAccountAuth(n) {
			continue
		}

		auth, err := c.auther.GetAuth(c.sdkContext, n)
		if err != nil {
			return sdkerrors.Wrapf(err, "missing account %s auth", n)
		}

		if !c.isHasAuth(auth) {
			return sdkerrors.Wrapf(ErrMissingAuth, "missing auth %s by account %s", auth, n)
		}
	}

	return nil
}

// RequireAuth require account auth
func (c Context) RequireAuth(permissions ...AccountID) {
	for _, id := range permissions {
		if accAdd, ok := id.ToAccAddress(); ok {
			c.authStats.requireAuths = append(c.authStats.requireAuths, accAdd)
		}

		if name, ok := id.ToName(); ok {
			c.authStats.requireAccountAuths = append(c.authStats.requireAccountAuths, name)
		}
	}
}

// RequireAccountAuth require address auth
func (c Context) RequireAccountAuth(adds ...AccAddress) {
	c.authStats.requireAuths = append(c.authStats.requireAuths, adds...)
}

// RequireAuth require account auth
func (c Context) RequireAccount(account ...Name) {
	c.authStats.requireAccountAuths = append(c.authStats.requireAccountAuths, account...)
}

// Authorize make authorize for account to this msg, it call by handlers to allow kumsg can use this auth
func (c Context) Authorize(account ...Name) {
	// TODO: Now kuchain not support user-define code or contracts, so this no check if handler REALLY have auth to all account
	c.authStats.allowedAccountAuths = append(c.authStats.allowedAccountAuths, account...)
}

// RequireTransfer require transfer coin large then amount for to
func (c Context) RequireTransfer(to AccountID, amount Coins) error {
	if !c.msg.GetTo().Equals(&to) {
		return ErrTransfNotTo
	}

	if !c.msg.GetAmount().IsAllGTE(amount) {
		return ErrTransfNoEnough
	}

	return nil
}

// GetTransf get transfer info from msg
func (c Context) GetTransf() (from, to AccountID, amount Coins) {
	if c.msg == nil {
		return
	}

	from = c.msg.GetFrom()
	to = c.msg.GetTo()
	amount = c.msg.GetAmount()
	return
}
