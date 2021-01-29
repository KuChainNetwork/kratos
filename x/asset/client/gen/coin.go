package gen

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset"
	assetTypes "github.com/KuChainNetwork/kuchain/x/asset/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

// GenCoinCmd builds gen genesis coin type to genesis config
func GensisCoinCmd(ctx *server.Context, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-coin",
		Short: "Add a genesis coin type to chain",
		Args:  cobra.ExactArgs(2),
		Long: `This command add a genesis coin to chain'.

		It creates a genesis coin type, then put the data to genesis.json
	`,

		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			maxSupply, err := types.ParseCoin(args[0])
			if err != nil {
				return err
			}

			creator, symbol, err := types.CoinAccountsFromDenom(maxSupply.Denom)
			if err != nil {
				return err
			}

			genFile := config.GenesisFile()

			var genesis asset.GenesisState
			if err := types.LoadGenesisStateFromFile(cdc, genFile, asset.ModuleName, &genesis); err != nil {
				return err
			}

			if err = addCoin(cdc, &genesis, creator, symbol, maxSupply, args[1]); err != nil {
				return err
			}

			return types.SaveGenesisStateToFile(cdc, genFile, asset.ModuleName, genesis)
		},
	}

	return cmd
}

func addCoin(cdc *codec.LegacyAmino, state *asset.GenesisState, creator, symbol types.Name, maxSupply types.Coin, desc string) error {
	for _, g := range state.GenesisCoins {
		if g.GetCreator().Eq(creator) && g.GetSymbol().Eq(symbol) {
			return fmt.Errorf("the application state already contains coins")
		}
	}

	c := asset.NewGenesisCoin(&assetTypes.CoinStat{
		Creator:      creator,
		Symbol:       symbol,
		CreateHeight: 1,
		Supply:       types.NewCoin(maxSupply.Denom, types.NewInt(0)),
		MaxSupply:    maxSupply,
	}, []byte(desc))

	if err := c.Validate(); err != nil {
		return err
	}

	state.GenesisCoins = append(state.GenesisCoins, c)
	return nil
}
