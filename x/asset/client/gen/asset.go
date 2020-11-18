package gen

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

// GensisAccountAssetCmd builds gen genesis account asset to genesis config
func GensisAccountAssetCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-account-coin",
		Short: "Add a genesis coin for a account to chain",
		Args:  cobra.ExactArgs(2),
		Long: `This command add a genesis coin to chain'.

		It creates a some genesis coin for a account, then put the data to genesis.json
	`,

		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			accountID, err := types.NewAccountIDFromStr(args[0])
			if err != nil {
				return err
			}

			coins, err := types.ParseCoins(args[1])
			if err != nil {
				return err
			}

			genFile := config.GenesisFile()

			var genesis asset.GenesisState
			if err := types.LoadGenesisStateFromFile(cdc, genFile, asset.ModuleName, &genesis); err != nil {
				return err
			}

			if err = addAccountAsset(cdc, &genesis, accountID, coins); err != nil {
				return err
			}

			return types.SaveGenesisStateToFile(cdc, genFile, asset.ModuleName, genesis)
		},
	}

	return cmd
}

func addAccountAsset(cdc *codec.Codec, state *asset.GenesisState, accountID types.AccountID, coins types.Coins) error {
	for _, g := range state.GenesisAssets {
		if g.GetID().Eq(accountID) {
			return fmt.Errorf("the application state already contains account coins")
		}
	}

	state.GenesisAssets = append(state.GenesisAssets, asset.NewGenesisAsset(accountID, coins...))
	return nil
}
