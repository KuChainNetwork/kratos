package cli

import (
	"bufio"
	"fmt"

	"github.com/KuChain-io/kuchain/chain/client/flags"
	"github.com/KuChain-io/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/asset/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

// Issue will create a account create tx and sign it with the given key.
func Issue(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue [creator] [symbol] [amount]",
		Short: "Issue coin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			creator, err := chainTypes.NewName(args[0])
			if err != nil {
				return err
			}

			creatorID := types.NewAccountIDFromName(creator)

			ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(creatorID)
			auth, err := txutil.QueryAccountAuth(ctx, creatorID)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", creator)
			}

			symbol, err := chainTypes.NewName(args[1])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}

			if chainTypes.CoinDenom(creator, symbol) != amount.GetDenom() {
				return fmt.Errorf("coin denom should equal %s != %s",
					chainTypes.CoinDenom(creator, symbol), amount.GetDenom())
			}

			msg := types.NewMsgIssue(auth, creator, symbol, amount)
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}
