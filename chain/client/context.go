package client

import (
	"io"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Context cli context for kuchain, add account info
type Context struct {
	ctx context.CLIContext

	FromAccount AccountID
}

// NewKuCLICtx creates a new Context
func NewKuCLICtx(ctx context.CLIContext) Context {
	return Context{
		ctx: ctx,
	}.WithFromAccount(types.NewAccountIDFromAccAdd(ctx.GetFromAddress()))
}

func NewKuCLICtxNoFrom(ctx context.CLIContext) Context {
	return Context{
		ctx: ctx,
	}
}

// NewKuCLICtxByBuf creates a new Context with cmd
func NewKuCLICtxByBuf(cdc *codec.Codec, inBuf io.Reader) Context {
	return NewKuCLICtx(context.NewCLIContextWithInput(inBuf).WithCodec(cdc))
}

func NewKuCLICtxByBufNoFrom(cdc *codec.Codec, inBuf io.Reader) Context {
	ctx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	return NewKuCLICtxNoFrom(ctx)
}

func (k Context) Ctx() context.CLIContext {
	return k.ctx
}

type CtxQueryFunc func(string, []byte) ([]byte, int64, error)

func (k Context) Codec() *codec.Codec                { return k.ctx.Codec }
func (k Context) GetFromAddress() types.AccAddress   { return k.ctx.FromAddress }
func (k Context) GetFromName() string                { return k.ctx.GetFromName() }
func (k Context) Output() io.Writer                  { return k.ctx.Output }
func (k Context) GetQueryWithDataFunc() CtxQueryFunc { return k.ctx.QueryWithData }

func (k Context) BroadcastTx(txBytes []byte) (sdk.TxResponse, error) {
	return k.ctx.BroadcastTx(txBytes)
}

func (k Context) PrintOutput(t interface{}) error {
	return k.ctx.PrintOutput(t)
}

func (k Context) SkipConfirm() bool  { return k.ctx.SkipConfirm }
func (k Context) Simulate() bool     { return k.ctx.Simulate }
func (k Context) GenerateOnly() bool { return k.ctx.GenerateOnly }
func (k Context) Indent() bool       { return k.ctx.Indent }

// WithFromAccount with account from accountID
func (k Context) WithFromAccount(from AccountID) Context {
	k.FromAccount = from
	return k
}

// WithAccount with account name
func (k Context) WithAccount(name Name) Context {
	k.FromAccount = types.NewAccountIDFromName(name)
	return k
}

// GetAccountID get account id
func (k Context) GetAccountID() AccountID {
	return k.FromAccount
}

// WithOutput returns a copy of the context with an updated output writer (e.g. stdout).
func (k Context) WithOutput(w io.Writer) Context {
	k.ctx.Output = w
	return k
}

// GetAccountInfo get account info by from account id
func (k Context) GetAccountInfo() (exported.Account, error) {
	return NewAccountRetriever(k).GetAccount(k.GetAccountID())
}
