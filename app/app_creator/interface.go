package appcreator

import (
	"github.com/KuChainNetwork/kuchain/x/account"
	"github.com/KuChainNetwork/kuchain/x/asset"
	"github.com/KuChainNetwork/kuchain/x/dex"
	distr "github.com/KuChainNetwork/kuchain/x/distribution"
	"github.com/KuChainNetwork/kuchain/x/evidence"
	"github.com/KuChainNetwork/kuchain/x/genutil"
	"github.com/KuChainNetwork/kuchain/x/genutil/types"
	"github.com/KuChainNetwork/kuchain/x/gov"
	"github.com/KuChainNetwork/kuchain/x/mint"
	"github.com/KuChainNetwork/kuchain/x/params"
	"github.com/KuChainNetwork/kuchain/x/slashing"
	"github.com/KuChainNetwork/kuchain/x/staking"
	"github.com/KuChainNetwork/kuchain/x/supply"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DeliverTxfn = genutil.DeliverTxfn

type KuAppWithKeeper interface {
	Router() sdk.Router
	QueryRouter() sdk.QueryRouter

	GetDeliverTx() DeliverTxfn
	GetStakingFuncMng() types.StakingFuncManager

	AccountKeeper() account.Keeper
	AssetKeeper() asset.Keeper
	SupplyKeeper() supply.Keeper
	DistrKeeper() distr.Keeper
	MintKeeper() mint.Keeper
	ParamsKeeper() *params.Keeper
	StakingKeeper() *staking.Keeper
	SlashingKeeper() slashing.Keeper
	EvidenceKeeper() evidence.Keeper
	GovKeeper() gov.Keeper
	DexKeeper() dex.Keeper
}
