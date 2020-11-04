package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset/keeper"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type IAssetKeeper interface {
	LockCoins(ctx sdk.Context, account types.AccountID, unlockBlockHeight int64, coins types.Coins) error
	Approve(ctx sdk.Context, id, spender types.AccountID, amt types.Coins, isLock bool) error
	CoinsToPower(ctx sdk.Context, from, to types.AccountID, amt types.Coins) error
	UnLockFreezedCoins(ctx sdk.Context, account types.AccountID, coins types.Coins) error

	GetLockCoins(ctx sdk.Context, account types.AccountID) (types.Coins, []keeper.LockedCoins, error)
	GetApproveCoins(ctx sdk.Context, account, spender types.AccountID) (*keeper.ApproveData, error)
	GetCoins(ctx sdk.Context, account types.AccountID) (types.Coins, error)
	GetCoinStat(ctx sdk.Context, creator, symbol types.Name) (*assetTypes.CoinStat, error)
}
