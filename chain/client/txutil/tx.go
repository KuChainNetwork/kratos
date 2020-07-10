package txutil

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// GasEstimateResponse defines a response definition for tx gas estimation.
type GasEstimateResponse struct {
	GasEstimate uint64 `json:"gas_estimate" yaml:"gas_estimate"`
}

func (gr GasEstimateResponse) String() string {
	return fmt.Sprintf("gas estimate: %d", gr.GasEstimate)
}

// GenerateOrBroadcastMsgs creates a StdTx given a series of messages. If
// the provided context has generate-only enabled, the tx will only be printed
// to STDOUT in a fully offline manner. Otherwise, the tx will be signed and
// broadcasted.
func GenerateOrBroadcastMsgs(cliCtx KuCLIContext, txBldr TxBuilder, msgs []sdk.Msg) error {
	if cliCtx.GenerateOnly {
		return PrintUnsignedStdTx(txBldr, cliCtx, msgs)
	}

	return CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
}

// CompleteAndBroadcastTxCLI implements a utility function that facilitates
// sending a series of messages in a signed transaction given a TxBuilder and a
// QueryContext. It ensures that the account exists, has a proper number and
// sequence set. In addition, it builds and signs a transaction with the
// supplied messages. Finally, it broadcasts the signed transaction to a node.
func CompleteAndBroadcastTxCLI(txBldr TxBuilder, cliCtx KuCLIContext, msgs []sdk.Msg) error {
	txBldr, err := PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err = EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return err
		}

		gasEst := GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	if cliCtx.Simulate {
		return nil
	}

	if !cliCtx.SkipConfirm {
		stdSignMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			return err
		}

		var json []byte
		if viper.GetBool(flags.FlagIndentResponse) {
			json, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				panic(err)
			}
		} else {
			json = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return err
		}
	}

	// build and sign the transaction, get from name is use to sign
	txBytes, err := txBldr.BuildAndSign(cliCtx.GetFromName(), keys.DefaultKeyPass, msgs)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return cliCtx.PrintOutput(res)
}

// EnrichWithGas calculates the gas estimate that would be consumed by the
// transaction and set the transaction's respective value accordingly.
func EnrichWithGas(txBldr TxBuilder, cliCtx KuCLIContext, msgs []sdk.Msg) (TxBuilder, error) {
	_, adjusted, err := simulateMsgs(txBldr, cliCtx, msgs)
	if err != nil {
		return txBldr, err
	}

	return txBldr.WithGas(adjusted), nil
}

// CalculateGas simulates the execution of a transaction and returns
// both the estimate obtained by the query and the adjusted amount.
func CalculateGas(
	queryFunc func(string, []byte) ([]byte, int64, error), cdc *codec.Codec,
	txBytes []byte, adjustment float64,
) (estimate, adjusted uint64, err error) {

	// run a simulation (via /app/simulate query) to
	// estimate gas and update TxBuilder accordingly
	rawRes, _, err := queryFunc("/app/simulate", txBytes)
	if err != nil {
		return estimate, adjusted, err
	}

	estimate, err = parseQueryResponse(cdc, rawRes)
	if err != nil {
		return
	}

	adjusted = adjustGasEstimate(estimate, adjustment)
	return estimate, adjusted, nil
}

// PrintUnsignedStdTx builds an unsigned StdTx and prints it to os.Stdout.
func PrintUnsignedStdTx(txBldr TxBuilder, cliCtx KuCLIContext, msgs []sdk.Msg) error {
	stdTx, err := buildUnsignedStdTxOffline(txBldr, cliCtx, msgs)
	if err != nil {
		return err
	}

	var json []byte
	if viper.GetBool(flags.FlagIndentResponse) {
		json, err = cliCtx.Codec.MarshalJSONIndent(stdTx, "", "  ")
	} else {
		json, err = cliCtx.Codec.MarshalJSON(stdTx)
	}
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(cliCtx.Output, "%s\n", json)
	return nil
}

// SignStdTx appends a signature to a StdTx and returns a copy of it. If appendSig
// is false, it replaces the signatures already attached with the new signature.
// Don't perform online validation or lookups if offline is true.
func SignStdTx(
	txBldr TxBuilder, cliCtx KuCLIContext, name string,
	stdTx StdTx, appendSig bool, offline bool,
) (StdTx, error) {
	// TODO: now the params has `name` and `id`, there should be both to make a id to get data to auth

	var signedStdTx StdTx

	// get auth to sign
	info, err := txBldr.Keybase().Get(name)
	if err != nil {
		return signedStdTx, err
	}

	addr := info.GetPubKey().Address()

	accAddr := sdk.AccAddress(addr)
	accountID := types.NewAccountIDFromAccAdd(accAddr)

	// check whether the address is a signer
	if !isTxSigner(accAddr, stdTx.GetSigners()) {
		return signedStdTx, fmt.Errorf("%s: %s", errInvalidSigner, name)
	}

	if !offline {
		txBldr, err = populateAccountFromState(txBldr, cliCtx, accountID)
		if err != nil {
			return signedStdTx, err
		}
	}

	return txBldr.SignStdTx(name, keys.DefaultKeyPass, stdTx, appendSig)
}

