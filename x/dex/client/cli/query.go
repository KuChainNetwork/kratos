package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
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
		GetSigInCmd(cdc),
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
			cliCtx := client.NewCtxByCodec(cdc)
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
		Use:   "symbol [creator] [base creator] [base code] [quote creator] [quote code]",
		Short: "Query dex symbol",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := client.NewCtxByCodec(cdc)
			getter := types.NewDexRetriever(cliCtx)

			var creator types.Name
			creator, err = chainTypes.NewName(args[0])
			if err != nil {
				err = errors.Wrap(err, "creator")
				return
			}

			baseCreator, baseCode, quoteCreator, quoteCode := args[1], args[2], args[3], args[4]
			if 0 >= len(baseCreator) ||
				0 >= len(baseCode) ||
				0 >= len(quoteCreator) ||
				0 >= len(quoteCode) {
				err = errors.Errorf("base code or quote code is empty")
				return
			}

			var dex *types.Dex
			dex, _, err = getter.GetDexWithHeight(creator)
			if err != nil {
				return err
			}

			baseCreatorName, err := chainTypes.NewName(baseCreator)
			if err != nil {
				err = errors.Wrap(err, "base creator")
				return
			}
			baseCodeName, err := chainTypes.NewName(baseCode)
			if err != nil {
				err = errors.Wrap(err, "base code")
				return
			}
			baseCode = types.CoinDenom(baseCreatorName, baseCodeName)
			quoteCreatorName, err := chainTypes.NewName(quoteCreator)
			if err != nil {
				err = errors.Wrap(err, "quote creator")
				return
			}
			quoteCodeName, err := chainTypes.NewName(quoteCode)
			if err != nil {
				err = errors.Wrap(err, "quote code")
				return
			}
			quoteCode = types.CoinDenom(quoteCreatorName, quoteCodeName)
			symbol, ok := dex.Symbol(baseCreator, baseCode, quoteCreator, quoteCode)
			if ok {
				err = cliCtx.PrintOutput(symbol)
			}
			return
		},
	}
	return flags.GetCommands(cmd)[0]
}

// GetSignInCmd returns a query dex
func GetSigInCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSigIn [account] [dex]",
		Short: "Query sigIn status for account to dex",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCtxByCodec(cdc)
			getter := types.NewDexRetriever(cliCtx)

			acc, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return errors.Wrap(err, "acc")
			}

			dex, err := chainTypes.NewAccountIDFromStr(args[1])
			if err != nil {
				return errors.Wrap(err, "dex")
			}

			c, _, err := getter.GetSigInWithHeight(acc, dex)
			if err != nil {
				return errors.Wrap(err, "get")
			}

			return cliCtx.PrintOutput(c)
		},
	}

	return flags.GetCommands(cmd)[0]
}
