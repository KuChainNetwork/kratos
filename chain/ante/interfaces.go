package ante

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SigVerifiableTx defines a Tx interface for all signature verification decorators
type SigVerifiableTx interface {
	types.Tx
	GetSignatures() []types.StdSignature
	GetSigners() []types.AccAddress
}

// AssetKeeper
type AssetKeeper interface {
	PayFee(sdk.Context, types.AccountID, types.Coins) error
}

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, id AccountID) exported.Account
}
