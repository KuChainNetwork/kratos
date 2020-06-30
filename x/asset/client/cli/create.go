package cli

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/KuChain-io/kuchain/chain/client/flags"
	"github.com/KuChain-io/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/asset/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Create will create a account create tx and sign it with the given key.
func Create(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [creator] [symbol] [max_supply] [canIssue] [canLock] [issueToHeight] [initSupply] [desc]",
		Short: "Create coin, if canIssue is 1 or canLock is 1, the coin cannot issue or lock after 64 blocks",
		Args:  cobra.ExactArgs(8),
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
				return sdkerrors.Wrapf(err, "query account %s auth error", creatorID)
			}

			symbol, err := chainTypes.NewName(args[1])
			if err != nil {
				return err
			}

			maxSupply, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}

			isCanIssue := args[3] == "1"
			isCanLock := args[4] == "1"
			issueToHeight, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				return err
			}

			initSupply, err := sdk.ParseCoin(args[6])
			if err != nil {
				return errors.Wrap(err, "init supply parse error")
			}

			if chainTypes.CoinDenom(creator, symbol) != maxSupply.GetDenom() {
				return fmt.Errorf("coin denom should equal %s != %s",
					chainTypes.CoinDenom(creator, symbol), maxSupply.GetDenom())
			}

			if maxSupply.GetDenom() != initSupply.GetDenom() {
				return fmt.Errorf("init coin denom should equal %s != %s",
					initSupply.GetDenom(), maxSupply.GetDenom())
			}

			desc := args[7]
			if len(desc) > types.CoinDescriptionLen {
				return fmt.Errorf("coin desc too long, should be less than %d", types.CoinDescriptionLen)
			}

			msg := types.NewMsgCreate(auth, creator, symbol, maxSupply, isCanIssue, isCanLock, issueToHeight, initSupply, []byte(desc))
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}
