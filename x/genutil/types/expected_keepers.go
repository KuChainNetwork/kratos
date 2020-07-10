package types

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KuChainNetwork/kuchain/chain/types"
	accountexported "github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/KuChainNetwork/kuchain/x/asset"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakingKeeper defines the expected staking keeper (noalias)
type StakingKeeper interface {
	ApplyAndReturnValidatorSetUpdates(sdk.Context) (updates []abci.ValidatorUpdate)
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, accountexported.Account) accountexported.Account
	SetAccount(sdk.Context, accountexported.Account)
	IterateAccounts(ctx sdk.Context, process func(accountexported.Account) (stop bool))
}

// GenesisAccountsIterator defines the expected iterating genesis accounts object (noalias)
type GenesisAccountsIterator interface {
	IterateGenesisAccounts(
		cdc *codec.Codec,
		appGenesis types.AppGenesisState,
		cb func(accountexported.Account) (stop bool),
	)
}

// GenesisAccountsIterator defines the expected iterating genesis accounts object (noalias)
type GenesisBalancesIterator interface {
	IterateGenesisBalances(
		cdc *codec.Codec,
		appStat types.AppGenesisState,
		cb func(asset.GenesisAsset) (stop bool),
	)
}

// StakingFuncManager defines the expected staking functions
type StakingFuncManager interface {
	// function: validate stakingtypes.MsgCreateValidator
	MsgCreateValidatorVal(msg sdk.Msg) bool

	// function: validate stakingtypes.MsgCreateValidator and check balance, then return monitor
	MsgDelegateWithBalance(m sdk.Msg, balancesMap map[string]asset.GenesisAsset) error
	GetMsgCreateValidatorMoniker(msg sdk.Msg) (string, error)
	// function: decode appGenesisState and get bondDenom
	GetBondDenom(appGenesisState map[string]json.RawMessage) string
}
