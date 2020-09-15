package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	stakingexport "github.com/KuChainNetwork/kuchain/x/staking/exported"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateValidator(cdc),
		GetCmdEditValidator(cdc),
		GetCmdDelegate(cdc),
		GetCmdRedelegate(storeKey, cdc),
		GetCmdUnbond(storeKey, cdc),
	)...)

	return stakingTxCmd
}

// GetCmdCreateValidator implements the create validator command handler.
func GetCmdCreateValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator [validator-operator-account]",
		Short: "create a new validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := txutil.NewKuCLICtxByBuf(cdc, inBuf)

			validatorAccount, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "validator account id error")
			}

			authAccAddress, err := txutil.QueryAccountAuth(cliCtx, validatorAccount)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", validatorAccount)
			}

			txBldr, msg, err := BuildCreateValidatorMsg(cliCtx, txBldr, validatorAccount, authAccAddress)
			if err != nil {
				return err
			}
			cliCtx = cliCtx.WithFromAccount(validatorAccount)

			if txBldr.FeePayer().Empty() {
				txBldr = txBldr.WithPayer(args[0])
			}
			return txutil.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsPk)
	cmd.Flags().AddFlagSet(fsDescriptionCreate)
	cmd.Flags().AddFlagSet(FsCommissionCreate)

	cmd.Flags().String(FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(FlagNodeID, "", "The node's ID")

	cmd.MarkFlagRequired(flags.FlagFrom)
	cmd.MarkFlagRequired(FlagPubKey)
	cmd.MarkFlagRequired(FlagMoniker)

	return cmd
}

// GetCmdEditValidator implements the create edit validator command.
// TODO: add full description
func GetCmdEditValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-validator [validator-operator-account]",
		Short: "edit an existing validator account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := txutil.NewKuCLICtxByBuf(cdc, inBuf)

			valAccount, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "validator operator accountID error")
			}

			description := types.NewDescription(
				viper.GetString(FlagMoniker),
				viper.GetString(FlagIdentity),
				viper.GetString(FlagWebsite),
				viper.GetString(FlagSecurityContact),
				viper.GetString(FlagDetails),
			)

			var newRate *sdk.Dec

			commissionRate := viper.GetString(FlagCommissionRate)
			if commissionRate != "" {
				rate, err := sdk.NewDecFromStr(commissionRate)
				if err != nil {
					return fmt.Errorf("invalid new commission rate: %v", err)
				}

				newRate = &rate
			}
			valAccAddress, err := txutil.QueryAccountAuth(cliCtx, valAccount)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", valAccount)
			}

			msg := types.NewKuMsgEditValidator(valAccAddress, valAccount, description, newRate)
			cliCtx = cliCtx.WithFromAccount(valAccount)
			if txBldr.FeePayer().Empty() {
				txBldr = txBldr.WithPayer(args[0])
			}
			// build and sign the transaction, then broadcast to Tendermint
			return txutil.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsDescriptionEdit)
	cmd.Flags().AddFlagSet(fsCommissionUpdate)
	cmd.Flags().AddFlagSet(FsMinSelfDelegation)

	return cmd
}

// GetCmdDelegate implements the delegate command.
func GetCmdDelegate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delegate [delegate-account] [validator-account] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Delegate liquid tokens to a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Delegate an amount of liquid coins to a validator from your wallet.

Example:
$ %s tx kustaking delegate jack validator 1000stake --from jack
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := txutil.NewKuCLICtxByBuf(cdc, inBuf)

			amount, err := chainTypes.ParseCoin(args[2])
			if err != nil {
				return sdkerrors.Wrap(err, "amount parse error")
			}

			delAccountID, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "delegate accountID error")
			}
			valAccountID, err := chainTypes.NewAccountIDFromStr(args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "validator accountID error")
			}

			delAccAddress, err := txutil.QueryAccountAuth(cliCtx, delAccountID)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", delAccountID)
			}

			msg := types.NewKuMsgDelegate(delAccAddress, delAccountID, valAccountID, amount)
			cliCtx = cliCtx.WithFromAccount(delAccountID)
			if txBldr.FeePayer().Empty() {
				txBldr = txBldr.WithPayer(args[0])
			}
			return txutil.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdRedelegate the begin redelegation command.
