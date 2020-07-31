package cli

import (
	"bufio"
	"fmt"
	"strconv"

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

// LockCoin will create a account create tx and sign it with the given key.
func LockCoin(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock [accountID] [unlockBlockHeight] [amount]",
		Short: "Lock coin until the unlockBlockHeight is arrived",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			id, err := types.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrapf(err, "account id %s parse error", args[0])
			}

			unlockBlockHeight, err := strconv.Atoi(args[1])
			if err != nil {
				return sdkerrors.Wrapf(err, "unlock block height parse error")
			}

			amount, err := chainTypes.ParseCoins(args[2])
			if err != nil {
				return err
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(id)
			auth, err := txutil.QueryAccountAuth(ctx, id)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", id)
			}

			fmt.Printf("auth %s\n", auth.String())

			msg := types.NewMsgLockCoin(auth, id, amount, int64(unlockBlockHeight))
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]
	return cmd
}

// UnlockCoin will create a account create tx and sign it with the given key.
func UnlockCoin(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlock [accountID] [amount]",
		Short: "Unlock all coin can unlock",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			id, err := types.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrapf(err, "account id %s parse error", args[0])
			}

			amount, err := chainTypes.ParseCoins(args[1])
			if err != nil {
				return err
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(id)
			auth, err := txutil.QueryAccountAuth(ctx, id)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", id)
			}

			msg := types.NewMsgUnlockCoin(auth, id, amount)
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]
	return cmd
}

// GetCoinsLockedCmd returns a query coin locked
func GetCoinsLockedCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locked [account]",
		Short: "Query all coin powers for a account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accGetter := types.NewAssetRetriever(cliCtx)

			key, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "account")
			}

			res, _, err := accGetter.GetLockedCoins(key)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}

	return flags.GetCommands(cmd)[0]
}
