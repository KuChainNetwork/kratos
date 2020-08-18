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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Create will create a account create tx and sign it with the given key.
func Create(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [creator] [symbol] [max_supply] [canIssue] [canLock] [canBurn] [issueToHeight] [initSupply] [desc]",
		Short: "Create coin, for canIssue, canLock and canBurn, the 1 means true.",
		Args:  cobra.ExactArgs(9),
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

			maxSupply, err := chainTypes.ParseCoin(args[2])
			if err != nil {
				return err
			}

			isCanIssue := args[3] == "1"
			isCanLock := args[4] == "1"
			isCanBurn := args[5] == "1"
			issueToHeight, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return err
			}

			initSupply, err := chainTypes.ParseCoin(args[7])
			if err != nil {
				return errors.Wrap(err, "init supply parse error")
			}

			if chainTypes.CoinDenom(creator, symbol) != maxSupply.Denom {
				return fmt.Errorf("coin denom should equal %s != %s",
					chainTypes.CoinDenom(creator, symbol), maxSupply.Denom)
			}

			if maxSupply.Denom != initSupply.Denom {
				return fmt.Errorf("init coin denom should equal %s != %s",
					initSupply.Denom, maxSupply.Denom)
			}

			desc := args[8]
			if len(desc) > types.CoinDescriptionLen {
				return fmt.Errorf("coin desc too long, should be less than %d", types.CoinDescriptionLen)
			}

			msg := types.NewMsgCreate(auth, creator, symbol, maxSupply, isCanIssue, isCanLock, isCanBurn, issueToHeight, initSupply, []byte(desc))
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}
