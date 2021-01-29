package client

import (
	"io"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Context cli context for kuchain, add account info
type Context struct {
	ctx client.Context

	FromAccount AccountID
}

// NewKuCLICtx creates a new Context
func NewKuCLICtx(ctx client.Context) Context {
	res := Context{
		ctx: ctx,
	}

	address := ctx.GetFromAddress()
	if !address.Empty() {
		res = res.WithFromAccount(types.NewAccountIDFromAccAdd(address))
	}

	return res
}

func NewKuCLICtxNoFrom(ctx client.Context) Context {
	return Context{
		ctx: ctx,
	}
}

// NewKuCLICtxByBuf creates a new Context with cmd
func NewKuCLICtxByBuf(cdc *codec.Codec, inBuf io.Reader) Context {
	return NewKuCLICtx(client.NewCLIContextWithInput(inBuf).WithCodec(cdc))
}

func NewKuCLICtxByBufNoFrom(cdc *codec.Codec, inBuf io.Reader) Context {
	ctx := client.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	return NewKuCLICtxNoFrom(ctx)
}

func NewCtxByCodec(cdc *codec.Codec) Context {
	return NewKuCLICtx(client.NewCLIContext().WithCodec(cdc))
}

func (k Context) Ctx() client.Context {
	return k.ctx
}

type CtxQueryFunc func(string, []byte) ([]byte, int64, error)

func (k Context) Codec() *codec.Codec                { return k.ctx.Codec }
func (k Context) GetFromAddress() types.AccAddress   { return k.ctx.FromAddress }
func (k Context) GetFromName() string                { return k.ctx.GetFromName() }
func (k Context) Output() io.Writer                  { return k.ctx.Output }
func (k Context) GetQueryWithDataFunc() CtxQueryFunc { return k.ctx.QueryWithData }
func (k Context) GetHeight() int64                   { return k.ctx.Height }
func (k Context) GetClient() rpcclient.Client        { return k.ctx.Client }
func (k Context) GetOutputFormat() string            { return k.ctx.OutputFormat }

func (k Context) BroadcastTx(txBytes []byte) (sdk.TxResponse, error) {
	return k.ctx.BroadcastTx(txBytes)
}

func (k Context) PrintOutput(t interface{}) error {
	return k.ctx.PrintOutput(t)
}

func (k Context) GetNode() (rpcclient.Client, error) {
	return k.ctx.GetNode()
}

func (k Context) Verify(height int64) (tmtypes.SignedHeader, error) {
	return k.ctx.Verify(height)
}

func (k Context) QueryWithData(path string, data []byte) ([]byte, int64, error) {
	return k.ctx.QueryWithData(path, data)
}

func (k Context) SkipConfirm() bool  { return k.ctx.SkipConfirm }
func (k Context) Simulate() bool     { return k.ctx.Simulate }
func (k Context) GenerateOnly() bool { return k.ctx.GenerateOnly }
func (k Context) Indent() bool       { return k.ctx.Indent }
func (k Context) TrustNode() bool    { return k.ctx.TrustNode }

func (k Context) QuerySubspace(subspace []byte, storeName string) ([]sdk.KVPair, int64, error) {
	return k.ctx.QuerySubspace(subspace, storeName)
}

func (k Context) QueryStore(key tmbytes.HexBytes, storeName string) ([]byte, int64, error) {
	return k.ctx.QueryStore(key, storeName)
}

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

func (k Context) WithHeight(height int64) Context {
	k.ctx = k.ctx.WithHeight(height)
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

func (k Context) WithBroadcastMode(mode string) Context {
	k.ctx = k.ctx.WithBroadcastMode(mode)
	return k
}

// GetAccountInfo get account info by from account id
func (k Context) GetAccountInfo() (exported.Account, error) {
	return NewAccountRetriever(k).GetAccount(k.GetAccountID())
}
