package gen

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/app"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	flagClientHome = "home-client"
)

// GensisAccountCmd builds gen genesis account to genesis config
func GensisAccountCmd(ctx *server.Context, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-account",
		Short: "Add a genesis account to chain",
		Args:  cobra.ExactArgs(2),
		Long: `This command add a genesis account to chain'.

		It creates a genesis account which contains a name and auth, then put the data to genesis.json
	`,

		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			name, err := types.NewName(args[0])
			if err != nil {
				return err
			}

			addr, err := types.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			genFile := config.GenesisFile()

			var accountGenesis account.GenesisState
			if err := types.LoadGenesisStateFromFile(cdc, genFile, account.ModuleName, &accountGenesis); err != nil {
				return err
			}

			if err = addGenesisAccount(cdc, &accountGenesis, name, addr); err != nil {
				return err
			}

			return types.SaveGenesisStateToFile(cdc, genFile, account.ModuleName, accountGenesis)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, app.DefaultCLIHome, "client's home directory")
	return cmd
}

// GensisAddAccountCmd builds gen genesis account to genesis config
func GensisAddAccountCmd(ctx *server.Context, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-address",
		Short: "Add a genesis account to chain",
		Args:  cobra.ExactArgs(1),
		Long: `This command add a genesis account to chain'.

		It creates a genesis account which contains a name and auth, then put the data to genesis.json
	`,

		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr, err := types.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			genFile := config.GenesisFile()

			var accountGenesis account.GenesisState
			if err := types.LoadGenesisStateFromFile(cdc, genFile, account.ModuleName, &accountGenesis); err != nil {
				return err
			}

			if err = addGenesisAddAccount(cdc, &accountGenesis, addr); err != nil {
				return err
			}

			return types.SaveGenesisStateToFile(cdc, genFile, account.ModuleName, accountGenesis)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, app.DefaultCLIHome, "client's home directory")
	return cmd
}

func addGenesisAccount(cdc *codec.LegacyAmino, state *account.GenesisState, name types.Name, auth types.AccAddress) error {
	for _, stateAcc := range state.Accounts {
		if stateAcc.GetName().Eq(name) {
			return fmt.Errorf("the application state already contains account %s", name)
		}
	}

	newAccount := account.NewKuAccount(types.NewAccountIDFromName(name))

	newAccount.SetAuth(auth)
	newAccount.SetAccountNumber(uint64(len(state.Accounts) + 1))

	state.Accounts = append(state.Accounts, newAccount)

	return nil
}

func addGenesisAddAccount(cdc *codec.LegacyAmino, state *account.GenesisState, auth types.AccAddress) error {
	id := types.NewAccountIDFromAccAdd(auth)
	for _, stateAcc := range state.Accounts {
		if stateAcc.GetID().Eq(id) {
			return fmt.Errorf("the application state already contains account %s", id)
		}
	}

	newAccount := account.NewKuAccount(id)
	newAccount.SetAccountNumber(uint64(len(state.Accounts) + 1))

	state.Accounts = append(state.Accounts, newAccount)
	return nil
}
