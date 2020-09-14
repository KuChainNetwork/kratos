package txutil

import (
	"io"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

// KuCLIContext cli context for kuchain, add account info
type KuCLIContext struct {
	context.CLIContext

	FromAccount AccountID
}

// NewKuCLICtx creates a new KuCLIContext
func NewKuCLICtx(ctx context.CLIContext) KuCLIContext {
	return KuCLIContext{
		CLIContext: ctx,
	}.WithFromAccount(types.NewAccountIDFromAccAdd(ctx.GetFromAddress()))
}

func NewKuCLICtxNoFrom(ctx context.CLIContext) KuCLIContext {
	return KuCLIContext{
		CLIContext: ctx,
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
	k.Output = w
	return k
}

// GetAccountInfo get account info by from account id
func (k KuCLIContext) GetAccountInfo() (exported.Account, error) {
	return NewAccountRetriever(k).GetAccount(k.GetAccountID())
}
