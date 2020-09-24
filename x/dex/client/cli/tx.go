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
		CreateCurrency(cdc),
		UpdateCurrency(cdc),
		PauseCurrency(cdc),
		RestoreCurrency(cdc),
		ShutdownCurrency(cdc),
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

// CreateCurrency returns a create currency command
func CreateCurrency(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "create_currency [creator] [base_code] [base_name] [base_full_name] [base_icon_url] [base_tx_url]" +
			" [quote_code] [quote_name] [quote_full_name] [quote_icon_url] [base_tx_url] [domain_address]",
		Short: "Create dex currency",
		Args:  cobra.ExactArgs(12),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			baseCode,
				baseName,
				baseFullName,
				baseIconUrl,
				baseTxUrl,
				quoteCode,
				quoteName,
				quoteFullName,
				quoteIconUrl,
				quoteTxUrl,
				domainAddress := args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11]

			if 0 >= len(baseCode) ||
				0 >= len(baseName) ||
				0 >= len(baseFullName) ||
				0 >= len(baseIconUrl) ||
				0 >= len(baseTxUrl) ||
				0 >= len(quoteCode) ||
				0 >= len(quoteName) ||
				0 >= len(quoteFullName) ||
				0 >= len(quoteIconUrl) ||
				0 >= len(quoteTxUrl) ||
				0 >= len(domainAddress) {
				err = errors.Errorf("all update failed are empty")
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				err = errors.Wrapf(err, "create dex currency account %s auth error", creator)
				return
			}

			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgCreateCurrency(
					auth,
					creator,
					&types.BaseCurrency{
						CurrencyBase: types.CurrencyBase{
							Code:     baseCode,
							Name:     baseName,
							FullName: baseFullName,
							IconUrl:  baseIconUrl,
							TxUrl:    baseTxUrl,
						},
					},
					&types.QuoteCurrency{
						CurrencyBase: types.CurrencyBase{
							Code:     quoteCode,
							Name:     quoteName,
							FullName: quoteFullName,
							IconUrl:  quoteIconUrl,
							TxUrl:    quoteTxUrl,
						},
					},
					domainAddress,
				),
			})

			return
		},
	}
	return flags.PostCommands(cmd)[0]
}

// UpdateCurrency returns a update currency command
func UpdateCurrency(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "update_currency [creator] [base_code] [base_name] [base_full_name] [base_icon_url] [base_tx_url]" +
			" [quote_code] [quote_name] [quote_full_name] [quote_icon_url] [base_tx_url]",
		Short: "Update dex currency",
		Args:  cobra.ExactArgs(11),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			baseCode,
				baseName,
				baseFullName,
				baseIconUrl,
				baseTxUrl,
				quoteCode,
				quoteName,
				quoteFullName,
				quoteIconUrl,
				quoteTxUrl := args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10]

			if 0 >= len(baseCode) || 0 >= len(quoteCode) {
				err = errors.Errorf("base code or quote code is empty")
				return
			}

			if 0 >= len(baseName) &&
				0 >= len(baseFullName) &&
				0 >= len(baseIconUrl) &&
				0 >= len(baseTxUrl) &&
				0 >= len(quoteName) &&
				0 >= len(quoteFullName) &&
				0 >= len(quoteIconUrl) &&
				0 >= len(quoteTxUrl) {
				err = errors.Errorf("all update failed are empty")
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				err = errors.Wrapf(err, "update dex currency account %s auth error", creator)
				return
			}

			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgUpdateCurrency(
					auth,
					creator,
					&types.BaseCurrency{
						CurrencyBase: types.CurrencyBase{
							Code:     baseCode,
							Name:     baseName,
							FullName: baseFullName,
							IconUrl:  baseIconUrl,
							TxUrl:    baseTxUrl,
						},
					},
					&types.QuoteCurrency{
						CurrencyBase: types.CurrencyBase{
							Code:     quoteCode,
							Name:     quoteName,
							FullName: quoteFullName,
							IconUrl:  quoteIconUrl,
							TxUrl:    quoteTxUrl,
						},
					},
				),
			})
			return
		},
	}
	return flags.PostCommands(cmd)[0]
}

// PauseCurrency returns a pause currency command
func PauseCurrency(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause_currency [creator] [base_code] [quote_code]",
		Short: "Pause dex currency",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			baseCode, quoteCode := args[1], args[2]
			if 0 >= len(baseCode) || 0 >= len(quoteCode) {
				err = errors.Errorf("base code or quote code is empty")
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				err = errors.Wrapf(err, "pause currency account %s auth error", creator)
				return
			}
			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgPauseCurrency(auth, creator, baseCode, quoteCode),
			})
			return
		},
	}
	return flags.PostCommands(cmd)[0]
}

// RestoreCurrency returns a restore currency command
func RestoreCurrency(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore_currency [creator] [base_code] [quote_code]",
		Short: "Restore dex currency",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			baseCode, quoteCode := args[1], args[2]
			if 0 >= len(baseCode) || 0 >= len(quoteCode) {
				err = errors.Errorf("base code or quote code is empty")
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				err = errors.Wrapf(err, "restore currency account %s auth error", creator)
				return
			}
			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgRestoreCurrency(auth, creator, baseCode, quoteCode),
			})
			return
		},
	}
	return flags.PostCommands(cmd)[0]
}

// ShutdownCurrency returns a shutdown currency command
func ShutdownCurrency(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shutdown_currency [creator] [base_code] [quote_code]",
		Short: "Shutdown dex currency",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var creator types.Name
			if creator, err = chainTypes.NewName(args[0]); nil != err {
				return
			}

			baseCode, quoteCode := args[1], args[2]
			if 0 >= len(baseCode) || 0 >= len(quoteCode) {
				err = errors.Errorf("base code or quote code is empty")
				return
			}

			ctx := txutil.NewKuCLICtx(cliCtx).WithAccount(creator)
			var auth chainTypes.AccAddress
			if auth, err = txutil.QueryAccountAuth(ctx, chainTypes.NewAccountIDFromName(creator)); nil != err {
				err = errors.Wrapf(err, "shutdown currency account %s auth error", creator)
				return
			}
			err = txutil.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{
				types.NewMsgShutdownCurrency(auth, creator, baseCode, quoteCode),
			})
			return
		},
	}
	return flags.PostCommands(cmd)[0]
}
