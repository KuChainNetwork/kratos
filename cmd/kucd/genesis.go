package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"

	accountGen "github.com/KuChain-io/kuchain/x/account/client/gen"
	assetGen "github.com/KuChain-io/kuchain/x/asset/client/gen"
)

// AddGenesisCmds
func AddGenesisCmds(
	ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string,
) *cobra.Command {
	genCmd := &cobra.Command{
		Use:                        "genesis",
		Short:                      "make genesis config sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	genCmd.AddCommand(
		accountGen.GenGensisAccountCmd(ctx, cdc),
		accountGen.GenGensisAddAccountCmd(ctx, cdc),
		assetGen.GenGensisCoinCmd(ctx, cdc),
		assetGen.GenGensisAccountAssetCmd(ctx, cdc),
	)

	return genCmd
}
