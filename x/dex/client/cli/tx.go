package cli

import (
	"bufio"
	"time"

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
		CreateSymbol(cdc),
		UpdateSymbol(cdc),
		PauseSymbol(cdc),
		RestoreSymbol(cdc),
		ShutdownSymbol(cdc),
		SigInCmd(cdc),
		SigOutCmd(cdc),
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
				return errors.Wrapf(err, "update dex description account %s auth error", creator)
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
				err = errors.Wrapf(err, "destroy dex account %s auth error", creator)
				return
			}

			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgDestroyDex(auth, creator),
			})
			return
		},
	}
	return flags.PostCommands(cmd)[0]
}

// CreateSymbol returns a create symbol command
func CreateSymbol(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "create_symbol [creator] [base_creator] [base_code] [base_name] [base_full_name] [base_icon_url] [base_tx_url]" +
			" [quote_creator] [quote_code] [quote_name] [quote_full_name] [quote_icon_url] [quote_tx_url]",
		Short: "Create dex symbol",
		Args:  cobra.ExactArgs(13),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			baseCreator,
				baseCode,
				baseName,
				baseFullName,
				baseIconURL,
				baseTxURL,
				quoteCreator,
				quoteCode,
				quoteName,
				quoteFullName,
				quoteIconURL,
				quoteTxURL := args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12]

			if 0 >= len(baseCreator) ||
				0 >= len(baseCode) ||
				0 >= len(baseName) ||
				0 >= len(baseFullName) ||
				0 >= len(quoteCreator) ||
				0 >= len(baseIconURL) ||
				0 >= len(baseTxURL) ||
				0 >= len(quoteCode) ||
				0 >= len(quoteName) ||
				0 >= len(quoteFullName) ||
				0 >= len(quoteIconURL) ||
				0 >= len(quoteTxURL) {
				err = errors.Errorf("all update failed are empty")
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				err = errors.Wrapf(err, "create dex symbol account %s auth error", creator)
				return
			}

			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgCreateSymbol(
					auth,
					creator,
					&types.BaseCurrency{
						CurrencyBase: types.CurrencyBase{
							Creator:  baseCreator,
							Code:     baseCode,
							Name:     baseName,
							FullName: baseFullName,
							IconURL:  baseIconURL,
							TxURL:    baseTxURL,
						},
					},
					&types.QuoteCurrency{
						CurrencyBase: types.CurrencyBase{
							Creator:  quoteCreator,
							Code:     quoteCode,
							Name:     quoteName,
							FullName: quoteFullName,
							IconURL:  quoteIconURL,
							TxURL:    quoteTxURL,
						},
					},
					time.Time{}, // use server time
				),
			})

			return
		},
	}
	return flags.PostCommands(cmd)[0]
}

// UpdateSymbol returns a update symbol command
func UpdateSymbol(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "update_symbol [creator] [base_creator] [base_code] [base_name] [base_full_name] [base_icon_url] [base_tx_url]" +
			" [quote_creator] [quote_code] [quote_name] [quote_full_name] [quote_icon_url] [quote_tx_url]",
		Short: "Update dex symbol",
		Args:  cobra.ExactArgs(13),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			baseCreator,
				baseCode,
				baseName,
				baseFullName,
				baseIconURL,
				baseTxURL,
				quoteCreator,
				quoteCode,
				quoteName,
				quoteFullName,
				quoteIconURL,
				quoteTxURL := args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12]

			if 0 >= len(baseCode) || 0 >= len(quoteCode) {
				err = errors.Errorf("base code or quote code is empty")
				return
			}

			if 0 >= len(baseCreator) &&
				0 >= len(baseName) &&
				0 >= len(baseFullName) &&
				0 >= len(baseIconURL) &&
				0 >= len(baseTxURL) &&
				0 >= len(quoteCreator) &&
				0 >= len(quoteName) &&
				0 >= len(quoteFullName) &&
				0 >= len(quoteIconURL) &&
				0 >= len(quoteTxURL) {
				err = errors.Errorf("all update failed are empty")
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				err = errors.Wrapf(err, "update dex symbol account %s auth error", creator)
				return
			}

			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgUpdateSymbol(
					auth,
					creator,
					&types.BaseCurrency{
						CurrencyBase: types.CurrencyBase{
							Creator:  baseCreator,
							Code:     baseCode,
							Name:     baseName,
							FullName: baseFullName,
							IconURL:  baseIconURL,
							TxURL:    baseTxURL,
						},
					},
					&types.QuoteCurrency{
						CurrencyBase: types.CurrencyBase{
							Creator:  quoteCreator,
							Code:     quoteCode,
							Name:     quoteName,
							FullName: quoteFullName,
							IconURL:  quoteIconURL,
							TxURL:    quoteTxURL,
						},
					},
				),
			})
			return
		},
	}
	return flags.PostCommands(cmd)[0]
}

// PauseSymbol returns a pause symbol command
func PauseSymbol(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause_symbol [creator] [base_creator] [base_code] [quote_creator] [quote_code]",
		Short: "Pause dex symbol",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			baseCreator, baseCode, quoteCreator, quoteCode := args[1], args[2], args[3], args[4]
			if 0 >= len(baseCreator) ||
				0 >= len(baseCode) ||
				0 >= len(quoteCreator) ||
				0 >= len(quoteCode) {
				err = errors.Errorf("base creator code or quote creator code is empty")
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				err = errors.Wrapf(err, "pause symbol account %s auth error", creator)
				return
			}
			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgPauseSymbol(auth, creator, baseCreator, baseCode, quoteCreator, quoteCode),
			})
			return
		},
	}
	return flags.PostCommands(cmd)[0]
}

// RestoreSymbol returns a restore symbol command
func RestoreSymbol(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore_symbol [creator] [base_creator] [base_code] [quote_creator] [quote_code]",
		Short: "Restore dex symbol",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			baseCreator, baseCode, quoteCreator, quoteCode := args[1], args[2], args[3], args[4]
			if 0 >= len(baseCreator) ||
				0 >= len(baseCode) ||
				0 >= len(quoteCreator) ||
				0 >= len(quoteCode) {
				err = errors.Errorf("base creator code or quote creator code is empty")
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				err = errors.Wrapf(err, "restore symbol account %s auth error", creator)
				return
			}
			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgRestoreSymbol(auth, creator, baseCreator, baseCode, quoteCreator, quoteCode),
			})
			return
		},
	}
	return flags.PostCommands(cmd)[0]
}

// ShutdownSymbol returns a shutdown symbol command
func ShutdownSymbol(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shutdown_symbol [creator] [base_creator] [base_code] [quote_creator] [quote_code]",
		Short: "Shutdown dex symbol",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			baseCreator, baseCode, quoteCreator, quoteCode := args[1], args[2], args[3], args[4]
			if 0 >= len(baseCreator) ||
				0 >= len(baseCode) ||
				0 >= len(quoteCreator) ||
				0 >= len(quoteCode) {
				err = errors.Errorf("base creator code or quote creator code is empty")
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				err = errors.Wrapf(err, "shutdown symbol account %s auth error", creator)
				return
			}
			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgShutdownSymbol(auth, creator, baseCreator, baseCode, quoteCreator, quoteCode),
			})
			return
		},
	}
	return flags.PostCommands(cmd)[0]
}