// SignStdTxWithSignerAddress attaches a signature to a StdTx and returns a copy of a it.
// Don't perform online validation or lookups if offline is true, else
// populate account and sequence numbers from a foreign account.
func SignStdTxWithSignerAddress(
	txBldr TxBuilder, cliCtx KuCLIContext,
	addr sdk.AccAddress, name string, stdTx StdTx, offline bool,
) (signedStdTx StdTx, err error) {

	// check whether the address is a signer
	if !isTxSigner(addr, stdTx.GetSigners()) {
		return signedStdTx, fmt.Errorf("%s: %s", errInvalidSigner, name)
	}

	if !offline {
		addrAccountID := types.NewAccountIDFromAccAdd(addr)
		txBldr, err = populateAccountFromState(txBldr, cliCtx, addrAccountID)
		if err != nil {
			return signedStdTx, err
		}
	}

	return txBldr.SignStdTx(name, keys.DefaultKeyPass, stdTx, false)
}

// Read and decode a StdTx from the given filename.  Can pass "-" to read from stdin.
func ReadStdTxFromFile(cdc *codec.Codec, filename string) (stdTx StdTx, err error) {
	var bytes []byte

	if filename == "-" {
		bytes, err = ioutil.ReadAll(os.Stdin)
	} else {
		bytes, err = ioutil.ReadFile(filename)
	}

	if err != nil {
		return
	}

	if err = cdc.UnmarshalJSON(bytes, &stdTx); err != nil {
		return
	}

	return
}

func populateAccountFromState(
	txBldr TxBuilder, cliCtx KuCLIContext, id types.AccountID,
) (TxBuilder, error) {

	num, seq, err := NewAccountRetriever(cliCtx).GetAuthNumberSequence(id)
	if err != nil {
		return txBldr, err
	}

	return txBldr.WithAccountNumber(num).WithSequence(seq), nil
}

// GetTxEncoder return tx encoder from global sdk configuration if ones is defined.
// Otherwise returns encoder with default logic.
func GetTxEncoder(cdc *codec.Codec) (encoder sdk.TxEncoder) {
	encoder = sdk.GetConfig().GetTxEncoder()
	if encoder == nil {
		encoder = DefaultTxEncoder(cdc)
	}

	return encoder
}

// nolint
// SimulateMsgs simulates the transaction and returns the gas estimate and the adjusted value.
func simulateMsgs(txBldr TxBuilder, cliCtx KuCLIContext, msgs []sdk.Msg) (estimated, adjusted uint64, err error) {
	txBytes, err := txBldr.BuildTxForSim(msgs)
	if err != nil {
		return
	}

	estimated, adjusted, err = CalculateGas(cliCtx.QueryWithData, cliCtx.Codec, txBytes, txBldr.GasAdjustment())
	return
}

func adjustGasEstimate(estimate uint64, adjustment float64) uint64 {
	return uint64(adjustment * float64(estimate))
}

func parseQueryResponse(cdc *codec.Codec, rawRes []byte) (uint64, error) {
	var gasUsed uint64
	if err := cdc.UnmarshalBinaryLengthPrefixed(rawRes, &gasUsed); err != nil {
		return 0, err
	}

	return gasUsed, nil
}

// PrepareTxBuilder populates a TxBuilder in preparation for the build of a Tx.
func PrepareTxBuilder(txBldr TxBuilder, cliCtx KuCLIContext) (TxBuilder, error) {
	from := cliCtx.GetAccountID()

	accGetter := NewAccountRetriever(cliCtx)

	if _, ok := from.ToAccAddress(); !ok {
		// if id is a address no need create account first
		fmt.Printf("need check if %s exit\b", from.String())
		if err := accGetter.EnsureExists(from); err != nil {
			return txBldr, err
		}
	}

	// TODO: Check if sign auth is account 's auth

	txbldrAccNum, txbldrAccSeq := txBldr.AccountNumber(), txBldr.Sequence()
	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if txbldrAccNum == 0 || txbldrAccSeq == 0 {
		num, seq, err := NewAccountRetriever(cliCtx).GetAuthNumberSequence(from)
		if err != nil {
			return txBldr, err
		}

		if txbldrAccNum == 0 {
			txBldr = txBldr.WithAccountNumber(num)
		}
		if txbldrAccSeq == 0 {
			txBldr = txBldr.WithSequence(seq)
		}
	}

	return txBldr, nil
}

// QueryAccountAuth query account auth by id
func QueryAccountAuth(cliCtx KuCLIContext, id types.AccountID) (types.AccAddress, error) {
	if cliCtx.GenerateOnly {
		// if just gen tx, cmd will not connect to node to get info, all auth is from --from params
		return cliCtx.FromAddress, nil
	}

	if _, ok := id.ToName(); ok {
		getter := NewAccountRetriever(cliCtx)
		acc, err := getter.GetAccount(id)
		if err != nil {
			return types.AccAddress{}, err
		}

		if acc == nil {
			return types.AccAddress{}, errors.New("account not found")
		}

		return acc.GetAuth(), nil
	}

	if auth, ok := id.ToAccAddress(); ok {
		return auth, nil
	}

	return types.AccAddress{}, errors.New("accountID type no support")
}

func buildUnsignedStdTxOffline(txBldr TxBuilder, cliCtx KuCLIContext, msgs []sdk.Msg) (stdTx StdTx, err error) {
	if txBldr.SimulateAndExecute() {
		if cliCtx.GenerateOnly {
			return stdTx, errors.New("cannot estimate gas with generate-only")
		}

		txBldr, err = EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return stdTx, err
		}

		_, _ = fmt.Fprintf(os.Stderr, "estimated gas = %v\n", txBldr.Gas())
	}

	stdSignMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return stdTx, err
	}

	return NewStdTx(stdSignMsg.Msg, stdSignMsg.Fee, nil, stdSignMsg.Memo), nil
}

func isTxSigner(user sdk.AccAddress, signers []sdk.AccAddress) bool {
	for _, s := range signers {
		if bytes.Equal(user.Bytes(), s.Bytes()) {
			return true
		}
	}

	return false
}
