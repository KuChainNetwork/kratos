package cli

import (
	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the transaction commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the asset module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCoinCmd(cdc),
		GetCoinsCmd(cdc),
		GetCoinPowerCmd(cdc),
		GetCoinPowersCmd(cdc),
		GetCoinsLockedCmd(cdc),
		GetCoinStatCmd(cdc),
	)

	return cmd
}

// GetCoinCmd returns a query coin
func GetCoinCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coin [account] [creator] [symbol]",
		Short: "Query coin for a account",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accGetter := types.NewAssetRetriever(cliCtx)

			key, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "account")
			}

			creator, err := chainTypes.NewName(args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "creator")
			}

			symbol, err := chainTypes.NewName(args[2])
			if err != nil {
				return sdkerrors.Wrap(err, "symbol")
			}

			coin, _, err := accGetter.GetCoin(key, creator, symbol)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(coin)
		},
	}

	return flags.GetCommands(cmd)[0]
}

// GetCoinCmd returns a query coin
func GetCoinsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coins [account]",
		Short: "Query all coins for a account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accGetter := types.NewAssetRetriever(cliCtx)

			key, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "account")
			}

			coin, _, err := accGetter.GetCoins(key)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(coin)
		},
	}

	return flags.GetCommands(cmd)[0]
}

// GetCoinCmd returns a query coin
func GetCoinPowerCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coinpower [account] [creator] [symbol]",
		Short: "Query coin power for a account",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accGetter := types.NewAssetRetriever(cliCtx)

			key, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "account")
			}

			creator, err := chainTypes.NewName(args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "creator")
			}

			symbol, err := chainTypes.NewName(args[2])
			if err != nil {
				return sdkerrors.Wrap(err, "symbol")
			}

			coin, _, err := accGetter.GetCoinPower(key, creator, symbol)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(coin)
		},
	}

	return flags.GetCommands(cmd)[0]
}

// GetCoinCmd returns a query coin
func GetCoinPowersCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coinpowers [account]",
		Short: "Query all coin powers for a account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accGetter := types.NewAssetRetriever(cliCtx)

			key, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "account")
			}

			coin, _, err := accGetter.GetCoinPowers(key)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(coin)
		},
	}

	return flags.GetCommands(cmd)[0]
}

// GetCoinStatCmd returns a query coin
func GetCoinStatCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [creator] [symbol]",
		Short: "Query coin status for creator/symbol token",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accGetter := types.NewAssetRetriever(cliCtx)

			name, err := chainTypes.NewName(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "name")
			}

			symbol, err := chainTypes.NewName(args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "symbol")
			}

			res, _, err := accGetter.GetCoinStat(name, symbol)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}

	return flags.GetCommands(cmd)[0]
}