func GetCmdRedelegate(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "redelegate [delegate-account] [src-validator-account] [dst-validator-account] [amount]",
		Short: "Redelegate illiquid tokens from one validator to another",
		Args:  cobra.ExactArgs(4),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Redelegate an amount of illiquid staking tokens from one validator to another.

Example:
$ %s tx kustaking redelegate jack validator1 validator2 100stake --from jack
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := txutil.NewKuCLICtxByBuf(cdc, inBuf)

			delAccountID, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "delegate acccount error")
			}

			valSrcAccID, err := chainTypes.NewAccountIDFromStr(args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "src-validator error")
			}

			valDstAccID, err := chainTypes.NewAccountIDFromStr(args[2])
			if err != nil {
				return sdkerrors.Wrap(err, "dst-validator error")
			}

			amount, err := chainTypes.ParseCoin(args[3])
			if err != nil {
				return err
			}
			delAccAddress, err := txutil.QueryAccountAuth(cliCtx, delAccountID)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", delAccountID)
			}
			msg := types.NewKuMsgRedelegate(delAccAddress, delAccountID, valSrcAccID, valDstAccID, amount)
			cliCtx = cliCtx.WithFromAccount(delAccountID)
			if txBldr.FeePayer().Empty() {
				txBldr = txBldr.WithPayer(args[0])
			}
			return txutil.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdUnbond implements the unbond validator command.
func GetCmdUnbond(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unbond [delegate-account] [validator-account] [amount]",
		Short: "Unbond shares from a validator",
		Args:  cobra.ExactArgs(3),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unbond an amount of bonded shares from a validator.

Example:
$ %s tx kustaking unbond jack validator 100stake --from jack
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := txutil.NewKuCLICtxByBuf(cdc, inBuf)

			delAccountID, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "delegate account id error")
			}

			amount, err := chainTypes.ParseCoin(args[2])
			if err != nil {
				return err
			}

			valAddr, err := chainTypes.NewAccountIDFromStr(args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "val account id error")
			}
			delAccAddress, err := txutil.QueryAccountAuth(cliCtx, delAccountID)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", delAccountID)
			}

			msg := types.NewKuMsgUnbond(delAccAddress, delAccountID, valAddr, amount)
			cliCtx = cliCtx.WithFromAccount(delAccountID)
			if txBldr.FeePayer().Empty() {
				txBldr = txBldr.WithPayer(args[0])
			}
			return txutil.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

//__________________________________________________________

var (
	defaultTokens                  = stakingexport.TokensFromConsensusPower(100)
	defaultAmount                  = defaultTokens.String() + stakingexport.DefaultBondDenom
	defaultCommissionRate          = "0.1"
	defaultCommissionMaxRate       = "0.2"
	defaultCommissionMaxChangeRate = "0.01"
	defaultMinSelfDelegation       = "1"
)

// Return the flagset, particular flags, and a description of defaults
// this is anticipated to be used with the gen-tx
func CreateValidatorMsgHelpers(ipDefault string) (fs *flag.FlagSet, nodeIDFlag, pubkeyFlag, amountFlag, defaultsDesc string) {

	fsCreateValidator := flag.NewFlagSet("", flag.ContinueOnError)
	fsCreateValidator.String(FlagIP, ipDefault, "The node's public IP")
	fsCreateValidator.String(FlagNodeID, "", "The node's NodeID")
	fsCreateValidator.String(FlagWebsite, "", "The validator's (optional) website")
	fsCreateValidator.String(FlagSecurityContact, "", "The validator's (optional) security contact email")
	fsCreateValidator.String(FlagDetails, "", "The validator's (optional) details")
	fsCreateValidator.String(FlagIdentity, "", "The (optional) identity signature (ex. UPort or Keybase)")
	fsCreateValidator.AddFlagSet(FsCommissionCreate)
	fsCreateValidator.AddFlagSet(FsMinSelfDelegation)
	fsCreateValidator.AddFlagSet(FsAmount)
	fsCreateValidator.AddFlagSet(FsPk)

	defaultsDesc = fmt.Sprintf(`
	delegation amount:           %s
	commission rate:             %s
	commission max rate:         %s
	commission max change rate:  %s
	minimum self delegation:     %s
`, defaultAmount, defaultCommissionRate,
		defaultCommissionMaxRate, defaultCommissionMaxChangeRate,
		defaultMinSelfDelegation)

	return fsCreateValidator, FlagNodeID, FlagPubKey, FlagAmount, defaultsDesc
}

