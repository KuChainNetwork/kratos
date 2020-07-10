package cli

import (
	"bufio"

	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/types"
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
		Short:                      "Account transactions sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		CreateAccount(cdc),
		UpdateAccountAuth(cdc),
	)

	return txCmd
}

// CreateAccount will create a account create tx and sign it with the given key.
func CreateAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [creator] [new_account_name] [new_account_owner_auth]",
		Short: "Create and sign a account create trx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			creator, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return err
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(creator)
			auth, err := txutil.QueryAccountAuth(ctx, creator)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", creator)
			}

			accountName, err := chainTypes.NewName(args[1])
			if err != nil {
				return err
			}

			accountAuth, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateAccount(auth, creator, accountName, accountAuth)
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

// UpdateAccountAuth will update auth for a account
func UpdateAccountAuth(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateauth [account_name] [new_account_owner_auth]",
		Short: "update account auth for a account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			accountName, err := chainTypes.NewName(args[0])
			if err != nil {
				return err
			}

			id := chainTypes.NewAccountIDFromName(accountName)

			ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(id)
			auth, err := txutil.QueryAccountAuth(ctx, id)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", id)
			}
			accountAuth, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateAccountAuth(auth, accountName, accountAuth)
			return txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}
