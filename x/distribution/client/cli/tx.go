package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/client/common"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// flagOnlyFromValidator = "only-from-validator"
	// flagIsValidator       = "is-validator"
	flagCommission       = "commission"
	flagMaxMessagesPerTx = "max-msgs"
)

const (
	MaxMessagesPerTxDefault = 5
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	distTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Distribution transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	distTxCmd.AddCommand(flags.PostCommands(
		GetCmdWithdrawRewards(cdc),
		GetCmdSetWithdrawAddr(cdc),
		GetCmdWithdrawAllRewards(cdc, storeKey),
	)...)

	return distTxCmd
}

type generateOrBroadcastFunc func(txutil.KuCLIContext, txutil.TxBuilder, []sdk.Msg) error

func splitAndApply(
	generateOrBroadcast generateOrBroadcastFunc,
	cliCtx txutil.KuCLIContext,
	txBldr txutil.TxBuilder,
	msgs []sdk.Msg,
	chunkSize int) error {
	if chunkSize == 0 {
		return generateOrBroadcast(cliCtx, txBldr, msgs)
	}

	// split messages into slices of length chunkSize
	totalMessages := len(msgs)

	for i := 0; i < len(msgs); i += chunkSize {
		sliceEnd := i + chunkSize
		if sliceEnd > totalMessages {
			sliceEnd = totalMessages
		}

		msgChunk := msgs[i:sliceEnd]
		if err := generateOrBroadcast(cliCtx, txBldr, msgChunk); err != nil {
			return err
		}
	}

	return nil
}

// command to withdraw rewards
func GetCmdWithdrawRewards(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-rewards [validator]  --from delegator",
		Short: "Withdraw rewards from a given delegation address and validator commission from a validator operator.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw rewards from a given delegation name,
and optionally withdraw validator commission if the delegation address given is a validator operator.

Example:
$ %s tx kudistribution withdraw-rewards validator  Delegator --from jack
$ %s tx kudistribution withdraw-rewards validator  Delegator --from jack --commission
`,
				version.ClientName, version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := txutil.NewKuCLICtxByBuf(cdc, inBuf)

			delAddr := cliCtx.GetFromAddress()
			delID, _ := chainType.NewAccountIDFromStr(args[1])
			valID, err := chainType.NewAccountIDFromStr(args[0])
			if err != nil {
				return err
			}

			msgs := []sdk.Msg{types.NewMsgWithdrawDelegatorReward(delAddr, delID, valID)}
			if viper.GetBool(flagCommission) {
				msgs = append(msgs, types.NewMsgWithdrawValidatorCommission(delAddr, valID))
			}
			cliCtx = cliCtx.WithFromAccount(delID)

			fmt.Println(delID)
			return txutil.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
		},
	}
	cmd.Flags().Bool(flagCommission, false, "also withdraw validator's commission")
	return cmd
}

// command to withdraw all rewards
func GetCmdWithdrawAllRewards(cdc *codec.Codec, queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-all-rewards [delegator] --from [delegator]",
		Short: "withdraw all delegations rewards for a delegator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw all rewards for a single delegator.

Example:
$ %s tx kudistribution withdraw-all-rewards  --from jack
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := txutil.NewKuCLICtxByBuf(cdc, inBuf)

			delAddr := cliCtx.GetFromAddress()
			delID, _ := chainType.NewAccountIDFromStr(args[0])

			// The transaction cannot be generated offline since it requires a query
			// to get all the validators.
			if cliCtx.GenerateOnly {
				return fmt.Errorf("command disabled with the provided flag: %s", flags.FlagGenerateOnly)
			}

			msgs, err := common.WithdrawAllDelegatorRewards(cliCtx.CLIContext, delAddr, queryRoute, delID)
			if err != nil {
				return err
			}

			chunkSize := viper.GetInt(flagMaxMessagesPerTx)
			cliCtx = cliCtx.WithFromAccount(delID)
			fmt.Println("args:", args[0], "GetFromName:", cliCtx.GetFromName())
			return splitAndApply(txutil.GenerateOrBroadcastMsgs, cliCtx, txBldr, msgs, chunkSize)
		},
	}

	cmd.Flags().Int(flagMaxMessagesPerTx, MaxMessagesPerTxDefault, "Limit the number of messages per tx (0 for unlimited)")
	return cmd
}

// command to replace a delegator's withdrawal address
func GetCmdSetWithdrawAddr(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-withdraw [withdrawAccount] --from account",
		Short: "change the default withdraw name for rewards associated with an name ",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Set the withdraw name for rewards associated with a delegator .

Example:
$ %s tx kudistribution set-withdraw withdrawacc  account --from account
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := txutil.NewKuCLICtxByBuf(cdc, inBuf)

			delAddr := cliCtx.GetFromAddress()
			delID, _ := chainType.NewAccountIDFromStr(args[1])

			withdrawAccID, err := chainType.NewAccountIDFromStr(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetWithdrawAccountID(delAddr, delID, withdrawAccID)
			cliCtx = cliCtx.WithFromAccount(delID)
			return txutil.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdSubmitProposal implements the command to submit a community-pool-spend proposal
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community-pool-spend [proposer] [proposal-file]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit a community pool spend proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a community pool spend proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal community-pool-spend <path/to/proposal.json> --from=<key>

Where proposal.json contains:

{
  "title": "Community Pool Spend",
  "description": "Pay me some Atoms!",
  "recipient": "jack",
  "amount": [
    {
      "denom": "stake",
      "amount": "10000"
    }
  ],
  "deposit": [
    {
      "denom": "stake",
      "amount": "10000"
    }
  ]
}
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := txutil.NewKuCLICtxByBuf(cdc, inBuf)

			proposal, err := ParseCommunityPoolSpendProposalJSON(cdc, args[1])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewCommunityPoolSpendProposal(proposal.Title, proposal.Description, proposal.Recipient, proposal.Amount)
			proposerAccount, err := chainType.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "proposer account id error")
			}

			msg := types.GovTypesNewKuMsgSubmitProposal(from, content, proposal.Deposit, proposerAccount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			cliCtx = cliCtx.WithFromAccount(proposerAccount)
			return txutil.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
