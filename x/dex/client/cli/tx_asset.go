package cli

import (
	"bufio"

	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// SigInCmd will create a dex
func SigInCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sigIn [account] [dex] [amt]",
		Short: "Create and sign a sigIn msg",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			account, err := types.NewAccountIDFromStr(args[0])
			if err != nil {
				return err
			}

			dex, err := types.NewAccountIDFromStr(args[1])
			if err != nil {
				return err
			}

			amt, err := chainTypes.ParseCoins(args[2])
			if err != nil {
				return err
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(account)
			auth, err := txutil.QueryAccountAuth(ctx, account)
			if err != nil {
				return errors.Wrapf(err, "query account %s auth error", account)
			}

			msg := types.NewMsgDexSigIn(auth, account, dex, amt)
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

// SigOutCmd will create a dex
func SigOutCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sigOut [account] [dex] [amt]",
		Short: "Create and sign a sigOut msg",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			account, err := types.NewAccountIDFromStr(args[0])
			if err != nil {
				return err
			}

			dex, err := types.NewAccountIDFromStr(args[1])
			if err != nil {
				return err
			}

			amt, err := chainTypes.ParseCoins(args[2])
			if err != nil {
				return err
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(dex)
			auth, err := txutil.QueryAccountAuth(ctx, dex)
			if err != nil {
				return errors.Wrapf(err, "query account %s auth error", dex)
			}

			msg := types.NewMsgDexSigOut(auth, false, account, dex, amt)
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}
