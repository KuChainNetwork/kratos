package cli

import (
	"bufio"

	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Asset transactions sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		Transfer(cdc),
		Create(cdc),
		Issue(cdc),
		Burn(cdc),
		LockCoin(cdc),
		UnlockCoin(cdc),
		Exercise(cdc),
	)

	return txCmd
}

// Transfer will create a account create tx and sign it with the given key.
func Transfer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [from] [to] [coins]",
		Short: "Transfer coins and sign a trx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			from, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return err
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(from)
			authAddress, err := txutil.QueryAccountAuth(ctx, from)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", from)
			}

			to, err := chainTypes.NewAccountIDFromStr(args[1])
			if err != nil {
				return err
			}

			coin, err := chainTypes.ParseCoins(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgTransfer(authAddress, from, to, coin)
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}
