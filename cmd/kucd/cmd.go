package main

import (
	"encoding/json"
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tversion "github.com/tendermint/tendermint/version"
	"gopkg.in/yaml.v2"
)

// add server commands
func AddCommands(
	ctx *server.Context, cdc *codec.Codec,
	rootCmd *cobra.Command,
	appCreator server.AppCreator, appExport server.AppExporter) {
	rootCmd.PersistentFlags().String("log_level", ctx.Config.LogLevel, "Log level")

	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint subcommands",
	}

	tendermintCmd.AddCommand(
		server.ShowNodeIDCmd(ctx),
		server.ShowValidatorCmd(ctx),
		server.ShowAddressCmd(ctx),
		server.VersionCmd(ctx),
	)

	rootCmd.AddCommand(
		StartCmd(ctx, appCreator),
		server.UnsafeResetAllCmd(ctx),
		flags.LineBreak,
		tendermintCmd,
		server.ExportCmd(ctx, cdc, appExport),
		flags.LineBreak,
		version.Cmd,
	)
}

const (
	flagVersionSimple = "simple"
	flagJSONFormat    = "json"
)

func versionCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print kuchain version",
		Long: `Print kuchain version numbers
against which this app has been compiled.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				bs  []byte
				err error
			)

			if !viper.GetBool(flagVersionSimple) {
				vs := struct {
					Name                string `json:"name" yaml:"name"`
					Tendermint          string `json:"tendermint_version" yaml:"tendermint_version"`
					ABCI                string `json:"abci" yaml:"abci"`
					BlockProtocol       uint64 `json:"block_protocol" yaml:"block_protocol"`
					P2PProtocol         uint64 `json:"p2p_protocol" yaml:"p2p_protocol"`
					KuchainBuildVersion string `json:"version" yaml:"version"`
					KuchainBuildBranch  string `json:"branch" yaml:"branch"`
					KuchainBuildTime    string `json:"build_time" yaml:"build_time"`
					SDKVersion          string `json:"sdk_version" yaml:"sdk_version"`
					Commit              string `json:"commit" yaml:"commit"`
					BuildTags           string `json:"build_tags" yaml:"build_tags"`
				}{
					Name:                version.Name,
					Tendermint:          tversion.Version,
					ABCI:                tversion.ABCIVersion,
					BlockProtocol:       tversion.BlockProtocol.Uint64(),
					P2PProtocol:         tversion.P2PProtocol.Uint64(),
					KuchainBuildVersion: constants.KuchainBuildVersion,
					KuchainBuildBranch:  constants.KuchainBuildBranch,
					KuchainBuildTime:    constants.KuchainBuildTime,
					SDKVersion:          constants.KuchainBuildSDKVersion,
					Commit:              version.Commit,
					BuildTags:           version.BuildTags,
				}

				if viper.GetBool(flagJSONFormat) {
					bs, err = json.Marshal(&vs)
				} else {
					bs, err = yaml.Marshal(&vs)
				}
			} else {
				bs = []byte(constants.KuchainBuildVersion)
			}

			if err != nil {
				return err
			}

			fmt.Println(string(bs))
			return nil
		},
	}

	cmd.Flags().BoolP(flagVersionSimple, "s", false, "if just print simple version of kucd")
	cmd.Flags().BoolP(flagJSONFormat, "j", false, "print version info by json")

	return cmd
}
