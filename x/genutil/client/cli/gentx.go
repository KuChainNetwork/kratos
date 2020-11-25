package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/genutil"
	"github.com/KuChainNetwork/kuchain/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	flagSignGenesisTrx = "sign"
)

// StakingMsgBuildingHelpers helpers for message building gen-tx command
type StakingMsgBuildingHelpers interface {
	CreateValidatorMsgHelpers(ipDefault string) (fs *flag.FlagSet, nodeIDFlag, pubkeyFlag, amountFlag, defaultsDesc string)
	PrepareFlagsForTxCreateValidator(config *cfg.Config, nodeID, chainID string, valPubKey crypto.PubKey)

	BuildCreateValidatorMsg(cliCtx txutil.KuCLIContext, txBldr txutil.TxBuilder,
		operAccountID chainTypes.AccountID, authAddress sdk.AccAddress) (txutil.TxBuilder, sdk.Msg, error)
	BuildDelegateMsg(cliCtx txutil.KuCLIContext, txBldr txutil.TxBuilder,
		authAddress chainTypes.AccAddress, delAccountID, valAccountID chainTypes.AccountID) (txutil.TxBuilder, sdk.Msg, error)
}

func genTxRunE(ctx *server.Context, cdc *codec.Codec,
	mbm module.BasicManager, smbh StakingMsgBuildingHelpers,
	flagNodeID, flagPubKey string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		config := ctx.Config
		config.SetRoot(viper.GetString(flags.FlagHome))
		nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(ctx.Config)
		if err != nil {
			return errors.Wrap(err, "failed to initialize node validator files")
		}

		// Read --nodeID, if empty take it from priv_validator.json
		if nodeIDString := viper.GetString(flagNodeID); nodeIDString != "" {
			nodeID = nodeIDString
		}

		// Read --pubkey, if empty take it from priv_validator.json
		if valPubKeyString := viper.GetString(flagPubKey); valPubKeyString != "" {
			valPubKey, err = sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, valPubKeyString)
			if err != nil {
				return errors.Wrap(err, "failed to get consensus node public key")
			}
		}

		genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
		if err != nil {
			return errors.Wrapf(err, "failed to read genesis doc file %s", config.GenesisFile())
		}

		var genesisState map[string]json.RawMessage
		if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
			return errors.Wrap(err, "failed to unmarshal genesis state")
		}

		if err = mbm.ValidateGenesis(genesisState); err != nil {
			return errors.Wrap(err, "failed to validate genesis state")
		}

		inBuf := bufio.NewReader(cmd.InOrStdin())
		_, err = keys.NewKeyring(sdk.KeyringServiceName(),
			viper.GetString(flags.FlagKeyringBackend), viper.GetString(flagClientHome), inBuf)
		if err != nil {
			return errors.Wrap(err, "failed to initialize keybase")
		}
		authAccAddress := chainTypes.MustAccAddressFromBech32(args[1])

		// Set flags for creating gentx
		viper.Set(flags.FlagHome, viper.GetString(flagClientHome))
		smbh.PrepareFlagsForTxCreateValidator(config, nodeID, genDoc.ChainID, valPubKey)

		// Set the generate-only flag here after the CLI context has
		// been created. This allows the from name/key to be correctly populated.
		viper.Set(flags.FlagGenerateOnly, true)
		viper.Set(flags.FlagFrom, args[1])

		cliCtx := txutil.NewKuCLICtxNoFrom(context.NewCLIContextWithInput(inBuf).WithCodec(cdc))
		valAccountID, err := chainTypes.NewAccountIDFromStr(args[0])
		if err != nil {
			return errors.Wrap(err, "Invalid validator account ID")
		}

		// create a 'create-validator' message
		txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc))
		txBldr, msg, err := smbh.BuildCreateValidatorMsg(cliCtx, txBldr, valAccountID, authAccAddress)
		if err != nil {
			return errors.Wrap(err, "failed to build create-validator message")
		}

		txBldr, msgdelegator, err := smbh.BuildDelegateMsg(cliCtx, txBldr, authAccAddress, valAccountID, valAccountID)
		if err != nil {
			return errors.Wrap(err, "failed to build create-validator message")
		}

		// set payer in gentx
		txBldr = txBldr.WithPayer(args[0])

		// write the unsigned transaction to the buffer
		w := bytes.NewBuffer([]byte{})
		cliCtx = cliCtx.WithOutput(w)

		if err = txutil.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg, msgdelegator}); err != nil {
			return errors.Wrap(err, "failed to print unsigned std tx")
		}

		// read the transaction
		stdTx, err := readUnsignedGenTxFile(cdc, w)
		if err != nil {
			return errors.Wrap(err, "failed to read unsigned gen tx file")
		}

		if viper.GetBool(flagSignGenesisTrx) {
			// sign the transaction and write it to the output file
			signedTx, err := txutil.SignStdTx(txBldr, cliCtx, viper.GetString(flags.FlagName), stdTx, false, true)
			if err != nil {
				return errors.Wrap(err, "failed to sign std tx")
			}
			stdTx = signedTx
		}

		return writeGenTxToFile(cdc, config, nodeID, stdTx)
	}
}

