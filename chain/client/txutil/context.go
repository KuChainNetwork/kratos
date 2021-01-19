package txutil

import (
	"io"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// KuCLIContext cli context for kuchain, add account info
type KuCLIContext struct {
	ctx context.CLIContext

	FromAccount AccountID
}

// NewKuCLICtx creates a new KuCLIContext
func NewKuCLICtx(ctx context.CLIContext) KuCLIContext {
	return KuCLIContext{
		ctx: ctx,
	}.WithFromAccount(types.NewAccountIDFromAccAdd(ctx.GetFromAddress()))
}

func NewKuCLICtxNoFrom(ctx context.CLIContext) KuCLIContext {
	return KuCLIContext{
		ctx: ctx,
	}
}

// NewKuCLICtxByBuf creates a new KuCLIContext with cmd
func NewKuCLICtxByBuf(cdc *codec.Codec, inBuf io.Reader) KuCLIContext {
	return NewKuCLICtx(context.NewCLIContextWithInput(inBuf).WithCodec(cdc))
}

func NewKuCLICtxByBufNoFrom(cdc *codec.Codec, inBuf io.Reader) KuCLIContext {
	ctx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	return NewKuCLICtxNoFrom(ctx)
}

func (k KuCLIContext) Ctx() context.CLIContext {
	return k.ctx
}

type CtxQueryFunc func(string, []byte) ([]byte, int64, error)

func (k KuCLIContext) Codec() *codec.Codec                { return k.ctx.Codec }
func (k KuCLIContext) GetFromAddress() types.AccAddress   { return k.ctx.FromAddress }
func (k KuCLIContext) GetFromName() string                { return k.ctx.GetFromName() }
func (k KuCLIContext) Output() io.Writer                  { return k.ctx.Output }
func (k KuCLIContext) GetQueryWithDataFunc() CtxQueryFunc { return k.ctx.QueryWithData }

func (k KuCLIContext) BroadcastTx(txBytes []byte) (sdk.TxResponse, error) {
	return k.ctx.BroadcastTx(txBytes)
}

func (k KuCLIContext) PrintOutput(t interface{}) error {
	return k.ctx.PrintOutput(t)
}

func (k KuCLIContext) SkipConfirm() bool  { return k.ctx.SkipConfirm }
func (k KuCLIContext) Simulate() bool     { return k.ctx.Simulate }
func (k KuCLIContext) GenerateOnly() bool { return k.ctx.GenerateOnly }
func (k KuCLIContext) Indent() bool       { return k.ctx.Indent }

// WithFromAccount with account from accountID
func (k KuCLIContext) WithFromAccount(from AccountID) KuCLIContext {
	k.FromAccount = from
	return k
}

// WithAccount with account name
func (k KuCLIContext) WithAccount(name Name) KuCLIContext {
	k.FromAccount = types.NewAccountIDFromName(name)
	return k
}

// GetAccountID get account id
func (k KuCLIContext) GetAccountID() AccountID {
	return k.FromAccount
}

// WithOutput returns a copy of the context with an updated output writer (e.g. stdout).
func (k KuCLIContext) WithOutput(w io.Writer) KuCLIContext {
	k.ctx.Output = w
	return k
}

// GetAccountInfo get account info by from account id
func (k KuCLIContext) GetAccountInfo() (exported.Account, error) {
	return NewAccountRetriever(k).GetAccount(k.GetAccountID())
}
