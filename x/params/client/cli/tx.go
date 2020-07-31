package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	paramscutils "github.com/KuChainNetwork/kuchain/x/params/client/utils"
	"github.com/KuChainNetwork/kuchain/x/params/external"
	paramproposal "github.com/KuChainNetwork/kuchain/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

// GetCmdSubmitProposal implements a command handler for submitting a parameter
// change proposal transaction.
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "param-change [proposaler] [proposal-file]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit a parameter change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a parameter proposal along with an initial deposit.
The proposal details must be supplied via a JSON file. For values that contains
objects, only non-empty fields will be updated.

note:you can only change params in module kugov

Example:
$ %s tx kugov submit-proposal param-change <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Staking Param Change",
  "description": "Update voting period",
  "changes": [
    {
      "subspace": "kugov",
      "key": "votingparams",
      "value":{
        "voting_period": "1209800000000000"
      }
    }
  ],
  "deposit": "1000kuchain/kcs"
}
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
			cliCtx := txutil.NewKuCLICtxByBuf(cdc, inBuf)

			proposal, err := paramscutils.ParseParamChangeProposalJSON(cdc, args[1])
			if err != nil {
				return err
			}
			ProposalAccount, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "depositor account id error")
			}

			content := paramproposal.NewParameterChangeProposal(proposal.Title, proposal.Description, proposal.Changes.ToParamChanges())
			deposit, err := chainTypes.ParseCoins(proposal.Deposit)
			if err != nil {
				return err
			}
			// Get proposal address
			authAccAddress, err := txutil.QueryAccountAuth(cliCtx, ProposalAccount)
			if err != nil {
				return sdkerrors.Wrapf(err, "query account %s auth error", ProposalAccount)
			}
			msg := external.GovNewMsgSubmitProposal(authAccAddress, content, deposit, ProposalAccount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			cliCtx = cliCtx.WithFromAccount(ProposalAccount)
			if txBldr.FeePayer().Empty() {
				txBldr = txBldr.WithPayer(args[0])
			}
			return txutil.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
