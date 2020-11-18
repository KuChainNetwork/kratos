package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"

	accountGen "github.com/KuChainNetwork/kuchain/x/account/client/gen"
	assetGen "github.com/KuChainNetwork/kuchain/x/asset/client/gen"
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
		accountGen.GensisAccountCmd(ctx, cdc),
		accountGen.GensisAddAccountCmd(ctx, cdc),
		assetGen.GensisCoinCmd(ctx, cdc),
		assetGen.GensisAccountAssetCmd(ctx, cdc),
	)

	return genCmd
}
