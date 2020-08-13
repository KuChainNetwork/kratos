package cli

import (
	"bufio"

	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

// Exercise will commit a Exercise msg to chain
func Exercise(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exercise [accountID] [amount]",
		Short: "Exercise coins from coinpowers",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			id, err := types.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrapf(err, "account id %s parse error", args[0])
			}

			amount, err := chainTypes.ParseCoin(args[1])
			if err != nil {
				return err
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(id)
			auth, err := txutil.QueryAccountAuth(ctx, id)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", id)
			}

			msg := types.NewMsgExercise(auth, id, amount)
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]
	return cmd
}
