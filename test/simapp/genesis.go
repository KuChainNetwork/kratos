package simapp

import (
	"encoding/json"
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	accountExported "github.com/KuChainNetwork/kuchain/x/account/exported"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	"github.com/pkg/errors"
)

// The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

type GenesisAccounts struct {
	accounts []accountExported.GenesisAccount
	assets   []assetTypes.GenesisAsset
	coins    []assetTypes.GenesisCoin
}

func NewGenesisAccounts(rootAuth types.AccAddress, accounts ...SimGenesisAccount) *GenesisAccounts {
	res := &GenesisAccounts{
		accounts: make([]accountExported.GenesisAccount, 0, len(accounts)+1),
		assets:   make([]assetTypes.GenesisAsset, 0, len(accounts)+1),
		coins:    make([]assetTypes.GenesisCoin, 0, 8),
	}

	coins2genesis := make(map[string]types.Coin)
	defaultCoins := types.NewCoin(constants.DefaultBondDenom, types.NewInt(0))

	coins2genesis[defaultCoins.Denom] = defaultCoins

	var accountNumber uint64 = 1

	res.accounts = append(res.accounts, NewGenesisAccount(constants.SystemAccountID, rootAuth, accountNumber))
	res.assets = append(res.assets,
		assetTypes.NewGenesisAsset(
			constants.SystemAccountID,
			types.NewCoin(constants.DefaultBondDenom, types.NewInt(1000000000000000000))))
	accountNumber++

	for _, a := range accounts {
		res.accounts = append(res.accounts, NewGenesisAccount(a.ID, a.GetAuth(), accountNumber))
		res.assets = append(res.assets, assetTypes.NewGenesisAsset(a.ID, a.Assets...))

		for _, as := range a.Assets {
			coins2genesis[as.Denom] = as
		}

		accountNumber++
	}

	for _, c := range coins2genesis {
		createor, symbol, err := types.CoinAccountsFromDenom(c.Denom)
		if err != nil {
			panic(errors.Wrapf(err, "coin symbol err %s", c.Denom))
		}
		max := types.NewInt(1000000000000000000) // for default
		max = max.Mul(types.NewInt(1000000000000000000))

		if c.Denom == constants.DefaultBondDenom {
			max = types.NewInt(0)
		}

		res.coins = append(res.coins, assetTypes.NewGenesisCoin(createor, symbol, max, fmt.Sprintf("desc for %s", c.Denom)))
	}

	return res
}
