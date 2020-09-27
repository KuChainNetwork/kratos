package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetQueryCmd returns the transaction commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the dex module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetDexCmd(cdc),
		GetSymbol(cdc),
	)

	return cmd
}

// GetDexCmd returns a query dex
func GetDexCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [creator]",
		Short: "Query dex for creator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			getter := types.NewDexRetriever(cliCtx)

			creator, err := chainTypes.NewName(args[0])
			if err != nil {
				return errors.Wrap(err, "creator")
			}

			dex, _, err := getter.GetDexWithHeight(creator)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(dex)
		},
	}

	return flags.GetCommands(cmd)[0]
}

// GetSymbol returns a query symbol
func GetSymbol(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "symbol [creator] [base code] [quote code]",
		Short: "Query dex symbol",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			getter := types.NewDexRetriever(cliCtx)

			var creator types.Name
			creator, err = chainTypes.NewName(args[0])
			if err != nil {
				err = errors.Wrap(err, "creator")
				return
			}

			if 0 >= len(args[1]) || 0 >= len(args[2]) {
				err = errors.Errorf("base code or quote code is empty")
				return
			}

			var dex *types.Dex
			dex, _, err = getter.GetDexWithHeight(creator)
			if err != nil {
				return err
			}

			symbol, ok := dex.Symbol(args[1], args[2])
			if ok {
				err = cliCtx.PrintOutput(symbol)
			}
			return
		},
	}
	return flags.GetCommands(cmd)[0]
}
