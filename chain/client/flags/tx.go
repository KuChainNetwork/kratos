package flags

import (
	"github.com/KuChainNetwork/kuchain/chain/transaction"
	cosmosFlags "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// PostCommands adds common flags for commands to post tx
func PostCommands(cmds ...*cobra.Command) []*cobra.Command {
	for _, c := range cmds {
		c.Flags().String(transaction.FlagPayer, "", "fee payer for tx")
	}

	return cosmosFlags.PostCommands(cmds...)
}
