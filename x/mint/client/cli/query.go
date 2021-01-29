package cli

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	"github.com/KuChainNetwork/kuchain/x/mint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for the minting module.
func GetQueryCmd(cdc *codec.LegacyAmino) *cobra.Command {
	mintingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the minting module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	mintingQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryParams(cdc),
			GetCmdQueryInflation(cdc),
			GetCmdQueryAnnualProvisions(cdc),
		)...,
	)

	return mintingQueryCmd
}

// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryParams(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current minting parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCtxByCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
}

// GetCmdQueryInflation implements a command to return the current minting
// inflation value.
func GetCmdQueryInflation(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "inflation",
		Short: "Query the current minting inflation value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCtxByCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryInflation)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var inflation sdk.Dec
			if err := cdc.UnmarshalJSON(res, &inflation); err != nil {
				return err
			}

			return cliCtx.PrintOutput(inflation)
		},
	}
}

// GetCmdQueryAnnualProvisions implements a command to return the current minting
// annual provisions value.
func GetCmdQueryAnnualProvisions(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "annual-provisions",
		Short: "Query the current minting annual provisions value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCtxByCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAnnualProvisions)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var inflation sdk.Dec
			if err := cdc.UnmarshalJSON(res, &inflation); err != nil {
				return err
			}

			return cliCtx.PrintOutput(inflation)
		},
	}
}
