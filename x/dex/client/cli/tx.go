package cli

import (
	"bufio"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/dex/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Dex transactions sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		CreateDex(cdc),
		UpdateDexDescriptionCmd(cdc),
		DestroyDex(cdc),
	)

	return txCmd
}

// CreateDex will create a dex
func CreateDex(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [creator] [stakings] [description]",
		Short: "Create and sign a dex create trx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			creator, err := chainTypes.NewName(args[0])
			if err != nil {
				return err
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			auth, err := txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator))
			if err != nil {
				return errors.Wrapf(err, "query account %s auth error", creator)
			}

			stakings, err := chainTypes.ParseCoins(args[1])
			if err != nil {
				return err
			}

			if len(args[2]) > types.MaxDexDescriptorLen {
				return types.ErrDexDescTooLong
			}

			msg := types.NewMsgCreateDex(auth, creator, stakings, []byte(args[2]))
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

// UpdateDexDescriptionCmd returns a updated dex
func UpdateDexDescriptionCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update_description [creator] [description]",
		Short: "Update creator dex description",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				return errors.Wrapf(err, "update description account %s auth error", creator)
			}

			if len(args[1]) > types.MaxDexDescriptorLen {
				err = types.ErrDexDescTooLong
				return
			}

			msg := types.NewMsgUpdateDexDescription(auth, creator, []byte(args[1]))
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}
	return flags.PostCommands(cmd)[0]
}

// DestroyDex return a destroy dex command
func DestroyDex(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy [creator]",
		Short: "Destroy creator dex",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				err = errors.Wrapf(err, "destroy account %s auth error", creator)
				return
			}

			// get dex's staking
			getter := types.NewDexRetriever(cliCtx)
			var dex *types.Dex
			if dex, _, err = getter.GetDexWithHeight(creator); nil != err {
				return
			}

			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgDestroyDex(auth, creator, dex.Staking),
			})
			return
		},
	}
	return flags.PostCommands(cmd)[0]
}