func writeGenTxToFile(cdc *codec.Codec, config *cfg.Config, nodeID string, stdTx chainTypes.StdTx) error {
	// Fetch output file name
	var err error

	outputDocument := viper.GetString(flags.FlagOutputDocument)
	if outputDocument == "" {
		outputDocument, err = makeOutputFilepath(config.RootDir, nodeID)
		if err != nil {
			return errors.Wrap(err, "failed to create output file path")
		}
	}

	if err := writeSignedGenTx(cdc, outputDocument, stdTx); err != nil {
		return errors.Wrap(err, "failed to write signed gen tx")
	}

	fmt.Fprintf(os.Stderr, "Genesis transaction written to %q\n", outputDocument)
	return nil
}

// GenTxCmd builds the application's gentx command.
func GenTxCmd(ctx *server.Context, cdc *codec.Codec,
	mbm module.BasicManager, smbh StakingMsgBuildingHelpers,
	genBalIterator types.GenesisBalancesIterator, defaultNodeHome, defaultCLIHome string,
	stakingFuncManager types.StakingFuncManager) *cobra.Command {
	ipDefault, _ := server.ExternalIP()
	fsCreateValidator, flagNodeID, flagPubKey, _, defaultsDesc := smbh.CreateValidatorMsgHelpers(ipDefault)

	cmd := &cobra.Command{
		Use:   "gentx [validator-operator-account] [validator-account-auth-address]",
		Short: "Generate a genesis tx carrying a self delegation",
		Args:  cobra.ExactArgs(2),
		Long: fmt.Sprintf(`This command is an alias of the 'tx create-validator' command'.

		It creates a genesis transaction to create a validator. 
		The following default parameters are included: 
		    %s`, defaultsDesc),

		RunE: genTxRunE(ctx, cdc, mbm, smbh, flagNodeID, flagPubKey),
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultCLIHome, "client's home directory")
	cmd.Flags().String(flags.FlagName, "", "name of private key with which to sign the gentx")
	cmd.Flags().String(flags.FlagOutputDocument, "",
		"write the genesis transaction JSON document to the given file instead of the default location")
	cmd.Flags().AddFlagSet(fsCreateValidator)
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().Bool(flagSignGenesisTrx, false, "If sign genesis trx in boot")
	viper.BindPFlag(flags.FlagKeyringBackend, cmd.Flags().Lookup(flags.FlagKeyringBackend))

	cmd.MarkFlagRequired(flags.FlagName)
	return cmd
}

func makeOutputFilepath(rootDir, nodeID string) (string, error) {
	writePath := filepath.Join(rootDir, "config", "gentx")
	if err := tmos.EnsureDir(writePath, 0700); err != nil {
		return "", err
	}
	return filepath.Join(writePath, fmt.Sprintf("gentx-%v.json", nodeID)), nil
}

func readUnsignedGenTxFile(cdc *codec.Codec, r io.Reader) (txutil.StdTx, error) {
	var stdTx txutil.StdTx
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return stdTx, err
	}
	err = cdc.UnmarshalJSON(bytes, &stdTx)
	return stdTx, err
}

func writeSignedGenTx(cdc *codec.Codec, outputDocument string, tx txutil.StdTx) error {
	outputFile, err := os.OpenFile(outputDocument, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	json, err := cdc.MarshalJSON(tx)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(outputFile, "%s\n", json)
	return err
}