// prepare flags in config
func PrepareFlagsForTxCreateValidator(
	config *cfg.Config, nodeID, chainID string, valPubKey crypto.PubKey,
) {

	ip := viper.GetString(FlagIP)
	if ip == "" {
		fmt.Fprintf(os.Stderr, "couldn't retrieve an external IP; "+
			"the tx's memo field will be unset")
	}

	website := viper.GetString(FlagWebsite)
	securityContact := viper.GetString(FlagSecurityContact)
	details := viper.GetString(FlagDetails)
	identity := viper.GetString(FlagIdentity)

	viper.Set(flags.FlagChainID, chainID)
	viper.Set(flags.FlagFrom, viper.GetString(flags.FlagName))
	viper.Set(FlagNodeID, nodeID)
	viper.Set(FlagIP, ip)
	viper.Set(FlagPubKey, sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, valPubKey))
	viper.Set(FlagMoniker, config.Moniker)
	viper.Set(FlagWebsite, website)
	viper.Set(FlagSecurityContact, securityContact)
	viper.Set(FlagDetails, details)
	viper.Set(FlagIdentity, identity)

	if config.Moniker == "" {
		viper.Set(FlagMoniker, viper.GetString(flags.FlagName))
	}
	if viper.GetString(FlagAmount) == "" {
		viper.Set(FlagAmount, defaultAmount)
	}
	if viper.GetString(FlagCommissionRate) == "" {
		viper.Set(FlagCommissionRate, defaultCommissionRate)
	}
}

// BuildCreateValidatorMsg makes a new MsgCreateValidator.
func BuildCreateValidatorMsg(cliCtx txutil.KuCLIContext, txBldr txutil.TxBuilder, valAddr chainTypes.AccountID, authAddress sdk.AccAddress) (txutil.TxBuilder, sdk.Msg, error) {
	delAddr := chainTypes.NewAccountIDFromAccAdd(authAddress)
	pkStr := viper.GetString(FlagPubKey)

	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pkStr)
	if err != nil {
		return txBldr, nil, err
	}

	description := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagSecurityContact),
		viper.GetString(FlagDetails),
	)

	// get the initial validator commission parameters
	rateStr := viper.GetString(FlagCommissionRate)
	rate, err := sdk.NewDecFromStr(rateStr)

	if err != nil {
		return txBldr, nil, err
	}

	msg := types.NewKuMsgCreateValidator(authAddress,
		valAddr, pk, description, rate, delAddr,
	)

	if viper.GetBool(flags.FlagGenerateOnly) {
		ip := viper.GetString(FlagIP)
		nodeID := viper.GetString(FlagNodeID)
		if nodeID != "" && ip != "" {
			txBldr = txBldr.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
		}
	}

	return txBldr, msg, nil
}

func BuildDelegateMsg(cliCtx txutil.KuCLIContext, txBldr txutil.TxBuilder, authAddress chainTypes.AccAddress, delAccountID chainTypes.AccountID, valAccountID chainTypes.AccountID) (txutil.TxBuilder, sdk.Msg, error) {

	defaultAmount = stakingexport.TokensFromConsensusPower(1).String() + stakingexport.DefaultBondDenom
	amount, err := chainTypes.ParseCoin(defaultAmount)
	if err != nil {
		return txBldr, nil, err
	}

	msg := types.NewKuMsgDelegate(authAddress, delAccountID, valAccountID, amount)

	return txBldr, msg, nil
}
